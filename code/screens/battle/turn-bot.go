package battle

import "pure-game-kit/debug"

// States for the Bot, handled by the turn manager.

func (tm *turnManager) botTurn() {
	debug.Print(tm.curIndex+1, ": bot turn")
	tm.states.GoToState(tm.botThink)
}
func (tm *turnManager) botThink() {
	if tm.states.StateTimer() > 1 {
		tm.nextTurn()
	}
}
