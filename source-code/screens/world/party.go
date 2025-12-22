package world

import (
	"game/source-code/global"
	"game/source-code/unit"
	"pure-game-kit/execution/screens"
	"pure-game-kit/geometry"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/utility/angle"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/time"
)

type Party struct {
	x, y, speed, moveTargetX, moveTargetY        float32
	isPlayer, isGoingToSettlement, isGoingToRoad bool
	units                                        []*unit.Unit
	shape                                        *geometry.Shape
}

func NewParty(units []*unit.Unit, x, y float32, isPlayer bool) *Party {
	return &Party{x: x, y: y, moveTargetX: x, moveTargetY: y, isPlayer: isPlayer, units: units, speed: 20,
		shape: geometry.NewShapeRectangle(10, 10, 0.5, 0.5)}
}

//=================================================================

func (party *Party) Update() {
	var world = screens.Current().(*World)
	party.handleMovement()
	world.camera.DrawShapes(palette.Red, party.shape.CornerPoints()...)

	if party.isPlayer {
		party.handlePlayer()
	}
}

//=================================================================
// private

func (party *Party) handleMovement() {
	var px, py, tx, ty = party.x, party.y, party.moveTargetX, party.moveTargetY
	var angle = angle.BetweenPoints(px, py, tx, ty)
	var speed = party.speed * time.FrameDelta() * global.TimeScale
	var velX, velY = point.MoveAtAngle(0, 0, angle, speed)

	party.shape.X, party.shape.Y = party.x, party.y
	var newVelX, newVelY = party.collideWithBarrier(velX, velY)
	var newSpeed = point.DistanceToPoint(0, 0, velX, velY)
	party.x, party.y = party.x+newVelX, party.y+newVelY

	if point.DistanceToPoint(party.x, party.y, tx, ty) < newSpeed*3 {
		party.x, party.y = tx, ty
	}
}
func (party *Party) collideWithBarrier(velX, velY float32) (newVelX, newVelY float32) {
	var world = screens.Current().(*World)
	newVelX, newVelY = velX, velY
	var x, y = party.shape.Collide(velX, velY, world.solids...)
	newVelX, newVelY = newVelX+x, newVelY+y
	return newVelX, newVelY
}
func (party *Party) handlePlayer() {
	var world = screens.Current().(*World)
	world.camera.X, world.camera.Y = party.x, party.y

	if world.hud.IsAnyHovered(world.camera) {
		return
	}

	world.camera.Zoom *= 1 + 0.001*mouse.ScrollSmooth()
	world.camera.Zoom = number.Limit(world.camera.Zoom, 0.1, 8)

	var cx, cy = world.camera.MousePosition()
	var dist = point.DistanceToPoint(party.x, party.y, cx, cy)
	if mouse.IsButtonPressed(button.Left) && dist > 10 {
		party.moveTargetX, party.moveTargetY = cx, cy
	}
}
