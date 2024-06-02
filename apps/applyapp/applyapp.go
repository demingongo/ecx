package applyapp

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh/spinner"
	"github.com/demingongo/ecx/aws"
	"github.com/demingongo/ecx/globals"
	"github.com/spf13/viper"
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

	yamlFile, err := os.ReadFile("ecx.yaml")
	if err != nil {
		globals.Logger.Fatalf("%v", err)
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

	if viper.GetString("project") != "" {
		err := os.Chdir(viper.GetString("project"))
		if err != nil {
			logger.Fatalf("project: %v", err)
		}
	}

	logger.Debugf("ecx apply %s", viper.GetString("project"))

	config.loadConfig()

	logger.Debug(config)

	if config.Api != validApi {
		logger.Fatalf("Value for \"%s\" is not valid. Expected \"%s\".", "api", validApi)
	}
	if config.ApiVersion != "0.1" {
		logger.Fatalf("Value for \"%s\" is not valid. Expected \"%s\".", "apiVersion", validApiVersion)
	}

	// logGroups
	if len(config.LogGroups) > 0 {
		var err error
		for _, logGroup := range config.LogGroups {
			_ = spinner.New().Type(spinner.MiniDot).
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

	// taskDefinitions
	if len(config.TaskDefinitions) > 0 {
		var err error
		for _, taskDefinitionFile := range config.TaskDefinitions {
			_ = spinner.New().Type(spinner.MiniDot).
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

	// flows
	if len(config.Flows) > 0 {
		var err error
		for _, flow := range config.Flows {
			_ = spinner.New().Type(spinner.MiniDot).
				Title(fmt.Sprintf(" flow: %v", flow)).
				Action(func() {
					// @TODO create target group, rules and/or service
					time.Sleep(2000 * time.Millisecond)

					var (
						targetGroup   aws.TargetGroup
						containerName string
						containerPort int
					)

					// create target group
					if flow.TargetGroup != "" {
						targetGroup, err = aws.CreateTargetGroup(flow.TargetGroup)
					}
					if err != nil {
						return
					}

					// create rules
					if targetGroup.TargetGroupArn != "" && len(flow.Rules) > 0 {
						for _, rule := range flow.Rules {
							_, err = aws.CreateRule(rule, targetGroup.TargetGroupArn)
							if err != nil {
								break
							}
						}
					}
					if err != nil {
						return
					}

					// create service
					if flow.Service != "" {
						// get port mapping named "http"
						// or the first port mapping
						if targetGroup.TargetGroupArn != "" {
							var containers []aws.ContainerPortMapping
							serviceConf := viper.New()
							serviceConf.SetConfigFile(flow.Service)
							err = serviceConf.ReadInConfig()
							if err != nil {
								return
							}
							serviceName := serviceConf.GetString("serviceName")
							taskDefinition := serviceConf.GetString("taskDefinition")

							logger.Debugf("serviceName %s", serviceName)
							logger.Debugf("taskDefinition %s", taskDefinition)

							containers, err = aws.ListPortMapping(taskDefinition)
							if err != nil {
								return
							}

							for i, container := range containers {
								if container.PortMapping.Name == "http" {
									containerName = container.Name
									containerPort = container.PortMapping.ContainerPort
									break
								}
								if i == 0 {
									containerName = container.Name
									containerPort = container.PortMapping.ContainerPort
								}
							}
						}

						_, err = aws.CreateService(flow.Service, aws.ServiceLoadBalancer{
							TargetGroupArn: targetGroup.TargetGroupArn,
							ContainerName:  containerName,
							ContainerPort:  containerPort,
						})
					}
				}).
				Run()
			if err != nil {
				logger.Fatalf("flow: %v", err)
			}
			fmt.Printf("flow: %v\n", flow)
		}
	}

	fmt.Println("Done")
}
