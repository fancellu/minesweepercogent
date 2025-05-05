package main

import (
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/states"
	"cogentcore.org/core/tree"
	"fmt"
	"math/rand"
)

const (
	Empty = '0'
	Mine  = 'M'
)

type Info struct {
	content  rune
	revealed bool
	button   *MyButton
}

type Board struct {
	grid        *core.Frame
	info        [][]Info
	rows        int
	cols        int
	mines       int
	minesPlaced bool
	frozen      bool
}

func NewBoard(grid *core.Frame, rows, cols, mines int) *Board {
	board := &Board{
		grid:        grid,
		rows:        rows,
		cols:        cols,
		mines:       mines,
		info:        make([][]Info, rows),
		minesPlaced: false,
		frozen:      false,
	}

	grid.DeleteChildren()

	for y := range board.info {
		board.info[y] = make([]Info, cols)

		for x := range board.info[y] {
			board.info[y][x].content = Empty
			board.info[y][x].button = board.newButton(grid, y, x)
		}
	}
	//board.placeMines()

	grid.Update()
	return board
}

// We don't place mines until first click, and avoid the first click
func (b *Board) placeMines(avoidy, avoidx int) {
	//fmt.Printf("avoid x=%d, y=%d\n", avoidx, avoidy)
	for i := 0; i < b.mines; {
		y := rand.Intn(b.rows)
		x := rand.Intn(b.cols)
		if x == avoidx && y == avoidy {
			//fmt.Printf("avoided x=%d, y=%d\n", x, y)
			continue
		}
		if b.info[y][x].content == Empty {
			bt := b.info[y][x].button
			bt.Mine = true
			//fmt.Printf("placed mine at x=%d, y=%d\n", x, y)
			b.info[y][x].content = Mine
			i++
		}
	}
	b.calculateNumbers()
	b.printBoard()
	b.minesPlaced = true
}

func (b *Board) showMines() {
	for y := 0; y < b.rows; y++ {
		for x := 0; x < b.cols; x++ {
			bt := b.info[y][x].button
			if bt.Mine {
				bt.ShowMineIcon()
			}
		}
	}
	core.MessageSnackbar(b.grid.Scene, "You lost")
	b.frozen = true
}

func (b *Board) calculateNumbers() {
	for y := 0; y < b.rows; y++ {
		for x := 0; x < b.cols; x++ {
			if b.info[y][x].content != Mine {
				b.info[y][x].content = rune(b.countAdjacentMines(y, x) + Empty)
			}
		}
	}
}

func (b *Board) wonCheck() {
	revealed := 0
	for y := 0; y < b.rows; y++ {
		for x := 0; x < b.cols; x++ {
			if b.info[y][x].revealed {
				fmt.Printf("rev %d, %d ", y, x)
				revealed++
			} else if b.info[y][x].button.Flag {
				fmt.Printf("flag %d, %d ", y, x)
				revealed++
			} else {
				fmt.Printf("not rev %d, %d ", y, x)
			}
		}
	}
	if revealed == (b.rows * b.cols) {
		core.MessageSnackbar(b.grid.Scene, "You won")
		b.frozen = true
		return
	}
}

func (b *Board) newButton(grid *core.Frame, y, x int) *MyButton {
	bt := tree.New[MyButton](grid)
	bt.OnClick(func(e events.Event) {
		if b.frozen {
			return
		}
		if bt.Flag {
			return
		}
		if !b.minesPlaced {
			b.placeMines(y, x)
		}
		b.reveal(y, x, func(ry, rx int) {
			fmt.Println("____Revealed cell at", ry, rx)
			rune := b.info[ry][rx].content
			bt := b.info[ry][rx].button
			if rune == Mine {
				bt.ShowMineIcon()
			} else if rune != Empty {
				bt.SetText(fmt.Sprintf("%c", rune))
				bt.SetIcon(icons.None)
			} else { // it's a zero
				bt.SetText("")
				bt.SetState(true, states.Checked)
			}
			bt.Update()
		})
	})
	bt.On(events.ContextMenu, func(e events.Event) {
		if b.frozen {
			return
		}
		if !b.minesPlaced {
			b.placeMines(y, x)
		}
		// can't flag an already known cell
		revealed := b.info[y][x].revealed
		if revealed {
			return
		}
		if bt.Icon == icons.Blank || bt.Icon == "" {
			bt.ShowFlagIcon()
			b.wonCheck()
		} else {
			bt.SetIcon(icons.Blank)
			bt.Flag = false
		}

		bt.Update()
	})

	return bt
}

func (b *Board) countAdjacentMines(row, col int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			r, c := row+i, col+j
			if r >= 0 && r < b.rows && c >= 0 && c < b.cols && b.info[r][c].content == Mine {
				count++
			}
		}
	}
	return count
}

func (b *Board) reveal(row, col int, onReveal func(row, col int)) bool {

	// already revealed, so ignore click
	if b.info[row][col].revealed {
		return true
	}

	rune := b.info[row][col].content

	if rune == Mine {
		fmt.Println("reveal(): Mine found at", row, col)
		onReveal(row, col)
		b.showMines()
		return false
	}

	b.info[row][col].revealed = true
	onReveal(row, col)

	if rune == Empty {
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				if i == 0 && j == 0 {
					continue
				}
				r, c := row+i, col+j
				if r >= 0 && r < b.rows && c >= 0 && c < b.cols && b.info[r][c].content != Mine {
					b.reveal(r, c, onReveal)
				}
			}
		}
	}

	b.wonCheck()
	return true
}

func (b *Board) printBoard() {
	for y := 0; y < b.rows; y++ {
		for x := 0; x < b.cols; x++ {
			fmt.Printf("%c ", b.info[y][x].content)
		}
		fmt.Println()
	}
}

func main() {
	b := core.NewBody("minesweeper").SetTitle("Minesweeper")

	rows, cols, mines := 10, 10, 10

	grid := core.NewFrame(b)
	grid.Styler(func(s *styles.Style) {
		s.Display = styles.Grid
		s.Columns = cols
		s.Gap.Zero()
	})

	board := NewBoard(grid, rows, cols, mines)

	core.NewButton(b).SetText("Reset!").SetIcon(icons.Reset).OnClick(func(e events.Event) {
		board = NewBoard(grid, rows, cols, mines)
		board.printBoard()
	})

	grid.Scene.ContextMenus = nil

	b.RunMainWindow()
}
