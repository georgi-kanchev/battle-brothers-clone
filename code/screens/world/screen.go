package world

import (
	"game/code/global"
	"game/code/screens/battle"
	"game/code/screens/loading"
	"game/code/unit"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/execution/screens"
	"pure-game-kit/geometry"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	"pure-game-kit/tiled"
	"pure-game-kit/tiled/property"
)

type WorldScreen struct {
	path   string
	camera *graphics.Camera

	hud, inventory, settlement, currentPopup *gui.GUI
	resultingCursorNonGUI                    int

	time       float32
	timeCircle *graphics.Sprite

	parties     []*Party
	settlements *tiled.Layer

	tmap      *tiled.Map
	mapLayers []*tiled.Layer
	solids    []*geometry.Shape
	roads     [][2]float32
}

func New(path string) *WorldScreen {
	var world = &WorldScreen{path: path, camera: graphics.NewCamera(1), time: 60 * 3}
	world.parties = []*Party{NewParty(nil, 2250, 1530, true)}
	return world
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

	assets.LoadTexture("art/Character/head.PNG")
	assets.LoadTexture("art/Character/body.PNG")
	assets.LoadTexture("art/Character/plate.PNG")

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
}
func (ws *WorldScreen) OnUpdate() {
	ws.camera.SetScreenAreaToWindow()

	//world.tmap.Draw(world.camera)
	for _, m := range ws.mapLayers {
		m.Draw(ws.camera)
	}

	for _, party := range ws.parties {
		party.Update()
	}

	ws.handleDayNightCycle()

	ws.hud.UpdateAndDraw(ws.camera)
	if ws.currentPopup != nil {
		ws.currentPopup.UpdateAndDraw(ws.camera)
	}

	if ws.resultingCursorNonGUI != -1 {
		mouse.SetCursor(ws.resultingCursorNonGUI)
	}

	ws.handleInput()

	switch ws.currentPopup {
	case ws.settlement:
		ws.handleSettlementPopup()
	case ws.inventory:
		ws.handleInventoryPopup()
	}

	if ws.hud.IsButtonJustClicked("main-menu", ws.camera) {
		screens.Enter(global.ScreenMainMenu, false)
	}
}

func (ws *WorldScreen) OnExit() {
}

//=================================================================
// private

var teamA = []*unit.Unit{unit.New(), unit.New(), unit.New(), unit.New(), unit.New()}
var teamB = []*unit.Unit{unit.New(), unit.New(), unit.New()}

func (ws *WorldScreen) handleInput() {
	if keyboard.IsKeyJustPressed(key.I) {
		ws.currentPopup = global.TogglePopup(ws.hud, ws.currentPopup, ws.inventory)
	} else if keyboard.IsKeyJustPressed(key.B) {
		screens.Enter(global.ScreenBattle, false)
		var scr = screens.Current().(*battle.BattleScreen)
		scr.Prepare(teamA, teamB, true)
	} else if keyboard.IsKeyJustPressed(key.Escape) && ws.currentPopup != nil {
		ws.currentPopup = global.TogglePopup(ws.hud, ws.currentPopup, ws.currentPopup)
		ws.parties[0].goingToSettlement = nil
	}
}
