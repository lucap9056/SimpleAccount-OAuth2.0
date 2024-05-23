package Informations

import (
	"simple_account_oauth/app/Error"
	"simple_account_oauth/app/Server/Author"
	"simple_account_oauth/app/Server/Message"
)

func DeleteApp(context *Message.Context) (string, int, error) {
	appId := context.Url.Shift()
	author, err := Author.Get(context)
	if err != nil {
		return "", Error.SYSTEM, err
	}
	if author == nil {
		return "", Error.NOT_LOGGED_IN, nil
	}

	connect := context.Database.Connect()
	result, err := connect.Exec("DELETE FROM ThirdApp WHERE id=? AND client=?", appId, author.Id)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return "", Error.SYSTEM, err
	}

	return "", Error.NULL, nil
}
