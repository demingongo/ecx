package updateserviceapp

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	formmmodel "github.com/demingongo/ecx/bubbles/formmodel"
	"github.com/demingongo/ecx/globals"
)

func generateFormProcess() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("confirm").
				Title("").
				Negative("Cancel").
				Affirmative("Proceed").
				Inline(true),
		),
	).
		WithTheme(globals.Theme).
		WithWidth(globals.FormWidth)

	return form
}

func runFormProcess() *huh.Form {

	form := generateFormProcess()
	fModel := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:         form,
		InfoBubble:   info,
		VerticalMode: true,
	}).Width(globals.Width)

	tea.NewProgram(&fModel).Run()

	return form
}
