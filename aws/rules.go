package aws

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func CreateRule(filepath string, targetGroupArn string) (string, error) {
	var args []string
	args = append(args, "elbv2", "create-rule", "--cli-input-json", fmt.Sprintf("file://%s", filepath))
	if targetGroupArn != "" {
		args = append(args, "--action", fmt.Sprintf("Type=forward,TargetGroupArn=%s", targetGroupArn))
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
