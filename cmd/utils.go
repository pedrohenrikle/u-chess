package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Check if a position lies on the board
func (b *Board) InBounds(p Pos) bool {
	return p.X >= 0 && p.X < b.Width && p.Y >= 0 && p.Y < b.Height
}

// At returns the piece at p, or nil if out of bounds/empty.
func (b *Board) At(p Pos) Piece {
	if p.X < 0 || p.X >= b.Width || p.Y < 0 || p.Y >= b.Height {
		return nil
	}
	return b.Cells[p.Y][p.X]
}

// parseSquare converts “A1”→Pos{0,0}, “C5”→Pos{2,4}
func parseSquare(s string) (Pos, bool) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if len(s) < 2 {
		return Pos{}, false
	}
	file := s[0]
	rank, err := strconv.Atoi(s[1:])
	if err != nil || rank < 1 {
		return Pos{}, false
	}
	return Pos{int(file - 'A'), rank - 1}, true
}

func formatSquare(p Pos) string {
	return fmt.Sprintf("%c%d", 'A'+p.X, p.Y+1)
}

func formatSquares(list []Pos) []string {
	out := make([]string, len(list))
	for i, p := range list {
		out[i] = formatSquare(p)
	}
	return out
}

// Symbol returns a one‐character representation of a piece.
func Symbol(p Piece) rune {
	switch pt := p.(type) {
	case *ProductOwner:
		if pt.col == White {
			return '♔'
		}
		return '♚'
	case *Developer:
		if pt.col == White {
			return '♖'
		}
		return '♜'
	case *Designer:
		if pt.col == White {
			return '♘'
		}
		return '♞'
	default:
		return '?'
	}
}

func opposite(c Color) Color {
	if c == White {
		return Black
	}
	return White
}

func containsPos(list []Pos, p Pos) bool {
	for _, x := range list {
		if x == p {
			return true
		}
	}
	return false
}

// movePiece enforces turn, legality, captures, and updates the board
func (g *Game) movePiece(from, to Pos) error {
	piece := g.board.At(from)
	if piece == nil {
		return fmt.Errorf("no piece at %s", formatSquare(from))
	}
	if piece.Color() != g.turn {
		return fmt.Errorf("not your turn")
	}

	valid := piece.ValidMoves(from, g.board)
	if !containsPos(valid, to) {
		return fmt.Errorf("illegal move to %s", formatSquare(to))
	}

	// If this is the Developer, sweep and capture any enemy in the path
	if dev, ok := piece.(*Developer); ok {
		dx, dy := to.X-from.X, to.Y-from.Y
		stepX, stepY := sign(dx), sign(dy)
		cur := Pos{from.X + stepX, from.Y + stepY}

		for cur != to {
			if p := g.board.At(cur); p != nil && p.Color() != dev.col {
				g.board.Cells[cur.Y][cur.X] = nil // capture
			}
			cur = Pos{cur.X + stepX, cur.Y + stepY}
		}
	}

	// finally, move the piece
	g.board.Cells[to.Y][to.X] = piece
	g.board.Cells[from.Y][from.X] = nil
	return nil
}

func sign(x int) int {
	switch {
	case x < 0:
		return -1
	case x > 0:
		return 1
	default:
		return 0
	}
}
