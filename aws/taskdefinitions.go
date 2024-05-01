package aws

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type ContainerPortMapping struct {
	Name        string
	PortMapping PortMapping
}

type PortMapping struct {
	ContainerPort int    `json:"containerPort"`
	HostPort      int    `json:"hostPort"`
	Name          string `json:"name"`
}

type ContainerDefinition struct {
	Name         string        `json:"name"`
	Image        string        `json:"image"`
	PortMappings []PortMapping `json:"portMappings"`
}

type TaskDefinition struct {
	TaskDefinitionArn    string                `json:"taskDefinitionArn"`
	ContainerDefinitions []ContainerDefinition `json:"containerDefinitions"`
}

type describeTaskDefinitionResponse struct {
	TaskDefinition TaskDefinition `json:"taskDefinition"`
}

func DescribeTaskDefinition(taskDefinitionArn string) (TaskDefinition, error) {
	result := TaskDefinition{}
	var args []string
	args = append(args, "ecs", "describe-task-definition", "--output", "json", "--no-paginate", "--task-definition", taskDefinitionArn)
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(2)
		result = TaskDefinition{
			TaskDefinitionArn: taskDefinitionArn,
			ContainerDefinitions: []ContainerDefinition{
				{
					Name:  "dmz-web",
					Image: "xxx/repository-dmz-web:tag",
					PortMappings: []PortMapping{
						{
							ContainerPort: 8080,
							HostPort:      0,
							Name:          "http",
						},
					},
				},
			},
		}
		return result, nil
	}

	var resp describeTaskDefinitionResponse
	_, err := execAWS(args, &resp)
	if err != nil {
		return result, err
	}
	log.Debug(resp)

	result = resp.TaskDefinition

	return result, nil
}

func ListPortMapping(taskDefinitionArn string) ([]ContainerPortMapping, error) {
	result := []ContainerPortMapping{}
	td, err := DescribeTaskDefinition(taskDefinitionArn)
	if err != nil {
		return result, err
	}
	for _, cd := range td.ContainerDefinitions {
		for _, pm := range cd.PortMappings {
			result = append(result, ContainerPortMapping{
				Name:        cd.Name,
				PortMapping: pm,
			})
		}
	}

	return result, nil
}
