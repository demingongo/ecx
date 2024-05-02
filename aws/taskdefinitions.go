package aws

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type ContainerPortMapping struct {
	Name        string
	PortMapping PortMapping
}

type PortMapping struct {
	ContainerPort      int    `json:"containerPort"`
	HostPort           int    `json:"hostPort"`
	Protocol           string `json:"protocol"`
	Name               string `json:"name"`
	AppProtocol        string `json:"appProtocol"`
	ContainerPortRange string `json:"containerPortRange"`
}

type ContainerDefinition struct {
	Name              string        `json:"name"`
	Image             string        `json:"image"`
	Cpu               int           `json:"cpu"`
	Memory            int           `json:"memory"`
	MemoryReservation int           `json:"memoryReservation"`
	PortMappings      []PortMapping `json:"portMappings"`
	Essential         bool          `json:"essential"`

	// not caring about the following props for now, but need to
	// define them to register a new task definition revision from one revision

	Links                  []any `json:"links"`
	RepositoryCredentials  any   `json:"repositoryCredentials"`
	EntryPoint             []any `json:"entryPoint"`
	Command                []any `json:"command"`
	Environment            []any `json:"environment"`
	EnvironmentFiles       []any `json:"environmentFiles"`
	MountPoints            []any `json:"mountPoints"`
	VolumesFrom            []any `json:"volumesFrom"`
	LinuxParameters        any   `json:"linuxParameters"`
	Secrets                []any `json:"secrets"`
	DependsOn              []any `json:"dependsOn"`
	StartTimeout           any   `json:"startTimeout"`
	StopTimeout            any   `json:"stopTimeout"`
	Hostname               any   `json:"hostname"`
	User                   any   `json:"user"`
	WorkingDirectory       any   `json:"workingDirectory"`
	DisableNetworking      bool  `json:"disableNetworking"`
	Privileged             bool  `json:"privileged"`
	ReadonlyRootFilesystem bool  `json:"readonlyRootFilesystem"`
	DnsServers             []any `json:"dnsServers"`
	DnsSearchDomains       []any `json:"dnsSearchDomains"`
	ExtraHosts             []any `json:"extraHosts"`
	DockerSecurityOptions  []any `json:"dockerSecurityOptions"`
	Interactive            bool  `json:"interactive"`
	PseudoTerminal         bool  `json:"pseudoTerminal"`
	DockerLabels           any   `json:"dockerLabels"`
	Ulimits                []any `json:"ulimits"`
	LogConfiguration       any   `json:"logConfiguration"`
	HealthCheck            any   `json:"healthCheck"`
	SystemControls         []any `json:"systemControls"`
	ResourceRequirements   []any `json:"resourceRequirements"`
	FirelensConfiguration  any   `json:"firelensConfiguration"`
	CredentialSpecs        []any `json:"credentialSpecs"`
}

type TaskDefinition struct {
	TaskDefinitionArn    string                `json:"taskDefinitionArn"`
	Family               string                `json:"family"`
	TaskRoleArn          string                `json:"taskRoleArn"`
	ExecutionRoleArn     string                `json:"executionRoleArn"`
	NetworkMode          string                `json:"networkMode"`
	ContainerDefinitions []ContainerDefinition `json:"containerDefinitions"`

	// not caring about the following props for now, but need to
	// define them to register a new task definition revision from one revision

	Volumes               []any `json:"volumes,omitempty"`
	PlacementConstraints  []any `json:"placementConstraints,omitempty"`
	RequiresCompabilities []any `json:"requiresCompabilities,omitempty"`
	Cpu                   any   `json:"cpu,omitempty"`
	Memory                any   `json:"memory,omitempty"`
	Tags                  []any `json:"tags,omitempty"`
	PidMode               any   `json:"pidMode,omitempty"`
	IpcMode               any   `json:"ipcMode,omitempty"`
	ProxyConfiguration    any   `json:"proxyConfiguration,omitempty"`
	InferenceAccelerators []any `json:"inferenceAccelerators,omitempty"`
	EphemeralStorage      any   `json:"ephemeralStorage,omitempty"`
	RuntimePlatform       any   `json:"runtimePlatform,omitempty"`
}

// "taskDefinition" argument is the family, family:revision or full ARN
func DescribeTaskDefinition(taskDefinition string) (TaskDefinition, error) {
	result := TaskDefinition{}
	var args []string
	args = append(args, "ecs", "describe-task-definition", "--output", "json", "--no-paginate", "--task-definition", taskDefinition)
	args = append(args, "--query", "{family: taskDefinition.family, taskDefinitionArn: taskDefinition.taskDefinitionArn, containerDefinitions: taskDefinition.containerDefinitions[*].{name: name, image: image, portMappings: portMappings}}")
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(2)
		result = TaskDefinition{
			TaskDefinitionArn: taskDefinition,
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

	_, err := execAWS(args, &result)
	if err != nil {
		return result, err
	}
	log.Debug(result)

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

func ExtractFamilyFromRevision(taskdefArn string) string {
	var result string = taskdefArn
	arnSuffix := ":task-definition/"
	arnSuffixPos := strings.LastIndex(result, arnSuffix)
	if arnSuffixPos > -1 {
		result = result[arnSuffixPos+len(arnSuffix):]
	}
	revisionPos := strings.LastIndex(result, ":")
	if revisionPos > -1 {
		result = result[:revisionPos]
	}
	return result
}
