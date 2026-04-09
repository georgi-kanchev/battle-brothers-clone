package battle

import (
	"game/code/global"
	"game/code/unit"
	"pure-game-kit/debug"
	"pure-game-kit/execution/condition"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
)

func (bs *BattleScreen) startBattle() {
	bs.curTurnIsTeam1 = bs.team1GoesFirst()
	bs.order = bs.calculateTurnOrder()
	bs.curTurnIndex = -1

	bs.nextTurn()
}

func (bs *BattleScreen) nextTurn() {
	bs.curTurnIndex++

	var newRound = bs.curTurnIndex >= len(bs.order)
	if newRound {
		bs.curTurnIndex = 0
		for _, u := range bs.order {
			u.MovePoints = u.BaseMovePoints
		}
		debug.Print("new round - - - - - - - - - - - - - - -")
	}

	bs.curTurnIsTeam1 = collection.Contains(bs.team1, bs.actingUnit())
	bs.states.GoToState(condition.If(bs.isPlayerTurn(), bs.playerTurn, bs.botTurn))
	var cx, cy = bs.actingUnit().Cell()

	bs.recalculatePathMap()
	bs.curMoveRangeCells = bs.calculateRangeCells(cx, cy, float32(bs.actingUnit().MovePoints)/10)
}

//=================================================================

func (bs *BattleScreen) isPlayerTurn() bool {
	return bs.playerIsTeam1 && bs.curTurnIsTeam1
}
func (bs *BattleScreen) team1GoesFirst() bool {
	var initiativesTeam1, initiativesTeam2 []float32

	for _, u := range bs.team1 {
		initiativesTeam1 = append(initiativesTeam1, float32(u.Initiative))
	}
	for _, u := range bs.team2 {
		initiativesTeam2 = append(initiativesTeam2, float32(u.Initiative))
	}

	var avg1, avg2 = number.Average(initiativesTeam1...), number.Average(initiativesTeam2...)
	if avg1 == avg2 {
		return random.Pick(true, false)
	}

	return avg1 > avg2
}
func (bs *BattleScreen) actingUnit() *unit.Unit {
	return bs.order[bs.curTurnIndex]
}

func (bs *BattleScreen) alliedUnits() []*unit.Unit {
	if bs.playerIsTeam1 {
		return bs.team1
	}
	return bs.team2
}
func (bs *BattleScreen) botUnits() []*unit.Unit {
	if bs.playerIsTeam1 {
		return bs.team2
	}
	return bs.team1
}
func (bs *BattleScreen) isBotUnit(u *unit.Unit) bool {
	return collection.Contains(bs.botUnits(), u)
}

func (bs *BattleScreen) calculateTurnOrder() []*unit.Unit {
	var allUnits = collection.Join(bs.team1, bs.team2)
	collection.SortByField(allUnits, func(u *unit.Unit) int { return u.Initiative })
	collection.Reverse(allUnits)
	return allUnits
}
func (bs *BattleScreen) calculateRangeCells(cellX, cellY, distance float32) [][2]int {
	return bs.pathMap.Range(int(cellX), int(cellY), distance, true)
}
func (bs *BattleScreen) calculateMovePoints(path []float32) (possible, target int) {
	if len(path) < 2 {
		return 0, 0
	}

	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	var outOfRange = false
	for i := 2; i < len(path); i += 2 {
		var prevX, prevY = int(path[i-2] / tw), int(path[i-1] / th)
		var currX, currY = int(path[i] / tw), int(path[i+1] / th)
		var dx, dy = number.Absolute(currX - prevX), number.Absolute(currY - prevY)
		var diagonal = dx > 0 && dy > 0

		var pts = condition.If(diagonal, 15, 10)
		target += pts
		if possible+pts > bs.actingUnit().MovePoints {
			outOfRange = true
		}
		if !outOfRange {
			possible += pts
		}
	}

	return possible, target
}
func (bs *BattleScreen) calculateMovePath(targetX, targetY float32) (possiblePts, targetPts int, path []float32) {
	var ux, uy = bs.actingUnit().Position()
	path = bs.pathMap.FindPathDiagonally(ux, uy, targetX, targetY, false)

	if len(path) < 2 {
		return 0, 0, nil
	}

	var inRange = path
	for i := 1; i < len(path); i += 2 {
		var crop = path[:i+1]
		possiblePts, targetPts = bs.calculateMovePoints(crop)
		if targetPts > bs.actingUnit().MovePoints {
			inRange = path[:i-1]
			break
		}
	}
	_, targetPts = bs.calculateMovePoints(path)
	return possiblePts, targetPts, inRange
}

//=================================================================

func (bs *BattleScreen) moveUnit() {
	var unitActing = bs.actingUnit()
	var targetX = bs.curMovePath[bs.curMoveIndex]
	var targetY = bs.curMovePath[bs.curMoveIndex+1]
	unitActing.MoveTo(targetX, targetY)
	var ux, uy = unitActing.Position()

	if ux == targetX && uy == targetY {
		bs.curMoveIndex += 2
	}
	if bs.curMoveIndex >= len(bs.curMovePath) {
		var cx, cy = unitActing.Cell()
		bs.recalculatePathMap()
		bs.curMoveRangeCells = bs.calculateRangeCells(cx, cy, float32(bs.actingUnit().MovePoints)/10)
		bs.curMovePath = nil
		bs.curMoveIndex = 0
		bs.states.GoToState(condition.If(bs.isPlayerTurn(), bs.waitForAction, bs.botThink))
	}
}
