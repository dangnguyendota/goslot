package goslot

import (
	"fmt"
	"math"
)

const InvalidReelsPenalty = 10000

type SlotMachine struct {
	conf  *Conf
	reels [][]int // reel hiện tại
	stops []int   // vị trí của line ăn đang xét hiện tại
	model Model
}

func NewMachine(config *Conf, model Model) *SlotMachine {
	return &SlotMachine{
		conf:  config,
		model: model,
		reels: make([][]int, 0),
		stops: make([]int, 0),
	}
}

func (s *SlotMachine) Reels() [][]int {
	return s.reels
}

func (s *SlotMachine) Stops() []int {
	return s.stops
}

func (s *SlotMachine) Model() Model {
	return s.model
}

func (s *SlotMachine) Conf() *Conf {
	return s.conf
}

// tính tỉ lệ lệch so với mục đích muốn tỉ lệ ăn
func (s *SlotMachine) Evaluate(reels [][]int) float64 {
	s.load(reels)
	s.init()
	parameters := s.calculate()
	var sum float64
	for i := 0; i < len(s.conf.Targets) && i < len(parameters); i++ {
		println(fmt.Sprintf("%f", parameters[i]))
		sum += (s.conf.Targets[i] - parameters[i]) * (s.conf.Targets[i] - parameters[i])
	}

	return sum
}

// số trường hợp có thể xảy ra
func (s *SlotMachine) combinations() int64 {
	var result int64 = 1
	for i := 0; i < len(s.reels); i++ {
		result *= int64(len(s.reels[i]))
	}
	return result
}

// tải reel vào
func (s *SlotMachine) load(reels [][]int) {
	s.reels = make([][]int, len(reels))
	for i := 0; i < len(s.reels); i++ {
		s.reels[i] = make([]int, len(reels[i]))
		copy(s.reels[i], reels[i])
	}
}

// khởi tạo các ô đang dừng ở đó
// và hàng ăn
func (s *SlotMachine) init() {
	s.stops = make([]int, len(s.reels))
	for i := 0; i < len(s.reels); i++ {
		s.stops[i] = 0
	}
}

// chuyển sang line khác
func (s *SlotMachine) next() {
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
func (s *SlotMachine) calculate() []float64 {
	result := make([]float64, len(s.conf.Targets))
	for g := s.combinations() - 1; g >= 0; g-- {
		r := s.model.Result(s)
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

func (s *SlotMachine) Compute(reels [][]int) map[int64][]float64 {
	s.load(reels)
	s.init()
	m := map[int64][]float64{}
	for g := s.combinations() - 1; g >= 0; g-- {
		m[s.StopsNumber()] = s.model.Result(s)
		s.next()
	}
	return m
}

func (s *SlotMachine) StopsNumber() int64 {
	var n int64
	for i := 0; i < s.conf.ColsSize; i++ {
		n += int64(s.stops[s.conf.ColsSize-1-i]) * int64(math.Pow10(2*i))
	}
	return n
}
