package starterapp

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	formmmodel "github.com/demingongo/ecx/bubbles/formmodel"
)

func selectRuleJSON(info string, title string, dir string) string {
	value := selectJSONFile(title, dir, info)
	return value
}

func generateFormRules() *huh.Form {

	confirm := true

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("confirm").
				Title("Create rules?").
				Value(&confirm),
		),
	).
		WithTheme(theme).
		WithWidth(formWidth)

	return form
}

func runFormRules() *huh.Form {

	form := generateFormRules()
	fModel := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:       form,
		InfoBubble: info,
	}).Width(width)

	tea.NewProgram(&fModel).Run()

	return form
}
