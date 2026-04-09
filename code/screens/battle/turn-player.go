package battle

import (
	"pure-game-kit/debug"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/motion/curve"
)

// States for the Player, handled by the turn logic.

func (bs *BattleScreen) playerTurn() {
	debug.Print(bs.curTurnIndex+1, ": player turn")
	bs.states.GoToState(bs.waitForAction)
}
func (bs *BattleScreen) waitForAction() {
	if mouse.IsButtonJustPressed(button.Left) {
		var mx, my = bs.camera.MousePosition()
		var pts, _, path = bs.calculateMovePath(mx, my)

		if pts > 0 && bs.hoveredUnit == nil {
			bs.actingUnit().MovePoints -= pts
			bs.curMovePath = curve.StraightenPath(path...)
			bs.curMovePath = curve.SmoothPath(path...)
			bs.curMoveIndex = 0
			bs.curMoveRangeCells = nil
			bs.states.GoToState(bs.moveUnit)
		}
	}
	if keyboard.IsKeyJustPressed(key.A) {
		bs.states.GoToState(bs.nextTurn)
	}
}
