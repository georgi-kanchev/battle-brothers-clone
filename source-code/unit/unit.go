package unit

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
)

type Unit struct {
	Initiative, Movement int

	x, y float32

	head, body, plate *graphics.Sprite
}

func New() *Unit {
	return &Unit{Movement: 50, Initiative: random.Range(30, 100)}
}

//=================================================================

func (u *Unit) Spawn(x, y float32, flip bool) {
	u.head = graphics.NewSprite(assets.LoadTexture("art/Character/head.PNG"), 0, 0)
	u.body = graphics.NewSprite(assets.LoadTexture("art/Character/body.PNG"), 0, 0)
	u.plate = graphics.NewSprite(assets.LoadTexture("art/Character/plate.PNG"), 0, 0)
	u.head.PivotY = 0.85
	u.body.PivotY = 0.85
	u.plate.PivotY = 0.85

	assets.SetTextureSmoothness(u.head.AssetId, true)
	assets.SetTextureSmoothness(u.body.AssetId, true)
	assets.SetTextureSmoothness(u.plate.AssetId, true)

	u.x, u.y = x, y

	if flip {
		u.head.ScaleX = -1
		u.body.ScaleX = -1
		u.plate.ScaleX = -1
	}
}

func (u *Unit) Draw(camera *graphics.Camera, tileW, tileH float32) {
	var x, y = u.x*tileW + (tileW * 0.5), u.y*tileH + (tileH * 0.6)

	u.plate.X, u.plate.Y = x, y
	u.body.X, u.body.Y = x, y
	u.head.X, u.head.Y = x, y

	camera.DrawSprites(u.plate, u.body, u.head)
}

//=================================================================

func (u *Unit) Cell() (x, y float32) {
	return u.x, u.y
}
func (u *Unit) Position(tileW, tileH float32) (x, y float32) {
	var cx, cy = u.Cell()
	return cx*tileW + tileW/2, cy*tileH + tileH/2
}

func (u *Unit) IsHovered(camera *graphics.Camera, mouseCellX, mouseCellY float32) bool {
	var hoversX = number.IsWithin(mouseCellX, u.x+0.5, 0.5)
	var hoversY = number.IsWithin(mouseCellY, u.y+0.5, 0.5)
	return hoversX && hoversY
}

func (u *Unit) CalculateMovementPoints(path [][2]float32, tileW, tileH float32) int {
	if len(path) < 2 {
		return 0
	}

	var totalPoints = 0
	for i := 1; i < len(path); i++ {
		var currX, currY = int(path[i][0] / tileW), int(path[i][1] / tileH)
		var prevX, prevY = int(path[i-1][0] / tileW), int(path[i-1][1] / tileH)
		var dx, dy = number.Absolute(currX - prevX), number.Absolute(currY - prevY)
		var diagonal = dx > 0 && dy > 0

		totalPoints += condition.If(diagonal, 15, 10)
	}

	return totalPoints
}
