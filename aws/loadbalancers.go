package aws

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type LoadBalancer struct {
	LoadBalancerName string `json:"LoadBalancerName"`
	Type             string `json:"Type"`
	LoadBalancerArn  string `json:"LoadBalancerArn"`
}

func CreateLoadBalancer(filepath string) (LoadBalancer, error) {
	var args []string
	args = append(args, "elbv2", "create-load-balancer", "--cli-input-json", fmt.Sprintf("file://%s", filepath), "--output", "json")
	args = append(args, "--query", "LoadBalancers[0].{LoadBalancerName: LoadBalancerName, Type: Type, LoadBalancerArn: LoadBalancerArn}")
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(1)
		return LoadBalancer{
			Type:             "application",
			LoadBalancerArn:  "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-load-balancer/50dc6c495c0c9188",
			LoadBalancerName: "my-load-balancer",
		}, nil
	}

	var resp LoadBalancer
	_, err := execAWS(args, &resp)

	return resp, err
}
