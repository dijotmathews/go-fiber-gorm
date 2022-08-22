package routes

import (
	"errors"

	"github.com/dijotmathews/go-fiber-gorm/database"
	"github.com/dijotmathews/go-fiber-gorm/models"
	"github.com/gofiber/fiber/v2"
)

// User ..
type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// CreateResponseUser ...
func CreateResponseUser(userModel models.User) User {
	return User{ID: userModel.ID, FirstName: userModel.FirstName, LastName: userModel.LastName}
}

// CreateUser ...
func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(&user)
	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)
}

// GetUsers ...
func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}

	database.Database.Db.Find(&users)

	responseUsers := []User{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)

}

func findUser(id int, user *models.User) error {
	database.Database.Db.Find(&user, "id=?", id)
	if user.ID == 0 {
		return errors.New("user does not exist")
	}

	return nil
}

// GetUser ...
func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User
	if err != nil {
		return c.Status(400).JSON("Please ensure the is is an integer")
	}

	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)

}

// UpdateUser ...
func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Please ensure the is is an integer")

	}
	var user models.User
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	var updateData UpdateUser

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if updateData.FirstName != "" {
		user.FirstName = updateData.FirstName
	}

	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}

	database.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)
	return c.Status(200).JSON(responseUser)

}

// DeleteUser ...
func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Please ensure the is is an integer")

	}
	var user models.User
	if err := findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&user).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("user deleted")

}
