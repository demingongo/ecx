package updateserviceapp

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/demingongo/ecx/aws"
	formmmodel "github.com/demingongo/ecx/bubbles/formmodel"
	"github.com/demingongo/ecx/globals"
)

type ComparableContainerDefinition struct {
	Name  string
	Image string
}

func generateFormSelectContainers(list []aws.ContainerDefinition) *huh.Form {
	options := []huh.Option[ComparableContainerDefinition]{}

	for _, containerDef := range list {
		options = append(options, huh.NewOption(containerDef.Name, ComparableContainerDefinition{
			Name:  containerDef.Name,
			Image: containerDef.Image,
		}))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[ComparableContainerDefinition]().
				Title("Select containers to update:").
				Key("containers").
				Options(
					options...,
				).Height(6),

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
		WithTheme(globals.Theme).
		WithWidth(globals.FormWidth)

	return form
}

func runFormSelectContainers(list []aws.ContainerDefinition) *huh.Form {

	form := generateFormSelectContainers(list)
	fModel := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:       form,
		InfoBubble: info,
	}).Width(globals.Width)

	tea.NewProgram(&fModel).Run()

	return form
}
