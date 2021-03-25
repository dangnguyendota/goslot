package goslot

const InvalidReelsPenalty = 10000

type Model interface {
	ValidLine(line []int)
}

type SlotModel struct {
	conf  *Conf
	reels [][]int // reel hiện tại
	stops []int   // vị trí của line ăn đang xét hiện tại
	//line  []int    // line ăn hiện tại
}

func NewSlotModel(config *Conf) *SlotModel {
	return &SlotModel{
		conf:  config,
		reels: make([][]int, 0),
		stops: make([]int, 0),
		//line:        make([]int, 0),
	}
}

// tải reel vào
func (s *SlotModel) Load(reels [][]int) {
	s.reels = make([][]int, len(reels))
	for i := 0; i < len(s.reels); i++ {
		s.reels[i] = make([]int, len(reels[i]))
		copy(s.reels[i], reels[i])
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
	//s.line = make([]int, len(s.reels))
	//for i := 0; i < len(s.reels); i++ {
	//	s.line[i] = s.reels[i][s.stops[i]]
	//}
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

	//s.line = make([]int, len(s.reels))
	//for i := 0; i < len(s.reels); i++ {
	//	s.line[i] = s.reels[i][s.stops[i]]
	//}
}

// trả về tiền ăn được
func (s *SlotModel) Win() int {
	win := 0
	for _, payLine := range s.conf.PayLines {
		line := make([]int, s.conf.ReelsSize)
		for i := 0; i < s.conf.ReelsSize; i++ {
			line[i] = s.reels[i][s.stops[i]] + payLine[i]
		}

		symbol := line[0]
		for i := 0; i < len(line); i++ {
			if s.conf.Types[symbol] != WILD {
				break
			}
			symbol = line[i]
		}
		for i := 0; i < len(line); i++ {
			if line[i] == 1 {
				line[i] = symbol
			}
		}

		counter := 0
		for i := 0; i < len(line); i++ {
			if line[i] == symbol {
				counter++
			} else {
				break
			}
		}
		win += s.conf.PayTable[counter][symbol]
	}
	return win
}

// trả về số lượng scatter
func (s *SlotModel) Scatters() int {
	counter := 0
	for i := 0; i < len(s.reels); i++ {
		for j := 0; j < s.conf.RowsSize; j++ {
			if s.conf.Types[s.reels[i][(s.stops[i]+j)%len(s.reels[i])]] == SCATTER {
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
		for j := 0; j < s.conf.RowsSize; j++ {
			if s.conf.Types[s.reels[i][(s.stops[i]+j)%len(s.reels[i])]] == BONUS {
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
			result[0] += float64(s.conf.Multipliers[3])
		case 4:
			result[2]++
			result[0] += float64(s.conf.Multipliers[4])
		case 5:
			result[3]++
			result[0] += float64(s.conf.Multipliers[5])
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

// tính tỉ lệ lệch so với mục đích muốn tỉ lệ ăn
func (s *SlotModel) Evaluate(reels [][]int) float64 {
	s.Load(reels)
	s.Init()
	parameters := s.calculate()
	sum := 0.0
	for i := 0; i < len(s.conf.Targets) && i < len(parameters); i++ {
		sum += (s.conf.Targets[i] - parameters[i]) * (s.conf.Targets[i] - parameters[i])
	}

	return sum
}
