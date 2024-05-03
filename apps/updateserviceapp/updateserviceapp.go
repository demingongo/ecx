package updateserviceapp

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/demingongo/ecx/aws"
	"github.com/demingongo/ecx/globals"
	"github.com/spf13/viper"
)

type containerUpdate struct {
	Name     string
	OldImage string
	NewImage string
}

type Config struct {
	cluster        string
	service        aws.Service
	taskDefinition aws.TaskDefinition

	containersToUpdate []containerUpdate

	taskDefinitionInfoDescription string

	serviceLogo        string
	taskDefinitionLogo string
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

func (m Config) findContainerDefinition(containerName string) *aws.ContainerDefinition {
	var result *aws.ContainerDefinition
	if len(m.taskDefinition.ContainerDefinitions) > 0 {
		for i, containerDef := range (&m.taskDefinition).ContainerDefinitions {
			if containerDef.Name == containerName {
				result = &(&m.taskDefinition).ContainerDefinitions[i]
				break
			}
		}
	}
	return result
}

func (m Config) addContainerToUpdate(containerName string, oldImage string, newImage string) {
	config.containersToUpdate = append(config.containersToUpdate, containerUpdate{
		Name:     containerName,
		OldImage: oldImage,
		NewImage: newImage,
	})
}

var (
	config Config

	info string

	subtle  = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	special = lipgloss.AdaptiveColor{Light: "230", Dark: "#010102"}

	notifText = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render

	subtleText = lipgloss.NewStyle().Foreground(subtle).Render

	// Titles.

	titleStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("7")).
			Foreground(special)

	subtitleStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle).
			Foreground(lipgloss.Color("6"))

	// Info block.

	infoStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("7")).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true).
			Width(globals.InfoWidth)
)

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

func generateInfo() string {

	var (
		serviceInfo        string
		taskDefinitionInfo string
	)

	if config.service.ServiceName != "" {
		serviceInfo = config.service.ServiceName
	} else {
		serviceInfo = config.service.ServiceArn
	}

	if config.taskDefinitionInfoDescription != "" {
		taskDefinitionInfo = config.taskDefinitionInfoDescription
	} else {
		taskDefinitionInfo = config.taskDefinition.Family
	}

	if len(serviceInfo) == 0 {
		serviceInfo = subtleText("-")
	}

	if len(taskDefinitionInfo) == 0 {
		taskDefinitionInfo = subtleText("-")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("SUMMARY"),
		subtitleStyle.Render("Cluster "),
		config.cluster,
		subtitleStyle.Render("Service "+config.serviceLogo),
		serviceInfo,
		subtitleStyle.Render("Task definition "+config.taskDefinitionLogo),
		taskDefinitionInfo,
	)

	infoWidth := globals.InfoWidth

	if len(config.containersToUpdate) > 0 {
		var containersInfo = []string{
			content,
			subtitleStyle.Render("Containers "),
		}
		for _, ctu := range config.containersToUpdate {
			containersInfo = append(containersInfo, "・"+ctu.OldImage+notifText(" » "))
			containersInfo = append(containersInfo, notifText(ctu.NewImage))
		}
		content = lipgloss.JoinVertical(lipgloss.Left,
			containersInfo...,
		)

		infoWidth = globals.InfoWidth * 2
	}

	return infoStyle.Width(infoWidth).Render(content)
}

func selectImage(containerName string, ecrRepositoryName string, images []aws.Image) {
	form := runFormSelectImage(
		fmt.Sprintf("Container \"%s\", repository \"%s\":", containerName, ecrRepositoryName),
		images,
	)
	if form.State == huh.StateCompleted {
		if image := form.Get("image").(aws.Image); image.ImageTag != "" {
			containerDefinition := config.findContainerDefinition(containerName)
			newURI := aws.ChangeImageTagFromURI(containerDefinition.Image, image.ImageTag)
			if newURI != "" {
				config.addContainerToUpdate(
					containerDefinition.Name,
					containerDefinition.Image,
					newURI,
				)
				containerDefinition.Image = newURI
			}
		}
	}
}

