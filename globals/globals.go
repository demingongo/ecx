package globals

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

var (
	Theme  *huh.Theme = huh.ThemeBase()
	Logger *log.Logger
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

	// Create logger
	styles := log.DefaultStyles()
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("50")).
		Foreground(lipgloss.Color("0"))
	// Add a custom style for key `err`
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("50"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)
	Logger = log.New(os.Stderr)
	Logger.SetStyles(styles)
	Logger.SetLevel(log.GetLevel())
}
