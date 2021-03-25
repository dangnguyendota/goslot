package goslot


func Example() {
	conf := &Conf{
		ColsSize:                5,
		ReelSize:                63,
		RowsSize:                3,
		NumberOfNodes:           2,
		LocalPopulationSize:     37,
		LocalOptimizationEpochs: 100,
		NumberOfLifeCircle:      67,
		Targets:                 []float64{},
		Symbols:                 []string{},
		Types:                   []SymbolType{},
		PayTable:                [][]int{},
		PayLines:                [][]int{},
	}
	conf.Validate()

	multipliers := []int{0, 0, 0, 1, 2, 3}

	generator := NewGenerator(conf, func(model Model) []float64 {
		result := make([]float64, 7)
		result[0] += float64(model.Win())

		switch model.Scatters() {
		case 0:
		case 1:
		case 2:
		case 3:
			result[1]++
			result[0] += float64(multipliers[3])
		case 4:
			result[2]++
			result[0] += float64(multipliers[4])
		case 5:
			result[3]++
			result[0] += float64(multipliers[5])
		default:
			result[0] += InvalidReelsPenalty
		}

		switch model.Bonus() {
		case 0:
		case 1:
		case 2:
		case 3:
			result[4]++
		case 4:
			result[5]++
		case 5:
			result[6]++
		default:
			result[0] += InvalidReelsPenalty
		}

		return result
	})
	generator.Start()
	println(generator.GetBestChromosome())
}
