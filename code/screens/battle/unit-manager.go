package battle

import (
	"game/code/global"
	"game/code/unit"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/utility/collection"
	col "pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/text"
)

func (bs *BattleScreen) spawnUnits(spawns []float32, units []*unit.Unit) {
	if len(units)*2 > len(spawns) {
		return
	}
	var pointIndex = 0
	for _, u := range units {
		if number.IsNaN(spawns[pointIndex]) {
			pointIndex += 2
		}
		u.Spawn(spawns[pointIndex], spawns[pointIndex+1])
		pointIndex += 2
	}
}

func (bs *BattleScreen) updateUnits() {
	var ySortedUnits = bs.unitsSortedByY()
	var acting = bs.actingUnit()

	bs.hoveredUnit = nil
	for _, u := range bs.units {
		if u.IsHovered(bs.camera) {
			bs.hoveredUnit = u
			break
		}
	}

	bs.states.UpdateCurrentState()
	bs.drawTurnIndicators()

	for _, u := range ySortedUnits {
		var x, y = u.Position()
		var scX = condition.If[float32](bs.isBotUnit(u), -1, 1)
		u.UpdateAndDraw(x, y, scX, 1, bs.camera)
	}

	if bs.hoveredUnit != nil {
		bs.drawUnitStats("Hovered Unit", bs.hoveredUnit)
	} else {
		bs.drawUnitStats("Unit taking turn", acting)
	}
}

func (bs *BattleScreen) drawTurnIndicators() {
	var bw, bh = global.BattleColumns, global.BattleRows
	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	var acting = bs.actingUnit()
	var ux, uy = acting.Position()
	var mx, my = bs.camera.MousePosition()

	// Range cells (fills then frames) + acting unit highlight — one batch
	var quads []*graphics.Quad
	if bs.hoveredUnit == nil {
		quads = append(quads, bs.rangeCellQuads(bs.curMoveRangeCells, palette.Green)...)
	}
	quads = append(quads, indicatorQuad(ux-tw/2, uy-th/2, tw, th, palette.Azure))
	if bs.hoveredUnit == nil {
		quads = append(quads, bs.rangeCellFrameQuads(bs.curMoveRangeCells, -1, palette.Green)...)
	}
	bs.camera.DrawQuads(quads...)

	if mx < 0 || mx > bw*tw || my < 0 || my > bh*th {
		return
	}

	var moving = len(bs.curMovePath) > 0
	var canMove = acting.MovePoints >= 10
	var cx, cy = float32(int(mx / tw)), float32(int(my / th))
	bs.camera.DrawQuads(cellFrameQuads(cx*tw, cy*th, tw, th, -2, palette.White)...)

	if bs.hoveredUnit != nil && !moving {
		var hx, hy = bs.hoveredUnit.Position()
		var hoverQuads = []*graphics.Quad{indicatorQuad(hx-tw/2, hy-th/2, 64, 64, col.FadeOut(palette.White, 0.75))}
		if mouse.IsButtonPressed(button.Right) {
			var moveRange = bs.calculateRangeCells(cx, cy, float32(bs.hoveredUnit.BaseMovePoints/10))
			hoverQuads = append(hoverQuads, bs.rangeCellQuads(moveRange, col.FadeOut(palette.Yellow, 0.5))...)
			hoverQuads = append(hoverQuads, bs.rangeCellQuads(bs.hoveredUnit.AttackRangeCells(), col.FadeOut(palette.Red, 0.5))...)
			hoverQuads = append(hoverQuads, bs.rangeCellFrameQuads(moveRange, -1, col.FadeOut(palette.Yellow, 0.5))...)
			hoverQuads = append(hoverQuads, bs.rangeCellFrameQuads(bs.hoveredUnit.AttackRangeCells(), -1, col.FadeOut(palette.Red, 0.5))...)
		} else {
			hoverQuads = append(hoverQuads, bs.rangeCellQuads(acting.AttackRangeCells(), palette.Red)...)
			hoverQuads = append(hoverQuads, bs.rangeCellFrameQuads(acting.AttackRangeCells(), -1, palette.Red)...)
		}
		bs.camera.DrawQuads(hoverQuads...)
	}

	if bs.hoveredUnit != nil || moving || !canMove {
		return
	}

	var pts, tPts, inRange = bs.calculateMovePath(mx, my)
	if pts == 0 {
		return
	}

	var tx, ty = point.Snap(mx, my, tw, th)
	var tText = text.New(tPts, "/", acting.MovePoints)
	var tColor = condition.If(tPts <= acting.MovePoints, palette.White, palette.Red)
	var outsideRange = tPts > acting.MovePoints

	var pathQuads []*graphics.Quad
	if !outsideRange {
		pathQuads = append(pathQuads, indicatorQuad(tx, ty, tw, th, palette.Azure))
	}
	if outsideRange {
		var x, y = inRange[len(inRange)-2] - tw/2, inRange[len(inRange)-1] - th/2
		pathQuads = append(pathQuads, indicatorQuad(x, y, tw, th, palette.Azure))
	}
	bs.camera.DrawQuads(pathQuads...)

	bs.camera.DrawTextAdvanced("", tText, tx+2, ty, 20, 1, 0, 0, palette.Black)
	bs.camera.DrawTextAdvanced("", tText, tx+2, ty, 20, 1, 0, 0, tColor)
	if outsideRange {
		var x, y = inRange[len(inRange)-2] - tw/2, inRange[len(inRange)-1] - th/2
		var txt = text.New(pts, "/", acting.MovePoints)
		bs.camera.DrawTextAdvanced("", txt, x+2, y, 20, 1, 0, 0, palette.Black)
		bs.camera.DrawTextAdvanced("", txt, x+2, y, 20, 1, 0, 0, palette.White)
	}
}

