package battle

import (
	"game/source-code/unit"
	"pure-game-kit/execution/flow"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color/palette"
)

type turnManager struct {
	turns          *flow.Sequence
	unitTakingTurn *unit.Unit

	team1IsPlayer, team1Takingturn bool

	pathMap      *geometry.ShapeGrid
	team1, team2 []*unit.Unit
}

func newTurnManager() *turnManager {
	var result = &turnManager{turns: flow.NewSequence()}
	result.turns.SetSteps(
		flow.NowDo(result.newTurnStart),
		flow.NowDoAndKeepRepeating(result.waitForAction),
	)
	return result
}

//=================================================================

func (tm *turnManager) startBattle(teamA, teamB []*unit.Unit, playerAttacks bool, pathMap *geometry.ShapeGrid) {
	tm.team1, tm.team2 = teamA, teamB
	tm.team1IsPlayer = playerAttacks
	tm.pathMap = pathMap
	tm.team1Takingturn = true
	tm.unitTakingTurn = teamA[0]
}
func (tm *turnManager) update(camera *graphics.Camera, tileW, tileH int) {
	var x, y = tm.unitTakingTurn.Position()
	x, y = x*float32(tileW)+float32(tileW/2), y*float32(tileH)+float32(tileH/2)
	camera.DrawCircle(x, y, 32, palette.White)
	var mx, my = camera.MousePosition()

	var path = tm.pathMap.FindPathSmoothly(x, y, mx, my, false)
	camera.DrawLinesPath(10, palette.Azure, path...)
	camera.DrawPoints(5, palette.Red, path...)

	tm.turns.Update()
}

func (tm *turnManager) newTurnStart() {
	tm.unitTakingTurn.RecalculateWalkRange(tm.pathMap)
}
func (tm *turnManager) waitForAction() {

}
func (tm *turnManager) isPlayerTurn() bool {
	return tm.team1IsPlayer && tm.team1Takingturn
}
