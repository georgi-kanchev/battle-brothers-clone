package loading

import (
	"game/code/global"
	"pure-game-kit/execution/screens"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/random"
	"pure-game-kit/window"
)

type LoadingScreen struct {
	camera *graphics.Camera

	textMid, textBot *graphics.TextBox
}

func New() *LoadingScreen {
	return &LoadingScreen{camera: graphics.NewCamera(1)}
}

func Show(message string) {
	msg = message
	screens.Enter(global.ScreenLoading, false)
}

//=================================================================

func (ls *LoadingScreen) OnLoad() {
	ls.textMid = graphics.NewTextBox("", 0, 0, "")
	ls.textBot = graphics.NewTextBox("", 0, 0, random.Pick("Tip: placeholder #1", "Tip: placeholder #2", "Tip: placeholder #3", "Tip: placeholder #4"))
	ls.textMid.AlignmentX, ls.textMid.AlignmentY = 0.5, 0.5
	ls.textBot.AlignmentX, ls.textBot.AlignmentY = 0.5, 1
	ls.textMid.LineHeight = 60
	ls.textBot.LineHeight = 40
}
func (ls *LoadingScreen) OnEnter() {
	ls.camera.SetScreenAreaToWindow()
	ls.textMid.Text = msg
	ls.textMid.Width, ls.textMid.Height = ls.camera.Size()
	ls.textBot.Width, ls.textBot.Height = ls.textMid.Width, ls.textMid.Height
	ls.camera.DrawTextBoxes(ls.textMid, ls.textBot)
	window.KeepOpen()
}
func (ls *LoadingScreen) OnUpdate() {}
func (ls *LoadingScreen) OnExit()   {}

// =================================================================
// private

var msg = ""
