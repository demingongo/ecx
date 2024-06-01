package starterapp

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/demingongo/ecx/aws"
	"github.com/demingongo/ecx/bubbles/filepickermodel"
	"github.com/demingongo/ecx/globals"
	"github.com/spf13/viper"
)

type TargetGroupConfig struct {
	New      bool   // must create a target group if it's true
	Filepath string // must be filled if New=true
	Arn      string // target group's arn
	Name     string // target group's name
}

func (tgc TargetGroupConfig) IsComplete() bool {
	return (tgc.New && tgc.Filepath != "") || tgc.Arn != ""
}

func (tgc TargetGroupConfig) IsNew() bool {
	return (tgc.New && tgc.Filepath != "")
}

type ServiceConfig struct {
	Filepath       string // must be filled if New=true
	Name           string // service's name
	TaskDefinition string // task definition (containers)
}

type Config struct {
	targetGroup   TargetGroupConfig
	rules         []string
	service       ServiceConfig
	containerName string
	containerPort int

	targetGroupDescription string
	rulesDescription       string
	serviceDescription     string

	targetGroupLogo string
	rulesLogo       string
	serviceLogo     string
}

type filepickerStyleStruct struct {
	cursor    lipgloss.Style
	directory lipgloss.Style
	file      lipgloss.Style
	selected  lipgloss.Style
	symlink   lipgloss.Style
}

var (

	// General.

	config Config
	info   string

	subtle  = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	special = lipgloss.AdaptiveColor{Light: "230", Dark: "#010102"}

	subtleText = lipgloss.NewStyle().Foreground(subtle).Render

	// Titles.

	titleStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("7")).
			Foreground(special)

	subtitleStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle).
			Foreground(lipgloss.Color("6"))

	// Info block.

	infoStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("7")).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true).
			Width(globals.InfoWidth)

	// filepicker
	filepickerStyle = filepickerStyleStruct{
		cursor:    lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
		directory: lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Underline(true),
		file:      lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		selected:  lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
		symlink:   lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Italic(true),
	}
)

func isProcessable() bool {
	return config.targetGroup.IsComplete() || len(config.rules) > 0 || config.service.Filepath != ""
}

func selectJSONFile(title string, currentDirectory string, info string) string {
	m := filepickermodel.NewFilepickerModel(filepickermodel.FilepickerModelConfig{
		AllowedTypes:     []string{".json"},
		CurrentDirectory: currentDirectory,
		EnableFastSelect: true,
		Title:            title,
		InfoBubble:       info,
	}).
		ShowPermissions(false).
		ShowSize(false).
		Height(8).
		Width(globals.Width).
		FilepickerWidth(globals.FormWidth).
		StyleDirectory(filepickerStyle.directory).
		StyleFile(filepickerStyle.file).
		StyleSymlink(filepickerStyle.symlink).
		StyleCursor(filepickerStyle.cursor).
		StyleSelected(filepickerStyle.selected)

	tm, _ := tea.NewProgram(&m).Run()

	mm := tm.(filepickermodel.FilepickerModel)

	return mm.SelectedFile
}

func generateInfo() string {

	var (
		tgInfo      string
		rulesInfo   string
		serviceInfo string
	)

	if config.targetGroupDescription != "" {
		tgInfo = config.targetGroupDescription
	} else {
		if config.targetGroup.New {
			tgInfo = config.targetGroup.Filepath
		} else {
			tgInfo = config.targetGroup.Name
		}
	}

	if config.rulesDescription != "" {
		rulesInfo = config.rulesDescription
	} else {
		rulesInfo = strings.Join(config.rules, ", ")
	}

	if config.serviceDescription != "" {
		serviceInfo = config.serviceDescription
	} else {
		serviceInfo = config.service.Filepath
	}

	if len(tgInfo) == 0 {
		tgInfo = subtleText("-")
	}
	if len(rulesInfo) == 0 {
		rulesInfo = subtleText("-")
	}
	if len(serviceInfo) == 0 {
		serviceInfo = subtleText("-")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("SUMMARY"),
		subtitleStyle.Render("Target group "+config.targetGroupLogo),
		tgInfo,
		subtitleStyle.Render("Rules "+config.rulesLogo),
		rulesInfo,
		subtitleStyle.Render("Service "+config.serviceLogo),
		serviceInfo,
	)

	return infoStyle.Render(content)
}

