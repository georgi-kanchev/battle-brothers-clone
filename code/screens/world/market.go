package world

import "pure-game-kit/gui/field"

func (ws *WorldScreen) handleMarketPopup() {
	ws.market.SetField("title-bgr", field.Text, "Marketplace")
	ws.tryExitPopup(ws.market, ws.settlement, nil)

	var allShopIds = ws.market.WidgetIdsOfContainer("shops")
	for _, id := range allShopIds {
		if ws.market.IsButtonJustClicked(id, ws.camera) {
			var text = ws.market.Field(id, field.Text, ws.camera)
			ws.market.SetField("shops", field.Hidden, "1")
			ws.market.SetField("shop-menu-dropdown", field.Text, text)
		}
	}
}
