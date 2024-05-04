package updateserviceapp

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/demingongo/ecx/aws"
	formmmodel "github.com/demingongo/ecx/bubbles/formmodel"
	"github.com/demingongo/ecx/globals"
)

func generateFormSelectImage(description string, list []aws.Image) *huh.Form {
	options := []huh.Option[aws.Image]{
		huh.NewOption("(leave unchanged)", aws.Image{}),
	}

	for _, image := range list {
		options = append(options, huh.NewOption(image.ImageTag, image))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[aws.Image]().
				Title("Select an image:").
				Description(description).
				Key("image").
				Options(
					options...,
				).Height(10),
		),
	).
		WithTheme(globals.Theme).
		WithWidth(globals.FormWidth)

	return form
}

func runFormSelectImage(description string, list []aws.Image) *huh.Form {

	form := generateFormSelectImage(description, list)
	fModel := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:       form,
		InfoBubble: info,
	}).Width(globals.Width)

	tea.NewProgram(&fModel).Run()

	return form
}

func generateFormInputImage(description string, placeholder string) *huh.Form {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Placeholder(placeholder).
				Suggestions([]string{placeholder}).
				Title("Select an image:").
				Description(description).
				Key("image").WithHeight(6),
		),
	).
		WithTheme(globals.Theme).
		WithWidth(globals.FormWidth)

	return form
}

func runFormInputImage(description string, placeholder string) *huh.Form {

	form := generateFormInputImage(description, placeholder)
	fModel := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:       form,
		InfoBubble: info,
	}).Width(globals.Width)

	tea.NewProgram(&fModel).Run()

	return form
}
