package aws

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type Listener struct {
	ListenerArn string `json:"ListenerArn"`
}

func CreateListener(filepath string, loadBalancerArn string, targetGroupArn string) (Listener, error) {
	var args []string
	args = append(args, "elbv2", "create-listener", "--cli-input-json", filepath, "--output", "json")
	args = append(args, "--query", "Listeners[0].{ListenerArn: ListenerArn}")
	if loadBalancerArn != "" {
		args = append(args, "--load-balancer-arn", loadBalancerArn)
	}
	if targetGroupArn != "" {
		args = append(args, "--default-actions", fmt.Sprintf("Type=forward,TargetGroupArn=%s", targetGroupArn))
	}
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(1)
		return Listener{
			ListenerArn: "arn:aws:elasticloadbalancing:us-east-1:850631746142:listener/gwy/my-agw-lb-example2/e0f9b3d5c7f7d3d6/afc127db15f925de",
		}, nil
	}

	var resp Listener
	_, err := execAWS(args, &resp)

	return resp, err
}
