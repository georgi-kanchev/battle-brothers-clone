package main

import (
	"game/source-code/screens"
	"game/source-code/screens/world"
	"pure-game-kit/data/assets"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/window"
)

const Version = "v0.0.3"

func main() {
	window.Title = "Battle Brothers Clone"

	assets.LoadDefaultAtlasIcons()
	assets.LoadDefaultFont()
	assets.LoadDefaultTexture()

	screens.New(nil, world.New(), nil)
	screens.LoadAll()

	for window.KeepOpen() {
		screens.UpdateCurrent()

		if keyboard.IsKeyJustPressed(key.F5) {
			assets.ReloadAll()
			screens.LoadAll()
		}
	}
}
