package goslot

import (
	"fmt"
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
	g.global = NewGeneticAlgorithm(g.conf)
	g.global.addRandomReels(NewMachine(g.conf, g.model), g.conf.LocalPopulationSize*g.conf.NumberOfNodes)
	g.StoreChromosome(g.global.getBestChromosome(), true)
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
		g.StoreChromosome(ga.getBestChromosome(), true)
	} else {
		g.StoreChromosome(ga.getBestChromosome(), false)
	}
}

func (g *Generator) StoreChromosome(c *Chromosome, best bool) {
	data := fmt.Sprintf("\n%s\n", ChromosomeString(c, g.conf.Symbols))
	if best {
		data = "BEST ==>" + data
	}
	println(data)
	if err := g.WriteFile([]byte(data)); err != nil {
		panic(err)
	}
}

func (g *Generator) start(rank int) {
	var counter int
	for {
		ga := g.getGA(rank)
		ga.optimize(NewMachine(g.conf, g.model), g.conf.LocalOptimizationEpochs)
		g.setGA(rank, ga)
		counter++
		println("rank:", rank, "counter:", counter)
		if counter > g.conf.NumberOfLifeCircle {
			break
		}
	}
}

func (g *Generator) WriteFile(data []byte) error {
	f, err := os.OpenFile(g.conf.OutputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
