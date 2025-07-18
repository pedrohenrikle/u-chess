package main

type Developer struct{ col Color }

func (d *Developer) Color() Color { return d.col }

func (d *Developer) ValidMoves(from Pos, b *Board) []Pos {
	var moves []Pos
	// 8 straight & diagonal directions
	directions := []Pos{
		{1, 0}, {0, 1}, {-1, 0}, {0, -1},
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
	}

	for _, dir := range directions {
		// 1) Normal moves (1–3 squares), path must be clear
		for step := 1; step <= 3; step++ {
			tgt := Pos{from.X + dir.X*step, from.Y + dir.Y*step}
			if !b.InBounds(tgt) {
				break
			}
			if p := b.At(tgt); p != nil {
				// can’t go past any piece if not capturing
				if p.Color() != d.col {
					// allow landing on enemy (simple capture)
					moves = append(moves, tgt)
				}
				break
			}
			moves = append(moves, tgt)
		}

		// 2) Jump-capture: exactly 2 squares
		//    over an enemy piece landing on empty square
		jumpOver := Pos{from.X + dir.X, from.Y + dir.Y}
		land := Pos{from.X + dir.X*2, from.Y + dir.Y*2}
		if b.InBounds(land) &&
			b.At(jumpOver) != nil &&
			b.At(jumpOver).Color() != d.col &&
			b.At(land) == nil {

			moves = append(moves, land)
		}
	}

	return moves
}

type Designer struct{ col Color }

func (d *Designer) Color() Color { return d.col }

func (d *Designer) ValidMoves(from Pos, b *Board) []Pos {
	var moves []Pos
	deltas := []Pos{
		{1, 2}, {2, 1}, {2, -1}, {1, -2},
		{-1, -2}, {-2, -1}, {-2, 1}, {-1, 2},
	}

	for _, delta := range deltas {
		tgt := Pos{from.X + delta.X, from.Y + delta.Y}
		if !b.InBounds(tgt) {
			continue
		}
		targetPiece := b.At(tgt)
		if targetPiece == nil || targetPiece.Color() != d.col {
			moves = append(moves, tgt)
		}
	}

	return moves
}

type ProductOwner struct{ col Color }

func (k *ProductOwner) Color() Color { return k.col }

func (k *ProductOwner) ValidMoves(from Pos, b *Board) []Pos {
	var moves []Pos
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			tgt := Pos{from.X + dx, from.Y + dy}
			if !b.InBounds(tgt) {
				continue
			}
			if p := b.At(tgt); p == nil || p.Color() != k.col {
				moves = append(moves, tgt)
			}
		}
	}
	return moves
}