func generateDescription(name string, filepath string) string {

	r := name

	filepathMaxSize := 57

	if filepath != "" {
		if len(filepath) > filepathMaxSize {
			r += " (..." + filepath[len(filepath)-filepathMaxSize:] + ")"
		} else {
			r += " (" + filepath + ")"
		}
	}

	return r
}

func process(logger *log.Logger) {
	if config.targetGroup.IsNew() {
		logger.Debug(fmt.Sprintf("create target group \"%s\"", config.targetGroup.Name))
		var (
			result aws.TargetGroup
			err    error
		)
		_ = spinner.New().Type(spinner.Meter).
			Title(fmt.Sprintf(" Creating target group \"%s\"...", config.targetGroup.Name)).
			Action(func() {
				result, err = aws.CreateTargetGroup(config.targetGroup.Filepath)
			}).
			Run()

		if err != nil {
			config.targetGroupLogo = globals.LogoError
			info = generateInfo()
			fmt.Println(info)
			logger.Fatal("CreateTargetGroup", err)
		}
		config.targetGroup.Arn = result.TargetGroupArn
		config.targetGroupLogo = globals.LogoSuccess
	}

	if len(config.rules) > 0 {
		logger.Debug(fmt.Sprintf("create rules for target group \"%s\"", config.targetGroup.Name))
		for i, v := range config.rules {
			var err error
			_ = spinner.New().Type(spinner.Meter).
				Title(fmt.Sprintf(" Creating rules (%d/%d)...", i+1, len(config.rules))).
				Action(func() {
					_, err = aws.CreateRule(v, config.targetGroup.Arn)
				}).
				Run()
			if err != nil {
				config.rulesLogo = globals.LogoError
				info = generateInfo()
				fmt.Println(info)
				logger.Fatal("CreateRule", err)
			}
		}
		config.rulesLogo = globals.LogoSuccess
	}

	if len(config.service.Filepath) > 0 {
		logger.Debug(fmt.Sprintf("create service \"%s\"", config.service.Name))
		var err error
		_ = spinner.New().Type(spinner.Meter).
			Title(fmt.Sprintf(" Creating service \"%s\"...", config.service.Name)).
			Action(func() {
				_, err = aws.CreateService(config.service.Filepath, aws.ServiceLoadBalancer{
					TargetGroupArn: config.targetGroup.Arn,
					ContainerName:  config.containerName,
					ContainerPort:  config.containerPort,
				})
			}).
			Run()
		if err != nil {
			config.serviceLogo = globals.LogoError
			info = generateInfo()
			fmt.Println(info)
			logger.Fatal("CreateService", err)
		}
		config.serviceLogo = globals.LogoSuccess
	}

	info = generateInfo()
	fmt.Println(info)
}

