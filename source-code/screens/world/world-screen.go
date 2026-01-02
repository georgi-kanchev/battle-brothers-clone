package world

import (
	"game/source-code/global"
	"game/source-code/screens/battle"
	"game/source-code/screens/loading"
	"game/source-code/unit"
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

func (w *WorldScreen) OnLoad() {
	loading.Show("Loading:\nWorld Map...")
	w.tmap = tiled.NewMap(assets.LoadTiledMap(w.path), global.Project)
	loading.Show("Loading:\nWorld GUI...")
	w.hud = gui.NewFromXMLs(file.LoadText("data/gui/world-hud.xml"), global.PopupDimGUI, global.ThemesGUI)
	w.inventory = gui.NewFromXMLs(file.LoadText("data/gui/world-inventory.xml"), global.ThemesGUI)
	w.settlement = gui.NewFromXMLs(file.LoadText("data/gui/world-settlement.xml"), global.ThemesGUI)
	w.currentPopup = nil

	var sc = global.Options.ScaleUI.Master
	w.hud.Scale = global.Options.ScaleUI.World.HUD * sc
	w.inventory.Scale = global.Options.ScaleUI.World.Inventory * sc
	w.settlement.Scale = global.Options.ScaleUI.World.Settlement * sc

	loading.Show("Loading:\nWorld images...")
	var timeCircle = assets.LoadTexture("art/UI/Time/time_circle.PNG")
	w.timeCircle = graphics.NewSprite(timeCircle, 0, 0)

	assets.LoadTexture("art/UI/Time/time_top.PNG")
	assets.LoadTexture("art/UI/Buttons/btn.PNG")
	assets.LoadTexture("art/UI/Buttons/btn_pause.PNG")
	assets.LoadTexture("art/UI/Buttons/btn_play.PNG")
	assets.LoadTexture("art/UI/Buttons/btn_playx2.PNG")
	assets.LoadTexture("art/UI/Buttons/btn_playx3.PNG")

	loading.Show("Processing:\nWorld data...")
	var solidLayers = w.tmap.FindLayersBy(property.LayerClass, "WorldSolids")
	var roadLayers = w.tmap.FindLayersBy(property.LayerClass, "WorldRoads")
	var settlements = w.tmap.FindLayersBy(property.LayerClass, "WorldSettlements")
	w.mapLayers = w.tmap.FindLayersBy(property.LayerClass, "WorldMap")
	for _, s := range solidLayers {
		w.solids = append(w.solids, s.ExtractShapes()...)
	}
	for _, r := range roadLayers {
		w.roads = append(w.roads, r.ExtractLines()...)
	}
	if len(settlements) > 0 {
		w.settlements = settlements[0]
	}

	for _, id := range assets.LoadedTextureIds() {
		assets.SetTextureSmoothness(id, true)
	}
}
func (w *WorldScreen) OnEnter() {
}
func (w *WorldScreen) OnUpdate() {
	w.camera.SetScreenAreaToWindow()

	//world.tmap.Draw(world.camera)
	for _, m := range w.mapLayers {
		m.Draw(w.camera)
	}

	for _, party := range w.parties {
		party.Update()
	}

	w.handleDayNightCycle()

	w.hud.UpdateAndDraw(w.camera)
	if w.currentPopup != nil {
		w.currentPopup.UpdateAndDraw(w.camera)
	}

	if w.resultingCursorNonGUI != -1 {
		mouse.SetCursor(w.resultingCursorNonGUI)
	}

	w.handleInput()

	if w.currentPopup == w.settlement {
		w.handleSettlementPopup()
	}
}

func (w *WorldScreen) OnExit() {
}

//=================================================================
// private

var teamA = []*unit.Unit{unit.New(), unit.New(), unit.New(), unit.New(), unit.New()}
var teamB = []*unit.Unit{unit.New(), unit.New(), unit.New()}

func (w *WorldScreen) handleInput() {
	if keyboard.IsKeyJustPressed(key.I) {
		w.currentPopup = global.TogglePopup(w.hud, w.currentPopup, w.inventory)
	} else if keyboard.IsKeyJustPressed(key.B) {
		screens.Enter(global.ScreenBattle, false)

		var scr = screens.Current().(*battle.BattleScreen)
		scr.Prepare(teamA, teamB, true)
	} else if keyboard.IsKeyJustPressed(key.Escape) {
		if w.currentPopup == nil {
			screens.Enter(global.ScreenMainMenu, false)
		} else {
			w.currentPopup = global.TogglePopup(w.hud, w.currentPopup, w.currentPopup)
		}
	}
}
