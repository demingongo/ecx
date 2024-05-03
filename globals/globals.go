package globals

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

var (
	Theme *huh.Theme = huh.ThemeBase()
)

const (
	FormWidth = 60
	InfoWidth = 38

	Width = 100

	LogoEmpty   = "" //"ᶻ 𝗓 𐰁"
	LogoSuccess = "✔️"
	LogoError   = "❌"
	LogoInfo    = "" //"🛈"
)

func LoadGlobals() {
	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}
	if viper.GetBool("colors") {
		Theme = huh.ThemeDracula()
	}
}
