package aws

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

type Image struct {
	ImageDigest string `json:"imageDigest"`
	ImageTag    string `json:"imageTag"`
}

type ListImagesOutput struct {
	ImageIds []Image `json:"imageIds"`
}

func ListImages(ecrRepositoryName string) ([]Image, error) {
	var result []Image
	var args []string
	args = append(args, "ecr", "list-images", "--output", "json", "--repository-name", ecrRepositoryName, "--no-paginate", "--filter", "tagStatus=TAGGED")
	log.Debug(args)
	if viper.GetBool("dummy") {
		sleep(1)
		for i := 1; i <= 8; i += 1 {
			ImageDigest := "sha256:b5a2c96250612366ea272ffac6d9744aaf4b45aacd96aa7cfcb931ee3b558259"
			ImageTag := "dummy1.13." + strconv.Itoa(2+i)
			result = append(result, Image{ImageDigest, ImageTag})
		}
	} else {
		var resp ListImagesOutput
		_, err := execAWS(args, &resp)
		if err != nil {
			return result, err
		}

		result = resp.ImageIds
	}

	// reversa array
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

func ExtractNameFromURI(ecrRepositoryUri string) string {
	var result string
	result1 := strings.Index(ecrRepositoryUri, ".ecr.")
	result2 := strings.Index(ecrRepositoryUri, ".amazonaws.")
	result3 := strings.Index(ecrRepositoryUri, "/")
	if result1 > 0 && result1 < result2 && result2 < result3 {
		result = ecrRepositoryUri[result3+1:]
		imageTagIdx := strings.Index(result, ":")
		if imageTagIdx > -1 {
			result = result[:imageTagIdx]
		}
	}
	return result
}
