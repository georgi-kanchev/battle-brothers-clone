package world

import (
	"pure-game-kit/gui/field"
	"pure-game-kit/tiled/property"
)

func (ws *WorldScreen) handleSettlementPopup() {
	var name = ws.playerParty.goingToSettlement.Properties[property.ObjectName].(string)
	ws.settlement.SetField("title-bgr", field.Text, "Town of \""+name+"\"")

	ws.tryExitPopup(ws.settlement, nil, func() { ws.playerParty.goingToSettlement = nil })

	if ws.settlement.IsButtonJustClicked("rest", ws.camera) {
		ws.playerParty.isResting = true
		ws.currentPopup = nil
	} else if ws.settlement.IsButtonJustClicked("market", ws.camera) {
		ws.currentPopup = ws.market
	} else if ws.settlement.IsButtonJustClicked("favors", ws.camera) {
		ws.currentPopup = ws.favors
	} else if ws.settlement.IsButtonJustClicked("recruit", ws.camera) {
		ws.currentPopup = ws.recruit
	} else if ws.settlement.IsButtonJustClicked("tavern", ws.camera) {
		ws.currentPopup = ws.tavern
	}
}
