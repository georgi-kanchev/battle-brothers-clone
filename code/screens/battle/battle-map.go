package battle

import (
	"game/code/global"
	"pure-game-kit/geometry"
)

func (bs *BattleScreen) recalculatePathMap() {
	var tw, th = global.BattleTileWidth, global.BattleTileHeight
	var w, h = bs.above.Size()
	bs.pathMap = geometry.NewShapeGrid(int(tw), int(th))
	for i := range h {
		for j := range w {
			var pts = bs.above.PointsAtCell(j, i)
			if len(pts) > 0 {
				bs.pathMap.SetAtCell(j, i, geometry.NewShapeCorners(pts...))
			}
		}
	}

	for _, unit := range bs.units {
		if unit != bs.actingUnit() {
			var ux, uy = unit.Position()
			bs.pathMap.SetAtCell(int(ux/tw), int(uy/th), geometry.NewShapeQuad(tw/2, th/2, 0.5, 0.5))
		}
	}
}
