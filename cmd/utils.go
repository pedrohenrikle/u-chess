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

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// maybeRemoveJumped checks for a single‐step jump in straight/diag.
// If true, removes the intervening piece.
func maybeRemoveJumped(from, to Pos, d *Developer, b *Board) {
	dx, dy := to.X-from.X, to.Y-from.Y
	if abs(dx) == 2 || abs(dy) == 2 {
		mid := Pos{from.X + sign(dx), from.Y + sign(dy)}
		if p := b.At(mid); p != nil && p.Color() != d.col {
			b.Cells[mid.Y][mid.X] = nil
		}
	}
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
	if dev, ok := piece.(*Developer); ok {
		maybeRemoveJumped(from, to, dev, g.board)
	}
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
