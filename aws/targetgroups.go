package aws

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type TargetGroup struct {
	TargetGroupArn  string
	TargetGroupName string
}

type describeTargetGroupsOutput struct {
	TargetGroups []TargetGroup
}

func DescribeTargetGroups() ([]TargetGroup, error) {
	result := []TargetGroup{}
	var args []string
	args = append(args, "elbv2", "describe-target-groups", "--output", "json", "--no-paginate")
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(2)
		for i := 1; i <= 10; i += 1 {
			arn := "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/dummy-" + strconv.Itoa(i) + "/73e2d6bc24d8a067"
			name := "dummy-" + strconv.Itoa(i)
			result = append(result, TargetGroup{TargetGroupArn: arn, TargetGroupName: name})
		}
		return result, nil
	}

	var resp describeTargetGroupsOutput
	_, err := execAWS(args, &resp)
	if err != nil {
		return result, err
	}

	result = resp.TargetGroups

	return result, nil
}

func CreateTargetGroup(filepath string) (TargetGroup, error) {
	var result TargetGroup
	var args []string
	args = append(args, "elbv2", "create-target-group", "--output", "json", "--cli-input-json", fmt.Sprintf("file://%s", filepath))
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(1)
		return TargetGroup{
			TargetGroupArn:  "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/dummy-created/73e2d6bc24d8a067",
			TargetGroupName: "dummy-created",
		}, nil
	}

	var resp describeTargetGroupsOutput
	_, err := execAWS(args, &resp)
	if err != nil {
		return result, err
	}

	if len(resp.TargetGroups) > 0 {
		result = resp.TargetGroups[0]
	}

	return result, nil
}
