package Informations

import (
	"encoding/json"
	"simple_account_oauth/app/AccountStruct"
	"simple_account_oauth/app/Error"
	"simple_account_oauth/app/Server/Author"
	"simple_account_oauth/app/Server/Message"
)

func GetAppInformation(context *Message.Context) (string, int, error) {
	name := context.Url.Shift()

	if name == "" {
		return GetOwnApps(context)
	} else {
		return GetApp(context, name)
	}
}

func GetOwnApps(context *Message.Context) (string, int, error) {
	author, err := Author.Get(context)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	if author == nil {
		return "", Error.NOT_LOGGED_IN, err
	}

	include := "id,name,callback,description,permissions"
	connect := context.Database.Connect()
	rows, err := connect.Query("SELECT "+include+" FROM ThirdApp WHERE client=?", author.Id)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	apps := []AccountStruct.ThirdApp{}
	for rows.Next() {
		var app AccountStruct.ThirdApp
		columns := app.MappingTable(include)
		err = rows.Scan(columns...)
		if err != nil {
			return "", Error.SYSTEM, err
		}
		apps = append(apps, app)
	}

	bytes, err := json.Marshal(apps)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	return string(bytes), Error.NULL, nil
}

func GetApp(context *Message.Context, id string) (string, int, error) {
	include := "name,callback,description,permissions"
	connect := context.Database.Connect()
	rows, err := connect.Query("SELECT "+include+" FROM ThirdApp WHERE id=? LIMIT 0,1", id)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	if rows.Next() {
		app := AccountStruct.ThirdApp{}
		columns := app.MappingTable(include)
		err = rows.Scan(columns...)
		if err != nil {
			return "", Error.SYSTEM, err
		}

		bytes, err := json.Marshal(app)
		if err != nil {
			return "", Error.SYSTEM, err
		}

		return string(bytes), Error.NULL, nil
	}

	return "", Error.APP_NOT_EXIST, nil
}