func (bs *BattleScreen) rangeCellQuads(cells [][2]int, color uint) []*graphics.Quad {
	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	var bw, bh = global.BattleColumns, global.BattleRows
	var quads = make([]*graphics.Quad, 0, len(cells))
	for _, cell := range cells {
		var cx, cy = cell[0], cell[1]
		if cx < 0 || cx >= int(bw) || cy < 0 || cy >= int(bh) {
			continue
		}
		quads = append(quads, indicatorQuad(float32(cx)*tw, float32(cy)*th, tw, th, col.FadeOut(color, 0.8)))
	}
	return quads
}
func (bs *BattleScreen) rangeCellFrameQuads(cells [][2]int, frameSize float32, color uint) []*graphics.Quad {
	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	var bw, bh = global.BattleColumns, global.BattleRows
	var quads = make([]*graphics.Quad, 0, len(cells)*4)
	for _, cell := range cells {
		var cx, cy = cell[0], cell[1]
		if cx < 0 || cx >= int(bw) || cy < 0 || cy >= int(bh) {
			continue
		}
		quads = append(quads, cellFrameQuads(float32(cx)*tw, float32(cy)*th, tw, th, frameSize, color)...)
	}
	return quads
}

func (bs *BattleScreen) drawUnitStats(description string, u *unit.Unit) {
	var lineHeight = 80 / bs.camera.Zoom
	var txt = text.New(
		description, "\n",
		"Initiative: ", u.Initiative, "\n",
		"Movement: ", u.MovePoints, "/", u.BaseMovePoints, "\n",
	)
	var lines = len(text.Split(txt, "\n"))
	var blx, bly = bs.camera.PointFromEdge(0, 1)
	var x, y = blx + 50/bs.camera.Zoom, bly - lineHeight*float32(lines)
	bs.camera.DrawTextAdvanced("", txt, x, y, lineHeight, 0, 0, 0, palette.Black)
	bs.camera.DrawTextAdvanced("", txt, x, y, lineHeight, 0, 0, 0, palette.White)
}

//=================================================================

func (bs *BattleScreen) unitsSortedByY() []*unit.Unit {
	var ySorted = make(map[float32][]*unit.Unit, len(bs.units))

	for _, u := range bs.units {
		var _, y = u.Cell()
		ySorted[y] = append(ySorted[y], u)
	}

	var keys = collection.MapKeys(ySorted)
	var result = make([]*unit.Unit, 0, len(ySorted))

	collection.SortNumbers(keys...)
	for _, key := range keys {
		result = append(result, ySorted[key]...)
	}

	return result
}

//=================================================================

func indicatorQuad(x, y, w, h float32, tint uint) *graphics.Quad {
	return &graphics.Quad{Area: graphics.Area{X: x, Y: y, Width: w, Height: h}, ScaleX: 1, ScaleY: 1, Tint: tint}
}
func cellFrameQuads(x, y, tw, th, frameSize float32, tint uint) []*graphics.Quad {
	if frameSize < 0 {
		var t = -frameSize
		return []*graphics.Quad{
			indicatorQuad(x, y, tw, t, tint),            // top    (covers corners)
			indicatorQuad(x, y+th-t, tw, t, tint),       // bottom (covers corners)
			indicatorQuad(x, y+t, t, th-2*t, tint),      // left   (inner height)
			indicatorQuad(x+tw-t, y+t, t, th-2*t, tint), // right  (inner height)
		}
	}
	var t = frameSize
	return []*graphics.Quad{
		indicatorQuad(x-t, y-t, tw+2*t, t, tint),  // top    (extended to cover corners)
		indicatorQuad(x-t, y+th, tw+2*t, t, tint), // bottom (extended to cover corners)
		indicatorQuad(x-t, y, t, th, tint),        // left
		indicatorQuad(x+tw, y, t, th, tint),       // right
	}
}
