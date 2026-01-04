package battle

import (
	"game/source-code/unit"
	"pure-game-kit/debug"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/flow"
	"pure-game-kit/execution/screens"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
)

type turnManager struct {
	pathMap       *geometry.ShapeGrid
	team1, team2  []*unit.Unit
	playerIsTeam1 bool

	turnOrder       []*unit.Unit
	turns           *flow.StateMachine
	turnIndex       int
	turnTakingTeam1 bool
}

func newTurnManager() *turnManager {
	return &turnManager{turns: flow.NewStateMachine()}
}

//=================================================================

func (tm *turnManager) startBattle(teamA, teamB []*unit.Unit, playerIsTeamA bool, pathMap *geometry.ShapeGrid) {
	tm.team1, tm.team2 = teamA, teamB
	tm.playerIsTeam1 = playerIsTeamA
	tm.pathMap = pathMap
	tm.turnTakingTeam1 = tm.isFirstTeam1()
	tm.turnOrder = tm.calculateTurnOrder()
	tm.turnIndex = -1

	tm.nextTurn()
}
func (tm *turnManager) update(camera *graphics.Camera, ySortedUnits []*unit.Unit) {
	tm.turns.UpdateCurrentState()

	if tm.turnIndex == -1 {
		return
	}

	var battle = screens.Current().(*BattleScreen)
	var tileW = battle.tmap.Properties[property.MapTileWidth].(int)
	var tileH = battle.tmap.Properties[property.MapTileHeight].(int)
	var unitTakingTurn = tm.turnOrder[tm.turnIndex]
	var x, y = unitTakingTurn.Position(tileW, tileH)
	var mx, my = camera.MousePosition()

	var cx, cy = unitTakingTurn.Cell()
	var walkRangeCells = tm.pathMap.MovementRange(int(cx), int(cy), float32(unitTakingTurn.Movement)/10)
	for _, cell := range walkRangeCells {
		var x, y = float32(cell[0] * tileW), float32(cell[1] * tileH)
		camera.DrawQuad(x, y, float32(tileW), float32(tileH), 0, color.FadeOut(palette.Red, 0.5))
	}

	var path = tm.pathMap.FindPathDiagonally(x, y, mx, my, false)
	camera.DrawLinesPath(6, palette.Azure, path...)
	camera.DrawPoints(3, palette.Red, path...)

	var mcx, mcy = battle.mouseCell()
	for i := len(ySortedUnits) - 1; i >= 0; i-- {
		if ySortedUnits[i].IsHovered(camera, mcx, mcy) {
			var ux, uy = ySortedUnits[i].Position(tileW, tileH)
			camera.DrawQuad(ux-float32(tileW)/2, uy-float32(tileH)/2, 64, 64, 0, palette.White)
			break
		}
	}

	camera.DrawCircle(x, y, 32, palette.Azure)
}

func (tm *turnManager) nextTurn() {
	tm.turnIndex++

	var newRound = tm.turnIndex >= len(tm.turnOrder)
	if newRound {
		tm.turnIndex = 0
		debug.Print("new round - - - - - - - - - - - - - - -")
	}

	tm.turnTakingTeam1 = collection.Contains(tm.team1, tm.turnOrder[tm.turnIndex])
	tm.turns.GoToState(condition.If(tm.isPlayerTurn(), tm.playerTurn, tm.botTurn))
}

//=================================================================

func (tm *turnManager) isPlayerTurn() bool {
	return tm.playerIsTeam1 && tm.turnTakingTeam1
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

func (tm *turnManager) calculateTurnOrder() []*unit.Unit {
	var allUnits = collection.Join(tm.team1, tm.team2)
	collection.SortByField(allUnits, func(u *unit.Unit) int { return u.Initiative })
	collection.Reverse(allUnits)
	return allUnits
}
