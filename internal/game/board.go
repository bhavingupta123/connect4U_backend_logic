package game

import consts "ludo_backend_refactored/internal/config"

type Board struct {
	grid [consts.Cols][consts.Rows]int
}

func NewBoard() *Board {
	return &Board{}
}

func (b *Board) IsValidMove(col int) bool {
	return col >= 0 && col < consts.Cols && b.grid[col][consts.Rows-1] == 0
}

func (b *Board) ApplyMove(col, playerID int) int {
	if !b.IsValidMove(col) {
		return -1
	}
	for r := 0; r < consts.Rows; r++ {
		if b.grid[col][r] == 0 {
			b.grid[col][r] = playerID
			return r
		}
	}
	return -1
}

func (b *Board) ApplyTempMove(col, playerID int) int {
	for r := 0; r < consts.Rows; r++ {
		if b.grid[col][r] == 0 {
			b.grid[col][r] = playerID
			defer func() { b.grid[col][r] = 0 }()
			return r
		}
	}
	return -1
}

func (b *Board) Reset() {
	for c := 0; c < consts.Cols; c++ {
		for r := 0; r < consts.Rows; r++ {
			b.grid[c][r] = 0
		}
	}
}

func (b *Board) GetCell(col, row int) int {
	if col < 0 || col >= consts.Cols || row < 0 || row >= consts.Rows {
		return -1
	}
	return b.grid[col][row]
}

func (b *Board) CheckWin(col, row, pid int) bool {
	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for _, d := range directions {
		count := 1
		for _, sign := range []int{1, -1} {
			x, y := col, row
			for {
				x += d[0] * sign
				y += d[1] * sign
				if x < 0 || x >= consts.Cols || y < 0 || y >= consts.Rows || b.grid[x][y] != pid {
					break
				}
				count++
			}
		}
		if count >= 4 {
			return true
		}
	}
	return false
}

func (b *Board) ResetCell(col, row int) {
	if col >= 0 && col < consts.Cols && row >= 0 && row < consts.Rows {
		b.grid[col][row] = 0
	}
}

func (b *Board) HasAnyWin(playerID int) bool {
	for col := 0; col < consts.Cols; col++ {
		for row := 0; row < consts.Rows; row++ {
			if b.grid[col][row] == playerID {
				if b.CheckWin(col, row, playerID) {
					return true
				}
			}
		}
	}
	return false
}

func (b *Board) GetWinningCells(col, row, pid int) [][2]int {
	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}

	for _, d := range directions {
		cells := [][2]int{{col, row}}

		x, y := col, row
		for {
			x += d[0]
			y += d[1]
			if x < 0 || x >= consts.Cols || y < 0 || y >= consts.Rows || b.grid[x][y] != pid {
				break
			}
			cells = append(cells, [2]int{x, y})
		}

		x, y = col, row
		for {
			x -= d[0]
			y -= d[1]
			if x < 0 || x >= consts.Cols || y < 0 || y >= consts.Rows || b.grid[x][y] != pid {
				break
			}
			cells = append(cells, [2]int{x, y})
		}

		if len(cells) >= 4 {
			return cells
		}
	}
	return nil
}
