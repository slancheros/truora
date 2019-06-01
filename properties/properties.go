package properties

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"truora/models"
)

func ObtainConfig() models.Config {
	var config models.Config
	if _, err := toml.DecodeFile("properties.toml", &config); err != nil {
		fmt.Println(err)
	}
	return config
}
