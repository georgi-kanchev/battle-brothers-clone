package battle

import (
	"game/code/global"
	"game/code/screens/loading"
	"game/code/unit"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/screens"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/number"
)

type BattleScreen struct {
	path   string // cached for reload
	camera *graphics.Camera

	hud, currentPopup, loot *gui.GUI

	ground, above *graphics.TileMap
	pathMap       *geometry.ShapeGrid

	spawnsA, spawnsB []float32

	unitManager *unitManager
}

func New(mapPath string) *BattleScreen {
	return &BattleScreen{path: mapPath, camera: graphics.NewCamera(1)}
}

//=================================================================

func (bs *BattleScreen) OnLoad() {
	loading.Show("Loading:\nBattle GUI...")
	bs.hud = gui.NewFromXMLs(bs.camera, file.LoadText("data/gui/battle-hud.xml"), global.ThemesGUI)
	bs.loot = gui.NewFromXMLs(bs.camera, file.LoadText("data/gui/battle-loot.xml"), global.ThemesGUI)
	bs.currentPopup = nil

	var sc = global.Opts.ScaleUI
	bs.hud.Scale = global.Opts.ScaleBattleHUD * sc
	bs.loot.Scale = global.Opts.ScaleBattleLoot * sc

	loading.Show("Loading:\nBattle Tiles...")
	var mapPath = "data/battlegrounds/test/map.tmx"
	var tileSetId, _ = assets.LoadTiledLayers(mapPath)

	bs.ground = graphics.NewTileMap(tileSetId, mapPath+"/Ground")
	bs.above = graphics.NewTileMap(tileSetId, mapPath+"/Above")

	bs.ground.PivotX, bs.ground.PivotY = 0, 0
	bs.above.PivotX, bs.above.PivotY = 0, 0

	var spawnsA = graphics.NewTileMap(tileSetId, mapPath+"/SpawnsA")
	spawnsA.PivotX, spawnsA.PivotY = 0, 0
	bs.spawnsA = spawnsA.Points()

	var spawnsB = graphics.NewTileMap(tileSetId, mapPath+"/SpawnsB")
	spawnsB.PivotX, spawnsB.PivotY = 0, 0
	bs.spawnsB = spawnsB.Points()

	var w, h = bs.ground.Size()
	global.BattleTileWidth, global.BattleTileHeight = bs.ground.SizeTile()
	global.BattleColumns, global.BattleRows = float32(w), float32(h)

	bs.pathMap = geometry.NewShapeGrid(int(global.BattleTileWidth), int(global.BattleTileHeight))
}
func (bs *BattleScreen) OnEnter() {
	bs.camera.X = global.BattleColumns / 2 * global.BattleTileWidth
	bs.camera.Y = global.BattleRows / 2 * global.BattleTileHeight
	bs.camera.Zoom = 0.8

	for _, id := range assets.LoadedTextureIds() { // probably shouldn't be here
		assets.SetTextureSmoothness(id, true)
	}
}
func (bs *BattleScreen) OnUpdate() {
	if bs.currentPopup == nil {
		bs.camera.MouseDragAndZoomSmoothly()
		bs.camera.Zoom = number.Limit(bs.camera.Zoom, 0.5, 10)
	}

	bs.camera.DrawTileMaps(bs.ground, bs.above)

	bs.unitManager.update()

	bs.hud.UpdateAndDraw()
	if bs.currentPopup != nil {
		bs.currentPopup.UpdateAndDraw()
	}

	bs.handleInput()

	global.TryShowFPS(bs.camera)
}
func (bs *BattleScreen) OnExit() {
}

func (bs *BattleScreen) Prepare(teamA, teamB []*unit.Unit, playerIsTeamA bool) {
	bs.unitManager = newUnitManager(teamA, teamB)
	bs.unitManager.spawnAll(bs.spawnsA, teamA)
	bs.unitManager.spawnAll(bs.spawnsB, teamB)
	bs.unitManager.turnManager.startBattle(teamA, teamB, playerIsTeamA)
}

//=================================================================
// private

func (bs *BattleScreen) handleInput() {
	if keyboard.IsKeyJustPressed(key.Escape) {
		screens.Enter(global.ScreenWorld, false)
	} else if keyboard.IsKeyJustPressed(key.L) {
		bs.currentPopup = condition.If(bs.currentPopup == bs.loot, nil, bs.loot)
	}
}
func (bs *BattleScreen) recalculatePathMap() {
	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	var w, h = bs.above.Size()
	bs.pathMap = geometry.NewShapeGrid(int(tw), int(th))
	for i := range h {
		for j := range w {
			var pts = bs.above.PointsAtCell(j, i)
			if len(pts) > 0 {
				bs.pathMap.SetAtCell(j, i, geometry.NewShapeCorners(pts...))
			}
		}
	}

	for _, unit := range bs.unitManager.units {
		if unit != bs.unitManager.turnManager.unitActing() {
			var ux, uy = unit.Position()
			bs.pathMap.SetAtCell(int(ux/tw), int(uy/th), geometry.NewShapeQuad(tw/2, th/2, 0.5, 0.5))
		}
	}
}
