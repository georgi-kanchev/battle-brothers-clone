package unit

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/random"
	"pure-game-kit/utility/text"
)

type Unit struct {
	x, y, MaxMoveCells float32

	Initiative int

	head, body, plate *graphics.Sprite
}

func New() *Unit {
	return &Unit{MaxMoveCells: 5, Initiative: random.Range(30, 100)}
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

func (u *Unit) Draw(camera *graphics.Camera, tileWidth, tileHeight int) {
	var tw, th = float32(tileWidth), float32(tileHeight)
	var x, y = u.x*tw + (tw / 2), u.y*th + (th / 2)

	u.plate.X, u.plate.Y = x, y
	u.body.X, u.body.Y = x, y
	u.head.X, u.head.Y = x, y

	camera.DrawSprites(u.plate, u.body, u.head)
	camera.DrawText("", text.New(u.Initiative), x, y-96, 32, palette.White)
}

//=================================================================

func (u *Unit) Position() (x, y float32) {
	return u.x, u.y
}
