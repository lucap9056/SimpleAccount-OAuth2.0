package Server

import (
	"context"
	"fmt"
	"net/http"
	"simple_account_oauth/app/Auths"
	"simple_account_oauth/app/Database"
	"simple_account_oauth/app/Logger"
	"simple_account_oauth/app/Server/Informations"
	"simple_account_oauth/app/Server/Message"
	"simple_account_oauth/app/Server/ThirdApp"
	"simple_account_oauth/app/Url"
)

type Config struct {
	ExtensionChannelSource string
	Port                   int
	AllowCreateApps        bool
	AppOwnershipLimits     int
}

type API struct {
	Config Config
	Server http.Server
	DB     *Database.API
	Auths  *Auths.Manager
	Logger *Logger.Manager
}

func New(config Config, db *Database.API, auths *Auths.Manager, logger *Logger.Manager) (*API, error) {
	api := API{
		DB:     db,
		Auths:  auths,
		Config: config,
		Logger: logger,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/oauth_api/", api.MainHandler)

	api.Server = http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
	}

	return &api, nil
}

func (api *API) MainHandler(w http.ResponseWriter, r *http.Request) {
	url := Url.New(r.URL)
	url.Shift()
	context := &Message.Context{
		Writer:                 w,
		Request:                r,
		Database:               api.DB,
		ExtensionChannelSource: api.Config.ExtensionChannelSource,
		AllowCreateApps:        api.Config.AllowCreateApps,
		AppOwnershipLimits:     api.Config.AppOwnershipLimits,
		Logger:                 api.Logger,
		Auths:                  api.Auths,
		Url:                    &url,
	}

	switch url.Shift() {
	case "info":
		Informations.Handler(context)
	case "app":
		ThirdApp.Handler(context)
	}
}

func (api *API) Start() {
	err := api.Server.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func (api *API) Exit(ctx context.Context) error {
	return api.Server.Shutdown(ctx)
}
