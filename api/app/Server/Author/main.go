package Author

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"simple_account_oauth/app/AccountStruct"
	"simple_account_oauth/app/Server/Message"
)

func Get(context *Message.Context) (*AccountStruct.User, error) {
	request := context.Request

	cookie, err := request.Cookie("secret")
	if err != nil {
		return nil, err
	}

	secret := cookie.Value
	token := request.Header.Get("Authorization")

	bodyJson := fmt.Sprintf(`{"secret": "%s","token":"%s"}`, secret, token)
	body := bytes.NewBuffer([]byte(bodyJson))
	resp, err := http.Post(context.ExtensionChannelSource+"/get_user", "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 讀取響應內容
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Message.Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New(response.Result)
	}

	user := &AccountStruct.User{}
	err = json.Unmarshal([]byte(response.Result), user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
