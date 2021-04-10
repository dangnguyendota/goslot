package goslot

import "fmt"

type Model interface {
	Win(machine *SlotMachine) int          // trả về số tiền ăn được (tỉ lệ so với tiền cược là 1)
	Scatters(machine *SlotMachine) int     // trả về số lượng Scatters trong view
	Bonus(machine *SlotMachine) int        // trả về số lượng Bonus trong view
	Result(machine *SlotMachine) []float64 // trả về tiền ăn cho đường ăn hiện tại
}

type model struct {
	conf        *Conf
	paylines    [][]int
	paytable    [][]int
	multipliers []float64
}

// base model dành cho các loại slot phổ thông với target là mảng 7 phần từ 
// target = [RPT, 3, 4, 5 Scatter, 3, 4, 5 Bonus] 
func NewBaseModel(conf *Conf, paylines [][]int, paytable [][]int, multipliers []float64) Model {
	if paylines == nil || len(paylines) == 0 {
		panic("invalid pay table")
	}

	for i := range paylines {
		if paylines[i] == nil || len(paylines[i]) != conf.ColsSize {
			panic(fmt.Sprintf("invalid pay lines or row size at %d is not %d", i, conf.ColsSize))
		}
		for j := range paylines[i] {
			if paylines[i][j] < 0 || paylines[i][j] >= conf.RowsSize {
				panic(fmt.Sprintf("invalid pay lines value, must be positive and less than %d", conf.RowsSize))
			}
		}
	}

	if paytable == nil || len(paytable) != conf.ColsSize+1 {
		panic(fmt.Sprintf("invalid pay table or paytable size (paytable size = number of columns + 1)"))
	}

	for i := range paytable {
		if paytable[i] == nil || len(paytable[i]) != len(conf.Symbols) {
			panic(fmt.Sprintf("invalid pay table at %d (size must equals number of symbols)", i))
		}
	}
	return &model{
		conf:        conf,
		paylines:    paylines,
		paytable:    paytable,
		multipliers: multipliers,
	}
}

func (m *model) Win(machine *SlotMachine) int {
	win := 0
	for _, payLine := range m.paylines {
		// lấy line tương ứng với payline này
		line := make([]int, m.conf.ColsSize)
		for i := 0; i < m.conf.ColsSize; i++ {
			line[i] = machine.reels[i][(machine.stops[i]+payLine[i])%m.conf.ReelSize]
		}

		// lấy biểu tượng đầu tiên (từ trái qua phải) khác WILD
		symbol := line[0]
		for i := 0; i < len(line); i++ {
			if m.conf.Types[symbol] != WILD {
				break
			}
			symbol = line[i]
		}

		// thay tất cả các WILD thành biểu tượng tìm được
		for i := 0; i < len(line); i++ {
			if m.conf.Types[line[i]] == WILD {
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
		win += m.paytable[counter][symbol]
	}
	return win
}

func (m *model) Scatters(machine *SlotMachine) int {
	counter := 0
	for i := 0; i < len(machine.reels); i++ {
		for j := 0; j < m.conf.RowsSize; j++ {
			if m.conf.Types[machine.reels[i][(machine.stops[i]+j)%len(machine.reels[i])]] == SCATTER {
				counter++
			}
		}
	}
	return counter
}

func (m *model) Bonus(machine *SlotMachine) int {
	counter := 0
	for i := 0; i < len(machine.reels); i++ {
		for j := 0; j < m.conf.RowsSize; j++ {
			if m.conf.Types[machine.reels[i][(machine.stops[i]+j)%len(machine.reels[i])]] == BONUS {
				counter++
			}
		}
	}
	return counter
}

func (m *model) Result(machine *SlotMachine) []float64 {
	result := make([]float64, 7)
	result[0] += float64(m.Win(machine))

	switch m.Scatters(machine) {
	case 0:
	case 1:
	case 2:
	case 3:
		result[1]++
		result[0] += m.multipliers[3]
	case 4:
		result[2]++
		result[0] += m.multipliers[4]
	case 5:
		result[3]++
		result[0] += m.multipliers[5]
	default:
		result[0] += InvalidReelsPenalty
	}

	switch m.Bonus(machine) {
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

	return result
}
