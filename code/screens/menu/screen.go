package menu

import (
	"game/code/global"
	"game/code/screens/loading"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/execution/condition"
	"pure-game-kit/execution/screens"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	"pure-game-kit/gui/field"
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

func (ms *MenuScreen) OnLoad() {
	loading.Show("Loading:\nMain Menu GUI...")
	var hud = file.LoadText("data/gui/menu-hud.xml")
	var options = file.LoadText("data/gui/menu-options.xml")
	var narrow = global.PopupNarrowGUI
	ms.hud = gui.NewFromXMLs(hud, global.DimGUI, global.ThemesGUI)
	ms.options = gui.NewFromXMLs(global.DimGUI, narrow, options, global.TitleGUI, global.XBtnGUI, global.ThemesGUI)
	ms.currentPopup = nil

	loading.Show("Loading:\nMain Menu images...")
	var bgr = assets.LoadTexture("art/UI/Titlescreen/bgr.png")
	var logo = assets.LoadTexture("art/UI/Titlescreen/logo.PNG")
	ms.bgr = graphics.NewSprite(bgr, 0, 0)
	ms.logo = graphics.NewSprite(logo, 0, 0)
	assets.SetTextureSmoothness(bgr, true)
	assets.SetTextureSmoothness(logo, true)

	ms.logo.ScaleX, ms.logo.ScaleY = 0.8, 0.8
	ms.logo.PivotX, ms.logo.PivotY = 1, 1

	var sc = global.Options.ScaleUI
	ms.options.Scale = global.Options.ScaleMenuOptions * sc

	for _, id := range assets.LoadedTextureIds() {
		assets.SetTextureSmoothness(id, true)
	}
}
func (ms *MenuScreen) OnEnter() {
	ms.updateOptionsGUI()
}
func (ms *MenuScreen) OnUpdate() {
	ms.camera.SetScreenAreaToWindow()

	ms.makeBackground()
	ms.handleInput()

	ms.hud.SetField("popup-dim-bgr", field.Hidden, condition.If(ms.currentPopup == nil, "1", ""))

	ms.hud.UpdateAndDraw(ms.camera)
	if ms.currentPopup != nil {
		ms.currentPopup.UpdateAndDraw(ms.camera)
	}

	switch ms.currentPopup {
	case ms.options:
		ms.handleOptionsPopup()
	}
}

func (ms *MenuScreen) OnExit() {
}

//=================================================================
// private

func (ms *MenuScreen) makeBackground() {
	var rx, ry = ms.camera.PointFromEdge(1, 0.5)
	var sc = ms.bgr.ScaleX

	ms.bgr.CameraFit(ms.camera)
	ms.bgr.X -= 100 * sc
	ms.bgr.ScaleX, ms.bgr.ScaleY = ms.bgr.ScaleX*1.1, ms.bgr.ScaleY*1.1
	ms.logo.X, ms.logo.Y = rx-80*sc, ry
	ms.logo.ScaleX, ms.logo.ScaleY = sc, sc
	ms.hud.Scale = sc * 1.5
	ms.camera.DrawColor(color.RGB(8, 3, 4))
	ms.camera.DrawSprites(ms.bgr, ms.logo)
}
func (ms *MenuScreen) handleInput() {
	if ms.hud.IsButtonJustClicked("new", ms.camera) {
		screens.Enter(global.ScreenWorld, false)
	} else if ms.hud.IsButtonJustClicked("options", ms.camera) {
		ms.currentPopup = ms.options
	} else if ms.hud.IsButtonJustClicked("quit", ms.camera) {
		window.Close()
	}
}

func (ms *MenuScreen) tryExitPopup(from *gui.GUI, to *gui.GUI) {
	if from.IsButtonJustClicked("exit-btn", ms.camera) ||
		from.IsButtonJustClicked("popup-dim-bgr", ms.camera) {
		ms.currentPopup = to
	}
}
