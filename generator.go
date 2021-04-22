package goslot

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
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
	model       Model
}

func NewGenerator(conf *Conf, model Model) *Generator {
	return &Generator{
		conf:        conf,
		model:       model,
		populations: make(map[int]*GeneticAlgorithm),
	}
}

func (g *Generator) Start() {
	ga := NewGeneticAlgorithm(g.conf)
	ga.addRandomReels(NewMachine(g.conf, g.model), g.conf.LocalPopulationSize*g.conf.NumberOfNodes)
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
	return g.global.getBestChromosome()
}

func (g *Generator) getGA(rank int) *GeneticAlgorithm {
	g.mux.Lock()
	defer g.mux.Unlock()
	if _, ok := g.populations[rank]; !ok || rand.Intn(int(g.conf.NumberOfLifeCircle/10)) == 0 {
		ga := NewGeneticAlgorithm(g.conf)
		g.global.subset(ga, g.conf.LocalPopulationSize)
		g.populations[rank] = ga
	}
	return g.populations[rank]
}

func (g *Generator) setGA(rank int, ga *GeneticAlgorithm) {
	g.mux.Lock()
	defer g.mux.Unlock()
	g.populations[rank] = ga
	if ga.getBestFitness() < g.global.getBestFitness() {
		g.global.addChromosome(ga.getBestChromosome())
		data := []byte(fmt.Sprintf("Found best chromosome with fitness: %f\n\n%s\n\n",
			ga.getBestChromosome().fitness,
			ChromosomeString(ga.getBestChromosome(), g.conf.Symbols)))
		if err := ioutil.WriteFile(g.conf.OutputFile, data, os.ModeAppend); err != nil {
			panic(err)
		}
	}
}

func (g *Generator) start(rank int) {
	//model := NewModel(g.conf, g.computeFunc)
	var counter int
	for {
		ga := g.getGA(rank)
		ga.optimize(NewMachine(g.conf, g.model), g.conf.LocalOptimizationEpochs)
		g.setGA(rank, ga)
		counter++
		//println("rank:", rank, "counter:", counter)
		if counter > g.conf.NumberOfLifeCircle {
			break
		}
	}
}
