package unit

import (
	"game/code/global"
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/random"
	"pure-game-kit/utility/time"
)

type Unit struct {
	BaseInitiative, Initiative int
	BaseMovePoints, MovePoints int
	BaseMoveSpeed, MoveSpeed   int

	x, y float32

	head, body, plate *graphics.Sprite
}

func New() *Unit {
	var intv = random.Range(30, 100)
	return &Unit{
		BaseInitiative: intv, Initiative: intv,
		BaseMovePoints: 50, MovePoints: 50,
		BaseMoveSpeed: 100, MoveSpeed: 100,
	}
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

func (u *Unit) UpdateAndDraw(camera *graphics.Camera) {
	// var tw, th = global.BattleTileWidth, global.BattleTileHeight
	var x, y = u.x, u.y //*tw + (tw * 0.5), u.y*th + (th * 0.6)

	u.plate.X, u.plate.Y = x, y
	u.body.X, u.body.Y = x, y
	u.head.X, u.head.Y = x, y

	camera.DrawSprites(u.plate, u.body, u.head)
}

func (u *Unit) MoveTo(targetX, targetY float32) {
	u.x, u.y = point.MoveToPoint(u.x, u.y, targetX, targetY, time.FrameDelta()*float32(u.MoveSpeed*2))
}

//=================================================================

func (u *Unit) Cell() (x, y float32) {
	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	return u.x / tw, u.y / th
}
func (u *Unit) Position() (x, y float32) {
	return u.x, u.y
}
func (u *Unit) IsHovered(camera *graphics.Camera) bool {
	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	var mx, my = camera.MousePosition()
	var cx, cy = u.Cell()
	var hoversX = number.IsWithin(mx/tw, cx, 0.5)
	var hoversY = number.IsWithin(my/th, cy, 0.5)
	return hoversX && hoversY
}
func (u *Unit) AttackRangeCells() [][2]int {
	var x, y = u.Cell()
	var cx, cy = int(x), int(y)
	var meleeRange = [][2]int{
		{cx - 1, cy - 1}, {cx, cy - 1}, {cx + 1, cy - 1},
		{cx - 1, cy + 0}, {cx, cy + 0}, {cx + 1, cy + 0},
		{cx - 1, cy + 1}, {cx, cy + 1}, {cx + 1, cy + 1},
	}
	return meleeRange
}
