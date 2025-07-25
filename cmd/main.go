package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// pos is a board coordinate
type Pos struct{ X, Y int }

// color of a piece
type Color int

// String method returns "White" or "Black"
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

// Piece represents any unit on the board
type Piece interface {
	Color() Color
	ValidMoves(from Pos, board *Board) []Pos
}

// Board holds the grid of pieces
type Board struct {
	Width, Height int
	Cells         [][]Piece
}

// HasProductOwner checks if the given color’s product owner is still on board
func (b *Board) HasProductOwner(col Color) bool {
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			if po, ok := b.Cells[y][x].(*ProductOwner); ok && po.col == col {
				return true
			}
		}
	}
	return false
}

// Game holds the current game state
type Game struct {
	board    *Board
	selected *Pos
	turn     Color
	gameOver bool
}

func main() {
	println("------------------------")
	println("Welcome to uChess")
	println("------------------------")

	width := GetDimension("width")
	height := GetDimension("height")
	fmt.Printf("Starting a %dx%d board...\n", width, height)

	game := NewGame(width, height)
	game.Run()
}

// NewGame initializes a new game
func NewGame(w, h int) *Game {
	return &Game{
		board:    NewBoard(w, h),
		selected: nil,
		turn:     White,
	}
}

// GetDimension reads an integer between 6 and 12 from stdin
func GetDimension(name string) int {
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

// NewBoard creates an empty board and places each side's pieces
func NewBoard(w, h int) *Board {
	cells := make([][]Piece, h)
	for y := range cells {
		cells[y] = make([]Piece, w)
	}
	b := &Board{Width: w, Height: h, Cells: cells}

	// place white pieces on bottom row
	b.Cells[0][0] = &ProductOwner{White}
	b.Cells[0][1] = &Developer{White}
	b.Cells[0][2] = &Designer{White}

	// place black pieces on top row
	top := h - 1
	b.Cells[top][w-1] = &ProductOwner{Black}
	b.Cells[top][w-2] = &Developer{Black}
	b.Cells[top][w-3] = &Designer{Black}

	return b
}

// Display prints the board, using x for highlights and . for empty houses
func (b *Board) Display(highlights map[Pos]bool) {
	// column headers
	fmt.Print("   ")
	for x := 0; x < b.Width; x++ {
		fmt.Printf(" %c", 'A'+x)
	}
	fmt.Println()

	// rows from top down
	for y := b.Height - 1; y >= 0; y-- {
		fmt.Printf("%2d ", y+1)
		for x := 0; x < b.Width; x++ {
			p := b.At(Pos{x, y})
			switch {
			case p != nil:
				// occupied cell
				fmt.Printf(" %c", Symbol(p))
			case highlights[Pos{x, y}]:
				// valid move highlight
				fmt.Print(" x")
			default:
				// empty cell
				fmt.Print(" .")
			}
		}
		fmt.Println()
	}
}

// PrintHelp shows available commands
func (g *Game) PrintHelp() {
	fmt.Println(`Commands:
  help                 show this help
  exit                 quit the game
  restart              pick size & restart
  select <sq>          choose a piece (e.g. select A1)
  move <from> <to>     move a piece (e.g. move A1 B3)`)
}

// Run method starts the main input loop
func (g *Game) Run() {
	reader := bufio.NewReader(os.Stdin)

	// initial draw
	g.board.Display(nil)
	fmt.Println("\nType `help` for commands.")

	for {
		// show turn unless game is over
		if !g.gameOver {
			fmt.Printf("\nTurn: %s\n", g.turn)
		}
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
			g.PrintHelp()

		case "exit":
			fmt.Println("Goodbye!")
			os.Exit(0)

		case "restart":
			// start a new game with fresh dimensions
			w, h := GetDimension("width"), GetDimension("height")
			*g = *NewGame(w, h)
			g.gameOver = false
			fmt.Printf("\nTurn: %s\n", g.turn)
			g.board.Display(nil)

		case "select":
			// choose a piece and highlight its valid moves
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

			moves := piece.ValidMoves(pos, g.board)
			hl := make(map[Pos]bool, len(moves))
			for _, m := range moves {
				hl[m] = true
			}

			g.board.Display(hl)
			fmt.Printf("Valid moves for %c at %s: %v\n",
				Symbol(piece), strings.ToUpper(parts[1]), formatSquares(moves))

		case "move":
			// attempt to move a piece, with full validation and reporting
			if len(parts) != 3 {
				fmt.Println("Usage: move <from> <to>")
				continue
			}
			fromS, toS := strings.ToUpper(parts[1]), strings.ToUpper(parts[2])
			from, ok1 := parseSquare(fromS)
			to, ok2 := parseSquare(toS)
			if !ok1 || !ok2 || !g.board.InBounds(from) || !g.board.InBounds(to) {
				fmt.Printf("Invalid Move: %s or %s is out of bounds\n", fromS, toS)
				continue
			}
			if from == to {
				fmt.Println("Invalid Move: Destination must be different from origin")
				continue
			}

			piece := g.board.At(from)
			if piece == nil {
				fmt.Printf("Invalid Move: No piece at %s\n", fromS)
				continue
			}
			if piece.Color() != g.turn {
				fmt.Println("Invalid Move: You can't move opponent's piece")
				continue
			}

			valid := piece.ValidMoves(from, g.board)
			if !containsPos(valid, to) {
				fmt.Printf("Invalid Move: %c piece can't move to %s\n", Symbol(piece), toS)
				continue
			}

			// detect captures before executing the move
			var captured []Piece
			if _, isDev := piece.(*Developer); !isDev {
				if p := g.board.At(to); p != nil {
					captured = append(captured, p)
				}
			} else {
				dx, dy := to.X-from.X, to.Y-from.Y
				stepX, stepY := sign(dx), sign(dy)
				cur := Pos{from.X + stepX, from.Y + stepY}
				for cur != to {
					if p := g.board.At(cur); p != nil && p.Color() != g.turn {
						captured = append(captured, p)
					}
					cur = Pos{cur.X + stepX, cur.Y + stepY}
				}
			}

			// perform the move and remove captures
			if err := g.movePiece(from, to); err != nil {
				fmt.Println("Move error:", err)
				continue
			}

			// redraw board
			g.board.Display(nil)

			// report move outcome
			icon := Symbol(piece)
			if len(captured) > 0 {
				icons := make([]string, len(captured))
				for i, cp := range captured {
					icons[i] = string(Symbol(cp))
				}
				fmt.Printf("\nMoved %c from %s to %s. Captured %s.\n",
					icon, fromS, toS, strings.Join(icons, ", "))
			} else {
				fmt.Printf("\nMoved %c from %s to %s.\n",
					icon, fromS, toS)
			}

			// switch turns
			g.turn = opposite(g.turn)

		default:
			fmt.Println("Unknown command:", cmd)
		}
	}
}
