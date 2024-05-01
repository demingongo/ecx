package updateserviceapp

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func Run() {
	log.Info(fmt.Sprintf("cluster: %s", viper.GetString("cluster")))
	log.Info(fmt.Sprintf("service: %s", viper.GetString("service")))
}
