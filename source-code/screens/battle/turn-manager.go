package battle

import (
	"game/source-code/unit"
	"pure-game-kit/execution/flow"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
)

type turnManager struct {
	turns          *flow.StateMachine
	unitTakingTurn *unit.Unit

	team1IsPlayer, team1Takingturn bool

	pathMap      *geometry.ShapeGrid
	team1, team2 []*unit.Unit

	currentTurnOrder []*unit.Unit
	walkRangeCells   [][2]int // relative to the unit cell position
}

func newTurnManager() *turnManager {
	return &turnManager{turns: flow.NewStateMachine()}
}

//=================================================================

func (tm *turnManager) startBattle(teamA, teamB []*unit.Unit, playerAttacks bool, pathMap *geometry.ShapeGrid) {
	tm.team1, tm.team2 = teamA, teamB
	tm.team1IsPlayer = playerAttacks
	tm.pathMap = pathMap
	tm.team1Takingturn = true
	tm.currentTurnOrder = tm.calculateTurnOrder()
	tm.unitTakingTurn = tm.currentTurnOrder[0]
	tm.turns.GoToState(tm.newTurn)
}
func (tm *turnManager) update(camera *graphics.Camera, tileW, tileH int) {
	var x, y = tm.unitTakingTurn.Position()
	x, y = x*float32(tileW)+float32(tileW/2), y*float32(tileH)+float32(tileH/2)
	camera.DrawCircle(x, y, 32, palette.White)
	var mx, my = camera.MousePosition()

	for _, cell := range tm.walkRangeCells {
		var x, y = float32(cell[0] * tileW), float32(cell[1] * tileH)
		camera.DrawQuad(x, y, float32(tileW), float32(tileH), 0, color.FadeOut(palette.Red, 0.5))
	}

	var path = tm.pathMap.FindPathSmoothly(x, y, mx, my, false)
	camera.DrawLinesPath(10, palette.Azure, path...)
	camera.DrawPoints(5, palette.Red, path...)

	tm.turns.UpdateCurrentState()
}

//=================================================================

func (tm *turnManager) isPlayerTurn() bool {
	return tm.team1IsPlayer && tm.team1Takingturn
}

func (tm *turnManager) calculateTurnOrder() []*unit.Unit {
	var allUnits = collection.Join(tm.team1, tm.team2)
	collection.SortByField(allUnits, func(u *unit.Unit) int { return u.Initiative })
	collection.Reverse(allUnits)
	return allUnits
}

//=================================================================
// states

func (tm *turnManager) newTurn() {
	var x, y = tm.unitTakingTurn.Position()
	tm.walkRangeCells = tm.pathMap.MovementRange(int(x), int(y), tm.unitTakingTurn.MaxMoveCells)
	tm.currentTurnOrder = tm.calculateTurnOrder()
	tm.turns.GoToState(tm.waitForAction)
}
func (tm *turnManager) waitForAction() {

}
