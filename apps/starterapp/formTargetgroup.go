package starterapp

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/demingongo/ecx/aws"
	formmmodel "github.com/demingongo/ecx/bubbles/formmodel"
)

func selectTargetGroupJSON(info string) string {
	value := selectJSONFile("Pick a target group (.json):", "", info)
	return value
}

func generateFormTargetgroup(list []aws.TargetGroup) *huh.Form {
	options := []huh.Option[aws.TargetGroup]{
		huh.NewOption("(None)", aws.TargetGroup{}),
	}

	arnTextMaxSize := 12

	for _, tg := range list {
		var arnText string

		if tg.TargetGroupArn != "" {
			if len(tg.TargetGroupArn) > arnTextMaxSize {
				arnText += " (..." + tg.TargetGroupArn[len(tg.TargetGroupArn)-arnTextMaxSize:] + ")"
			} else {
				arnText += " (" + tg.TargetGroupArn + ")"
			}
		}
		options = append(options, huh.NewOption(tg.TargetGroupName+arnText, tg))
	}

	confirm := true

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[aws.TargetGroup]().
				Title("Select a target group:").
				Key("targetgroup").
				Options(
					options...,
				).Height(6),

			huh.NewConfirm().
				Key("confirm").
				Title("Are you sure?").
				Value(&confirm).
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

func runFormTargetgroup(list []aws.TargetGroup) *huh.Form {

	form := generateFormTargetgroup(list)
	fModel := formmmodel.NewModel(formmmodel.ModelConfig{
		Form:       form,
		InfoBubble: info,
	}).Width(width)

	tea.NewProgram(&fModel).Run()

	return form
}
