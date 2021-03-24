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
	population  []*Chromosome
	resultIndex int
	firstIndex  int
	secondIndex int
	bestIndex   int
	worstIndex  int
}

func NewGeneticAlgorithm(populationSize int) *GeneticAlgorithm {
	return &GeneticAlgorithm{
		population:  make([]*Chromosome, populationSize),
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

func (g *GeneticAlgorithm) GetResultIndex() int {
	return g.resultIndex
}

func (g *GeneticAlgorithm) GetBestIndex() int {
	return g.bestIndex
}

func (g *GeneticAlgorithm) SetChromosome(chromosome *Chromosome, index int) {
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

func (g *GeneticAlgorithm) GetChromosome(index int) *Chromosome {
	if len(g.population) <= index || index <= -1 {
		panic("invalid index")
	}
	return g.population[index]
}

func (g *GeneticAlgorithm) GetBestChromosome() *Chromosome {
	return g.population[g.bestIndex]
}

func (g *GeneticAlgorithm) GetRandomChromosome() *Chromosome {
	return g.population[rand.Intn(len(g.population))]
}

func (g *GeneticAlgorithm) GetWorstChromosome() *Chromosome {
	return g.population[g.worstIndex]
}

func (g *GeneticAlgorithm) ReplaceWorst(chromosome *Chromosome) {
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

func (g *GeneticAlgorithm) SetFitness(fitness float64, index int) {
	if index < 0 {
		index = len(g.population) - 1
	}

	if len(g.population) <= index {
		panic("invalid index")
	}

	g.population[index].fitness = fitness
	if fitness < g.population[g.bestIndex].fitness {
		g.bestIndex = index
	}
	if fitness > g.population[g.worstIndex].fitness {
		g.worstIndex = index
	}
}

func (g *GeneticAlgorithm) GetFitness(index int) float64 {
	return g.population[index].fitness
}

func (g *GeneticAlgorithm) GetBestFitness() float64 {
	return g.population[g.bestIndex].fitness
}

func (g *GeneticAlgorithm) Size() int {
	return len(g.population)
}

func (g *GeneticAlgorithm) Subset(ga *GeneticAlgorithm, size int) {
	if len(g.population) <= 0 {
		return
	}
	for i := 0; i < size; i++ {
		ga.SetChromosome(g.GetRandomChromosome(), -1)
	}
}

func (g *GeneticAlgorithm) Selection() {
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

func (g *GeneticAlgorithm) Crossover() {
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

func (g *GeneticAlgorithm) Mutation() {
	index := rand.Intn(len(g.population))
	i := rand.Intn(len(g.population[g.resultIndex].reels))
	j := rand.Intn(len(g.population[g.resultIndex].reels[i]))
	g.population[g.resultIndex].reels[i][j] = g.population[index].reels[i][j]
	g.population[g.resultIndex].fitness = InvalidFitNessValue
}

func (g *GeneticAlgorithm) AddRandomReels(model *SlotModel, targets []float64, populationSize int) {
	for p := 0; p < populationSize; p++ {
		reels := make([][]int, model.config.ReelsSize)
		for i := 0; i < model.config.ReelsSize; i++ {
			for j := 0; j < model.config.ReelSize; j++ {
				value := rand.Intn(len(model.symbols) - 1) + 1
				reels[i] = append(reels[i], value)
			}
		}

		g.SetChromosome(NewChromosome(reels, InvalidReelsPenalty), -1)
		g.SetFitness(evaluate(model, targets, reels), -1)
	}
}

func (g *GeneticAlgorithm) Optimize(model *SlotModel, target []float64, epochs int64) {
	var e int64
	for e = 0; e < epochs*int64(g.Size()); e++ {
		g.Selection()
		g.Crossover()
		g.Mutation()
		index := g.GetResultIndex()
		g.SetFitness(evaluate(model, target, g.GetChromosome(index).reels), index)
	}
}

func (g *GeneticAlgorithm) String(model *SlotModel) string {
	result := ""
	result += fmt.Sprintf("number of populations %d\n\n", len(g.population))
	for p := 0; p < len(g.population); p++ {
		result += fmt.Sprintf("===> population: %d\n", p)
		for i := 0; i < len(g.population[p].reels); i++ {
			for j := 0; j < len(g.population[p].reels[i]); j++ {
				result += fmt.Sprintf("%s", model.symbols[g.population[p].reels[i][j]])
				result += " "
			}
			result += "\n"
		}

		result += "\n"
	}

	result = result[:len(result) - 1]
	return result
}

// tính tỉ lệ lệch so với mục đích muốn tỉ lệ ăn
func evaluate(model *SlotModel, target []float64, reels [][]int) float64 {
	model.Load(reels)
	model.Init()
	parameters := model.calculate()
	sum := 0.0
	for i := 0; i < len(target) && i < len(parameters); i++ {
		sum += (target[i] - parameters[i]) * (target[i] - parameters[i])
	}

	return sum
}