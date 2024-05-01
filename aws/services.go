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
