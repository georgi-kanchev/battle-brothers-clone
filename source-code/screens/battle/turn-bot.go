package battle

import "pure-game-kit/debug"

func (tm *turnManager) botTurn() {
	debug.Print(tm.turnIndex+1, ": bot turn")
	tm.turns.GoToState(tm.botThink)
}
func (tm *turnManager) botThink() {
	if tm.turns.StateTimer() > 3 {
		tm.nextTurn()
	}
}
