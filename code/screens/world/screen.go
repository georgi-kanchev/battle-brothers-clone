package world

import (
	"game/code/global"
	"game/code/screens/battle"
	"game/code/screens/loading"
	"game/code/unit"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/data/folder"
	"pure-game-kit/data/path"
	"pure-game-kit/execution/screens"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	"pure-game-kit/tiled"
	"pure-game-kit/tiled/property"
	"pure-game-kit/utility/text"
)

type WorldScreen struct {
	path   string
	camera *graphics.Camera

	hud, inventory, settlement, currentPopup *gui.GUI
	resultingCursorNonGUI                    int

	time       float32
	timeCircle *graphics.Sprite

	playerParty  *Party
	otherParties []*Party
	settlements  *tiled.Layer

	tmap      *tiled.Map
	mapLayers []*tiled.Layer
	solids    []*geometry.Shape
	roads     [][2]float32
}

func New(path string) *WorldScreen {
	return &WorldScreen{path: path, camera: graphics.NewCamera(1), time: 60 * 3}
}

//=================================================================

func (ws *WorldScreen) OnLoad() {
	loading.Show("Loading:\nWorld Map...")
	ws.tmap = tiled.NewMap(assets.LoadTiledMap(ws.path), global.Project)
	loading.Show("Loading:\nWorld GUI...")
	ws.hud = gui.NewFromXMLs(file.LoadText("data/gui/world-hud.xml"), global.PopupDimGUI, global.ThemesGUI)
	ws.inventory = gui.NewFromXMLs(file.LoadText("data/gui/world-inventory.xml"), global.ThemesGUI)
	ws.settlement = gui.NewFromXMLs(file.LoadText("data/gui/world-settlement.xml"), global.ThemesGUI)
	ws.currentPopup = nil

	var sc = global.Options.ScaleUI.Master
	ws.hud.Scale = global.Options.ScaleUI.World.HUD * sc
	ws.inventory.Scale = global.Options.ScaleUI.World.Inventory * sc
	ws.settlement.Scale = global.Options.ScaleUI.World.Settlement * sc

	loading.Show("Loading:\nWorld images...")
	var timeCircle = assets.LoadTexture("art/UI/Time/time_circle.PNG")
	ws.timeCircle = graphics.NewSprite(timeCircle, 0, 0)

	assets.LoadTexture("art/UI/Time/time_top.PNG")
	assets.LoadTexture("art/UI/Buttons/btn.PNG")
	assets.LoadTexture("art/UI/Buttons/btn_pause.PNG")
	assets.LoadTexture("art/UI/Buttons/btn_play.PNG")
	assets.LoadTexture("art/UI/Buttons/btn_playx2.PNG")
	assets.LoadTexture("art/UI/Buttons/btn_playx3.PNG")

	var allAssets []string
	allAssets = append(allAssets, folder.Content("art/Character", true)...)
	allAssets = append(allAssets, folder.Content("art/Character/hair", true)...)
	allAssets = append(allAssets, folder.Content("art/Character/body_armor", true)...)
	allAssets = append(allAssets, folder.Content("art/Character/head_armor", true)...)
	for _, filePath := range allAssets {
		if path.IsFile(filePath) {
			assets.LoadTexture(filePath)
		}
	}

	unit.Names = text.Split(file.LoadText("data/names.txt"), " ")
	unit.Nicknames = text.Split(file.LoadText("data/nicknames.txt"), " ")

	loading.Show("Processing:\nWorld data...")
	var solidLayers = ws.tmap.FindLayersBy(property.LayerClass, "WorldSolids")
	var roadLayers = ws.tmap.FindLayersBy(property.LayerClass, "WorldRoads")
	var settlements = ws.tmap.FindLayersBy(property.LayerClass, "WorldSettlements")
	ws.mapLayers = ws.tmap.FindLayersBy(property.LayerClass, "WorldMap")
	for _, s := range solidLayers {
		ws.solids = append(ws.solids, s.ExtractShapes()...)
	}
	for _, r := range roadLayers {
		ws.roads = append(ws.roads, r.ExtractLines()...)
	}
	if len(settlements) > 0 {
		ws.settlements = settlements[0]
	}

	for _, id := range assets.LoadedTextureIds() {
		assets.SetTextureSmoothness(id, true)
	}
}
func (ws *WorldScreen) OnEnter() {
	var units = []*unit.Unit{unit.New(), unit.New(), unit.New(), unit.New(), unit.New(), unit.New()}
	ws.playerParty = NewParty(units, 2250, 1530, true)
}
func (ws *WorldScreen) OnUpdate() {
	ws.camera.SetScreenAreaToWindow()

	//world.tmap.Draw(world.camera)
	for _, m := range ws.mapLayers {
		m.Draw(ws.camera)
	}

	ws.handleResting()

	ws.playerParty.Update()
	for _, party := range ws.otherParties {
		party.Update()
	}

	ws.handleDayNightCycle()

	ws.hud.UpdateAndDraw(ws.camera)
	if ws.currentPopup != nil {
		ws.currentPopup.UpdateAndDraw(ws.camera)
	} else if ws.resultingCursorNonGUI != -1 {
		mouse.SetCursor(ws.resultingCursorNonGUI)
	}

	ws.handleInput()

	switch ws.currentPopup {
	case ws.settlement:
		ws.handleSettlementPopup()
	case ws.inventory:
		ws.handleInventoryPopup()
	}
}

func (ws *WorldScreen) OnExit() {
}

//=================================================================
// private

var teamA = []*unit.Unit{}
var teamB = []*unit.Unit{}

func (ws *WorldScreen) handleInput() {
	if (ws.currentPopup == nil || ws.currentPopup == ws.inventory) && keyboard.IsKeyJustPressed(key.I) {
		ws.currentPopup = global.TogglePopup(ws.hud, ws.currentPopup, ws.inventory)
	}

	if keyboard.IsKeyJustPressed(key.B) {
		screens.Enter(global.ScreenBattle, false)
		var scr = screens.Current().(*battle.BattleScreen)
		scr.Prepare(teamA, teamB, true)
	}
	if ws.hud.IsButtonJustClicked("main-menu", ws.camera) {
		screens.Enter(global.ScreenMainMenu, false)
	}
}
