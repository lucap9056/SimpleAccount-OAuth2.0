package ThirdApp

import (
	"encoding/base64"
	"encoding/json"
	"simple_account_oauth/app/AccountStruct"
	"simple_account_oauth/app/Auths"
	"simple_account_oauth/app/Error"
	message "simple_account_oauth/app/Server/Message"
)

func Verify(context *message.Context) (string, int, error) {
	code := context.Url.Shift()
	author := context.Request.Header.Get("Authorization")

	request := context.Auths.Get(code)
	if request == nil {
		return "", Error.VERIFICATION_CODE_INVALID, nil
	}

	app := request.App

	secret, err := base64.StdEncoding.DecodeString(author)
	if err != nil {
		return "", Error.AUTHORIZATION_INVALID, nil
	}

	salt, err := base64.StdEncoding.DecodeString(app.Salt)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	if Auths.Hash(secret, salt) != app.Hash {
		return "", Error.AUTHORIZATION_INVALID, nil
	}

	context.Auths.Remove(code)

	connect := context.Database.Connect()

	include := request.GetReadPermissions()
	rows, err := connect.Query("SELECT "+include+" FROM User WHERE id=?", request.User.Id)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	if !rows.Next() {
		return "", Error.USER_NOT_EXIST, nil
	}

	var user AccountStruct.User
	columns := user.MappingTable(include)
	err = rows.Scan(columns...)
	if err != nil {
		return "", Error.SYSTEM, err
	}
	user.MoveTempToFinal()

	userBytes, err := json.Marshal(user)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	return string(userBytes), Error.NULL, nil
}
