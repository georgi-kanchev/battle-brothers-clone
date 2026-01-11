package battle

import (
	"pure-game-kit/debug"
	"pure-game-kit/execution/screens"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/motion/curve"
)

// States for the Player, handled by the turn manager.

func (tm *turnManager) playerTurn() {
	debug.Print(tm.curIndex+1, ": player turn")
	tm.states.GoToState(tm.waitForAction)
}
func (tm *turnManager) waitForAction() {
	if mouse.IsButtonJustPressed(button.Left) {
		var battle = screens.Current().(*BattleScreen)
		var mx, my = battle.camera.MousePosition()
		var pts, _, path = tm.calculateMovePath(mx, my)

		if pts > 0 && battle.unitManager.hoveredUnit == nil {
			tm.unitActing().MovePoints -= pts
			tm.curMovePath = curve.StraightenPath(path)
			tm.curMovePath = curve.SmoothPath(path)
			tm.curMoveIndex = 0
			tm.curMoveRangeCells = nil
			tm.states.GoToState(tm.moveUnit)
		}
	}
	if keyboard.IsKeyJustPressed(key.A) {
		tm.states.GoToState(tm.nextTurn)
	}
}
