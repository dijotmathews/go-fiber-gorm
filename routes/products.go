package routes

import (
	"errors"

	"github.com/dijotmathews/go-fiber-gorm/database"
	"github.com/dijotmathews/go-fiber-gorm/models"
	"github.com/gofiber/fiber/v2"
)

//Product ...
type Product struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

func CreateResponseProduct(productModel models.Product) Product {
	return Product{ID: productModel.ID, Name: productModel.Name, SerialNumber: productModel.SerialNumber}
}

// CreateProduct ...
func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	database.Database.Db.Create(&product)
	responseProduct := CreateResponseProduct(product)

	return c.Status(fiber.StatusOK).JSON(responseProduct)
}

// GetProducts ...
func GetProducts(c *fiber.Ctx) error {
	var products []models.Product

	database.Database.Db.Find(&products)

	var responseProducts []Product

	for _, product := range products {
		responseProduct := CreateResponseProduct(product)
		responseProducts = append(responseProducts, responseProduct)
	}

	return c.Status(fiber.StatusOK).JSON(responseProducts)
}

func findProduct(id int, product *models.Product) error {
	database.Database.Db.Find(&product, "id=?", id)
	if product.ID == 0 {
		return errors.New("product does not exist")
	}

	return nil
}

// GetProduct ...
func GetProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("the id: is required")

	}
	var product models.Product

	if err = findProduct(id, &product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}
	responseProduct := CreateResponseProduct(product)

	return c.Status(fiber.StatusOK).JSON(responseProduct)
}

// UpdateProduct ...
func UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("the :id is required and needs to be an integer")

	}

	var product models.Product
	if err = findProduct(id, &product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	type UpdateProduct struct {
		Name         string `json:"name"`
		SerialNumber string `json:"serial_number"`
	}

	var updateData UpdateProduct

	if err = c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if updateData.Name != "" {
		product.Name = updateData.Name
	}

	if updateData.SerialNumber != "" {
		product.SerialNumber = updateData.SerialNumber
	}

	database.Database.Db.Save(product)

	responseProduct := CreateResponseProduct(product)
	return c.Status(fiber.StatusOK).JSON(responseProduct)

}

// DeleteProduct ...
func DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("the :id is required and needs to be an integer")

	}

	var product models.Product
	if err = findProduct(id, &product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&product).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString("product deleted")
}
