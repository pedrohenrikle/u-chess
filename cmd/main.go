package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Pos is a coordinate on the board
type Pos struct{ X, Y int }

// Color of a piece
type Color int

func (c Color) String() string {
	if c == White {
		return "White"
	}
	return "Black"
}

const (
	White Color = iota
	Black
)

// Piece is any board‐occupying unit
type Piece interface {
	Color() Color
	ValidMoves(from Pos, board *Board) []Pos
}

// Board holds a 2D slice of Piece (nil = empty)
type Board struct {
	Width, Height int
	Cells         [][]Piece
}

// Game holds state for the REPL
type Game struct {
	board    *Board
	selected *Pos
	turn     Color
}

func main() {
	println("------------------------")
	println("Welcome to Unvoid Chess")
	println("------------------------")

	width := getDimension("width")
	height := getDimension("height")
	fmt.Printf("Starting a %dx%d board...\n", width, height)

	game := NewGame(width, height)
	game.Run()
}

// NewGame creates a fresh game
func NewGame(w, h int) *Game {
	return &Game{
		board:    NewBoard(w, h),
		selected: nil,
		turn:     White,
	}
}

// getDimension prompts the user for a value between 6 and 12.
func getDimension(name string) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Enter %s (6–12): ", name)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		v, err := strconv.Atoi(line)
		if err != nil || v < 6 || v > 12 {
			fmt.Println("Invalid. Please enter an integer between 6 and 12.")
			continue
		}
		return v
	}
}

// NewBoard allocates the grid and places the three pieces per side.
func NewBoard(w, h int) *Board {
	cells := make([][]Piece, h)
	for y := range cells {
		cells[y] = make([]Piece, w)
	}
	b := &Board{Width: w, Height: h, Cells: cells}

	// White on row 0, cols 0,1,2
	b.Cells[0][0] = &ProductOwner{White}
	b.Cells[0][1] = &Developer{White}
	b.Cells[0][2] = &Designer{White}

	// Black on top row, rightmost cols
	top := h - 1
	b.Cells[top][w-1] = &ProductOwner{Black}
	b.Cells[top][w-2] = &Developer{Black}
	b.Cells[top][w-3] = &Designer{Black}

	return b
}

// Display prints the board with “.” for empty and piece symbols for occupied.
func (b *Board) Display() {
	// Column headers
	fmt.Print("   ")
	for x := 0; x < b.Width; x++ {
		fmt.Printf(" %c", 'A'+x)
	}
	fmt.Println()

	// Rows top→bottom
	for y := b.Height - 1; y >= 0; y-- {
		fmt.Printf("%2d ", y+1)
		for x := 0; x < b.Width; x++ {
			if p := b.At(Pos{x, y}); p == nil {
				fmt.Print(" .")
			} else {
				fmt.Printf(" %c", Symbol(p))
			}
		}
		fmt.Println()
	}
}

// printHelp lists commands
func (g *Game) printHelp() {
	fmt.Println(`Commands:
  help                 show this help
  exit                 quit the game
  restart              pick size & restart
  select <sq>          choose a piece (e.g. select A1)
  move <from> <to>     move a piece (e.g. move A1 B3)`)
}

// Run starts the REPL loop
func (g *Game) Run() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n")
	g.board.Display()

	fmt.Println("\nType `help` for commands.")

	for {
		// Show whose turn it is:
		fmt.Printf("\nTurn: %s\n", g.turn)
		fmt.Print("> ")
		raw, _ := reader.ReadString('\n')
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		cmd := strings.ToLower(parts[0])

		switch cmd {
		case "help":
			g.printHelp()

		case "exit":
			fmt.Println("Goodbye!")
			os.Exit(0)

		case "restart":
			w, h := getDimension("width"), getDimension("height")
			*g = *NewGame(w, h)
			g.board.Display()

		case "select":
			if len(parts) != 2 {
				fmt.Println("Usage: select <square>")
				continue
			}
			pos, ok := parseSquare(parts[1])
			if !ok || !g.board.InBounds(pos) {
				fmt.Println("Invalid square:", parts[1])
				continue
			}
			piece := g.board.At(pos)
			if piece == nil || piece.Color() != g.turn {
				fmt.Println("No", g.turn, "piece at", parts[1])
				continue
			}
			g.selected = &pos
			moves := piece.ValidMoves(pos, g.board)
			g.board.Display()
			fmt.Printf("Valid moves for %c at %s: %v\n",
				Symbol(piece), parts[1], formatSquares(moves))

		case "move":
			if len(parts) != 3 {
				fmt.Println("Usage: move <from> <to>")
				continue
			}
			from, ok1 := parseSquare(parts[1])
			to, ok2 := parseSquare(parts[2])
			if !ok1 || !ok2 || !g.board.InBounds(from) || !g.board.InBounds(to) {
				fmt.Println("Invalid squares:", parts[1], parts[2])
				continue
			}
			if err := g.movePiece(from, to); err != nil {
				fmt.Println("Move error:", err)
				continue
			}
			g.turn = opposite(g.turn)
			g.selected = nil
			g.board.Display()
			fmt.Println(g.turn, "to move.")

		default:
			fmt.Println("Unknown command:", cmd)
		}
	}
}
