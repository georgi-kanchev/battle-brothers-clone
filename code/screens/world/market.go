package world

import "pure-game-kit/gui/field"

func (ws *WorldScreen) handleMarketPopup() {
	if ws.market.IsButtonJustClicked("exit-btn", ws.camera) ||
		ws.market.IsButtonJustClicked("popup-dim-bgr", ws.camera) {
		ws.currentPopup = ws.settlement
	}

	var allShopIds = ws.market.WidgetIdsOfContainer("shops")
	for _, id := range allShopIds {
		if ws.market.IsButtonJustClicked(id, ws.camera) {
			var text = ws.market.Field(id, field.Text, ws.camera)
			ws.market.SetField("shops", field.Hidden, "1")
			ws.market.SetField("shop-menu-dropdown", field.Text, text)
		}
	}
}
