package main

import (
	"cc_DavidGayle_BackendAPI/internal/app"
	"cc_DavidGayle_BackendAPI/internal/app/common"
	"cc_DavidGayle_BackendAPI/internal/app/core"
	"cc_DavidGayle_BackendAPI/internal/app/model"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	if len(os.Args[1:]) != 1 {
		usage()
		return
	}

	var config *model.Config
	var err error
	if config, err = getConfig(os.Args[1]); err != nil {
		fmt.Printf("Encountered error while reading config file: %v", err)
	}

	// Setup DB
	if err = common.NewDb(config).InitialMigration(); err != nil {
		log.Fatal("Error creating users table", err)
	}

	svcHandlers := core.NewUserHandler(config)
	svcRouter := app.NewSvcRouter(svcHandlers)
	svcRouter.Start()
	fmt.Println("Yes, it runs")
}

func getConfig(cfgFile string) (*model.Config, error) {
	viper.SetConfigType("json")
	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config model.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func usage() {
	fmt.Println("Path and config file name must be passed in as arguments")
	fmt.Println("    example:    <executable> <file>")
}
