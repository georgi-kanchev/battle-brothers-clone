package loading

import (
	"game/code/global"
	"pure-game-kit/data/assets"
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

func (l *LoadingScreen) OnLoad() {
	assets.LoadDefaultFont()
	l.textMid = graphics.NewTextBox("", 0, 0, "")
	l.textBot = graphics.NewTextBox("", 0, 0, random.Pick("Tip: placeholder #1", "Tip: placeholder #2", "Tip: placeholder #3", "Tip: placeholder #4"))
	l.textMid.AlignmentX, l.textMid.AlignmentY = 0.5, 0.5
	l.textBot.AlignmentX, l.textBot.AlignmentY = 0.5, 1
	l.textMid.LineHeight = 60
	l.textBot.LineHeight = 40
}
func (l *LoadingScreen) OnEnter() {
	l.camera.SetScreenAreaToWindow()
	l.textMid.Text = msg
	l.textMid.Width, l.textMid.Height = l.camera.Size()
	l.textBot.Width, l.textBot.Height = l.textMid.Width, l.textMid.Height
	l.camera.DrawTextBoxes(l.textMid, l.textBot)
	window.KeepOpen()
}
func (l *LoadingScreen) OnUpdate() {}
func (l *LoadingScreen) OnExit()   {}

// =================================================================
// private

var msg = ""
