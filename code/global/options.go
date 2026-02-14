package global

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/data/storage"
	"pure-game-kit/utility/number"
	"pure-game-kit/window"
)

type Options struct {
	WindowState   int  `yaml:"graphics-window-state"`
	Monitor       int  `yaml:"graphics-monitor"`
	Antialiasing  bool `yaml:"graphics-antialiasing"`
	TextureFilter bool `yaml:"graphics-texture-filter"`
	VSync         bool `yaml:"graphics-vsync"`
	LimitFPS      int  `yaml:"graphics-limit-fps"`

	ShowFPS                      bool    `yaml:"ui-show-fps"`
	ScaleUI                      float32 `yaml:"ui-scale"`
	ScaleMenuOptions             float32 `yaml:"ui-scale-menu-options"`
	ScaleWorldHUD                float32 `yaml:"ui-scale-world-hud"`
	ScaleWorldInventory          float32 `yaml:"ui-scale-world-inventory"`
	ScaleWorldEvents             float32 `yaml:"ui-scale-world-events"`
	ScaleWorldSettlement         float32 `yaml:"ui-scale-world-settlement"`
	ScaleWorldSettlementMarket   float32 `yaml:"ui-scale-world-settlement-market"`
	ScaleWorldSettlementFavors   float32 `yaml:"ui-scale-world-settlement-favors"`
	ScaleWorldSettlementRecruits float32 `yaml:"ui-scale-world-settlement-recruits"`
	ScaleWorldSettlementTavern   float32 `yaml:"ui-scale-world-settlement-tavern"`
	ScaleBattleHUD               float32 `yaml:"ui-scale-battle-hud"`
	ScaleBattleLoot              float32 `yaml:"ui-scale-battle-loot"`

	AudioVolume      float32 `yaml:"audio-volume"`
	AudioVolumeMusic float32 `yaml:"audio-volume-music"`
	AudioVolumeSound float32 `yaml:"audio-volume-sound"`
}

func LoadOptions() {
	var opts Options
	storage.FromYAML(file.LoadText("data/options.yaml"), &opts)
	Opts = &opts
}

func ApplyOptions() {
	window.IsVSynced = Opts.VSync
	window.IsAntialiased = Opts.Antialiasing
	window.FrameRateLimit = byte(number.Limit(Opts.LimitFPS, 0, 250))
	window.ApplyState(Opts.WindowState)
	window.MoveToMonitor(Opts.Monitor)

	var allTextures = assets.LoadedTextureIds()
	for _, tex := range allTextures {
		assets.SetTextureSmoothness(tex, Opts.TextureFilter)
	}
}
