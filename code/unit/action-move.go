package unit

import (
	"game/code/global"
	"pure-game-kit/execution/condition"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

type ActionMove struct {
	BasePoints, Points int
}

func NewActionMove() *ActionMove {
	return &ActionMove{BasePoints: 50, Points: 50}
}

//=================================================================

func (a *ActionMove) CalculatePoints(path [][2]float32) int {
	if len(path) < 2 {
		return 0
	}

	var totalPoints = 0
	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	for i := 1; i < len(path); i++ {
		var currX, currY = int(path[i][0] / tw), int(path[i][1] / th)
		var prevX, prevY = int(path[i-1][0] / tw), int(path[i-1][1] / th)
		var dx, dy = number.Absolute(currX - prevX), number.Absolute(currY - prevY)
		var diagonal = dx > 0 && dy > 0

		totalPoints += condition.If(diagonal, 15, 10)
	}

	return totalPoints
}

func (a *ActionMove) DrawPathToMouse(camera *graphics.Camera, unit *Unit, pathMap *geometry.ShapeGrid) {
	var mx, my = camera.MousePosition()
	var ux, uy = unit.Position()
	var path = pathMap.FindPathDiagonally(ux, uy, mx, my, false)

	if len(path) < 2 {
		return
	}

	if keyboard.IsKeyPressed(key.W) {
		print()
	}

	var inRange = path
	var outOfRange [][2]float32
	for i := 2; i < len(path); i++ {
		var crop = path[:i+1]
		var pts = a.CalculatePoints(crop)
		if pts > a.Points {
			inRange = path[:i]
			outOfRange = path[i-1:]
			break
		}
	}

	// var smooth = curve.StraightenPath(curve.SmoothPath(path))
	camera.DrawLinesPath(4, palette.Black, outOfRange...)
	camera.DrawPoints(2, palette.Black, outOfRange...)
	camera.DrawLinesPath(6, palette.Black, inRange...)
	camera.DrawPoints(3, palette.Black, inRange...)

	camera.DrawLinesPath(2, palette.Red, outOfRange...)
	camera.DrawPoints(1, palette.Red, outOfRange...)
	camera.DrawLinesPath(4, palette.Green, inRange...)
	camera.DrawPoints(2, palette.Green, inRange...)

	var x, y = path[len(path)-1][0], path[len(path)-1][1]
	var pts = a.CalculatePoints(path)
	var txt = text.New(pts, "/", unit.ActionMove.Points)
	var color = condition.If(pts <= a.Points, palette.Green, palette.Red)
	var height = 60 / camera.Zoom
	camera.DrawText("", txt, x, y, height, 0.95, palette.Black)
	camera.DrawText("", txt, x, y, height, 0.5, color)
}
