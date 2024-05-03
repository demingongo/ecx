package updateserviceapp

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/demingongo/ecx/aws"
	formmmodel "github.com/demingongo/ecx/bubbles/formmodel"
	"github.com/demingongo/ecx/globals"
)

type ComparableService struct {
	ServiceArn  string
	ServiceName string
}

func generateFormService(list []aws.Service) *huh.Form {
	options := []huh.Option[ComparableService]{
		huh.NewOption("(None)", ComparableService{}),
	}

	for _, s := range list {
		if s.ServiceArn != "" {
			var text string
			if s.ServiceName != "" {
				text = s.ServiceName
			} else {
				text = s.ServiceArn
			}
			options = append(options, huh.NewOption(text, ComparableService{
				ServiceArn:  s.ServiceArn,
				ServiceName: s.ServiceName,
			}))
		}
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[ComparableService]().
				Title("Select a service:").
				Key("service").
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

func runFormService(list []aws.Service) *huh.Form {

	form := generateFormService(list)
	fModel := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:       form,
		InfoBubble: info,
	}).Width(globals.Width)

	tea.NewProgram(&fModel).Run()

	return form
}
