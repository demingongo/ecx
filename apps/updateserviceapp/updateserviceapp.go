package updateserviceapp

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/demingongo/ecx/aws"
	"github.com/spf13/viper"
)

func Run() {
	log.Info(fmt.Sprintf("cluster: %s", viper.GetString("cluster")))
	log.Info(fmt.Sprintf("service: %s", viper.GetString("service")))

	log.Info(fmt.Sprintf("ECR repo name: %s", aws.ExtractNameFromURI("xxx.dkr.ecr.us-west-2.amazonaws.com/repo/dummy")))

}
