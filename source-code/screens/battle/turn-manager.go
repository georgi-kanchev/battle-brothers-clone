package battle

import (
	"game/source-code/unit"
	"pure-game-kit/debug"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/flow"
	"pure-game-kit/execution/screens"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
)

type turnManager struct {
	team1, team2  []*unit.Unit
	playerIsTeam1 bool

	order  []*unit.Unit
	states *flow.StateMachine

	curIndex int
	curTeam1 bool

	curWalkRangeCells [][2]int
}

func newTurnManager() *turnManager {
	return &turnManager{states: flow.NewStateMachine()}
}

//=================================================================

func (tm *turnManager) startBattle(teamA, teamB []*unit.Unit, playerIsTeamA bool) {
	tm.team1, tm.team2 = teamA, teamB
	tm.playerIsTeam1 = playerIsTeamA
	tm.curTeam1 = tm.isFirstTeam1()
	tm.order = tm.calculateTurnOrder()
	tm.curIndex = -1

	tm.nextTurn()
}

func (tm *turnManager) nextTurn() {
	tm.curIndex++

	var newRound = tm.curIndex >= len(tm.order)
	if newRound {
		tm.curIndex = 0
		debug.Print("new round - - - - - - - - - - - - - - -")
	}

	tm.curTeam1 = collection.Contains(tm.team1, tm.unit())
	tm.states.GoToState(condition.If(tm.isPlayerTurn(), tm.playerTurn, tm.botTurn))
	var cx, cy = tm.unit().Cell()
	var battle = screens.Current().(*BattleScreen)
	battle.recalculatePathMap()
	tm.curWalkRangeCells = battle.pathMap.Range(int(cx), int(cy), float32(tm.unit().Movement)/10, true)
}

//=================================================================

func (tm *turnManager) isPlayerTurn() bool {
	return tm.playerIsTeam1 && tm.curTeam1
}
func (tm *turnManager) isFirstTeam1() bool {
	var initiativesTeam1, initiativesTeam2 []float32

	for _, unit := range tm.team1 {
		initiativesTeam1 = append(initiativesTeam1, float32(unit.Initiative))
	}
	for _, unit := range tm.team2 {
		initiativesTeam2 = append(initiativesTeam2, float32(unit.Initiative))
	}

	var avg1, avg2 = number.Average(initiativesTeam1...), number.Average(initiativesTeam2...)
	if avg1 == avg2 {
		return random.Pick(true, false)
	}

	return avg1 > avg2
}
func (tm *turnManager) unit() *unit.Unit {
	return tm.order[tm.curIndex]
}

func (tm *turnManager) calculateTurnOrder() []*unit.Unit {
	var allUnits = collection.Join(tm.team1, tm.team2)
	collection.SortByField(allUnits, func(u *unit.Unit) int { return u.Initiative })
	collection.Reverse(allUnits)
	return allUnits
}
