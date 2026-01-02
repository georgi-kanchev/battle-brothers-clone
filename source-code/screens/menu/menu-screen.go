package menu

import (
	"game/source-code/global"
	"game/source-code/screens/loading"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/execution/screens"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/color"
	"pure-game-kit/window"
)

type MenuScreen struct {
	camera *graphics.Camera

	bgr, logo *graphics.Sprite

	hud, play, load, options, currentPopup *gui.GUI
}

func New() *MenuScreen {
	var menu = &MenuScreen{camera: graphics.NewCamera(1)}
	return menu
}

//=================================================================

func (m *MenuScreen) OnLoad() {
	loading.Show("Loading:\nMain Menu GUI...")
	m.hud = gui.NewFromXMLs(file.LoadText("data/gui/menu-hud.xml"), global.PopupDimGUI, global.ThemesGUI)
	m.options = gui.NewFromXMLs(file.LoadText("data/gui/menu-options.xml"), global.ThemesGUI)
	m.currentPopup = nil

	loading.Show("Loading:\nMain Menu images...")
	var bgr = assets.LoadTexture("art/UI/Titlescreen/bgr.png")
	var logo = assets.LoadTexture("art/UI/Titlescreen/logo.PNG")
	m.bgr = graphics.NewSprite(bgr, 0, 0)
	m.logo = graphics.NewSprite(logo, 0, 0)
	assets.SetTextureSmoothness(bgr, true)
	assets.SetTextureSmoothness(logo, true)

	m.logo.ScaleX, m.logo.ScaleY = 0.8, 0.8
	m.logo.PivotX, m.logo.PivotY = 1, 1

	var sc = global.Options.ScaleUI.Master
	m.options.Scale = global.Options.ScaleUI.Menu.Options * sc

	for _, id := range assets.LoadedTextureIds() {
		assets.SetTextureSmoothness(id, true)
	}
}
func (m *MenuScreen) OnEnter() {
}
func (m *MenuScreen) OnUpdate() {
	m.camera.DrawColor(color.RGB(8, 3, 4))
	m.camera.SetScreenAreaToWindow()

	m.positionElements()

	m.hud.UpdateAndDraw(m.camera)
	if m.currentPopup != nil {
		m.currentPopup.UpdateAndDraw(m.camera)
	}

	m.handleInput()
}

func (m *MenuScreen) OnExit() {
}

//=================================================================
// private

func (m *MenuScreen) positionElements() {
	var rx, ry = m.camera.PointFromEdge(1, 0.5)
	var sc = m.bgr.ScaleX

	m.bgr.CameraFit(m.camera)
	m.bgr.X -= 100 * sc
	m.bgr.ScaleX, m.bgr.ScaleY = m.bgr.ScaleX*1.1, m.bgr.ScaleY*1.1
	m.logo.X, m.logo.Y = rx-80*sc, ry
	m.logo.ScaleX, m.logo.ScaleY = sc, sc
	m.hud.Scale = sc * 1.5
	m.camera.DrawSprites(m.bgr, m.logo)
}
func (m *MenuScreen) handleInput() {
	if m.hud.IsButtonJustClicked("new", m.camera) {
		screens.Enter(global.ScreenWorld, false)
	} else if m.hud.IsButtonJustClicked("options", m.camera) {
		m.currentPopup = global.TogglePopup(m.hud, m.currentPopup, m.options)
	} else if m.hud.IsButtonJustClicked("quit", m.camera) {
		window.Close()
	} else if m.currentPopup != nil && keyboard.IsKeyJustPressed(key.Escape) {
		m.currentPopup = global.TogglePopup(m.hud, m.currentPopup, m.currentPopup)
	}
}
