package Message

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"simple_account_oauth/app/Auths"
	"simple_account_oauth/app/Database"
	"simple_account_oauth/app/Logger"
	"simple_account_oauth/app/Url"
	"strconv"
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

func GetUser(context *Context) int {
	cookie, err := context.Request.Cookie("secret")
	if err != nil {
		return 0
	}
	secret := cookie.Value
	token := context.Request.Header.Get("Authorization")
	if secret == "" || token == "" {
		return 0
	}

	data := Data{
		Secret: secret,
		Token:  token,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return 0
	}

	body := bytes.NewBuffer(dataBytes)

	res, err := http.Post(context.ExtensionChannelSource+"/get_user", "application/json", body)
	if err != nil {
		return 0
	}
	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return 0
	}

	userId, err := strconv.ParseInt(string(resBytes), 10, 0)
	if err != nil {
		return 0
	}

	return int(userId)
}
