package world

import (
	"game/source-code/global"
	"pure-game-kit/gui/field"
	"pure-game-kit/tiled/property"
)

func (world *World) handleSettlementPopup() {
	var player = world.parties[0]
	var name = player.goingToSettlement.Properties[property.ObjectName].(string)
	world.settlement.SetField("settlement-title-label", field.Text, name)

	if world.settlement.IsButtonJustClicked("settlement-exit-btn", world.camera) {
		world.currentPopup = global.TogglePopup(world.hud, world.currentPopup, world.settlement)
		player.goingToSettlement = nil
	}
}
