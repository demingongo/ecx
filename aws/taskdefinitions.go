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
	Protocol           string `json:"protocol"`
	Name               string `json:"name"`
	HostPort           int    `json:"hostPort"`
	AppProtocol        string `json:"appProtocol,omitempty"`
	ContainerPortRange string `json:"containerPortRange,omitempty"`
}

type ContainerDefinition struct {
	Name              string        `json:"name"`
	Image             string        `json:"image"`
	Cpu               int           `json:"cpu,omitempty"`
	Memory            int           `json:"memory,omitempty"`
	MemoryReservation int           `json:"memoryReservation,omitempty"`
	PortMappings      []PortMapping `json:"portMappings"`
	Essential         bool          `json:"essential"`

	// not caring about the following props for now, but need to
	// define them to register a new task definition revision from one revision

	Links                 []any `json:"links,omitempty"`
	RepositoryCredentials any   `json:"repositoryCredentials,omitempty"`
	EntryPoint            []any `json:"entryPoint,omitempty"`
	Command               []any `json:"command,omitempty"`
	Environment           []any `json:"environment,omitempty"`
	// @TODO check if it's registered
	EnvironmentFiles       []any `json:"environmentFiles,omitempty"`
	MountPoints            []any `json:"mountPoints,omitempty"`
	VolumesFrom            []any `json:"volumesFrom,omitempty"`
	LinuxParameters        any   `json:"linuxParameters,omitempty"`
	Secrets                []any `json:"secrets,omitempty"`
	DependsOn              []any `json:"dependsOn,omitempty"`
	StartTimeout           any   `json:"startTimeout,omitempty"`
	StopTimeout            any   `json:"stopTimeout,omitempty"`
	Hostname               any   `json:"hostname,omitempty"`
	User                   any   `json:"user,omitempty"`
	WorkingDirectory       any   `json:"workingDirectory,omitempty"`
	DisableNetworking      bool  `json:"disableNetworking,omitempty"`
	Privileged             bool  `json:"privileged,omitempty"`
	ReadonlyRootFilesystem bool  `json:"readonlyRootFilesystem,omitempty"`
	DnsServers             []any `json:"dnsServers,omitempty"`
	DnsSearchDomains       []any `json:"dnsSearchDomains,omitempty"`
	ExtraHosts             []any `json:"extraHosts,omitempty"`
	DockerSecurityOptions  []any `json:"dockerSecurityOptions,omitempty"`
	Interactive            bool  `json:"interactive,omitempty"`
	PseudoTerminal         bool  `json:"pseudoTerminal,omitempty"`
	DockerLabels           any   `json:"dockerLabels,omitempty"`
	// @TODO check if it's registered
	Ulimits               []any `json:"ulimits,omitempty"`
	LogConfiguration      any   `json:"logConfiguration,omitempty"`
	HealthCheck           any   `json:"healthCheck,omitempty"`
	SystemControls        []any `json:"systemControls,omitempty"`
	ResourceRequirements  []any `json:"resourceRequirements,omitempty"`
	FirelensConfiguration any   `json:"firelensConfiguration,omitempty"`
	CredentialSpecs       []any `json:"credentialSpecs,omitempty"`
}

type TaskDefinition struct {
	TaskDefinitionArn    string                `json:"taskDefinitionArn"`
	Family               string                `json:"family"`
	TaskRoleArn          string                `json:"taskRoleArn"`
	ExecutionRoleArn     string                `json:"executionRoleArn"`
	NetworkMode          string                `json:"networkMode"`
	ContainerDefinitions []ContainerDefinition `json:"containerDefinitions,omitempty"`

	// not caring about the following props for now, but need to
	// define them to register a new task definition revision from one revision

	Volumes                 []any `json:"volumes,omitempty"`
	PlacementConstraints    []any `json:"placementConstraints,omitempty"`
	RequiresCompatibilities []any `json:"requiresCompatibilities,omitempty"`
	Cpu                     any   `json:"cpu,omitempty"`
	Memory                  any   `json:"memory,omitempty"`
	Tags                    []any `json:"tags,omitempty"`
	PidMode                 any   `json:"pidMode,omitempty"`
	IpcMode                 any   `json:"ipcMode,omitempty"`
	ProxyConfiguration      any   `json:"proxyConfiguration,omitempty"`
	InferenceAccelerators   []any `json:"inferenceAccelerators,omitempty"`
	EphemeralStorage        any   `json:"ephemeralStorage,omitempty"`
	RuntimePlatform         any   `json:"runtimePlatform,omitempty"`
}

type describeTaskDefinitionOutput struct {
	TaskDefinition TaskDefinition `json:"taskDefinition,omitempty"`
	Tags           []any          `json:"tags,omitempty"`
}

// "taskDefinition" argument is the family, family:revision or full ARN
func DescribeTaskDefinition(taskDefinition string) (TaskDefinition, error) {
	result := TaskDefinition{}
	var args []string
	args = append(args, "ecs", "describe-task-definition", "--output", "json", "--no-paginate", "--include", "TAGS", "--task-definition", taskDefinition)
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(2)
		result = TaskDefinition{
			TaskDefinitionArn: taskDefinition,
			ContainerDefinitions: []ContainerDefinition{
				{
					Name:  "dmz-web",
					Image: "xxx.dkr.ecr.us-west-2.amazonaws.com/repository-dummy:tag",
					PortMappings: []PortMapping{
						{
							ContainerPort: 8080,
							HostPort:      0,
							Name:          "http",
						},
					},
				},
			},
			Family: ExtractFamilyFromRevision(taskDefinition),
		}
		return result, nil
	}

	var output describeTaskDefinitionOutput
	_, err := execAWS(args, &output)
	if err != nil {
		return result, err
	}

	result = output.TaskDefinition
	result.Tags = output.Tags

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

func RegisterTaskDefinition(inputJson string) (string, error) {
	var args []string
	args = append(args, "ecs", "register-task-definition", "--cli-input-json", inputJson)
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(1)
		return strings.Join(args, " "), nil
	}

	var resp any
	stdout, err := execAWS(args, &resp)

	return string(stdout), err
}
