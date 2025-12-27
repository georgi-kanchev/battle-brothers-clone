package world

import (
	"game/source-code/global"
	"pure-game-kit/gui/field"
	"pure-game-kit/tiled/property"
)

func (w *WorldScreen) handleSettlementPopup() {
	var player = w.parties[0]
	var name = player.goingToSettlement.Properties[property.ObjectName].(string)
	w.settlement.SetField("settlement-title-label", field.Text, name)

	if w.settlement.IsButtonJustClicked("settlement-exit-btn", w.camera) {
		w.currentPopup = global.TogglePopup(w.hud, w.currentPopup, w.settlement)
		player.goingToSettlement = nil
	}
}
