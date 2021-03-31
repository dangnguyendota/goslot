package goslot

import (
	"fmt"
	"math/rand"
)

const KeepElite bool = true
const (
	CrossoverResultIntoBestPercent   int = 1
	CrossoverResultIntoMiddlePercent     = 9
	CrossoverResultIntoWorstPercent      = 90
	InvalidFitNessValue                  = 2147483647
)

type GeneticAlgorithm struct {
	conf        *Conf
	population  []*Chromosome
	resultIndex int
	firstIndex  int
	secondIndex int
	bestIndex   int
	worstIndex  int
}

func NewGeneticAlgorithm(conf *Conf) *GeneticAlgorithm {
	return &GeneticAlgorithm{
		conf: conf,
		population:  []*Chromosome{},
		resultIndex: 0,
		firstIndex:  0,
		secondIndex: 0,
		bestIndex:   0,
		worstIndex:  0,
	}
}

func (g *GeneticAlgorithm) selectRandom() {
	for {
		g.resultIndex = rand.Int() % len(g.population)
		g.firstIndex = rand.Int() % len(g.population)
		g.secondIndex = rand.Int() % len(g.population)
		if !(g.resultIndex == g.firstIndex || g.resultIndex == g.secondIndex ||
			(g.resultIndex == g.bestIndex && KeepElite)) {
			break
		}
	}
}

func (g *GeneticAlgorithm) getResultIndex() int {
	return g.resultIndex
}

func (g *GeneticAlgorithm) getBestIndex() int {
	return g.bestIndex
}

func (g *GeneticAlgorithm) addChromosome(chromosome *Chromosome) {
	g.setChromosome(chromosome, -1)
}

func (g *GeneticAlgorithm) setChromosome(chromosome *Chromosome, index int) {
	if index < 0 {
		g.population = append(g.population, chromosome)
		index = len(g.population) - 1
	} else if index < len(g.population) {
		g.population[index] = chromosome
	}

	if g.population[index].fitness < g.population[g.bestIndex].fitness {
		g.bestIndex = index
	}
	if g.population[index].fitness > g.population[g.worstIndex].fitness {
		g.worstIndex = index
	}
}

func (g *GeneticAlgorithm) getChromosome(index int) *Chromosome {
	if len(g.population) <= index || index <= -1 {
		panic("invalid index")
	}
	return g.population[index]
}

func (g *GeneticAlgorithm) getBestChromosome() *Chromosome {
	return g.population[g.bestIndex]
}

func (g *GeneticAlgorithm) getRandomChromosome() *Chromosome {
	return g.population[rand.Intn(len(g.population))]
}

func (g *GeneticAlgorithm) getWorstChromosome() *Chromosome {
	return g.population[g.worstIndex]
}

func (g *GeneticAlgorithm) replaceWorst(chromosome *Chromosome) {
	g.population[g.worstIndex] = chromosome
	g.bestIndex = 0
	g.worstIndex = 0
	for i := 0; i < len(g.population); i++ {
		if g.population[i].fitness < g.population[g.bestIndex].fitness {
			g.bestIndex = i
		}
		if g.population[i].fitness > g.population[g.worstIndex].fitness {
			g.worstIndex = i
		}
	}
}

func (g *GeneticAlgorithm) addFitness(fitness float64) {
	g.setFitness(fitness, len(g.population)-1)
}
func (g *GeneticAlgorithm) setFitness(fitness float64, index int) {
	g.population[index].fitness = fitness
	if fitness < g.population[g.bestIndex].fitness {
		g.bestIndex = index
	}
	if fitness > g.population[g.worstIndex].fitness {
		g.worstIndex = index
	}
}

func (g *GeneticAlgorithm) getFitness(index int) float64 {
	return g.population[index].fitness
}

func (g *GeneticAlgorithm) getBestFitness() float64 {
	return g.population[g.bestIndex].fitness
}

func (g *GeneticAlgorithm) size() int {
	return len(g.population)
}

