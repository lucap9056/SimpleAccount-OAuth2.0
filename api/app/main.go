package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"simple_account_oauth/app/Auths"
	"simple_account_oauth/app/Database"
	"simple_account_oauth/app/Logger"
	"simple_account_oauth/app/Server"
	"syscall"
	"time"
)

type Config struct {
	MainAPISource      string `json:"main_API_source"`
	Port               int    `json:"port"`
	AllowCreateApps    bool   `json:"allow_creating_new_apps"`
	AppOwnershipLimits int    `json:"app_ownership_limits"`
	LogsPath           string `json:"logs_dir_path"`
	Database           struct {
		SourceName string `json:"source_name"`
	} `json:"database"`
	Auths struct {
		AuthorizationCodeDuration int `json:"authorization_code_duration"`
	} `json:"auths"`
}

func main() {
	configBytes, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	logger, err := Logger.New(config.LogsPath)
	if err != nil {
		panic(err)
	}

	db, err := Database.New(Database.Config{
		SourceName: config.Database.SourceName,
	})
	if err != nil {
		panic(err)
	}

	auths := Auths.New(Auths.Config{
		AuthorizationCodeDuration: time.Duration(config.Auths.AuthorizationCodeDuration) * time.Minute,
	})

	server_config := Server.Config{
		ExtensionChannelSource: config.MainAPISource,
		Port:                   config.Port,
		AllowCreateApps:        config.AllowCreateApps,
		AppOwnershipLimits:     config.AppOwnershipLimits,
	}
	server, err := Server.New(server_config, db, auths, logger)
	if err != nil {
		panic(err)
	}

	go server.Start()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer server.Exit(ctx)

	fmt.Println("START")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
