package game

import (
	consts "ludo_backend_refactored/internal/config"
	"math"
)

const maxDepth = 4

func (g *Game) BotBestMoveMiniMax() int {
	bestScore := math.MinInt
	bestCol := -1

	for col := 0; col < consts.Cols; col++ {
		if !g.Board.IsValidMove(col) {
			continue
		}
		row := g.Board.ApplyMove(col, 2) // Bot = 2
		score := minimax(g.Board, maxDepth-1, false)
		g.Board.ResetCell(col, row)

		if score > bestScore {
			bestScore = score
			bestCol = col
		}
	}

	if bestCol == -1 {
		for col := 0; col < consts.Cols; col++ {
			if g.Board.IsValidMove(col) {
				return col
			}
		}
	}

	return bestCol
}

func minimax(board *Board, depth int, maximizing bool) int {
	if board.HasAnyWin(2) {
		return 100000 + depth
	}
	if board.HasAnyWin(1) {
		return -100000 - depth
	}
	if depth == 0 {
		return evaluateBoard(board)
	}

	if maximizing {
		maxEval := math.MinInt
		for col := 0; col < consts.Cols; col++ {
			if !board.IsValidMove(col) {
				continue
			}
			row := board.ApplyMove(col, 2)
			score := minimax(board, depth-1, false)
			board.ResetCell(col, row)
			maxEval = max(maxEval, score)
		}
		return maxEval
	} else {
		minEval := math.MaxInt
		for col := 0; col < consts.Cols; col++ {
			if !board.IsValidMove(col) {
				continue
			}
			row := board.ApplyMove(col, 1)
			score := minimax(board, depth-1, true)
			board.ResetCell(col, row)
			minEval = min(minEval, score)
		}
		return minEval
	}
}

func evaluateBoard(board *Board) int {
	score := 0
	for col := 0; col < consts.Cols; col++ {
		for row := 0; row < consts.Rows; row++ {
			cell := board.GetCell(col, row)
			if cell == 2 {
				score += positionScore(col)
			} else if cell == 1 {
				score -= positionScore(col)
			}
		}
	}
	return score
}

func positionScore(col int) int {
	center := consts.Cols / 2
	return consts.Cols - abs(col-center)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
