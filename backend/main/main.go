package main

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/diakovliev/mesap/backend/fake_database"
	"github.com/diakovliev/mesap/backend/ifaces"
	"github.com/diakovliev/mesap/backend/models"
)

var (
	user0    models.User
	user1    models.User
	database ifaces.Database
)

func init() {
	user0 = models.User{Login: "Test login 0"}
	user1 = models.User{Login: "Test login 1"}

	database = fake_database.NewDatabase()

	log.Println("initialize")
}

func main() {
	if err := database.Open(); err != nil {
		panic(err)
	}
	defer database.Close()

	users, err := database.Users()
	if err != nil {
		panic(err)
	}

	newId, err := users.Insert(user0)
	if err != nil {
		panic(err)
	}
	log.Printf("%d", newId)

	newId, err = users.Insert(user1)
	if err != nil {
		panic(err)
	}
	log.Printf("%d", newId)

	users.Each(
		func(record models.User) bool {
			writer := bytes.NewBufferString("")
			encoder := json.NewEncoder(writer)

			if err := encoder.Encode(record); err != nil {
				panic(err)
			}

			dataStr := writer.String()

			log.Printf("[ENC] str: %s", dataStr)

			var user_in models.User

			decoder := json.NewDecoder(strings.NewReader(dataStr))
			if err := decoder.Decode(&user_in); err != nil {
				panic(err)
			}

			log.Printf("[DEC] User.Id: %d, User.Login: '%s'", user_in.GetId(), user_in.Login)

			return true
		},
	)

	users.Delete(0)
	users.Delete(1)
	// err = users.Delete(100)
	// if err != nil {
	// 	panic(err)
	// }
}
