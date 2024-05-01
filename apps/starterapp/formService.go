package starterapp

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	formmmodel "github.com/demingongo/ecx/bubbles/formmodel"
)

func selectServiceJSON(info string) string {
	value := selectJSONFile("Pick a service (.json):", "", info)
	return value
}

func generateFormService() *huh.Form {

	confirm := true

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("confirm").
				Title("Create a service?").
				Value(&confirm),
		),
	).
		WithTheme(theme).
		WithWidth(formWidth)

	return form
}

func runFormService() *huh.Form {

	form := generateFormService()
	fModel := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:       form,
		InfoBubble: info,
	}).Width(width)

	tea.NewProgram(&fModel).Run()

	return form
}
