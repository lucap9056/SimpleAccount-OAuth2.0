package Informations

import (
	"encoding/json"
	"io"
	"simple_account_oauth/app/AccountStruct"
	"simple_account_oauth/app/Auths"
	"simple_account_oauth/app/Error"
	"simple_account_oauth/app/Server/Author"
	"simple_account_oauth/app/Server/Message"
	"strings"
	"unicode/utf8"
)

func SetAppInformation(context *Message.Context) (string, int, error) {
	if !context.AllowCreateApps {
		return "", Error.NOT_ALLOWED_TO_CREATE_APPS, nil
	}

	author, err := Author.Get(context)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	if author == nil {
		return "", Error.NOT_LOGGED_IN, nil
	}

	body, err := io.ReadAll(context.Request.Body)
	if err != nil {
		return "", Error.CLIENT_INVALID_REQUEST, nil
	}

	app := AccountStruct.ThirdApp{}
	err = json.Unmarshal(body, &app)
	if err != nil {
		return "", Error.CLIENT_INVALID_REQUEST, nil
	}

	if app.Name == "" {
		return "", Error.APP_NAME_IS_EMPTY, nil
	}

	if utf8.RuneCountInString(app.Name) > 255 {
		return "", Error.APP_NAME_TOO_LONG, nil
	}

	if app.Callback == "" {
		return "", Error.APP_CALLBACK_IS_EMPTY, nil
	} else {

		if strings.Index(app.Callback, "http://") == 0 {
			app.Callback = "https://" + app.Callback[7:]
		}

		if strings.Index(app.Callback, "https://") != 0 {
			app.Callback = "https://" + app.Callback
		}
	}

	if utf8.RuneCountInString(app.Callback) > 255 {
		return "", Error.CALLBACK_TOO_LONG, nil
	}

	if len(app.Description) > 256 {
		return "", Error.APP_DESCRIPTION_TOO_LONG, nil
	}

	connect := context.Database.Connect()
	rows, err := connect.Query("SELECT id FROM ThirdApp WHERE client=?", author.Id)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	ownedApps := 0
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return "", Error.SYSTEM, nil
		}
		if id == app.Id {
			return Update(context, app)
		}
		ownedApps++
		if ownedApps >= context.AppOwnershipLimits {
			return "", Error.OWNED_APPS_REACHED_LIMIT, nil
		}
	}

	return Create(context, app, author)
}

func Create(context *Message.Context, app AccountStruct.ThirdApp, author *AccountStruct.User) (string, int, error) {
	id, err := context.Database.UUID()
	if err != nil {
		return "", Error.SYSTEM, err
	}

	secret, salt, hash := Auths.GenerateSecret()

	connect := context.Database.Connect()

	insertQuery := `INSERT INTO ThirdApp(id,name,client,salt,hash,callback,description,permissions)
	VALUE(?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := connect.Prepare(insertQuery)
	if err != nil {
		return "", Error.SYSTEM, err
	}
	_, err = stmt.Exec(id, app.Name, author.Id, salt, hash, app.Callback, app.Description, app.Permissions)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	newApp := AccountStruct.ThirdApp{
		Id:     id,
		Secret: secret,
	}

	result, err := json.Marshal(newApp)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	return string(result), Error.NULL, nil
}

func Update(context *Message.Context, app AccountStruct.ThirdApp) (string, int, error) {
	secret, salt, hash := Auths.GenerateSecret()
	connect := context.Database.Connect()
	_, err := connect.Exec(`UPDATE ThirdApp SET 
	name=?, salt=?, hash=?, callback=?, description=?, permissions=? 
	WHERE id=?`,
		app.Name, salt, hash, app.Callback, app.Description, app.Permissions, app.Id)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	newApp := AccountStruct.ThirdApp{
		Id:     app.Id,
		Secret: secret,
	}

	result, err := json.Marshal(newApp)
	if err != nil {
		return "", Error.SYSTEM, err
	}

	return string(result), Error.NULL, nil
}
