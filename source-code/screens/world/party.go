package world

import (
	"game/source-code/global"
	"game/source-code/unit"
	"pure-game-kit/execution/screens"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/button"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/time"
)

type Party struct {
	x, y, speed,
	moveTargetX, moveTargetY float32
	isPlayer bool
	units    []*unit.Unit
}

func NewParty(units []*unit.Unit, x, y float32, isPlayer bool) *Party {
	return &Party{x: x, y: y, moveTargetX: x, moveTargetY: y, isPlayer: isPlayer, units: units, speed: 50}
}

//=================================================================

func (party *Party) Update() {
	var world = screens.Current().(*World)
	var px, py, tx, ty = party.x, party.y, party.moveTargetX, party.moveTargetY
	party.x, party.y = point.MoveToPoint(px, py, tx, ty, party.speed*time.FrameDelta()*global.TimeScale)

	if party.isPlayer {
		world.camera.X, world.camera.Y = party.x, party.y
	}

	world.camera.DrawTexture("", party.x-15, party.y-15, 30, 30, 0, palette.Cyan)

	if !party.isPlayer || world.hud.IsAnyHovered(world.camera) {
		return
	}

	world.camera.Zoom *= 1 + 0.001*mouse.ScrollSmooth()
	world.camera.Zoom = number.Limit(world.camera.Zoom, 0.1, 8)

	if mouse.IsButtonPressed(button.Left) {
		party.moveTargetX, party.moveTargetY = world.camera.MousePosition()
	}
	if mouse.IsButtonJustPressed(button.Right) {
		party.x, party.y = world.camera.MousePosition()
		party.moveTargetX, party.moveTargetY = party.x, party.y
	}
}
