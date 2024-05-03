package updateserviceapp

import (
	"encoding/json"
	"fmt"

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

var (
	config Config
)

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

func removeJSONKey(taskDefinition aws.TaskDefinition, key string) ([]byte, error) {
	// marshal to []byte
	var jsonByte []byte
	var err error
	var output map[string]interface{}
	if jsonByte, err = json.Marshal(taskDefinition); err != nil {
		log.Error("json.Marshal(taskDefinition)")
		return jsonByte, err
	}
	// unmarshal to map[string]interface{}
	if err := json.Unmarshal(jsonByte, &output); err != nil {
		log.Error("json.Unmarshal(jsonByte, &output)")
		return jsonByte, err
	}
	// remove key
	delete(output, key)
	// marshal the updated map[string]interface{} to []byte
	if jsonByte, err = json.Marshal(output); err != nil {
		log.Error("json.Marshal(output)")
		return jsonByte, err
	}

	return jsonByte, err
}

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

	log.Info(fmt.Sprintf("TaskDefinitionArn: %s", config.taskDefinition.TaskDefinitionArn))

	// marshal to []byte
	var jsonByte []byte
	var err error
	if jsonByte, err = removeJSONKey(config.taskDefinition, "taskDefinitionArn"); err != nil {
		log.Fatal("removeJSONKey", err)
	}

	_, err = aws.RegisterTaskDefinition(string(jsonByte))
	if err != nil {
		log.Fatal("RegisterTaskDefinition", err)
	}

	fmt.Println("Done")
}