func (g *GeneticAlgorithm) subset(ga *GeneticAlgorithm, size int) {
	if len(g.population) <= 0 {
		return
	}
	for i := 0; i < size; i++ {
		ga.setChromosome(g.getRandomChromosome(), -1)
	}
}

func (g *GeneticAlgorithm) selection() {
	var percent int
	percent = rand.Int() % (CrossoverResultIntoWorstPercent +
		CrossoverResultIntoMiddlePercent +
		CrossoverResultIntoBestPercent)
	if percent < CrossoverResultIntoWorstPercent {
		for {
			g.selectRandom()
			if !(g.population[g.resultIndex].fitness < g.population[g.firstIndex].fitness ||
				g.population[g.resultIndex].fitness < g.population[g.secondIndex].fitness) {
				break
			}
		}
	} else if percent < CrossoverResultIntoWorstPercent+CrossoverResultIntoMiddlePercent {
		for {
			g.selectRandom()
			if !(g.population[g.resultIndex].fitness < g.population[g.firstIndex].fitness ||
				g.population[g.resultIndex].fitness < g.population[g.secondIndex].fitness) {
				break
			}
		}
	} else if percent < CrossoverResultIntoWorstPercent+CrossoverResultIntoMiddlePercent+CrossoverResultIntoBestPercent {
		for {
			g.selectRandom()
			if !(g.population[g.resultIndex].fitness > g.population[g.firstIndex].fitness ||
				g.population[g.resultIndex].fitness > g.population[g.secondIndex].fitness) {
				break
			}
		}
	}
}

func (g *GeneticAlgorithm) crossover() {
	a := g.population[g.firstIndex].reels
	b := g.population[g.secondIndex].reels
	c := g.population[g.resultIndex].reels
	for i := 0; i < len(a) && i < len(b) && i < len(c); i++ {
		for j := 0; j < len(a[i]) && j < len(b[i]) && j < len(c[i]); j++ {
			if rand.Int()%2 == 0 {
				c[i][j] = a[i][j]
			} else {
				c[i][j] = b[i][j]
			}
		}
	}

	g.population[g.resultIndex].fitness = InvalidFitNessValue
}

func (g *GeneticAlgorithm) mutation() {
	index := rand.Intn(len(g.population))
	i := rand.Intn(len(g.population[g.resultIndex].reels))
	j := rand.Intn(len(g.population[g.resultIndex].reels[i]))
	g.population[g.resultIndex].reels[i][j] = g.population[index].reels[i][j]
	g.population[g.resultIndex].fitness = InvalidFitNessValue
}

func (g *GeneticAlgorithm) addRandomReels(model Model, populationSize int) {
	for p := 0; p < populationSize; p++ {
		reels := make([][]int, g.conf.ColsSize)
		for i := 0; i < g.conf.ColsSize; i++ {
			for j := 0; j < g.conf.ReelSize; j++ {
				value := rand.Intn(len(g.conf.Symbols))
				reels[i] = append(reels[i], value)
			}
		}

		g.addChromosome(NewChromosome(reels, InvalidReelsPenalty))
		g.addFitness(model.Evaluate(reels))
	}
}

func (g *GeneticAlgorithm) optimize(model Model, epochs int64) {
	var e int64
	for e = 0; e < epochs*int64(g.size()); e++ {
		g.selection()
		g.crossover()
		g.mutation()
		index := g.getResultIndex()
		g.setFitness(model.Evaluate(g.getChromosome(index).reels), index)
	}
}

func (g *GeneticAlgorithm) String() string {
	result := ""
	result += fmt.Sprintf("number of populations %d\n\n", len(g.population))
	for p := 0; p < len(g.population); p++ {
		result += fmt.Sprintf("===> population: %d\n", p)
		for i := 0; i < len(g.population[p].reels); i++ {
			for j := 0; j < len(g.population[p].reels[i]); j++ {
				result += fmt.Sprintf("%s", g.conf.Symbols[g.population[p].reels[i][j]])
				result += " "
			}
			result += "\n"
		}

		result += "\n"
	}

	result = result[:len(result)-1]
	return result
}
