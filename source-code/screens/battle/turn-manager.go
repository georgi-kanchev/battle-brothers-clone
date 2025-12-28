package battle

import (
	"game/source-code/unit"
	"pure-game-kit/execution/flow"
	"pure-game-kit/geometry"
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
	var result = &turnManager{turns: flow.NewSequence()}
	result.turns.SetSteps(false,
		flow.NowDoAndKeepRepeating(result.waitForAction),
	)
	return result
}

//=================================================================

func (tm *turnManager) startBattle(teamA, teamB []*unit.Unit, playerAttacks bool, pathMap *geometry.ShapeGrid) {
	tm.teamA, tm.teamB = teamA, teamB
	tm.playerIsTeamA = playerAttacks
	tm.pathMap = pathMap
	tm.takingTurnTeamA = true
	tm.takingTurnUnit = teamA[0]

	tm.turns.Run()
}
func (tm *turnManager) waitForAction() {
}
func (tm *turnManager) isPlayerTurn() bool {
	return tm.playerIsTeamA && tm.takingTurnTeamA
}
