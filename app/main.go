package main

import (
	"log"
	"fmt"
	"encoding/json"

	"github.com/gmkvaal/wtd/app/shared/database"
	"github.com/gmkvaal/wtd/app/shared/server"
	"github.com/gmkvaal/wtd/app/controller"
	"github.com/gmkvaal/wtd/app/shared/readjson"
	//"github.com/gmkvaal/wtd/app/shared/database"

)



func main() {

	readjson.Load("config.json", config)
	database.Connect(config.Database.Postgres)

	if database.CheckConnection() == true {
		log.Println("Successfully connected to DB")
	} else {
		log.Println("Unable to connect to DB")
	}

	fmt.Println(config)

	fmt.Println(server.HttpAddress(config.Server))
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

var config = &configuration{}

type configuration struct {
	Database database.Database
	Server server.Server
	OAuth controller.OAuthConfig
}




