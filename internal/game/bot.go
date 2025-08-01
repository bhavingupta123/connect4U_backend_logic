package game

import (
	"math"
)

func (g *Game) BotBestMove() int {
	bestScore := math.MinInt
	bestCol := -1

	for col := 0; col < Cols; col++ {
		if !g.Board.IsValidMove(col) {
			continue
		}

		row := g.Board.ApplyMove(col, 1)
		if g.Board.CheckWin(col, row, 1) {
			g.Board.ResetCell(col, row)
			return col
		}
		g.Board.ResetCell(col, row)

		rowOpp := g.Board.ApplyMove(col, 0)
		if g.Board.CheckWin(col, rowOpp, 0) {
			g.Board.ResetCell(col, rowOpp)
			return col
		}
		g.Board.ResetCell(col, rowOpp)

		score := evaluateColumn(col)
		if score > bestScore {
			bestScore = score
			bestCol = col
		}
	}

	if bestCol == -1 {
		for col := 0; col < Cols; col++ {
			if g.Board.IsValidMove(col) {
				return col
			}
		}
	}
	return bestCol
}

func evaluateColumn(col int) int {
	center := Cols / 2
	return 3 - abs(col-center)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
