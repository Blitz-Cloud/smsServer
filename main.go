package main

import (
	"log"

	dbTypes "github.com/Blitz-Cloud/smsServer/db"
	"github.com/davecgh/go-spew/spew"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// app := fiber.New()
	db, err := gorm.Open(sqlite.Open("test.db"),&gorm.Config{})
	db.AutoMigrate(&dbTypes.User{},&dbTypes.Message{},&dbTypes.Node{})
	if err!=nil{
		log.Fatal(err)
	}
	
	// newUser:= dbTypes.User{
	// 	Name: "Ionut",
	// 	Email: "ionut@blitzcloud.me",
	// 	Password: "test",
	// }

	// newNode := dbTypes.Node{
	// 	UserId: 1,
	// 	Name: "rico",
	// 	MacAddress: "AAAA:BBBB:CCCC:DDDD",
	// }

	// db.Create(&newUser)
	// // spew.Dump(result)
	// db.Create(&newNode)
	

	foundUser := dbTypes.User{}
	db.Preload("Nodes").First(&foundUser, 1)
	// db.First(&foundUser,1)
	spew.Dump(&foundUser)





	

	// db.Create(&newUser)


	// app.Get("/", func(c *fiber.Ctx) error {

	// 	return c.SendString("Hello, World!")
	// })

	// app.Listen(":3000")
}