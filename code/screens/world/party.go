package world

import (
	"game/code/global"
	"game/code/unit"
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
	isPlayer, isUsingRoads, isResting     bool

	goingToSettlement *tiled.Object

	units  []*unit.Unit
	hitbox *geometry.Shape

	path [][2]float32
}

func NewParty(units []*unit.Unit, x, y float32, isPlayer bool) *Party {
	return &Party{x: x, y: y, moveTargetX: x, moveTargetY: y, isPlayer: isPlayer, units: units, speed: 20,
		hitbox: geometry.NewShapeQuad(10, 10, 0.5, 0.5)}
}

//=================================================================

func (p *Party) Update() {
	var world = screens.Current().(*WorldScreen)
	var isInRoadRange = p.isInRoadRange()
	p.handleMovement(isInRoadRange)

	if p.isPlayer {
		p.handlePlayer()
	}

	p.tryEnterSettlement()

	if !p.isResting {
		world.camera.DrawShapes(palette.Red, p.hitbox.CornerPoints()...)
	}
}

//=================================================================
// private

func (p *Party) handleMovement(isInRoadRange bool) {
	if p.isResting {
		return
	}

	if p.isUsingRoads && len(p.path) > 0 {
		p.moveTargetX, p.moveTargetY = p.path[0][0], p.path[0][1]
	}

	var px, py, tx, ty = p.x, p.y, p.moveTargetX, p.moveTargetY
	var angle = angle.BetweenPoints(px, py, tx, ty)
	var speed = p.speed * time.FrameDelta() * global.TimeScale

	if isInRoadRange {
		speed *= 2
	}

	var velX, velY = point.MoveAtAngle(0, 0, angle, speed)
	p.hitbox.X, p.hitbox.Y = p.x, p.y
	var newVelX, newVelY = p.collideWithSolid(velX, velY)
	var newSpeed = point.DistanceToPoint(0, 0, velX, velY)
	p.x, p.y = p.x+newVelX, p.y+newVelY
	var dist = point.DistanceToPoint(p.x, p.y, tx, ty)

	if dist < newSpeed*3 {
		p.x, p.y = tx, ty

		if p.isUsingRoads {
			p.path = collection.RemoveAt(p.path, 0)
		}
	}
}
func (p *Party) collideWithSolid(velX, velY float32) (newVelX, newVelY float32) {
	var world = screens.Current().(*WorldScreen)
	newVelX, newVelY = velX, velY
	var x, y = p.hitbox.Collide(velX, velY, world.solids...)
	newVelX, newVelY = newVelX+x, newVelY+y
	return newVelX, newVelY
}
func (p *Party) tryEnterSettlement() {
	var world = screens.Current().(*WorldScreen)
	if p.isResting || p.goingToSettlement == nil || world.currentPopup != nil {
		return
	}

	for _, s := range world.settlements.Objects {
		if p.goingToSettlement == s && p.hitbox.IsOverlappingShapes(s.ExtractShapes()...) {
			p.moveTargetX, p.moveTargetY = p.x, p.y
			p.path = nil
			world.resultingCursorNonGUI = -1
			world.currentPopup = global.TogglePopup(world.hud, world.currentPopup, world.settlement)
		}
	}
}

func (party *Party) handlePlayer() {
	var world = screens.Current().(*WorldScreen)
	world.camera.X, world.camera.Y = party.x, party.y

	if !party.isResting {
		var col = palette.White
		var p = party.path
		if len(p) > 0 {
			world.camera.DrawLine(party.x, party.y, p[0][0], p[0][1], 2, col)
			world.camera.DrawLinesPath(2, col, p...)
		}
		var mx, my = party.lastPathPoint()
		world.camera.DrawPoints(4, col, [2]float32{mx, my})
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

	if world.hud.IsAnyHovered(world.camera) {
		return
	}
	world.resultingCursorNonGUI = -1

	var mx, my = world.camera.MousePosition()
	var settlements = world.settlements.Objects
	for _, s := range settlements {
		var shape = s.ExtractShapes()[0]
		var hovering = shape.IsContainingPoint(mx, my)
		if hovering || shape.IsContainingPoint(party.moveTargetX, party.moveTargetY) {
			var pts = shape.CornerPoints()
			world.camera.DrawShapes(color.FadeOut(palette.White, 0.8), pts...)
			world.camera.DrawLinesPath(2, color.FadeOut(palette.White, 0.5), pts...)
		}

		if !hovering {
			continue
		}

		world.resultingCursorNonGUI = cursor.Hand
		if mouse.IsButtonJustPressed(button.Left) {
			party.goingToSettlement = s
		}
	}

	world.camera.Zoom *= 1 + 0.001*mouse.ScrollSmooth()
	world.camera.Zoom = number.Limit(world.camera.Zoom, 0.1, 8)

	var dist = point.DistanceToPoint(party.x, party.y, mx, my)
	if mouse.IsButtonPressed(button.Left) && dist > 10 {
		party.moveTargetX, party.moveTargetY = mx, my

		if party.isUsingRoads && mouse.IsButtonJustPressed(button.Left) {
			party.path = geometry.FollowPaths(party.x, party.y, party.moveTargetX, party.moveTargetY, world.roads...)
		}
	}
}

func (p *Party) isInRoadRange() bool {
	var world = screens.Current().(*WorldScreen)
	for i := 1; i < len(world.roads); i++ {
		var ax, ay = world.roads[i-1][0], world.roads[i-1][1]
		var bx, by = world.roads[i][0], world.roads[i][1]
		var line = geometry.NewLine(ax, ay, bx, by)
		var closestX, closestY = line.ClosestToPoint(p.x, p.y)
		var distance = point.DistanceToPoint(p.x, p.y, closestX, closestY)
		if distance < 15 {
			return true
		}
	}
	return false
}
func (p *Party) lastPathPoint() (x, y float32) {
	if len(p.path) > 0 {
		return p.path[len(p.path)-1][0], p.path[len(p.path)-1][1]
	}
	return p.moveTargetX, p.moveTargetY
}
