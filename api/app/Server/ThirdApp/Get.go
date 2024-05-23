package ThirdApp

import (
	"fmt"
	"simple_account_oauth/app/AccountStruct"
	"simple_account_oauth/app/Error"
	"simple_account_oauth/app/Server/Author"
	"simple_account_oauth/app/Server/Message"
	"strings"
)

func Get(context *Message.Context) (string, int, error) {
	author, err := Author.Get(context)
	if err != nil {
		fmt.Println(err)
	}
	if author == nil {
		return "", Error.NOT_LOGGED_IN, nil
	}

	appId := context.Url.Shift()

	connect := context.Database.Connect()
	include := "id,callback,permissions,salt,hash"
	query := "SELECT " + include + " FROM ThirdApp WHERE id=?"
	rows, err := connect.Query(query, appId)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	if !rows.Next() {
		return "", Error.APP_NOT_EXIST, nil
	}

	var app AccountStruct.ThirdApp
	columns := app.MappingTable(include)
	err = rows.Scan(columns...)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	verificationCode, err := context.Auths.Add(*author, app, app.Permissions)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	url := strings.Replace(app.Callback, "{code}", verificationCode, -1)

	return url, Error.NULL, nil
}
