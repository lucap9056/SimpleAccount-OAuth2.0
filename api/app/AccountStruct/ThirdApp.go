package AccountStruct

import "strings"

type ThirdApp struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Client      int    `json:"client"`
	Salt        string `json:"-"`
	Hash        string `json:"-"`
	Secret      string `json:"secret"`
	Callback    string `json:"callback"`
	Description string `json:"description"`
	Permissions int    `json:"permissions"`
}

func (app *ThirdApp) MappingTable(args ...string) []interface{} {

	if len(args) == 1 && strings.Contains(args[0], ",") {
		args = strings.Split(args[0], ",")
	}

	for i, arg := range args {
		if !strings.Contains(arg, ".") {
			args[i] = "app." + arg
		}
	}

	tableMap := map[string]interface{}{
		"app.id":          &app.Id,
		"app.name":        &app.Name,
		"app.client":      &app.Client,
		"app.salt":        &app.Salt,
		"app.hash":        &app.Hash,
		"app.callback":    &app.Callback,
		"app.description": &app.Description,
		"app.permissions": &app.Permissions,
	}

	var columns []interface{}
	for _, arg := range args {
		refer, exist := tableMap[arg]
		if !exist {
			var i interface{}
			refer = &i
		}
		columns = append(columns, refer)
	}
	return columns
}
