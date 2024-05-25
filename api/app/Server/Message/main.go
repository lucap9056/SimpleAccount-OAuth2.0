package Message

import (
	"net/http"
	"simple_account_oauth/app/Auths"
	"simple_account_oauth/app/Database"
	"simple_account_oauth/app/Logger"
	"simple_account_oauth/app/Url"
)

type Response struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
	Error   int    `json:"error"`
}

type Context struct {
	ExtensionChannelSource string
	AppOwnershipLimits     int
	AllowCreateApps        bool
	Writer                 http.ResponseWriter
	Request                *http.Request
	Url                    *Url.Url
	Database               *Database.API
	Logger                 *Logger.Manager
	Auths                  *Auths.Manager
}

type Data struct {
	Secret string `json:"secret"`
	Token  string `json:"token"`
}
