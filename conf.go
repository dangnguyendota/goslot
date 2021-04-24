package goslot

type SymbolType int

const (
	REGULAR SymbolType = iota
	WILD
	BONUS
	SCATTER
)

type Conf struct {
	ColsSize                   int          `json:"cols_size"`                     // số cột
	ReelSize                   int          `json:"reel_size"`                     // số hàng tối đa
	RowsSize                   int          `json:"rows_size"`                     // số hàng hiển thị
	NumberOfNodes              int          `json:"number_of_nodes"`               // số lượng node sẽ run
	LocalPopulationSize        int          `json:"local_population_size"`         // số lượng tế bào trong 1 node
	LocalOptimizationEpochs    int64        `json:"local_optimization_epochs"`     // số năm tiến hóa của 1 tế bào trong 1 vòng đời
	NumberOfLifeCircle         int          `json:"number_of_life_circle"`         // số vòng đời của 1 node
	Targets                    []float64    `json:"targets"`                       // bảng tỉ lệ ăn
	Symbols                    []string     `json:"symbols"`                       // list các symbol
	Types                      []SymbolType `json:"types"`                         // kiểu của kí tự
	OutputFile                 string       `json:"output_file"`                   // file lưu kết quả
}

func (c *Conf) Validate() {
	if c.ColsSize <= 0 {
		panic("columns size must be more than 0")
	}

	if c.RowsSize <= 0 {
		panic("rows size must be  more than 0")
	}

	if c.ReelSize <= c.RowsSize {
		panic("reels size must be more than rows size")
	}

	if c.NumberOfNodes <= 0 {
		panic("number of nodes must be more than 0")
	}

	if c.LocalPopulationSize <= 2 {
		panic("local population size must be more than 2")
	}

	if c.LocalOptimizationEpochs <= 0 {
		panic("local optimization epochs must be more than 0")
	}

	if c.NumberOfLifeCircle <= 0 {
		panic("number of life circle must be more than 0")
	}

	if c.Targets == nil || len(c.Targets) == 0 {
		panic("invalid target")
	}

	if c.Symbols == nil || len(c.Symbols) == 0 {
		panic("invalid symbols")
	}

	if c.Types == nil || len(c.Types) != len(c.Symbols) {
		panic("types is nil or types length is not as same as symbols length")
	}
}
