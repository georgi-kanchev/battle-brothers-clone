package world

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/gui/field"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/text"
)

var inventorySelectedUnitIndex int

func (ws *WorldScreen) handleInventoryPopup() {
	ws.tryExitPopup(ws.inventory, nil, nil)

	var x, y, _, _, _ = ws.inventory.Area("display", ws.camera)
	var sc = 1 / ws.camera.Zoom * ws.inventory.Scale
	var cx, cy = x + 35*sc, y + 125*sc
	var units = ws.playerParty.units
	var selectedUnit = ws.playerParty.units[inventorySelectedUnitIndex]
	var r, g, b, _ = color.Channels(selectedUnit.NameColor)

	ws.inventory.SetField("name-label", field.Text, selectedUnit.NickAndName())
	ws.inventory.SetField("name-label", field.TextColor, text.New(r, " ", g, " ", b))
	ws.inventory.SetField("movement", field.Text, text.New(" Movement = ", selectedUnit.BaseMovePoints))
	ws.inventory.SetField("initiative", field.Text, text.New(" Initiative = ", selectedUnit.BaseInitiative))

	for i := range 20 {
		var hidden = condition.If(i < len(units), "", "1")
		var unitId = text.New("unit", i)
		var ux, uy, _, _, _ = ws.inventory.Area(unitId, ws.camera)
		ws.inventory.SetField(unitId, field.Hidden, hidden)

		if hidden == "" {
			var cx, cy, cw, ch, _ = ws.inventory.Area("units", ws.camera)
			var tlx, tly = ws.camera.PointToScreen(cx, cy)
			var brx, bry = ws.camera.PointToScreen(cx+cw, cy+ch)
			ws.camera.Mask(tlx, tly, brx-tlx, bry-tly)
			selectedUnit.UpdateAndDraw(ux+50*sc, uy+75*sc, sc*0.85, sc*0.85, ws.camera)
			ws.camera.SetScreenAreaToWindow()
		}

		if i < len(units) {
			var r, g, b, _ = color.Channels(units[i].NameColor)
			ws.inventory.SetField(unitId, field.Text, units[i].NickAndName()+"   \n@@@  ")
			ws.inventory.SetField(unitId, field.TextColor, text.New(r, " ", g, " ", b))
		}

		if ws.inventory.IsButtonJustClicked(unitId, ws.camera) {
			inventorySelectedUnitIndex = i
		}
	}

	selectedUnit.UpdateAndDraw(cx, cy, -sc*1.35, sc*1.35, ws.camera)
	// ws.camera.DrawQuad(x, y, w, h, a, palette.Red)
}