func inputImage(containerName string, currentImage string) {
	form := runFormInputImage(
		fmt.Sprintf("Container \"%s\":", containerName),
		currentImage,
	)
	if newURI := form.GetString("image"); form.State == huh.StateCompleted && newURI != "" {
		containerDefinition := config.findContainerDefinition(containerName)
		if containerDefinition.Name != "" {
			config.addContainerToUpdate(
				containerDefinition.Name,
				containerDefinition.Image,
				newURI,
			)
			containerDefinition.Image = newURI
		}
	}
}

func Run() {

	config.cluster = viper.GetString("cluster")
	config.service = aws.Service{
		ServiceArn: viper.GetString("service"),
	}

	log.Debug(fmt.Sprintf("cluster: %s", config.cluster))
	info = generateInfo()

	if config.service.ServiceArn != "" {
		var err error
		config.service, err = aws.DescribeService(config.cluster, config.service.ServiceArn)
		if err != nil {
			log.Fatal("DescribeService", err)
		}
	} else {
		list, err := aws.DescribeServices(config.cluster, "")
		if err != nil {
			log.Fatal("DescribeServices", err)
		}
		form := runFormService(list)
		if form.State == huh.StateCompleted && form.GetBool("confirm") {
			service := form.Get("service").(ComparableService)
			if service.ServiceArn != "" {
				for _, fullService := range list {
					if fullService.ServiceArn == service.ServiceArn {
						config.service = fullService
						break
					}
				}
			}
		}
	}

	log.Debug(fmt.Sprintf("service: %s", config.service.ServiceArn))
	info = generateInfo()

	if config.service.ServiceArn != "" {

		log.Info(fmt.Sprintf("ECR repo name: %s", aws.ExtractNameFromURI("xxx.dkr.ecr.us-west-2.amazonaws.com/repo/dummy")))

		// retrieve the last revision from aws
		if config.CurrentTaskDefinitionFamily() != "" {
			var err error
			config.taskDefinition, err = aws.DescribeTaskDefinition(config.CurrentTaskDefinitionFamily())
			if err != nil {
				log.Fatal(err)
			}

			log.Debug(fmt.Sprintf("TaskDefinitionArn: %s", config.taskDefinition.TaskDefinitionArn))
			info = generateInfo()

			// select containers to update
			var containersList []ComparableContainerDefinition
			if len(config.taskDefinition.ContainerDefinitions) > 0 {
				containersForm := runFormSelectContainers(config.taskDefinition.ContainerDefinitions)
				if containersForm.State == huh.StateCompleted && containersForm.GetBool("confirm") {
					containersList = containersForm.Get("containers").([]ComparableContainerDefinition)
					for _, container := range containersList {
						log.Info(fmt.Sprintf("you selected: %s", container.Name))
					}
				}
			}

			if len(containersList) > 0 {
				for _, container := range containersList {
					ecrRepositoryName := aws.ExtractNameFromURI(container.Image)
					if ecrRepositoryName != "" {
						images, err := aws.ListImages(ecrRepositoryName)
						if err != nil {
							log.Error(err)
						}

						if len(images) > 0 {
							// select an image
							selectImage(container.Name, ecrRepositoryName, images)
						} else {
							// write a new image
							inputImage(container.Name, container.Image)
						}
					} else {
						// write a new image
						inputImage(container.Name, container.Image)
					}
				}

				for _, c := range config.taskDefinition.ContainerDefinitions {
					log.Debug(fmt.Sprintf("image: %s", c.Image))
				}
				// @TODO
				// - if len(containersToUpdate) > 0
				// -- register a new revision for the task definition
				// -- update service with new revision
				/*
					// create new revision for task definition
					var jsonByte []byte
					if jsonByte, err = removeJSONKey(config.taskDefinition, "taskDefinitionArn"); err != nil {
						log.Fatal("removeJSONKey", err)
					}

					_, err = aws.RegisterTaskDefinition(string(jsonByte))
					if err != nil {
						log.Fatal("RegisterTaskDefinition", err)
					}
				*/
			} else {
				// @TODO force deployment
			}
		} else {
			// @TODO force redeployment
		}
	}

	info = generateInfo()
	fmt.Println(info)

	fmt.Println("Done")
}
