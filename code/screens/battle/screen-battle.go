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

	layers  []*graphics.TileMap
	pathMap *geometry.ShapeGrid

	unitManager *unitManager
}

func New(mapPath string) *BattleScreen {
	return &BattleScreen{path: mapPath, camera: graphics.NewCamera(1)}
}

//=================================================================

func (bs *BattleScreen) Prepare(teamA, teamB []*unit.Unit, playerIsTeamA bool) {
	bs.unitManager = newUnitManager(teamA, teamB)
	bs.unitManager.spawnAll(nil, teamA)
	bs.unitManager.spawnAll(nil, teamB)
	bs.unitManager.turnManager.startBattle(teamA, teamB, playerIsTeamA)
}

func (bs *BattleScreen) OnLoad() {
	loading.Show("Loading:\nBattle GUI...")
	bs.hud = gui.NewFromXMLs(file.LoadText("data/gui/battle-hud.xml"), global.ThemesGUI)
	bs.loot = gui.NewFromXMLs(file.LoadText("data/gui/battle-loot.xml"), global.ThemesGUI)
	bs.currentPopup = nil

	var sc = global.Opts.ScaleUI
	bs.hud.Scale = global.Opts.ScaleBattleHUD * sc
	bs.loot.Scale = global.Opts.ScaleBattleLoot * sc

	loading.Show("Loading:\nBattle Tiles...")
	var tileSetId, tileDataIds = assets.LoadTiledData("data/battlegrounds/test/map.tmx")
	var layers = make([]*graphics.TileMap, len(tileDataIds))

	for i, tileDataId := range tileDataIds {
		layers[i] = graphics.NewTileMap(tileSetId, tileDataId)
	}

	var tileWidth, tileHeight = assets.Size(tileSetId)
	var columns, rows = assets.Size(tileDataIds[0])
	global.BattleTileWidth, global.BattleTileHeight = float32(tileWidth), float32(tileHeight)
	global.BattleColumns, global.BattleRows = float32(columns), float32(rows)
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
	bs.camera.SetScreenAreaToWindow()

	if bs.currentPopup == nil {
		bs.camera.MouseDragAndZoomSmoothly()
		bs.camera.Zoom = number.Limit(bs.camera.Zoom, 0.5, 10)
	}

	bs.camera.DrawTileMaps(bs.layers...)

	bs.unitManager.update()

	bs.hud.UpdateAndDraw(bs.camera)
	if bs.currentPopup != nil {
		bs.currentPopup.UpdateAndDraw(bs.camera)
	}

	bs.handleInput()

	global.TryShowFPS(bs.camera)
}

func (bs *BattleScreen) OnExit() {
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
	//var pathMapLayers = bs.layers.FindLayersBy(property.LayerClass, "BattlePathMap")
	// if len(pathMapLayers) > 0 {
	// 	bs.pathMap = pathMapLayers[0].ExtractShapeGrid()
	// }

	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	for _, unit := range bs.unitManager.units {
		if unit != bs.unitManager.turnManager.unitActing() {
			var ux, uy = unit.Position()
			bs.pathMap.SetAtCell(int(ux/tw), int(uy/th), geometry.NewShapeQuad(tw/2, th/2, 0.5, 0.5))
		}
	}
}
