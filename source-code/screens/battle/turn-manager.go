package battle

import (
	"game/source-code/unit"
	"pure-game-kit/execution/flow"
	"pure-game-kit/geometry"
	"pure-game-kit/utility/number"
)

type turnManager struct {
	playerIsTeamA   bool
	turns           *flow.Sequence
	takingTurnUnit  *unit.Unit
	takingTurnTeamA bool

	pathMap      *geometry.ShapeGrid
	teamA, teamB []*unit.Unit
}

func newTurnManager() *turnManager {
	return &turnManager{}
}

//=================================================================

func (tm *turnManager) startBattle(teamA, teamB []*unit.Unit, playerAttacks bool, pathMap *geometry.ShapeGrid) {
	tm.teamA, tm.teamB = teamA, teamB
	tm.playerIsTeamA = playerAttacks
	tm.pathMap = pathMap
	tm.takingTurnTeamA = true
	tm.takingTurnUnit = teamA[0]

	tm.turns = flow.NewSequence()
	tm.turns.SetSteps(false,
		flow.NowDoLoop(number.ValueMaximum[int](), func(i int) { tm.waitForAction() }),
	)
}
func (tm *turnManager) waitForAction() {
	if tm.isPlayerTurn() {
		return
	}
}
func (tm *turnManager) isPlayerTurn() bool {
	return tm.playerIsTeamA && tm.takingTurnTeamA
}
