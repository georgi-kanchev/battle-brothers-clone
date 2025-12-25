package world

import (
	"game/source-code/global"
	"game/source-code/unit"
	"pure-game-kit/execution/screens"
	"pure-game-kit/geometry"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/tiled"
	"pure-game-kit/utility/angle"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/time"
)

type Party struct {
	x, y, speed, moveTargetX, moveTargetY float32
	isPlayer, isUsingRoads                bool

	goingToSettlement *tiled.Object

	units  []*unit.Unit
	hitbox *geometry.Shape

	path [][2]float32
}

func NewParty(units []*unit.Unit, x, y float32, isPlayer bool) *Party {
	return &Party{x: x, y: y, moveTargetX: x, moveTargetY: y, isPlayer: isPlayer, units: units, speed: 20,
		hitbox: geometry.NewShapeRectangle(10, 10, 0.5, 0.5)}
}

//=================================================================

func (party *Party) Update() {
	var world = screens.Current().(*World)
	var isInRoadRange = party.isInRoadRange()
	party.handleMovement(isInRoadRange)
	party.tryEnterSettlement()

	if party.isPlayer {
		party.handlePlayer()
	}

	world.camera.DrawShapes(palette.Red, party.hitbox.CornerPoints()...)
}

//=================================================================
// private

func (party *Party) handleMovement(isInRoadRange bool) {
	if party.isUsingRoads && len(party.path) > 0 {
		party.moveTargetX, party.moveTargetY = party.path[0][0], party.path[0][1]
	}

	var px, py, tx, ty = party.x, party.y, party.moveTargetX, party.moveTargetY
	var angle = angle.BetweenPoints(px, py, tx, ty)
	var speed = party.speed * time.FrameDelta() * global.TimeScale

	if isInRoadRange {
		speed *= 2
	}

	var velX, velY = point.MoveAtAngle(0, 0, angle, speed)
	party.hitbox.X, party.hitbox.Y = party.x, party.y
	var newVelX, newVelY = party.collideWithSolid(velX, velY)
	var newSpeed = point.DistanceToPoint(0, 0, velX, velY)
	party.x, party.y = party.x+newVelX, party.y+newVelY
	var dist = point.DistanceToPoint(party.x, party.y, tx, ty)

	if dist < newSpeed*3 {
		party.x, party.y = tx, ty

		if party.isUsingRoads {
			party.path = collection.RemoveAt(party.path, 0)
		}
	}
}
func (party *Party) collideWithSolid(velX, velY float32) (newVelX, newVelY float32) {
	var world = screens.Current().(*World)
	newVelX, newVelY = velX, velY
	var x, y = party.hitbox.Collide(velX, velY, world.solids...)
	newVelX, newVelY = newVelX+x, newVelY+y
	return newVelX, newVelY
}
func (party *Party) tryEnterSettlement() {
	var world = screens.Current().(*World)
	if party.goingToSettlement == nil || world.currentPopup == world.settlement {
		return
	}

	for _, s := range world.settlements.Objects {
		if party.goingToSettlement == s && party.hitbox.IsOverlappingShapes(s.ExtractShapes()...) {
			party.moveTargetX, party.moveTargetY = party.x, party.y
			party.path = nil
			party.goingToSettlement = s
			world.resultingCursorNonGUI = -1
			world.currentPopup = global.TogglePopup(world.hud, world.currentPopup, world.settlement)
		}
	}
}

func (party *Party) handlePlayer() {
	var world = screens.Current().(*World)
	world.camera.X, world.camera.Y = party.x, party.y

	var col = palette.White
	var p = party.path
	if len(p) > 0 {
		world.camera.DrawLine(party.x, party.y, p[0][0], p[0][1], 2, col)
		world.camera.DrawLinesPath(2, col, p...)
	}
	var mx, my = party.lastPathPoint()
	world.camera.DrawPoints(4, col, [2]float32{mx, my})

	if world.hud.IsAnyHovered(world.camera) {
		return
	}

	world.resultingCursorNonGUI = -1

	var settlements = world.settlements.Objects
	var cx, cy = world.camera.MousePosition()
	for _, s := range settlements {
		var shape = s.ExtractShapes()[0]
		if shape.IsContainingPoint(cx, cy) || shape.IsContainingPoint(party.moveTargetX, party.moveTargetY) {
			var pts = shape.CornerPoints()
			world.camera.DrawShapes(color.FadeOut(palette.White, 0.8), pts...)
			world.camera.DrawLinesPath(2, color.FadeOut(palette.White, 0.5), pts...)
			world.resultingCursorNonGUI = cursor.Hand

			if mouse.IsButtonJustPressed(button.Left) {
				party.goingToSettlement = s
			}
		}
	}

	world.camera.Zoom *= 1 + 0.001*mouse.ScrollSmooth()
	world.camera.Zoom = number.Limit(world.camera.Zoom, 0.1, 8)

	var dist = point.DistanceToPoint(party.x, party.y, cx, cy)
	if mouse.IsButtonPressed(button.Left) && dist > 10 {
		party.moveTargetX, party.moveTargetY = cx, cy

		if party.isUsingRoads && mouse.IsButtonJustPressed(button.Left) {
			party.path = geometry.FollowPaths(party.x, party.y, party.moveTargetX, party.moveTargetY, world.roads...)
		}
	}
	if keyboard.IsKeyJustPressed(key.Enter) {
		party.isUsingRoads = !party.isUsingRoads
		var standingStill = party.x == party.moveTargetX && party.y == party.moveTargetY

		if party.isUsingRoads && !standingStill {
			party.path = geometry.FollowPaths(party.x, party.y, party.moveTargetX, party.moveTargetY, world.roads...)
		} else if !party.isUsingRoads && len(party.path) > 0 {
			party.moveTargetX, party.moveTargetY = party.path[len(party.path)-1][0], party.path[len(party.path)-1][1]
			party.path = nil
		}
	}
}

func (party *Party) isInRoadRange() bool {
	var world = screens.Current().(*World)
	for i := 1; i < len(world.roads); i++ {
		var ax, ay = world.roads[i-1][0], world.roads[i-1][1]
		var bx, by = world.roads[i][0], world.roads[i][1]
		var line = geometry.NewLine(ax, ay, bx, by)
		var closestX, closestY = line.ClosestToPoint(party.x, party.y)
		var distance = point.DistanceToPoint(party.x, party.y, closestX, closestY)
		if distance < 15 {
			return true
		}
	}
	return false
}
func (party *Party) lastPathPoint() (x, y float32) {
	if len(party.path) > 0 {
		return party.path[len(party.path)-1][0], party.path[len(party.path)-1][1]
	}
	return party.moveTargetX, party.moveTargetY
}
