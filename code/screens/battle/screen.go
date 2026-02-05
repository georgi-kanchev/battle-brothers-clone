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
	"pure-game-kit/tiled"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

type BattleScreen struct {
	path   string // cached for reload
	camera *graphics.Camera

	hud, currentPopup, loot *gui.GUI

	tmap    *tiled.Map
	tiles   []*graphics.Sprite
	pathMap *geometry.ShapeGrid

	unitManager *unitManager
}

func New(mapPath string) *BattleScreen {
	return &BattleScreen{path: mapPath, camera: graphics.NewCamera(1)}
}

//=================================================================

func (b *BattleScreen) Prepare(teamA, teamB []*unit.Unit, playerIsTeamA bool) {
	b.unitManager = newUnitManager(teamA, teamB)
	b.unitManager.spawnAll(b.tmap, teamA, "BattleSpawnsTeamA")
	b.unitManager.spawnAll(b.tmap, teamB, "BattleSpawnsTeamB")
	b.unitManager.turnManager.startBattle(teamA, teamB, playerIsTeamA)
}

func (b *BattleScreen) OnLoad() {
	loading.Show("Loading:\nBattle Map...")
	b.tmap = tiled.NewMap(assets.LoadTiledMap(b.path), global.Project)
	loading.Show("Loading:\nBattle GUI...")
	b.hud = gui.NewFromXMLs(file.LoadText("data/gui/battle-hud.xml"), global.PopupDimGUI, global.ThemesGUI)
	b.loot = gui.NewFromXMLs(file.LoadText("data/gui/battle-loot.xml"), global.ThemesGUI)
	b.currentPopup = nil

	var sc = global.Options.ScaleUI
	b.hud.Scale = global.Options.ScaleBattleHUD * sc
	b.loot.Scale = global.Options.ScaleBattleLoot * sc

	loading.Show("Processing:\nBattle data...")
	var layers = b.tmap.FindLayersBy(property.LayerClass, "BattleMap")
	layers = append(layers, b.tmap.FindLayersBy(property.LayerClass, "BattlePathMap")...)
	for _, l := range layers {
		b.tiles = append(b.tiles, l.ExtractSprites()...)
	}

	var atlasId = b.tmap.Tilesets[0].Properties["atlasId"].(string)
	for i := range 31 {
		var tileId = text.New(atlasId, "/", i)
		assets.SetTextureAtlasTile(atlasId, tileId, float32(i), 1, 1, 1, 0, false)
	}
}
func (b *BattleScreen) OnEnter() {
	global.BattleTileWidth = float32(b.tmap.Properties[property.MapTileWidth].(int))
	global.BattleTileHeight = float32(b.tmap.Properties[property.MapTileHeight].(int))
	global.BattleTileColumns = float32(b.tmap.Properties[property.MapColumns].(int))
	global.BattleTileRows = float32(b.tmap.Properties[property.MapRows].(int))

	b.camera.X = global.BattleTileColumns / 2 * global.BattleTileWidth
	b.camera.Y = global.BattleTileRows / 2 * global.BattleTileHeight
	b.camera.Zoom = 0.8

	for _, id := range assets.LoadedTextureIds() { // probably shouldn't be here
		assets.SetTextureSmoothness(id, true)
	}
}
func (b *BattleScreen) OnUpdate() {
	b.camera.SetScreenAreaToWindow()

	if b.currentPopup == nil {
		b.camera.MouseDragAndZoomSmoothly()
		b.camera.Zoom = number.Limit(b.camera.Zoom, 0.5, 10)
	}

	b.camera.DrawSprites(b.tiles...)

	b.unitManager.update()

	b.hud.UpdateAndDraw(b.camera)
	if b.currentPopup != nil {
		b.currentPopup.UpdateAndDraw(b.camera)
	}

	b.handleInput()
}

func (b *BattleScreen) OnExit() {
}

//=================================================================
// private

func (b *BattleScreen) handleInput() {
	if keyboard.IsKeyJustPressed(key.Escape) {
		screens.Enter(global.ScreenWorld, false)
	} else if keyboard.IsKeyJustPressed(key.L) {
		b.currentPopup = condition.If(b.currentPopup == b.loot, nil, b.loot)
	}
}
func (b *BattleScreen) recalculatePathMap() {
	var pathMapLayers = b.tmap.FindLayersBy(property.LayerClass, "BattlePathMap")
	if len(pathMapLayers) > 0 {
		b.pathMap = pathMapLayers[0].ExtractShapeGrid()
	}

	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	for _, unit := range b.unitManager.units {
		if unit != b.unitManager.turnManager.unitActing() {
			var ux, uy = unit.Position()
			b.pathMap.SetAtCell(int(ux/tw), int(uy/th), geometry.NewShapeQuad(tw/2, th/2, 0.5, 0.5))
		}
	}
}
