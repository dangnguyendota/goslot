package goslot

import "fmt"

type Chromosome struct {
	conf    *Conf
	fitness float64
	reels   [][]int
}

func NewChromosome(reels [][]int, fitness float64, conf *Conf) *Chromosome {
	return &Chromosome{
		conf:    conf,
		fitness: fitness,
		reels:   reels,
	}
}

func (c *Chromosome) String() string {
	result := fmt.Sprintf("fitness: %f\n", c.fitness)
	for i := 0; i < c.conf.ReelSize; i++ {
		for j := 0; j < c.conf.RowsSize; j++ {
			result += c.conf.Symbols[c.reels[j][i]]
			result += " "
		}
		result += "\n"
	}
	return result
}