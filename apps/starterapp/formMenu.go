package starterapp

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	formmmodel "github.com/demingongo/ecx/bubbles/formmodel"
)

func generateFormMenu() *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What to do:").
				Description("Steps to create a service").
				Key("operation").
				Options(
					huh.NewOption("Create a target group", "create-targetgroup"),
					huh.NewOption("Select a target group", "select-targetgroup"),
					huh.NewOption("Create a service", "create-service"),
					huh.NewOption("Nothing", "none"),
				),

			huh.NewConfirm().
				Key("confirm").
				Title("Are you sure?").
				Validate(func(b bool) error {
					if !b {
						return errors.New("waiting till you confirm")
					}
					return nil
				}),
		),
	).
		WithTheme(theme).
		WithWidth(formWidth)

	return form
}

func runFormMenu() *huh.Form {

	form := generateFormMenu()
	fModel := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:       form,
		InfoBubble: info,
	}).Width(width)

	tea.NewProgram(&fModel).Run()

	return form
}
