package battle

import "pure-game-kit/debug"

// States for the Bot, handled by the turn logic.

func (bs *BattleScreen) botTurn() {
	debug.Print(bs.curTurnIndex+1, ": bot turn")
	bs.states.GoToState(bs.botThink)
}
func (bs *BattleScreen) botThink() {
	if bs.states.StateTimer() > 1 {
		bs.nextTurn()
	}
}
