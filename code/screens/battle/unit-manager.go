package battle

import (
	"game/code/global"
	"game/code/unit"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/screens"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/tiled"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/collection"
	col "pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/text"
)

// Handles all units on the battle screen, taking care of their spawning, updating gameplay & drawing.
type unitManager struct {
	units       []*unit.Unit
	hoveredUnit *unit.Unit
	turnManager *turnManager
}

func newUnitManager(teamA, teamB []*unit.Unit) *unitManager {
	return &unitManager{units: collection.Join(teamA, teamB), turnManager: newTurnManager()}
}

//=================================================================

func (um *unitManager) spawnAll(tmap *tiled.Map, units []*unit.Unit, flip bool, layerClass string) {
	var spawns = tmap.FindLayersBy(property.LayerClass, layerClass)[0].ExtractPoints()
	if len(units) > len(spawns) {
		return
	}
	for i, u := range units {
		u.Spawn(spawns[i][0], spawns[i][1], flip)
	}
}
func (um *unitManager) update() {
	var battle = screens.Current().(*BattleScreen)
	var ySortedUnits = um.ySortAll()
	var unitActing = um.turnManager.unitActing()

	um.hoveredUnit = nil
	for _, unit := range um.units {
		if unit.IsHovered(battle.camera) {
			um.hoveredUnit = unit
			break
		}
	}

	um.turnManager.states.UpdateCurrentState()
	um.drawIndicators()

	for _, unit := range ySortedUnits {
		unit.UpdateAndDraw(battle.camera)
	}

	if um.hoveredUnit != nil {
		um.drawStats("Hovered Unit", um.hoveredUnit)
	} else {
		um.drawStats("Unit taking turn", unitActing)
	}
}

func (um *unitManager) drawIndicators() {
	var battle = screens.Current().(*BattleScreen)
	var bw, bh = global.BattleTileColumns, global.BattleTileRows
	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	var unitActing = um.turnManager.unitActing()
	var ux, uy = unitActing.Position()
	var mx, my = battle.camera.MousePosition()

	if um.hoveredUnit == nil {
		um.drawRange(um.turnManager.curMoveRangeCells, -1, palette.Green)
	}

	battle.camera.DrawQuad(ux-tw/2, uy-th/2, tw, th, 0, palette.Azure)

	if mx < 0 || mx > bw*tw || my < 0 || my > bh*th {
		return
	}

	var moving = len(um.turnManager.curMovePath) > 0
	var canMove = unitActing.MovePoints >= 10
	var cx, cy = float32(int(mx / tw)), float32(int(my / th))
	battle.camera.DrawQuadFrame(cx*tw, cy*th, tw, th, 0, -2, palette.White)
	if um.hoveredUnit != nil && !moving {
		var ux, uy = um.hoveredUnit.Position()
		battle.camera.DrawQuad(ux-tw/2, uy-th/2, 64, 64, 0, col.FadeOut(palette.White, 0.75))

		if mouse.IsButtonPressed(button.Right) {
			var moveRange = um.turnManager.calculateRangeCells(cx, cy, float32(um.hoveredUnit.BaseMovePoints/10))
			um.drawRange(moveRange, -1, col.FadeOut(palette.Yellow, 0.5))
			um.drawRange(um.hoveredUnit.AttackRangeCells(), -1, col.FadeOut(palette.Red, 0.5))
		} else {
			um.drawRange(unitActing.AttackRangeCells(), -1, palette.Red)
		}
	}

	if um.hoveredUnit != nil || moving || !canMove {
		return
	}

	var pts, tPts, inRange = um.turnManager.calculateMovePath(mx, my)
	if pts == 0 {
		return
	}

	var tx, ty = point.Snap(mx, my, tw, th)
	var tText = text.New(tPts, "/", unitActing.MovePoints)
	var tColor = condition.If(tPts <= unitActing.MovePoints, palette.White, palette.Red)
	var outsideRange = tPts > unitActing.MovePoints

	if !outsideRange {
		battle.camera.DrawQuad(tx, ty, tw, th, 0, palette.Azure)
	}
	battle.camera.DrawText("", tText, tx+2, ty, 20, 0.95, 1, palette.Black)
	battle.camera.DrawText("", tText, tx+2, ty, 20, 0.5, 1, tColor)

	if outsideRange {
		var x, y = inRange[len(inRange)-1][0] - tw/2 + 2, inRange[len(inRange)-1][1] - th/2
		var txt = text.New(pts, "/", unitActing.MovePoints)
		battle.camera.DrawQuad(x, y, tw, th, 0, palette.Azure)
		battle.camera.DrawText("", txt, x, y, 20, 0.95, 1, palette.Black)
		battle.camera.DrawText("", txt, x, y, 20, 0.5, 1, palette.White)
	}
}
func (um *unitManager) drawRange(cells [][2]int, frameSize float32, color uint) {
	var battle = screens.Current().(*BattleScreen)
	var bw, bh = global.BattleTileColumns, global.BattleTileRows
	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	for _, cell := range cells {
		var cx, cy = cell[0], cell[1]
		if cx < 0 || cx >= int(bw) || cy < 0 || cy >= int(bh) {
			continue
		}

		var curCx, curCy = float32(cx) * tw, float32(cy) * th
		battle.camera.DrawQuad(curCx, curCy, float32(tw), float32(th), 0, col.FadeOut(color, 0.8))
		battle.camera.DrawQuadFrame(curCx, curCy, float32(tw), float32(th), 0, frameSize, color)
	}
}
func (um *unitManager) drawStats(description string, unit *unit.Unit) {
	var battle = screens.Current().(*BattleScreen)
	var lineHeight = 80 / battle.camera.Zoom
	var txt = text.New(
		description, "\n",
		"Initiative: ", unit.Initiative, "\n",
		"Movement: ", unit.MovePoints, "/", unit.BaseMovePoints, "\n",
	)
	var lines = len(text.Split(txt, "\n"))
	var blx, bly = battle.camera.PointFromEdge(0, 1)
	var x, y = blx + 50/battle.camera.Zoom, bly - lineHeight*float32(lines)
	battle.camera.DrawText("", txt, x, y, lineHeight, 0.95, 0, palette.Black)
	battle.camera.DrawText("", txt, x, y, lineHeight, 0.45, 0, palette.White)
}

//=================================================================

func (um *unitManager) ySortAll() []*unit.Unit {
	var ySorted = make(map[float32][]*unit.Unit, len(um.units))

	for _, unit := range um.units {
		var _, y = unit.Cell()
		ySorted[y] = append(ySorted[y], unit)
	}

	var keys = collection.MapKeys(ySorted)
	var result = make([]*unit.Unit, 0, len(ySorted))

	collection.SortNumbers(keys...)
	for _, key := range keys {
		result = append(result, ySorted[key]...)
	}

	return result
}