func Run() {

	logger := globals.Logger

	info = generateInfo()

	menuForm := runFormMenu()

	if menuForm.State == huh.StateCompleted && menuForm.GetString("operation") != "none" {

		operation := menuForm.GetString("operation")

		// create-targetgroup
		if operation == "create-targetgroup" {
			config.targetGroup.New = true
			config.targetGroup.Filepath = selectTargetGroupJSON(info)

			if config.targetGroup.Filepath != "" {
				tgConf := viper.New()
				tgConf.SetConfigFile(config.targetGroup.Filepath)
				err := tgConf.ReadInConfig()
				if err != nil {
					logger.Fatal("Could not read file:", err)
				}

				config.targetGroup.Name = tgConf.GetString("Name")

				config.targetGroupDescription = generateDescription(config.targetGroup.Name, config.targetGroup.Filepath)
			}
		}

		// select-targetgroup
		if operation == "select-targetgroup" {
			var (
				targetgroups []aws.TargetGroup
				err          error
			)
			_ = spinner.New().Type(spinner.Globe).
				Title(" Searching target groups...").
				Action(func() {
					targetgroups, err = aws.DescribeTargetGroups()
				}).
				Run()
			if err != nil {
				logger.Fatal(err)
			}

			targetGroupForm := runFormTargetgroup(targetgroups)
			if targetGroupForm.State == huh.StateCompleted {
				tg := targetGroupForm.Get("targetgroup").(aws.TargetGroup)
				if tg.TargetGroupArn != "" {
					config.targetGroup.Arn = tg.TargetGroupArn
					config.targetGroup.Name = tg.TargetGroupName
					config.targetGroupDescription = generateDescription(tg.TargetGroupName, tg.TargetGroupArn)
					config.targetGroupLogo = globals.LogoInfo
				}
			}
		}
		if config.targetGroupDescription == "" {
			config.targetGroupLogo = globals.LogoEmpty
		}
		info = generateInfo()

		// create rules
		if config.targetGroup.IsComplete() {
			rulesForm := runFormRules()
			if rulesForm.State == huh.StateCompleted && rulesForm.GetBool("confirm") {
				var searchDir string
				var maxRules = 10
				for len(config.rules) < 10 {
					title := fmt.Sprintf("Pick a rule (.json) (%d/%d):", len(config.rules), maxRules)
					file := selectRuleJSON(info, title, searchDir)
					if len(file) > 0 {
						if slices.Contains(config.rules, file) {
							break
						} else {
							config.rules = append(config.rules, file)
							searchDir = filepath.Dir(file)
							info = generateInfo()
						}
					} else {
						break
					}
				}
			}
		}
		if len(config.rules) == 0 {
			config.rulesLogo = globals.LogoEmpty
			info = generateInfo()
		}

		// create service
		if operation == "create-targetgroup" || operation == "select-targetgroup" {
			serviceForm := runFormService()
			if serviceForm.State == huh.StateCompleted && serviceForm.GetBool("confirm") {
				config.service.Filepath = selectServiceJSON(info)
			}
		} else if operation == "create-service" {
			config.service.Filepath = selectServiceJSON(info)
		}
		if config.service.Filepath == "" {
			config.serviceLogo = globals.LogoEmpty
		} else {
			tgConf := viper.New()
			tgConf.SetConfigFile(config.service.Filepath)
			err := tgConf.ReadInConfig()
			if err != nil {
				logger.Fatal("Could not read file:", err)
			}
			config.service.Name = tgConf.GetString("serviceName")
			config.service.TaskDefinition = tgConf.GetString("taskDefinition")
			config.serviceDescription = generateDescription(config.service.Name, config.service.Filepath)
		}
		info = generateInfo()

		// create load balancer for service
		if config.service.TaskDefinition != "" && config.targetGroup.IsComplete() {
			// select container and port
			var (
				containers []aws.ContainerPortMapping
				err        error
			)
			_ = spinner.New().Type(spinner.Points).
				Title(" Checking task definition containers...").
				Action(func() {
					containers, err = aws.ListPortMapping(config.service.TaskDefinition)
				}).
				Run()
			if err != nil {
				logger.Fatal(err)
			}
			if len(containers) > 0 {
				lbForm := runFormLoadBalancer(containers)
				if lbForm.State == huh.StateCompleted {
					container := lbForm.Get("loadbalancer").(aws.ContainerPortMapping)
					if container.Name != "" && container.PortMapping.ContainerPort > 0 {
						config.containerName = container.Name
						config.containerPort = container.PortMapping.ContainerPort
					}
				}
			}
		}

		//fmt.Println(info)

		if isProcessable() {
			if form := runFormProcess(); form.State == huh.StateCompleted && form.GetBool("confirm") {
				process(logger)
			}
		}
	}

	fmt.Println("Done")
}
