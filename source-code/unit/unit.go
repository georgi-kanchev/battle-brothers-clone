package unit

import gfx "pure-game-kit/graphics"

type Unit struct {
	x, y float32
}

func New() *Unit {
	return &Unit{}
}

//=================================================================

func (unit *Unit) Draw(camera *gfx.Camera) {
}

//=================================================================

func (unit *Unit) Position() (x, y float32) {
	return unit.x, unit.y
}
