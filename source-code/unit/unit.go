package unit

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
)

type Unit struct {
	x, y float32

	head, body, plate *graphics.Sprite

	mapColumns, mapRows, tileWidth, tileHeight int
}

func New(mapColumns, mapRows, tileWidth, tileHeight int) *Unit {
	return &Unit{mapColumns: mapColumns, mapRows: mapRows, tileWidth: tileWidth, tileHeight: tileHeight}
}

//=================================================================

func (unit *Unit) Spawn(x, y float32, flip bool) {
	unit.head = graphics.NewSprite(assets.LoadTexture("art/Character/head.PNG"), 0, 0)
	unit.body = graphics.NewSprite(assets.LoadTexture("art/Character/body.PNG"), 0, 0)
	unit.plate = graphics.NewSprite(assets.LoadTexture("art/Character/plate.PNG"), 0, 0)
	unit.head.PivotY = 0.85
	unit.body.PivotY = 0.85
	unit.plate.PivotY = 0.85

	assets.SetTextureSmoothness(unit.head.AssetId, true)
	assets.SetTextureSmoothness(unit.body.AssetId, true)
	assets.SetTextureSmoothness(unit.plate.AssetId, true)

	unit.x, unit.y = x, y

	if flip {
		unit.head.ScaleX = -1
		unit.body.ScaleX = -1
		unit.plate.ScaleX = -1
	}
}

func (unit *Unit) Draw(camera *graphics.Camera) {
	var tw, th = float32(unit.tileWidth), float32(unit.tileHeight)
	var x, y = unit.x*tw + (tw / 2), unit.y*th + (th / 2)

	unit.plate.X, unit.plate.Y = x, y
	unit.body.X, unit.body.Y = x, y
	unit.head.X, unit.head.Y = x, y
	camera.DrawSprites(unit.plate, unit.body, unit.head)
}

//=================================================================

func (unit *Unit) Position() (x, y float32) {
	return unit.x, unit.y
}
