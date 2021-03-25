package goslot

type SlotType int

const (
	REGULAR SlotType = iota
	WILD
	BONUS
	SCATTER
)

type Conf struct {
	ReelsSize               int        // số cột
	ReelSize                int        // số hàng tối đa
	RowsSize                int        // số hàng hiển thị
	NumberOfNodes           int        // số lượng node sẽ run
	LocalPopulationSize     int        // số lượng tế bào trong 1 node
	LocalOptimizationEpochs int64      // số năm tiến hóa của 1 tế bào trong 1 vòng đời
	NumberOfLifeCircle      int        // số vòng đời của 1 node
	Targets                 []float64  // bảng tỉ lệ ăn
	Symbols                 []string   // list các symbol
	Types                   []SlotType // kiểu của kí tự
	PayTable                [][]int    // bảng tỉ lệ ăn
	Multipliers             []int      // bảng tỉ lệ nhân với scatter
	PayLines                [][]int	   // hàng ăn
}
