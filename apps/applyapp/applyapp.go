package applyapp

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh/spinner"
	"github.com/demingongo/ecx/aws"
	"github.com/demingongo/ecx/globals"
	"gopkg.in/yaml.v2"
)

type LogGroup struct {
	Group     string `yaml:"group"`
	Retention int    `yaml:"retention"`
}

type Flow struct {
	Name        string   `yaml:"name"`
	Service     string   `yaml:"service"`
	TargetGroup string   `yaml:"targetGroup"`
	Rules       []string `yaml:"rules"`
}

type Config struct {
	Api             string     `yaml:"api"`
	ApiVersion      string     `yaml:"apiVersion"`
	LogGroups       []LogGroup `yaml:"logGroups"`
	TaskDefinitions []string   `yaml:"taskDefinitions"`
	Flows           []Flow     `yaml:"flows"`
}

func (c *Config) loadConfig() *Config {

	yamlFile, err := os.ReadFile("../../ecx-tests/project1/ecx.yaml")
	if err != nil {
		globals.Logger.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		globals.Logger.Fatalf("Unmarshal: %v", err)
	}

	return c
}

var (
	config Config

	validApi        = "ecx"
	validApiVersion = "0.1"
)

func Run() {
	logger := globals.Logger

	logger.Debug("ecx apply")

	config.loadConfig()

	logger.Debug(config)

	if config.Api != validApi {
		logger.Fatalf("Value for \"%s\" is not valid. Expected \"%s\".", "api", validApi)
	}
	if config.ApiVersion != "0.1" {
		logger.Fatalf("Value for \"%s\" is not valid. Expected \"%s\".", "apiVersion", validApiVersion)
	}

	if len(config.TaskDefinitions) > 0 {
		var err error
		for _, taskDefinitionFile := range config.TaskDefinitions {
			_ = spinner.New().Type(spinner.Meter).
				Title(fmt.Sprintf(" task definition: %s", taskDefinitionFile)).
				Action(func() {
					// create new revision for task definition
					_, err = aws.RegisterTaskDefinition(fmt.Sprintf("file://%s", taskDefinitionFile))
				}).
				Run()
			if err != nil {
				logger.Fatalf("RegisterTaskDefinition: %v", err)
			}
			fmt.Printf("task definition: %s\n", taskDefinitionFile)
		}
	}

	if len(config.LogGroups) > 0 {
		var err error
		for _, logGroup := range config.LogGroups {
			_ = spinner.New().Type(spinner.Meter).
				Title(fmt.Sprintf(" log group: %s", logGroup.Group)).
				Action(func() {
					// create log group
					aws.CreateLogGroup(logGroup.Group)
					if logGroup.Retention > 0 {
						// put retention policy in number of days
						_, err = aws.PutRetentionPolicy(logGroup.Group, logGroup.Retention)
					}
				}).
				Run()
			if err != nil {
				logger.Fatalf("CreateLogGroup: %v", err)
			}
			fmt.Printf("log group: %s\n", logGroup.Group)
		}
	}
}
