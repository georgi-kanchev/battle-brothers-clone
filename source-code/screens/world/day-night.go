package world

import (
	"game/source-code/global"
	"pure-game-kit/gui/field"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/time"
)

var overlayColors = []uint{
	color.RGBA(0, 10, 30, 160),   // midnight
	color.RGBA(0, 10, 30, 160),   // night
	color.RGBA(255, 120, 60, 51), // dawn
	0, 0, 0,                      // morning, noon, afternoon
	color.RGBA(255, 100, 40, 46), // dusk
	color.RGBA(0, 10, 30, 160),   // evening
}
var clock = []string{"Midnight", "Night", "Dawn", "Morning", "Noon", "Afternoon", "Dusk", "Evening"}

func (world *World) handleDayNightCycle() {
	var dayNightCycleDuration = time.FromMinutes(float32(len(clock)))
	var topX, topY = world.camera.PointFromEdge(0.5, 0)
	world.timeCircle.X, world.timeCircle.Y = topX, topY+(185*world.hud.Scale/world.camera.Zoom)
	world.timeCircle.ScaleX = -1 / world.camera.Zoom * 0.5 * world.hud.Scale
	world.timeCircle.ScaleY = 1 / world.camera.Zoom * 0.5 * world.hud.Scale

	var scrX, scrY = world.camera.PointToScreen(topX, topY)
	var scrW, scrH = 100 * world.hud.Scale, 150 * world.hud.Scale
	world.camera.Mask(int(float32(scrX)-50*world.hud.Scale), int(scrY), int(scrW), int(scrH))
	world.timeCircle.Angle = number.Map(world.time, 0, dayNightCycleDuration, 0, 360)

	world.camera.DrawSprites(world.timeCircle)
	world.camera.SetScreenAreaToWindow()

	var x0, x1 = world.hud.Field("x0", field.Value), world.hud.Field("x1", field.Value)
	var x2, x3 = world.hud.Field("x2", field.Value), world.hud.Field("x3", field.Value)
	if world.currentPopup != nil || x0 != "" {
		global.TimeScale = 0
	} else if x1 != "" {
		global.TimeScale = 1
	} else if x2 != "" {
		global.TimeScale = 5
	} else if x3 != "" {
		global.TimeScale = 20
	}

	world.time += time.FrameDelta() * global.TimeScale
	world.time = number.Wrap(world.time, 0, dayNightCycleDuration)
	var timeOfDay = number.Round(world.time/60, 0)
	var timeOfDayIndex = int(timeOfDay) % len(clock)
	world.hud.SetField("time-word", field.Text, clock[timeOfDayIndex])

	var curColor, nextColor = overlayColors[timeOfDayIndex], overlayColors[int(timeOfDay+1)%len(clock)]
	var progress = number.Map(float32(timeOfDayIndex)-world.time/60, -0.5, 0.5, 1, 0)
	var col = color.Fade(curColor, nextColor, progress)

	world.camera.DrawColor(col)
}
