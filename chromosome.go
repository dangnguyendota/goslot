package goslot

type Chromosome struct {
	fitness float64
	reels [][]int
}

func NewChromosome(reels [][]int, fitness float64) *Chromosome {
	return &Chromosome{
		fitness: fitness,
		reels:   reels,
	}
}

