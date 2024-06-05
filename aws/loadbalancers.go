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

func DescribeLoadBalancersWithNames(names []string) ([]LoadBalancer, error) {
	result := []LoadBalancer{}
	var args []string
	args = append(args, "elbv2", "describe-load-balancers", "--output", "json", "--no-paginate")
	args = append(args, "--query", "LoadBalancers[*].{LoadBalancerName: LoadBalancerName, Type: Type, LoadBalancerArn: LoadBalancerArn}")
	if len(names) > 0 {
		args = append(args, "--load-balancer-names")
		args = append(args, names...)
	}
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(2)
		if len(names) > 0 {
			name := names[0]
			arn := "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/" + name + "/50dc6c495c0c9188"
			result = append(result, LoadBalancer{LoadBalancerArn: arn, LoadBalancerName: name, Type: "application"})
		}
		return result, nil
	}

	_, err := execAWS(args, &result)

	return result, err
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
			LoadBalancerArn:  "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/dummy-load-balancer/50dc6c495c0c9188",
			LoadBalancerName: "my-load-balancer",
		}, nil
	}

	var resp LoadBalancer
	_, err := execAWS(args, &resp)

	return resp, err
}
