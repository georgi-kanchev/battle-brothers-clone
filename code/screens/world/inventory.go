package world

import (
	"pure-game-kit/utility/color/palette"
)

func (ws *WorldScreen) handleInventoryPopup() {
	var x, y, w, h, a = ws.inventory.Area("display", ws.camera)
	ws.camera.DrawQuad(x, y, w, h, a, palette.Red)
}
