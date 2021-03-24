package goslot

import (
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
//LocalPopulationSize     = 37
//LocalOptimizationEpochs = 100
//NumberOfBroadcasts      = 67
//RootNode                = 0
)

type NodeConfig struct {
	NumberOfNodes           int
	LocalPopulationSize     int
	LocalOptimizationEpochs int64
	NumberOfBroadcasts      int64
	RootNode                int
	Targets                 []float64
}

type NodeManager struct {
	nodes map[int]*Node
}

func NewNodeManager() *NodeManager {
	return &NodeManager{
		nodes: make(map[int]*Node),
	}
}

func (n *NodeManager) Register(node *Node) {
	n.nodes[node.rank] = node
	node.send = func(rank int, ga *GeneticAlgorithm) {
		println("[START] node", node.rank, "send to node", rank)
		if recvNode, ok := n.nodes[rank]; ok {
			recvNode.mux.Lock()
			if _, ok := recvNode.recvChan[node.rank]; !ok {
				recvNode.recvChan[node.rank] = make(chan *GeneticAlgorithm, 1)
			}
			recvChan := recvNode.recvChan[rank]
			recvNode.mux.Unlock()
			select {
			case recvChan <- ga:
			}
		}
		println("[STOP] node", node.rank, "send to node", rank)
	}

	node.recv = func(rank int) *GeneticAlgorithm {
		println("[START] node", node.rank, "recv to node", rank)
		node.mux.Lock()
		if _, ok := node.recvChan[rank]; !ok {
			node.recvChan[rank] = make(chan *GeneticAlgorithm, 1)
		}
		recvChan := node.recvChan[rank]
		node.mux.Unlock()
		select {
		case ga, ok := <-recvChan:
			if ok {
				return ga
			}
		}
		println("[STOP] node", node.rank, "recv to node", rank)
		return nil
	}
}

type Node struct {
	mux      sync.Mutex
	config   *NodeConfig
	rank     int
	model    *SlotModel
	recvChan map[int]chan *GeneticAlgorithm
	send     func(rank int, ga *GeneticAlgorithm)
	recv     func(rank int) *GeneticAlgorithm
}

func NewNode(rank int, config *NodeConfig, model *SlotModel) *Node {
	return &Node{
		rank:     rank,
		config:   config,
		model:    model,
		recvChan: make(map[int]chan *GeneticAlgorithm),
		send:     nil,
		recv:     nil,
	}
}

func (n *Node) Start() {
	n.master()
	n.slave()
}

func (n *Node) master() {
	var counter int64
	if n.rank != n.config.RootNode {
		return
	}

	var global = NewGeneticAlgorithm(0)
	var populations = make(map[int]*GeneticAlgorithm, 0)
	for {
		println("Round:", counter+1)
		for r := 0; r < n.config.NumberOfNodes; r++ {
			if r != n.config.RootNode {
				continue
			}
			if counter == 0 {
				global.AddRandomReels(n.model, n.config.Targets, n.config.LocalPopulationSize*n.config.NumberOfNodes)
				var ga = NewGeneticAlgorithm(0)
				global.Subset(ga, n.config.LocalPopulationSize)
				populations[r] = ga
			} else {
				if rand.Intn(int(n.config.NumberOfBroadcasts/10)) == 0 {
					var ga = NewGeneticAlgorithm(0)
					global.Subset(ga, n.config.LocalPopulationSize)
					populations[r] = ga
				}
			}
			n.send(r, populations[r])
		}

		for r := 0; r < n.config.NumberOfNodes; r++ {
			if r == n.config.RootNode {
				continue
			}
			var ga *GeneticAlgorithm
			ga = n.recv(r)
			populations[r] = ga
			if ga.GetBestFitness() < global.GetBestFitness() {
				global.SetChromosome(ga.GetBestChromosome(), -1)
			}
			println("Worker", r, ga.GetBestChromosome().fitness)
		}

		println("Global : ", global.GetBestChromosome().fitness)
		counter++
		if counter > n.config.NumberOfBroadcasts {
			break
		}
	}
	println("done master")
}

func (n *Node) slave() {
	var counter int64
	if n.rank == n.config.RootNode {
		return
	}
	for {
		var ga *GeneticAlgorithm
		ga = n.recv(n.config.RootNode)
		ga.Optimize(n.model, n.config.Targets, n.config.LocalOptimizationEpochs)
		n.send(n.config.RootNode, ga)
		counter++
		if counter > n.config.NumberOfBroadcasts {
			break
		}
	}
	println("done slave", n.rank)
}
