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
	debug.Print(tm.turnIndex+1, ": player turn")
	tm.turns.GoToState(tm.waitForAction)
}
func (tm *turnManager) waitForAction() {
	if mouse.IsButtonJustPressed(button.Left) {
		var battle = screens.Current().(*BattleScreen)
		var tileW = battle.tmap.Properties[property.MapTileWidth].(int)
		var tileH = battle.tmap.Properties[property.MapTileHeight].(int)
		var mx, my = battle.camera.MousePosition()
		var coord = [2]int{int(mx / float32(tileW)), int(my / float32(tileH))}
		var unit = tm.turnOrder[tm.turnIndex]
		var ux, uy = unit.Cell()
		var walkRangeCells = tm.pathMap.MovementRange(int(ux), int(uy), float32(unit.Movement)/10)

		if collection.Contains(walkRangeCells, coord) {
			debug.Print("in walk range")
		} else {
			debug.Print("out of walk range")
		}
	}
	if keyboard.IsKeyJustPressed(key.A) {
		tm.turns.GoToState(tm.nextTurn)
	}
}
