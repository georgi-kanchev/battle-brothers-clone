package world

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/gui/field"
	"pure-game-kit/tiled/property"
)

func (ws *WorldScreen) handleSettlementPopup() {
	var name = ws.playerParty.goingToSettlement.Properties[property.ObjectName].(string)
	ws.settlement.SetField("title-label", field.Text, "Town of \""+name+"\"")

	if ws.settlement.IsButtonJustClicked("exit-btn", ws.camera) ||
		ws.settlement.IsButtonJustClicked("popup-dim-bgr", ws.camera) {
		ws.currentPopup = condition.If(ws.currentPopup == ws.settlement, nil, ws.settlement)
		ws.playerParty.goingToSettlement = nil
	}

	if ws.settlement.IsButtonJustClicked("rest", ws.camera) {
		ws.playerParty.isResting = true
		ws.currentPopup = nil
	} else if ws.settlement.IsButtonJustClicked("market", ws.camera) {
		ws.currentPopup = ws.market
	} else if ws.settlement.IsButtonJustClicked("quests", ws.camera) {
		ws.currentPopup = ws.quests
	}
}
