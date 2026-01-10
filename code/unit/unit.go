package unit

import (
	"game/code/global"
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
)

type Unit struct {
	Initiative int

	ActionMove *ActionMove

	x, y float32

	head, body, plate *graphics.Sprite
}

func New() *Unit {
	return &Unit{ActionMove: NewActionMove(), Initiative: random.Range(30, 100)}
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

func (u *Unit) Draw(camera *graphics.Camera) {
	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	var x, y = u.x*tw + (tw * 0.5), u.y*th + (th * 0.6)

	u.plate.X, u.plate.Y = x, y
	u.body.X, u.body.Y = x, y
	u.head.X, u.head.Y = x, y

	camera.DrawSprites(u.plate, u.body, u.head)
}

//=================================================================

func (u *Unit) Cell() (x, y float32) {
	return u.x, u.y
}
func (u *Unit) Position() (x, y float32) {
	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	var cx, cy = u.Cell()
	return cx*tw + tw/2, cy*th + th/2
}

func (u *Unit) IsHovered(camera *graphics.Camera, mouseCellX, mouseCellY float32) bool {
	var hoversX = number.IsWithin(mouseCellX, u.x+0.5, 0.5)
	var hoversY = number.IsWithin(mouseCellY, u.y+0.5, 0.5)
	return hoversX && hoversY
}
