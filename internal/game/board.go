package game

const (
	Rows = 6
	Cols = 7
)

type Board struct {
	grid [Cols][Rows]int
}

func NewBoard() *Board {
	return &Board{}
}

func (b *Board) IsValidMove(col int) bool {
	return col >= 0 && col < Cols && b.grid[col][Rows-1] == 0
}

func (b *Board) ApplyMove(col, playerID int) int {
	if !b.IsValidMove(col) {
		return -1
	}
	for r := 0; r < Rows; r++ {
		if b.grid[col][r] == 0 {
			b.grid[col][r] = playerID
			return r
		}
	}
	return -1
}

func (b *Board) ApplyTempMove(col, playerID int) int {
	for r := 0; r < Rows; r++ {
		if b.grid[col][r] == 0 {
			b.grid[col][r] = playerID
			defer func() { b.grid[col][r] = 0 }()
			return r
		}
	}
	return -1
}

func (b *Board) Reset() {
	for c := 0; c < Cols; c++ {
		for r := 0; r < Rows; r++ {
			b.grid[c][r] = 0
		}
	}
}

func (b *Board) GetCell(col, row int) int {
	if col < 0 || col >= Cols || row < 0 || row >= Rows {
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
				if x < 0 || x >= Cols || y < 0 || y >= Rows || b.grid[x][y] != pid {
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
	if col >= 0 && col < Cols && row >= 0 && row < Rows {
		b.grid[col][row] = 0
	}
}

func (b *Board) HasAnyWin(playerID int) bool {
	for col := 0; col < Cols; col++ {
		for row := 0; row < Rows; row++ {
			if b.grid[col][row] == playerID {
				if b.CheckWin(col, row, playerID) {
					return true
				}
			}
		}
	}
	return false
}
