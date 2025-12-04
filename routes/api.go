package routes

import (
	"fmt"
	"log"

	"github.com/Blitz-Cloud/smsServer/db"
	"github.com/Blitz-Cloud/smsServer/types"
	"github.com/Blitz-Cloud/smsServer/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterApiRoutes(app *fiber.App) {
	apiG := app.Group("/api")

	apiG.Get("/status", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	apiG.Post("/signup/user", func(c *fiber.Ctx) error {
		DB, ok := c.Locals("db").(*gorm.DB)
		if !ok {

			log.Fatal("Failed to access db")
		}
		user := types.User{}
		foundUser := db.User{}
		// spew.Dump(c.BodyRaw())

		err := c.BodyParser(&user)
		if err != nil {
			log.Fatal(err)
		}
		result := DB.First(&foundUser, "email=?", user.Email)
		if result.Error == nil {
			// utilizator existent
			// urmand sa seted si un cookie pentru autentificare
			token, err := utils.GenerateToken(&foundUser)
			if err != nil {
				log.Fatal(err)
			}

			return c.SendString(token)
		}
		foundUser.Name = user.Name
		foundUser.Email = user.Email
		foundUser.Password, err = utils.HashPassword(user.Password)
		if err != nil {
			log.Fatal("Failed to hash password")
		}
		// foundUser.UUID = uuid.New().String()

		DB.Save(&foundUser)

		token, err := utils.GenerateToken(&foundUser)
		if err != nil {
			log.Fatal(err)
		}

		return c.SendString(token)
	})

	apiG.Post("/signup/node", func(c *fiber.Ctx) error {
		DB, ok := c.Locals("db").(*gorm.DB)
		if !ok {
			log.Fatal("Failed to access db")
		}
		nodeData := types.Node{}
		err := c.BodyParser(&nodeData)
		if err != nil {
			log.Fatal(err)
		}
		foundUser := db.User{}
		foundNode := db.Node{}
		result := DB.Preload("Nodes").Where("email = ? ",nodeData.Email).Where("id = ?",nodeData.UserId).First(&foundUser)
		fmt.Println("{POST:")
		spew.Dump(nodeData)
		fmt.Println("{DB}:")
		spew.Dump(foundUser)
		if result.Error == gorm.ErrRecordNotFound {
			// utilizatorul nu exista
			return c.Status(fiber.StatusUnauthorized).SendString("This user doesn't exists")
		}
		nodeFound := false
		for _, node := range foundUser.Nodes {
			if node.MacAddress == nodeData.MacAddress {
				nodeFound = true
				foundNode = node;
				break
			}
		}
		spew.Dump(nodeFound)
		
		if nodeFound {
			// 
			return c.SendString("Exists\nPlease follow this link to authorize and associate the Node with your account:\nhttp://localhost:3000/authorize/node?id=" + fmt.Sprintf("%d",int(foundNode.ID)))
		}
		newNode := db.Node{}
		newNode.Status = "Pending approval"
		newNode.UserID = uint(nodeData.UserId)
		newNode.MacAddress = nodeData.MacAddress
		result = DB.Save(&newNode)
		if result.Error !=nil{
			spew.Dump(result.Error)
		}
		foundUser.Nodes = append(foundUser.Nodes,newNode )
		result = DB.Save(&foundUser)
		if result.Error !=nil{
			spew.Dump(result.Error)
		}
		return c.SendString("Please follow this link to authorize and associate the Node with your account:\nhttp://localhost:3000/authorize/node?id=" + fmt.Sprintf("%d",int(newNode.ID)))
	})
	//
	// apiG.Use(middleware.RouteProtector)

}
