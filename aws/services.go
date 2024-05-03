package aws

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type ServiceLoadBalancer struct {
	TargetGroupArn string
	ContainerName  string
	ContainerPort  int
}

type Deployment struct {
	Id             string `json:"id"`
	TaskDefinition string `json:"taskDefinition"`
}

type Service struct {
	ServiceArn  string       `json:"serviceArn"`
	ServiceName string       `json:"serviceName"`
	Deployments []Deployment `json:"deployments"`
}

func CreateService(filepath string, loadBalancer ServiceLoadBalancer) (string, error) {
	var args []string
	args = append(args, "ecs", "create-service", "--output", "json", "--cli-input-json", fmt.Sprintf("file://%s", filepath))
	if loadBalancer.TargetGroupArn != "" && loadBalancer.ContainerName != "" {
		args = append(args, "--load-balancers", fmt.Sprintf(
			"targetGroupArn=%s,containerName=%s,containerPort=%d",
			loadBalancer.TargetGroupArn, loadBalancer.ContainerName, loadBalancer.ContainerPort,
		))
	}
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(1)
		return strings.Join(args, " "), nil
	}

	var resp any
	stdout, err := execAWS(args, &resp)

	return string(stdout), err
}

func DescribeService(cluster string, serviceArn string) (Service, error) {
	var result Service
	var args []string
	args = append(args, "ecs", "describe-services", "--output", "json", "--cluster", cluster, "--no-paginate", "--services", serviceArn)
	args = append(args, "--query", "services[0].{serviceArn: serviceArn, serviceName: serviceName, deployments: deployments[*].{id: id, taskDefinition: taskDefinition}}")

	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(1)
		return Service{
			ServiceArn:  serviceArn,
			ServiceName: "dummy-service",
			Deployments: []Deployment{
				{
					Id:             "ecs-svc/1234567890123456789",
					TaskDefinition: "arn:aws:ecs:us-east-1:053534965804:task-definition/dummy:5",
				},
			},
		}, nil
	}

	_, err := execAWS(args, &result)

	return result, err
}

func DescribeServices(cluster string, serviceArn string) ([]Service, error) {
	var result []Service
	var args []string
	args = append(args, "ecs", "describe-services", "--output", "json", "--cluster", cluster, "--no-paginate")
	if serviceArn != "" {
		args = append(args, "--services", serviceArn)
	}
	args = append(args, "--query", "services[*].{serviceArn: serviceArn, serviceName: serviceName, deployments: deployments[*].{id: id, taskDefinition: taskDefinition}}")

	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(1)
		return []Service{
			{
				ServiceArn:  "arn:aws:ecs:us-west-2:123456789012:service/dummy-service",
				ServiceName: "dummy-service",
				Deployments: []Deployment{
					{
						Id:             "ecs-svc/1234567890123456789",
						TaskDefinition: "arn:aws:ecs:us-east-1:053534965804:task-definition/dummy:5",
					},
				},
			},
			{
				ServiceArn:  "arn:aws:ecs:us-west-2:123456789012:service/dummy-service-2",
				ServiceName: "dummy-service-2",
				Deployments: []Deployment{
					{
						Id:             "ecs-svc/9876543210987654321",
						TaskDefinition: "arn:aws:ecs:us-east-1:053534965804:task-definition/dummy2:18",
					},
				},
			},
		}, nil
	}

	_, err := execAWS(args, &result)

	return result, err
}
