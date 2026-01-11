package battle

import (
	"game/code/global"
	"game/code/unit"
	"pure-game-kit/debug"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/flow"
	"pure-game-kit/execution/screens"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
)

// Handles all of the turn logic for the units during battle in the form of a state machine.
type turnManager struct {
	team1, team2  []*unit.Unit
	playerIsTeam1 bool

	order  []*unit.Unit
	states *flow.StateMachine

	curIndex int
	curTeam1 bool

	curMoveRangeCells [][2]int
	curMovePath       [][2]float32
	curMoveIndex      int
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
		for _, u := range tm.order {
			u.MovePoints = u.BaseMovePoints
		}
		debug.Print("new round - - - - - - - - - - - - - - -")
	}

	tm.curTeam1 = collection.Contains(tm.team1, tm.unitActing())
	tm.states.GoToState(condition.If(tm.isPlayerTurn(), tm.playerTurn, tm.botTurn))
	var battle = screens.Current().(*BattleScreen)
	var cx, cy = tm.unitActing().Cell()

	battle.recalculatePathMap()
	tm.curMoveRangeCells = tm.calculateRangeCells(cx, cy, float32(tm.unitActing().MovePoints)/10)
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
func (tm *turnManager) unitActing() *unit.Unit {
	return tm.order[tm.curIndex]
}

func (tm *turnManager) allies() []*unit.Unit {
	if tm.playerIsTeam1 {
		return tm.team1
	}
	return tm.team2
}
func (tm *turnManager) bots() []*unit.Unit {
	if tm.playerIsTeam1 {
		return tm.team2
	}
	return tm.team1
}
func (tm *turnManager) isBot(unit *unit.Unit) bool {
	return collection.Contains(tm.bots(), unit)
}

func (tm *turnManager) calculateTurnOrder() []*unit.Unit {
	var allUnits = collection.Join(tm.team1, tm.team2)
	collection.SortByField(allUnits, func(u *unit.Unit) int { return u.Initiative })
	collection.Reverse(allUnits)
	return allUnits
}
func (tm *turnManager) calculateRangeCells(cellX, cellY, distance float32) [][2]int {
	var battle = screens.Current().(*BattleScreen)
	return battle.pathMap.Range(int(cellX), int(cellY), distance, true)
}
func (tm *turnManager) calculateMovePoints(path [][2]float32) (possible, target int) {
	if len(path) < 2 {
		return 0, 0
	}

	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	var outOfRange = false
	for i := 1; i < len(path); i++ {
		var currX, currY = int(path[i][0] / tw), int(path[i][1] / th)
		var prevX, prevY = int(path[i-1][0] / tw), int(path[i-1][1] / th)
		var dx, dy = number.Absolute(currX - prevX), number.Absolute(currY - prevY)
		var diagonal = dx > 0 && dy > 0

		var pts = condition.If(diagonal, 15, 10)
		target += pts
		if possible+pts > tm.unitActing().MovePoints {
			outOfRange = true
		}
		if !outOfRange {
			possible += pts
		}
	}

	return possible, target
}
func (tm *turnManager) calculateMovePath(targetX, targetY float32) (possiblePts, targetPts int, path [][2]float32) {
	var battle = screens.Current().(*BattleScreen)
	var ux, uy = tm.unitActing().Position()
	path = battle.pathMap.FindPathDiagonally(ux, uy, targetX, targetY, false)

	if len(path) < 2 {
		return 0, 0, nil
	}

	var inRange = path
	for i := 1; i < len(path); i++ {
		var crop = path[:i+1]
		possiblePts, targetPts = tm.calculateMovePoints(crop)
		if targetPts > tm.unitActing().MovePoints {
			inRange = path[:i]
			break
		}
	}
	_, targetPts = tm.calculateMovePoints(path)
	return possiblePts, targetPts, inRange
}

//=================================================================

func (tm *turnManager) moveUnit() {
	var unitActing = tm.unitActing()
	var targetPos = tm.curMovePath[tm.curMoveIndex]
	unitActing.MoveTo(targetPos[0], targetPos[1])
	var ux, uy = unitActing.Position()

	if ux == targetPos[0] && uy == targetPos[1] {
		tm.curMoveIndex++
	}
	if tm.curMoveIndex >= len(tm.curMovePath) {
		var battle = screens.Current().(*BattleScreen)
		var cx, cy = unitActing.Cell()
		battle.recalculatePathMap()
		tm.curMoveRangeCells = tm.calculateRangeCells(cx, cy, float32(tm.unitActing().MovePoints)/10)
		tm.curMovePath = nil
		tm.curMoveIndex = 0
		tm.states.GoToState(condition.If(tm.isPlayerTurn(), tm.waitForAction, tm.botThink))
	}
}
