package menu

import (
	"game/code/global"
	"pure-game-kit/data/file"
	"pure-game-kit/data/storage"
	"pure-game-kit/execution/condition"
	"pure-game-kit/gui/field"
	"pure-game-kit/utility/text"
	"pure-game-kit/window"
)

func (ms *MenuScreen) handleOptionsPopup() {
	ms.tryExitPopup(ms.options, nil)
	ms.options.SetField("title-bgr", field.Text, "Options")

	var gfx = ms.options.Field("tab-graphics", field.Value, ms.camera)
	var ui = ms.options.Field("tab-ui", field.Value, ms.camera)
	var audio = ms.options.Field("tab-audio", field.Value, ms.camera)
	var ctrls = ms.options.Field("tab-controls", field.Value, ms.camera)
	var game = ms.options.Field("tab-gameplay", field.Value, ms.camera)
	ms.options.SetField("graphics", field.Hidden, condition.If(gfx == "", "1", ""))
	ms.options.SetField("ui", field.Hidden, condition.If(ui == "", "1", ""))
	ms.options.SetField("audio", field.Hidden, condition.If(audio == "", "1", ""))
	ms.options.SetField("controls", field.Hidden, condition.If(ctrls == "", "1", ""))
	ms.options.SetField("gameplay", field.Hidden, condition.If(game == "", "1", ""))

	ms.tryChangeWindowState("floating")
	ms.tryChangeWindowState("maximized")
	ms.tryChangeWindowState("fullscreen")
	ms.tryChangeWindowState("fullscreen-borderless")
	ms.tryChangeMonitor()
}

func (ms *MenuScreen) updateOptionsGUI() {
	var widgetStates = ms.options.WidgetIdsOfContainer("states")
	for _, widget := range widgetStates {
		var state = ms.options.Field(widget, "state", ms.camera)
		if state == text.New(global.Options.WindowState) {
			ms.options.SetField("window-state", field.Text, ms.options.Field(widget, field.Text, ms.camera))
			break
		}
	}

	var monitors, cur = window.Monitors()
	for i := range 10 {
		var hidden = condition.If(i >= len(monitors), "1", "")
		var id = text.New("m", i)
		ms.options.SetField(id, field.Hidden, hidden)
		if hidden == "" {
			var label = text.New("  ", i+1, ": ", monitors[i])
			ms.options.SetField(id, field.Text, label)

			if cur == i {
				ms.options.SetField("monitors-menu", field.Text, text.Trim(label))
			}
		}
	}
}

func (ms *MenuScreen) tryChangeWindowState(id string) {
	if ms.options.IsButtonJustClicked(id, ms.camera) {
		var state = int(ms.options.FieldNumber(id, "state", ms.camera))
		window.ApplyState(state)
		ms.options.SetField("window-state", field.Text, ms.options.Field(id, field.Text, ms.camera))
		global.Options.WindowState = state
		ms.saveFile()
	}
}
func (ms *MenuScreen) tryChangeMonitor() {
	for i := range 10 {
		var id = text.New("m", i)
		if ms.options.IsButtonJustClicked(id, ms.camera) {
			window.MoveToMonitor(i)
			ms.options.SetField("monitors-menu", field.Text, ms.options.Field(id, field.Text, ms.camera))
			global.Options.Monitor = i
			ms.saveFile()
		}
	}
}

func (ms *MenuScreen) saveFile() {
	var yaml = storage.ToYAML(&global.Options)
	var comment = "#0=Floating 1=FloatingBorderless 2=Fullscreen 3=FullscreenBorderless 4=Maximized 5=Minimized\n"
	file.SaveText("data/options.yaml", comment+yaml)
}
