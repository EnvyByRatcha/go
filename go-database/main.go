package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

var db *sql.DB

type Product struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	sdb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	db = sdb

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	app := fiber.New()

	app.Get("/products", getProductsHandler)
	app.Get("/products/:id", getProductHandler)
	app.Post("/products", createProductHandler)
	app.Put("/products/:id", updateProductHandler)
	app.Delete("products/:id", deleteProductHandler)

	app.Listen(":8080")

}

func getProductsHandler(c *fiber.Ctx) error {
	product, err := getProducts()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(product)
}

func getProductHandler(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	product, err := getProduct(productId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(product)
}

func createProductHandler(c *fiber.Ctx) error {
	p := new(Product)
	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := createProduct(p)
	if err != nil {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	return c.JSON(p)
}

func updateProductHandler(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	productUpdate := new(Product)
	if err := c.BodyParser(productUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	p, err := updateProduct(productId, productUpdate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(p)
}

func deleteProductHandler(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = deleteProduct(productId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
