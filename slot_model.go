package goslot

const (
	//ReelsSize           int = 3
	//ReelSize                = 6
	//RowsSize                = 1
	InvalidReelsPenalty     = 10000
)
type ModelConfig struct {
	ReelsSize int
	ReelSize int
	RowsSize int
}

type SlotModel struct {
	config *ModelConfig
	reels       [][]int
	stops       []int
	line        []int
	symbols     []string
	types       []string
	payTable    [][]int
	multipliers []int
	lines       [][]int
}

func NewSlotModel(config *ModelConfig) *SlotModel {
	return &SlotModel{
		config: config,
		reels:       make([][]int, 0),
		stops:       make([]int, 0),
		line:        make([]int, 0),
		symbols:     make([]string, 0),
		types:       make([]string, 0),
		payTable:    make([][]int, 0),
		multipliers: make([]int, 0),
		lines:       make([][]int, 0),
	}
}

func (s *SlotModel) SetSymbols(symbols []string) {
	s.symbols = make([]string, len(symbols))
	copy(s.symbols, symbols)
}

func (s *SlotModel) SetTypes(types []string) {
	s.types = make([]string, len(types))
	copy(s.types, types)
}

func (s *SlotModel) SetPayTable(payTable [][]int) {
	s.payTable = make([][]int, len(payTable))
	for i := 0; i < len(payTable); i++ {
		s.payTable[i] = make([]int, len(payTable[i]))
		copy(s.payTable[i], payTable[i])
	}
}

func (s *SlotModel) SetMultipliers(multiplier []int) {
	s.multipliers = make([]int, len(multiplier))
	copy(s.multipliers, multiplier)
}

func (s *SlotModel) SetLines(lines [][]int) {
	s.lines = make([][]int, len(lines))
	for i := 0; i < len(lines); i++ {
		s.lines[i] = make([]int, len(lines[i]))
		copy(s.lines[i], lines[i])
	}
}

// tải reel vào
func (s *SlotModel) Load(values [][]int) {
	s.reels = make([][]int, len(values))
	for i := 0; i < len(s.reels); i++ {
		s.reels[i] = make([]int, len(values[i]))
		copy(s.reels[i], values[i])
	}
}

// số ô trong reel
func (s *SlotModel) Combinations() int64 {
	var result int64 = 1
	for i := 0; i < len(s.reels); i++ {
		result *= int64(len(s.reels[i]))
	}
	return result
}

// khởi tạo các ô đang dừng ở đó
// và hàng ăn
func (s *SlotModel) Init() {
	s.stops = make([]int, len(s.reels))
	for i := 0; i < len(s.reels); i++ {
		s.stops[i] = 0
	}
	s.line = make([]int, len(s.reels))
	for i := 0; i < len(s.reels); i++ {
		s.line[i] = s.reels[i][s.stops[i]]
	}
}

// chuyển sang line khác
func (s *SlotModel) Next() {
	s.stops[len(s.reels)-1] += 1
	for i := len(s.reels) - 1; i > 0; i-- {
		if s.stops[i] >= len(s.reels[i]) {
			s.stops[i] = 0
			s.stops[i-1] += 1
		}
	}
	if s.stops[0] >= len(s.reels[0]) {
		s.stops[0] = 0
	}
	s.line = make([]int, len(s.reels))
	for i := 0; i < len(s.reels); i++ {
		s.line[i] = s.reels[i][s.stops[i]]
	}
}

// trả về tiền ăn được
func (s *SlotModel) Win() int {
	symbol := s.line[0]
	for i := 0; i < len(s.line); i++ {
		if symbol != 1 {
			break
		}
		symbol = s.line[i]
	}
	for i := 0; i < len(s.line); i++ {
		if s.line[i] == 1 {
			s.line[i] = symbol
		}
	}

	counter := 0
	for i := 0; i < len(s.line); i++ {
		if s.line[i] == symbol {
			counter++
		} else {
			break
		}
	}
	return s.payTable[counter][symbol]
}

// trả về số lượng scatter
func (s *SlotModel) Scatters() int {
	counter := 0
	for i := 0; i < len(s.reels); i++ {
		for j := 0; j < s.config.RowsSize; j++ {
			if s.reels[i][(s.stops[i]+j)%len(s.reels[i])] == 16 {
				counter++
			}
		}
	}
	return counter
}

//  tính số lượng bonus
func (s *SlotModel) Bonus() int {
	counter := 0
	for i := 0; i < len(s.reels); i++ {
		for j := 0; j < s.config.RowsSize; j++ {
			if s.reels[i][(s.stops[i]+j)%len(s.reels[i])] == 15 {
				counter++
			}
		}
	}
	return counter
}

// returns [ăn được bao nhiêu, số case ăn 1 scatter, số case ăn 2 scatter,
//			số case ăn 3 scatter, số case ăn 1 bonus, số case ăn 2 bonus,
//			số case ăn 3 bonus]
func (s *SlotModel) calculate() []float64 {
	result := make([]float64, 7)
	// RTP
	result[0] = 0
	// free games frequency (scatters as activator).
	// activation can be from 3, 4 or 5 scatters
	result[1] = 0
	result[2] = 0
	result[3] = 0
	// bonus game frequency (bonus symbols as activator)
	// activation can be from 3, 4 or 5 bonus symbols
	result[4] = 0
	result[5] = 0
	result[6] = 0
	for g := s.Combinations() - 1; g >= 0; g-- {
		result[0] += float64(s.Win())

		switch s.Scatters() {
		case 0:
		case 1:
		case 2:
		case 3:
			result[1]++
			result[0] += float64(s.multipliers[3])
		case 4:
			result[2]++
			result[0] += float64(s.multipliers[4])
		case 5:
			result[3]++
			result[0] += float64(s.multipliers[5])
		default:
			result[0] += InvalidReelsPenalty
		}

		switch s.Bonus() {
		case 0:
		case 1:
		case 2:
		case 3:
			result[4]++
		case 4:
			result[5]++
		case 5:
			result[6]++
		default:
			result[0] += InvalidReelsPenalty
		}

		s.Next()
	}

	for i := 0; i < len(result); i++ {
		result[i] /= float64(s.Combinations())
	}
	return result
}
