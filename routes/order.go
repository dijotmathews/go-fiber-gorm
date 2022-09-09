package routes

import (
	"errors"
	"time"

	"github.com/dijotmathews/go-fiber-gorm/database"
	"github.com/dijotmathews/go-fiber-gorm/models"
	"github.com/gofiber/fiber/v2"
)

// Order ...
type Order struct {
	ID        uint    `json:"id"`
	Product   Product `json:"product"`
	User      User    `json:"user"`
	CreatedAt time.Time
}

// CreateResponseOrder ...
func CreateResponseOrder(order models.Order, user User, product Product) Order {
	return Order{ID: order.ID, User: user, Product: product, CreatedAt: order.CreatedAt}
}

// CreateOrder ...
func CreateOrder(c *fiber.Ctx) error {
	var order models.Order

	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	var user models.User

	if err := findUser(order.UserRefer, &user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	var product models.Product
	if err := findProduct(order.ProductRefer, &product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())

	}

	database.Database.Db.Create(&order)

	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(fiber.StatusOK).JSON(responseOrder)
}

// GetOrders ...
func GetOrders(c *fiber.Ctx) error {
	var orders []models.Order

	database.Database.Db.Find(&orders)

	var responseOrders []Order
	for _, order := range orders {
		var user models.User
		if err := findUser(order.UserRefer, &user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		var product models.Product
		if err := findProduct(order.ProductRefer, &product); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		responseOrder := CreateResponseOrder(order, CreateResponseUser(user), CreateResponseProduct(product))
		responseOrders = append(responseOrders, responseOrder)

	}

	return c.Status(fiber.StatusOK).JSON(responseOrders)
}

func findOrder(id int, order *models.Order) error {
	database.Database.Db.Find(&order, "id=?", id)
	if order.ID == 0 {
		return errors.New("order does not exist")
	}

	return nil
}

// GetOrder ...
func GetOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Please ensure that :id is an integer")

	}

	var order models.Order
	if err := findOrder(id, &order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())

	}

	var user models.User
	if err := findUser(order.UserRefer, &user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	var product models.Product
	if err := findProduct(order.ProductRefer, &product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	responseOrder := CreateResponseOrder(order, CreateResponseUser(user), CreateResponseProduct(product))

	return c.Status(fiber.StatusOK).JSON(responseOrder)
}
