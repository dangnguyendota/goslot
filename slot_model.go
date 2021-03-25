package goslot

const InvalidReelsPenalty = 10000

type ComputeResultFunc func(model Model) []float64

type Model interface {
	Win() int                       // trả về số tiền ăn được (tỉ lệ so với tiền cược là 1)
	Scatters() int                  // trả về số lượng Scatters trong view
	Bonus() int                     // trả về số lượng Bonus trong view
	Evaluate(reels [][]int) float64 // trả về epsilon cho reel  càng gần 0 thì càng chuẩn
}

type SlotModel struct {
	conf        *Conf
	computeFunc ComputeResultFunc
	reels       [][]int // reel hiện tại
	stops       []int   // vị trí của line ăn đang xét hiện tại
}

func NewModel(config *Conf, computeFunc ComputeResultFunc) Model {
	return &SlotModel{
		conf:        config,
		computeFunc: computeFunc,
		reels:       make([][]int, 0),
		stops:       make([]int, 0),
	}
}

// trả về tiền ăn được
func (s *SlotModel) Win() int {
	win := 0
	for _, payLine := range s.conf.PayLines {
		// lấy line tương ứng với payline này
		line := make([]int, s.conf.ColsSize)
		for i := 0; i < s.conf.ColsSize; i++ {
			line[i] = s.reels[i][s.stops[i]] + payLine[i]
		}

		// lấy biểu tượng đầu tiên (từ trái qua phải) khác WILD
		symbol := line[0]
		for i := 0; i < len(line); i++ {
			if s.conf.Types[symbol] != WILD {
				break
			}
			symbol = line[i]
		}

		// thay tất cả các WILD thành biểu tượng tìm được
		for i := 0; i < len(line); i++ {
			if s.conf.Types[line[i]] == WILD {
				line[i] = symbol
			}
		}

		// đếm từ trái qua phải xem có bao nhiêu symbol liên tiếp
		counter := 0
		for i := 0; i < len(line); i++ {
			if line[i] == symbol {
				counter++
			} else {
				break
			}
		}
		// tính tiền số lượng symbol đó
		win += s.conf.PayTable[counter][symbol]
	}
	return win
}

// trả về số lượng scatter trong view
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

//  tính số lượng bonus trong view
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

// tính tỉ lệ lệch so với mục đích muốn tỉ lệ ăn
func (s *SlotModel) Evaluate(reels [][]int) float64 {
	s.load(reels)
	s.init()
	parameters := s.calculate()
	sum := 0.0
	for i := 0; i < len(s.conf.Targets) && i < len(parameters); i++ {
		sum += (s.conf.Targets[i] - parameters[i]) * (s.conf.Targets[i] - parameters[i])
	}

	return sum
}

// số trường hợp có thể xảy ra
func (s *SlotModel) combinations() int64 {
	var result int64 = 1
	for i := 0; i < len(s.reels); i++ {
		result *= int64(len(s.reels[i]))
	}
	return result
}

// tải reel vào
func (s *SlotModel) load(reels [][]int) {
	s.reels = make([][]int, len(reels))
	for i := 0; i < len(s.reels); i++ {
		s.reels[i] = make([]int, len(reels[i]))
		copy(s.reels[i], reels[i])
	}
}

// khởi tạo các ô đang dừng ở đó
// và hàng ăn
func (s *SlotModel) init() {
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
func (s *SlotModel) next() {
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
}

// returns sum of result from computeFunc
func (s *SlotModel) calculate() []float64 {
	result := make([]float64, len(s.conf.Targets))
	for g := s.combinations() - 1; g >= 0; g-- {
		r := s.computeFunc(s)
		for i := 0; i < len(result); i++ {
			result[i] += r[i]
		}
		s.next()
	}

	for i := 0; i < len(result); i++ {
		result[i] /= float64(s.combinations())
	}
	return result
}
