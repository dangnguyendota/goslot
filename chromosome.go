package goslot

import (
	"bytes"
	"fmt"
)

type Chromosome struct {
	fitness float64
	reels   [][]int
}

func NewChromosome(reels [][]int, fitness float64) *Chromosome {
	return &Chromosome{
		fitness: fitness,
		reels:   reels,
	}
}

func (c *Chromosome) Reels() [][]int {
	return c.reels
}

func ChromosomeString(c *Chromosome, symbols []string) string {
	result := fmt.Sprintf("===> fitness: %f\n", c.fitness)
	maxSymbolSize := 0
	for _, s := range symbols {
		if len(s) > maxSymbolSize {
			maxSymbolSize = len(s)
		}
	}

	for i := 0; i < len(c.reels[0]); i++ {
		result += path("-", maxSymbolSize, len(c.reels)) + "\n|"
		for j := 0; j < len(c.reels); j++ {
			result += centerString(symbols[c.reels[j][i]], maxSymbolSize)
			result += "|"
		}
		result += "\n"
	}
	result += path("-", maxSymbolSize, len(c.reels))
	result += "\n"
	for i := 0; i < len(c.reels); i++ {
		for j := 0; j < len(c.reels[i]); j++ {
			result += symbols[c.reels[i][j]]
			result+= "."
		}
	}
	result += "\n"
	return result
}

func (c *Chromosome) Code(symbols []string) string {
	result := ""
	for i := 0; i < len(c.reels); i++ {
		for j := 0; j < len(c.reels[i]); j++ {
			result += symbols[c.reels[i][j]]
			result+= "."
		}
	}
	return result
}

func path(s string, length int, size int) string {
	ss := ""
	for i := 0; i < length; i++ {
		ss += s
	}
	result := ""
	for i := 0; i < size; i++ {
		result += "+" + ss
	}
	result += "+"
	return result
}

func centerString(str string, totalFieldWidth int) string {

	strLen := len(str)
	spacesToPad := totalFieldWidth - strLen
	var tmpSpaces float64
	var lrSpaces int

	tmpSpaces = float64(spacesToPad) / 2
	lrSpaces = int(tmpSpaces)

	buffer := bytes.NewBufferString("")

	spacesRemaining := totalFieldWidth

	for i := 0; i < lrSpaces; i++ {
		buffer.WriteString(" ")
		spacesRemaining = spacesRemaining - 1
	}
	buffer.WriteString(str)
	spacesRemaining = spacesRemaining - strLen
	for i := spacesRemaining; i > 0; i-- {
		buffer.WriteString(" ")
	}

	return buffer.String()
}