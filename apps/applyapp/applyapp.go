package applyapp

import (
	"fmt"
	"os"

	"github.com/demingongo/ecx/globals"
	"gopkg.in/yaml.v2"
)

type LogGroup struct {
	Group     string `yaml:"group"`
	Retention int64  `yaml:"retention"`
}

type Flow struct {
	Name        string   `yaml:"name"`
	Service     string   `yaml:"service"`
	TargetGroup string   `yaml:"targetGroup"`
	Rules       []string `yaml:"rules"`
}

type conf struct {
	Api             string     `yaml:"api"`
	ApiVersion      string     `yaml:"apiVersion"`
	LogGroups       []LogGroup `yaml:"logGroups"`
	TaskDefinitions []string   `yaml:"taskDefinitions"`
	Flows           []Flow     `yaml:"flows"`
}

func (c *conf) getConf() *conf {

	yamlFile, err := os.ReadFile("../../ecx-tests/project1/ecx.yaml")
	if err != nil {
		globals.Logger.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		globals.Logger.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func Run() {
	fmt.Println("apply ecx")

	var c conf
	c.getConf()

	fmt.Println(c)
}
