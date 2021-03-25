package goslot

import (
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Generator struct {
	mux         sync.Mutex
	wait        sync.WaitGroup
	conf        *Conf
	global      *GeneticAlgorithm
	populations map[int]*GeneticAlgorithm
}

func NewGenerator(conf *Conf) *Generator {
	model := NewSlotModel(conf)
	ga := NewGeneticAlgorithm()
	ga.AddRandomReels(model, conf.LocalPopulationSize*conf.NumberOfNodes)

	return &Generator{
		conf:        conf,
		global:      ga,
		populations: make(map[int]*GeneticAlgorithm),
	}
}

func (g *Generator) Start() {
	for rank := 0; rank < g.conf.NumberOfNodes; rank++ {
		g.wait.Add(1)
		go func(rank int) {
			g.start(rank)
			g.wait.Done()
		}(rank)
	}
	g.wait.Wait()
}

func (g *Generator) GetBestChromosome() *Chromosome {
	return g.global.GetBestChromosome()
}

func (g *Generator) getGA(rank int) *GeneticAlgorithm {
	g.mux.Lock()
	defer g.mux.Unlock()
	if _, ok := g.populations[rank]; !ok || rand.Intn(int(g.conf.NumberOfLifeCircle/10)) == 0 {
		ga := NewGeneticAlgorithm()
		g.global.Subset(ga, g.conf.LocalPopulationSize)
		g.populations[rank] = ga
	}
	return g.populations[rank]
}

func (g *Generator) setGA(rank int, ga *GeneticAlgorithm) {
	g.mux.Lock()
	defer g.mux.Unlock()
	g.populations[rank] = ga
	if ga.GetBestFitness() < g.global.GetBestFitness() {
		g.global.AddChromosome(ga.GetBestChromosome())
	}
}

func (g *Generator) start(rank int) {
	model := NewSlotModel(g.conf)
	var counter int
	for {
		ga := g.getGA(rank)
		ga.Optimize(model, g.conf.LocalOptimizationEpochs)
		g.setGA(rank, ga)
		counter++
		if counter > g.conf.NumberOfLifeCircle {
			break
		}
	}
}
