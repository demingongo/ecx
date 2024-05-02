package updateserviceapp

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/demingongo/ecx/aws"
	"github.com/spf13/viper"
)

type Config struct {
	cluster string
	service aws.Service

	// the new task definition
	taskDefinition aws.TaskDefinition
}

func (m Config) CurrentTaskDefinitionArn() string {
	var result string
	if len(m.service.Deployments) > 0 {
		result = m.service.Deployments[0].TaskDefinition
	}
	return result
}

func (m Config) CurrentTaskDefinitionFamily() string {
	return aws.ExtractFamilyFromRevision(m.CurrentTaskDefinitionArn())
}

var (
	config Config
)

func Run() {

	config.cluster = viper.GetString("cluster")
	config.service = aws.Service{
		ServiceArn: viper.GetString("service"),
	}

	log.Info(fmt.Sprintf("cluster: %s", config.cluster))
	log.Info(fmt.Sprintf("service: %s", config.service.ServiceArn))

	log.Info(fmt.Sprintf("ECR repo name: %s", aws.ExtractNameFromURI("xxx.dkr.ecr.us-west-2.amazonaws.com/repo/dummy")))

	log.Info(fmt.Sprintf("Task def family: %s", aws.ExtractFamilyFromRevision("arn:aws:ecs:us-east-1:053534965804:task-definition/webserver:5")))

	if config.service.ServiceArn != "" {
		var err error
		config.service, err = aws.DescribeService(config.cluster, config.service.ServiceArn)
		if err != nil {
			log.Fatal("DescribeService", err)
		}
	}

	// retrieve the last revision
	if config.CurrentTaskDefinitionFamily() != "" {
		var err error
		config.taskDefinition, err = aws.DescribeTaskDefinition(config.CurrentTaskDefinitionFamily())
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Info(config.taskDefinition.ContainerDefinitions)

	dir, _ := os.Getwd()
	jsonByte, _ := json.Marshal(config.taskDefinition)
	fmt.Println(string(jsonByte))

	f, err := os.Create(dir + "/task_def.json")
	if err != nil {
		log.Fatal("os.Create", err)
	}
	defer f.Close()
	_, err = f.Write(jsonByte)
	if err != nil {
		log.Fatal("os.File.WriteString", err)
	}

	fmt.Println("Done")
}
