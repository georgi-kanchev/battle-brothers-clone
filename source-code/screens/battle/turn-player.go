package battle

import (
	"pure-game-kit/debug"
	"pure-game-kit/execution/screens"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/collection"
)

func (tm *turnManager) playerTurn() {
	debug.Print(tm.curIndex+1, ": player turn")
	tm.states.GoToState(tm.waitForAction)
}
func (tm *turnManager) waitForAction() {
	// var battle = screens.Current().(*BattleScreen)
	// var tileW, tileH = battle.tileSize()
	// var unit = tm.turnOrder[tm.turnIndex]
	// var mx, my = battle.camera.MousePosition()

	// for _, cell := range tm.turnWalkRangeCells {
	// 	var x, y = float32(cell[0]) * tileW, float32(cell[1]) * tileH
	// 	battle.camera.DrawQuad(x, y, float32(tileW), float32(tileH), 0, color.FadeOut(palette.Red, 0.5))

	// 	var isWalkCellHovered = int(mx/tileW) == cell[0] && int(my/tileH) == cell[1]
	// 	if isWalkCellHovered {
	// 		var path = tm.pathMap.FindPathDiagonally(x, y, mx, my, false)
	// 		battle.camera.DrawLinesPath(6, palette.Azure, path...)
	// 		battle.camera.DrawPoints(3, palette.Red, path...)
	// 		battle.camera.DrawCircle(x, y, 32, palette.Azure)

	// 		var x, y = path[len(path)-1][0], path[len(path)-1][1]
	// 		var txt = text.New(unit.CalculateMovementPoints(path, tileW, tileH), "/", unit.Movement)
	// 		battle.camera.DrawText("", txt, x, y-50/battle.camera.Zoom, 50/battle.camera.Zoom, palette.White)
	// 		break
	// 	}
	// }

	if mouse.IsButtonJustPressed(button.Left) {
		var battle = screens.Current().(*BattleScreen)
		var tileW = battle.tmap.Properties[property.MapTileWidth].(int)
		var tileH = battle.tmap.Properties[property.MapTileHeight].(int)
		var mx, my = battle.camera.MousePosition()
		var coord = [2]int{int(mx / float32(tileW)), int(my / float32(tileH))}

		if collection.Contains(tm.curWalkRangeCells, coord) {

		} else {
			debug.Print("out of walk range")
		}
	}
	if keyboard.IsKeyJustPressed(key.A) {
		tm.states.GoToState(tm.nextTurn)
	}
}
