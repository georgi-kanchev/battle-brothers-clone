package menu

import (
	"game/code/global"
	"pure-game-kit/data/file"
	"pure-game-kit/data/storage"
	"pure-game-kit/execution/condition"
	"pure-game-kit/gui/field"
	"pure-game-kit/utility/number"
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
	ms.tryChangeVsyncOrLimitFPS()
	ms.tryChangeVsyncOrTexFilter()
}

func (ms *MenuScreen) tryChangeWindowState(id string) {
	if ms.options.IsButtonJustClicked(id) {
		var state = int(ms.options.FieldNumber(id, "state", ms.camera))
		var _, cur = window.Monitors()

		global.Opts.WindowState = state
		global.Opts.Monitor = cur
		window.ApplyState(state)

		ms.updateOptionsGUI()
		ms.saveFile()
	}
}
func (ms *MenuScreen) tryChangeMonitor() {
	for i := range 10 {
		var id = text.New("m", i)
		if ms.options.IsButtonJustClicked(id) {
			global.Opts.Monitor = i

			global.ApplyOptions()
			condition.CallAfter(0.1, func() { ms.updateOptionsGUI() })
			ms.saveFile()
		}
	}
}
func (ms *MenuScreen) tryChangeVsyncOrLimitFPS() {
	if !ms.options.IsButtonJustClicked("vsync") && !ms.options.IsSliderJustSlid("limit-fps") {
		return
	}

	var vsync = ms.options.Field("vsync", field.Value, ms.camera)
	var slider = ms.options.FieldNumber("limit-fps", field.Value, ms.camera)
	var maxFPS = number.Snap(number.Map(slider, 0, 1, 0, 250), 10)

	global.Opts.VSync = vsync == "1"
	global.Opts.LimitFPS = condition.If(maxFPS == 250, 0, int(maxFPS))

	global.ApplyOptions()
	ms.updateOptionsGUI()
	ms.saveFile()
}
func (ms *MenuScreen) tryChangeVsyncOrTexFilter() {
	var clickAntialiasing = ms.options.IsButtonJustClicked("aa")
	var clickTexFilter = ms.options.IsButtonJustClicked("tex-filter")

	if clickAntialiasing {
		var antialiasing = ms.options.Field("aa", field.Value, ms.camera)
		global.Opts.Antialiasing = antialiasing == "1"
	}
	if clickTexFilter {
		var textureFilter = ms.options.Field("tex-filter", field.Value, ms.camera)
		global.Opts.TextureFilter = textureFilter == "1"
	}
	if clickAntialiasing || clickTexFilter {
		global.ApplyOptions()
		ms.updateOptionsGUI()
		ms.saveFile()
	}
}

func (ms *MenuScreen) updateOptionsGUI() {
	var widgetStates = ms.options.WidgetIdsOfContainer("states")
	for _, widget := range widgetStates {
		var state = ms.options.Field(widget, "state", ms.camera)
		if state == text.New(global.Opts.WindowState) {
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

	var vsync = condition.If(global.Opts.VSync, "1", "")
	var aa = condition.If(global.Opts.Antialiasing, "1", "")
	var texFilter = condition.If(global.Opts.TextureFilter, "1", "")
	var labelFPS = condition.If(global.Opts.LimitFPS == 0, "unlimited", text.New(global.Opts.LimitFPS))
	var slider = number.Map(float32(global.Opts.LimitFPS), 0, 250, 0, 1)

	if vsync == "1" {
		var hz = text.Remove(text.Between(monitors[cur], ", ", ")"), "Hz")
		labelFPS = text.New(hz, " (Monitor Hz)")
	}
	ms.options.SetField("vsync", field.Text, text.Replace(vsync, "1", "X"))
	ms.options.SetField("vsync", field.Value, vsync)
	ms.options.SetField("limit-fps-label", field.Hidden, vsync)
	ms.options.SetField("limit-fps", field.Hidden, vsync)
	ms.options.SetField("limit-fps", field.Value, text.New(slider))
	ms.options.SetField("spacing", field.Hidden, condition.If(vsync == "1", "", "1"))
	ms.options.SetField("fps", field.Text, text.New("Maximum FPS = ", labelFPS))
	ms.options.SetField("aa", field.Text, text.Replace(aa, "1", "X"))
	ms.options.SetField("aa", field.Value, text.Replace(aa, "1", "X"))
	ms.options.SetField("tex-filter", field.Text, text.Replace(texFilter, "1", "X"))
	ms.options.SetField("tex-filter", field.Value, text.Replace(texFilter, "1", "X"))
}
func (ms *MenuScreen) saveFile() {
	var yaml = storage.ToYAML(&global.Opts)
	var state = "#0=Floating 1=FloatingBorderless 2=Fullscreen 3=FullscreenBorderless 4=Maximized 5=Minimized\n"
	var vsync = "graphics-vsync: true #limit-fps is ignored, max FPS = monitor Hz"
	yaml = text.Replace(yaml, "graphics-vsync: true", vsync)
	file.SaveText("data/options.yaml", state+yaml)
}
