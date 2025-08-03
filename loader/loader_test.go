package loader_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/bxgstudio/appconfigloader/loader"
)

type MyConfig struct {
	AppHost       *string `yaml:"app_host" json:"app_host" env:"APP_HOST"`
	AppPort       *int    `yaml:"app_port" json:"app_port" env:"APP_PORT"`
	AppOnline     *bool   `yaml:"app_online" json:"app_online" env:"APP_ONLINE"`
	AppOtherParam string  `yaml:"app_other_param" json:"app_other_param" env:"APP_OTHER_PARAM"`
}

func (myConfig *MyConfig) Validate() error {
	if myConfig.AppHost == nil {
		return errors.New("field 'app_host' should be provided in configuration file")
	}

	if myConfig.AppPort == nil {
		return errors.New("field 'app_port' should be provided in configuration file")
	}

	if myConfig.AppOnline == nil {
		return errors.New("field 'app_online' should be provided in configuration file")
	}

	if myConfig.AppOtherParam == "" {
		return errors.New("field 'app_other_param' should not be empty")
	}

	return nil
}

func TestConfigLoader(t *testing.T) {

	// Create configuration
	myConfig := &MyConfig{}

	// Simulate envvar overriding config file & evaluate envvar mecanism
	os.Setenv("APP_HOST", "localhost")
	os.Setenv("APP_OTHER_PARAM", "${APP_HOST}_replica_0")

	// Load application config file in configuration
	err := loader.LoadConfig("./config.yaml", myConfig)
	if err != nil {
		fmt.Println("error while loading configuration file: ", err.Error())
		os.Unsetenv("APP_HOST")
		os.Unsetenv("APP_OTHER_PARAM")
		os.Exit(1)
	}
	os.Unsetenv("APP_HOST")
	os.Unsetenv("APP_OTHER_PARAM")

	// Validate configuration
	if err := myConfig.Validate(); err != nil {
		fmt.Println("validation failed: ", err.Error())
		os.Exit(1)
	} else {
		fmt.Println("validation succeed")
	}

	// Check that env override is successfull
	if *myConfig.AppHost == "localhost" {
		fmt.Println("override succeed")
	} else {
		fmt.Println("override failed")
	}

	// Enjoy your configuration
}
