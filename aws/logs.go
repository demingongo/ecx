package aws

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func CreateLogGroup(logGroupName string) (string, error) {
	var args []string
	args = append(args, "logs", "create-log-group", "--log-group-name", logGroupName)
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(1)
		return strings.Join(args, " "), nil
	}

	var resp any
	stdout, err := execAWS(args, &resp)

	return string(stdout), err
}

func PutRetentionPolicy(logGroupName string, retentionInDays int) (string, error) {
	var args []string
	args = append(args, "logs", "put-retention-policy", "--log-group-name", logGroupName, "--retention-in-days", strconv.Itoa(retentionInDays))
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(1)
		return strings.Join(args, " "), nil
	}

	var resp any
	stdout, err := execAWS(args, &resp)

	return string(stdout), err
}
