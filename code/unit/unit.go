package unit

import (
	"game/code/global"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/random"
	"pure-game-kit/utility/time"
)

// index for unit sprites
const main, secondary, helmet, hair, beard, head, armor, body, plate = 0, 1, 2, 3, 4, 5, 6, 7, 8

var Names, Nicknames []string // world loads them from files

type Unit struct {
	Name, Nickname string
	NameColor      uint

	BaseInitiative, Initiative int
	BaseMovePoints, MovePoints int
	BaseMoveSpeed, MoveSpeed   int

	x, y float32

	sprites []*graphics.Sprite
}

func New() *Unit {
	var initiative = random.Range(30, 100)
	var unit = &Unit{
		BaseInitiative: initiative, Initiative: initiative,
		BaseMovePoints: 50, MovePoints: 50,
		BaseMoveSpeed: 100, MoveSpeed: 100,
	}
	unit.Name = random.Pick(Names...)
	unit.Nickname = random.Pick(Nicknames...)
	unit.NameColor = color.RandomBright()
	if random.HasChance(5) {
		unit.Name, unit.Nickname = unit.Nickname, unit.Name
	}
	Names = collection.Remove(Names, unit.Name)
	Nicknames = collection.Remove(Nicknames, unit.Nickname)

	unit.sprites = make([]*graphics.Sprite, 9)
	unit.sprites[main] = graphics.NewSprite("none", 0, 0)
	unit.sprites[secondary] = graphics.NewSprite("none", 0, 0)
	unit.sprites[helmet] = graphics.NewSprite( /*"art/Character/head_armor/greathelm_03.PNG"*/ "", 0, 0)
	unit.sprites[hair] = graphics.NewSprite("art/Character/hair/hair_01.PNG", 0, 0)
	unit.sprites[beard] = graphics.NewSprite("art/Character/hair/beard_03.PNG", 0, 0)
	unit.sprites[head] = graphics.NewSprite("art/Character/head.PNG", 0, 0)
	unit.sprites[armor] = graphics.NewSprite("art/Character/body_armor/gambeson.PNG", 0, 0)
	unit.sprites[body] = graphics.NewSprite("art/Character/body.PNG", 0, 0)
	unit.sprites[plate] = graphics.NewSprite("art/Character/plate.PNG", 0, 0)
	for _, spr := range unit.sprites {
		spr.PivotY = 0.85
		spr.Width, spr.Height = 64, 128
	}
	return unit
}

//=================================================================

func (u *Unit) Spawn(x, y float32) {
	u.x, u.y = x, y
}

func (u *Unit) UpdateAndDraw(x, y, scaleX, scaleY float32, camera *graphics.Camera) {
	for i := len(u.sprites) - 1; i >= 0; i-- {
		u.sprites[i].X, u.sprites[i].Y = x, y
		u.sprites[i].ScaleX, u.sprites[i].ScaleY = scaleX, scaleY
		if u.sprites[i].AssetId != "" {
			camera.DrawSprites(u.sprites[i])
		}
	}
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

func (u *Unit) NickAndName() string {
	return u.Nickname + " " + u.Name
}
