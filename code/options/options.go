package options

type Options struct {
	WindowState                 int     `yaml:"graphics-window-state"`
	Monitor                     int     `yaml:"graphics-monitor"`
	VSync                       bool    `yaml:"graphics-vsync"`
	LimitFPS                    int     `yaml:"graphics-limit-fps"`
	ScaleUI                     float32 `yaml:"ui-scale"`
	ScaleMenuOptions            float32 `yaml:"ui-scale-menu-options"`
	ScaleWorldHUD               float32 `yaml:"ui-scale-world-hud"`
	ScaleWorldInventory         float32 `yaml:"ui-scale-world-inventory"`
	ScaleWorldEvents            float32 `yaml:"ui-scale-world-events"`
	ScaleWorldSettlement        float32 `yaml:"ui-scale-world-settlement"`
	ScaleWorldSettlementMarket  float32 `yaml:"ui-scale-world-settlement-market"`
	ScaleWorldSettlementFavors  float32 `yaml:"ui-scale-world-settlement-favors"`
	ScaleWorldSettlementRecruit float32 `yaml:"ui-scale-world-settlement-recruit"`
	ScaleWorldSettlementTavern  float32 `yaml:"ui-scale-world-settlement-tavern"`
	ScaleBattleHUD              float32 `yaml:"ui-scale-battle-hud"`
	ScaleBattleLoot             float32 `yaml:"ui-scale-battle-loot"`

	AudioVolume      float32 `yaml:"audio-volume"`
	AudioVolumeMusic float32 `yaml:"audio-volume-music"`
	AudioVolumeSound float32 `yaml:"audio-volume-sound"`
	// Controls     struct {
	// } `yaml:"controls"`
	// Gameplay struct {
	// } `yaml:"gameplay"`
}
