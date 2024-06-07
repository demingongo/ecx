package applyapp

import (
	"fmt"
	"os"
	"strings"

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

type FlowRule struct {
	Listener    string `yaml:"listener"`
	Priority    int    `yaml:"priority"`
	TargetGroup string `yaml:"targetGroup"`
	Value       string `yaml:"value"`
}

type Flow struct {
	Name                          string     `yaml:"name"`
	Service                       string     `yaml:"service"`
	TargetGroup                   string     `yaml:"targetGroup"`
	HealthCheckGracePeriodSeconds int        `yaml:"healthCheckGracePeriodSeconds"`
	Rules                         []FlowRule `yaml:"rules"`
}

type LoadBalancer struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Rule struct {
	Priority    int    `yaml:"priority"`
	TargetGroup string `yaml:"targetGroup"`
	Value       string `yaml:"value"`
}

type Listener struct {
	Key          string `yaml:"key"`
	Value        string `yaml:"value"`
	LoadBalancer string `yaml:"loadBalancer"`
	TargetGroup  string `yaml:"targetGroup"`
	Rules        []Rule `yaml:"rules"`
}

type TargetGroup struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Config struct {
	Api             string         `yaml:"api"`
	ApiVersion      string         `yaml:"apiVersion"`
	LogGroups       []LogGroup     `yaml:"logGroups"`
	TaskDefinitions []string       `yaml:"taskDefinitions"`
	Flows           []Flow         `yaml:"flows"`
	LoadBalancers   []LoadBalancer `yaml:"loadBalancers"`
	Listeners       []Listener     `yaml:"listeners"`
	TargetGroups    []TargetGroup  `yaml:"targetGroups"`
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

type ConfigRefs struct {
	LoadBalancers map[string]aws.LoadBalancer
	Listeners     map[string]aws.Listener
	TargetGroups  map[string]aws.TargetGroup
}

var (
	config Config

	validApi        = "ecx"
	validApiVersion = "0.1"
)

func createConfigRefs() ConfigRefs {
	var configRefs ConfigRefs
	configRefs.LoadBalancers = make(map[string]aws.LoadBalancer)
	configRefs.Listeners = make(map[string]aws.Listener)
	configRefs.TargetGroups = make(map[string]aws.TargetGroup)
	return configRefs
}

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

	// references
	refs := createConfigRefs()

	// @TODO targetGroups
	if len(config.TargetGroups) > 0 {
		logger.Debugf("TargetGroups: %v", config.TargetGroups)
		for _, targetGroup := range config.TargetGroups {
			if targetGroup.Value != "" {
				var (
					err  error
					resp aws.TargetGroup
				)
				// get name from file
				content := viper.New()
				content.SetConfigFile(targetGroup.Value)
				err = content.ReadInConfig()
				if err != nil {
					logger.Fatalf("checking target group %s: %v", targetGroup.Key, err)
				}
				name := content.GetString("Name")
				if name != "" {
					found := false
					_ = spinner.New().Type(spinner.Pulse).
						Title(fmt.Sprintf(" DescribeTargetGroupsWithNames %s", name)).
						Action(func() {
							// do not handle error below as aws cli
							// returns error if one name is not found
							results, _ := aws.DescribeTargetGroupsWithNames([]string{name})
							if len(results) > 0 {
								if targetGroup.Key != "" {
									refs.TargetGroups[targetGroup.Key] = results[0]
								}
								found = true
							}
						}).
						Run()
					if found {
						logger.Infof("target group named \"%s\" already exists", name)
						continue
					}
				}

				_ = spinner.New().Type(spinner.MiniDot).
					Title(fmt.Sprintf(" target group: %s", targetGroup.Key)).
					Action(func() {
						// create target group
						resp, err = aws.CreateTargetGroup(targetGroup.Value)
					}).
					Run()
				if err != nil {
					logger.Fatalf("CreateTargetGroup: %v", err)
				}
				if targetGroup.Key != "" && resp.TargetGroupArn != "" {
					refs.TargetGroups[targetGroup.Key] = resp
				}
				fmt.Printf("target group: %s\n", targetGroup.Key)
			}
		}
	}

	// @TODO loadBalancers
	if len(config.LoadBalancers) > 0 {
		logger.Debugf("LoadBalancers: %v", config.LoadBalancers)
		for _, loadBalancer := range config.LoadBalancers {
			if loadBalancer.Value != "" {
				var (
					err  error
					resp aws.LoadBalancer
				)
				// get name from file
				content := viper.New()
				content.SetConfigFile(loadBalancer.Value)
				err = content.ReadInConfig()
				if err != nil {
					logger.Fatalf("checking load balancer %s: %v", loadBalancer.Key, err)
				}
				name := content.GetString("Name")
				if name != "" {
					found := false
					_ = spinner.New().Type(spinner.Pulse).
						Title(fmt.Sprintf(" DescribeLoadBalancersWithNames %s", name)).
						Action(func() {
							// do not handle error below as aws cli
							// returns error if one name is not found
							results, _ := aws.DescribeLoadBalancersWithNames([]string{name})
							if len(results) > 0 {
								if loadBalancer.Key != "" {
									refs.LoadBalancers[loadBalancer.Key] = results[0]
								}
								found = true
							}
						}).
						Run()
					if found {
						logger.Infof("load balancer named \"%s\" already exists", name)
						continue
					}
				}

				_ = spinner.New().Type(spinner.MiniDot).
					Title(fmt.Sprintf(" load balancer: %s", loadBalancer.Key)).
					Action(func() {
						// create load balancer
						resp, err = aws.CreateLoadBalancer(loadBalancer.Value)
					}).
					Run()
				if err != nil {
					logger.Fatalf("CreateLoadBalancer: %v", err)
				}
				if loadBalancer.Key != "" && resp.LoadBalancerArn != "" {
					refs.LoadBalancers[loadBalancer.Key] = resp
				}
				fmt.Printf("load balancer: %s\n", loadBalancer.Key)
			}
		}
	}

	// @TODO listeners
	if len(config.Listeners) > 0 {
		logger.Debugf("Listeners: %v", config.Listeners)
		for _, listener := range config.Listeners {
			if listener.Value != "" {
				var (
					err     error
					resp    aws.Listener
					loading *spinner.Spinner
				)
				loading = spinner.New().Type(spinner.MiniDot).
					Title(fmt.Sprintf(" listener: %s", listener.Key)).
					Action(func() {
						var (
							lbArn string
							tgArn string
						)
						if listener.LoadBalancer != "" {
							if strings.HasPrefix(listener.LoadBalancer, "ref:") {
								key := listener.LoadBalancer[4:]
								lbArn = refs.LoadBalancers[key].LoadBalancerArn
								if lbArn == "" {
									logger.Fatalf("ecx - could not find load balancer reference \"%s\"", key)
								}
							} else {
								lbArn = listener.LoadBalancer
							}
						}
						if listener.TargetGroup != "" {
							if strings.HasPrefix(listener.TargetGroup, "ref:") {
								key := listener.TargetGroup[4:]
								tgArn = refs.TargetGroups[key].TargetGroupArn
								if tgArn == "" {
									logger.Fatalf("ecx - could not find target group reference \"%s\"", key)
								}
							} else {
								tgArn = listener.TargetGroup
							}
						}
						// create listener
						resp, err = aws.CreateListener(listener.Value, lbArn, tgArn)
						if err != nil {
							return
						}
						// create rules
						for _, rule := range listener.Rules {
							var ruleDestination string
							if strings.HasPrefix(rule.TargetGroup, "ref:") {
								key := rule.TargetGroup[4:]
								ruleDestination = refs.TargetGroups[key].TargetGroupArn
								if ruleDestination == "" {
									err = fmt.Errorf("ecx - could not find target group reference \"%s\"", key)
									break
								}
							} else {
								ruleDestination = rule.TargetGroup
							}
							loading.Title(fmt.Sprintf(" listener: %s - rule: %s", listener.Key, rule.Value))
							// create rule
							_, err = aws.CreateRule2(rule.Value, ruleDestination, rule.Priority, resp.ListenerArn)
							if err != nil {
								break
							}
						}
					})
				_ = loading.Run()
				if err != nil {
					logger.Fatalf("CreateListener: %v", err)
				}
				if listener.Key != "" && resp.ListenerArn != "" {
					refs.Listeners[listener.Key] = resp
				}
				fmt.Printf("listener: %s\n", listener.Key)
			}
		}
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

					var (
						targetGroup   aws.TargetGroup
						containerName string
						containerPort int
					)

					// create target group
					if flow.TargetGroup != "" {
						if strings.HasPrefix(flow.TargetGroup, "ref:") {
							key := flow.TargetGroup[4:]
							targetGroup = refs.TargetGroups[key]
							if targetGroup.TargetGroupArn == "" {
								err = fmt.Errorf("ecx - could not find target group reference \"%s\"", key)
								return
							}
						} else {
							targetGroup, err = aws.CreateTargetGroup(flow.TargetGroup)
						}
					}
					if err != nil {
						return
					}

					// create rules
					if targetGroup.TargetGroupArn != "" && len(flow.Rules) > 0 {
						logger.Info("letsgoooooo")
						for _, rule := range flow.Rules {
							var listenerArn string = rule.Listener
							if strings.HasPrefix(listenerArn, "ref:") {
								key := listenerArn[4:]
								listenerArn = refs.Listeners[key].ListenerArn
								if listenerArn == "" {
									err = fmt.Errorf("ecx - could not find listener reference \"%s\"", key)
									break
								}
							}
							_, err = aws.CreateRule2(rule.Value, targetGroup.TargetGroupArn, rule.Priority, listenerArn)
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
						}, flow.HealthCheckGracePeriodSeconds)
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
