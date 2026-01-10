package battle

import "pure-game-kit/debug"

func (tm *turnManager) botTurn() {
	debug.Print(tm.curIndex+1, ": bot turn")
	tm.states.GoToState(tm.botThink)
}
func (tm *turnManager) botThink() {
	if tm.states.StateTimer() > 1 {
		tm.nextTurn()
	}
}
