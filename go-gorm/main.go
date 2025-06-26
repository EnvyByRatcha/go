package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

func authRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtSecretKey := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&Book{}, &User{})
	fmt.Println("Database migration completed!")

	app := fiber.New()
	app.Use("/books", authRequired)

	app.Post("/register", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		err := createUser(db, user)
		if err != nil {
			return c.SendStatus(fiber.ErrBadRequest.Code)
		}

		return c.JSON(fiber.Map{
			"message": "Register successful",
		})
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		token, err := loginUser(db, user)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 72),
			HTTPOnly: true,
		})

		return c.JSON(fiber.Map{
			"message": "Login successful",
			"token":   token,
		})
	})

	app.Get("/books", func(c *fiber.Ctx) error {
		books, err := getBooks(db)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.JSON(books)
	})
	app.Get("/books/:id", func(c *fiber.Ctx) error {
		bookId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		book, err := getBook(db, uint(bookId))
		if err != nil {
			return c.SendStatus(fiber.ErrBadRequest.Code)
		}

		return c.JSON(book)
	})
	app.Post("/books", func(c *fiber.Ctx) error {
		book := new(Book)
		if err := c.BodyParser(book); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		err := createBook(db, book)
		if err != nil {
			return c.SendStatus(fiber.ErrBadRequest.Code)
		}

		return c.JSON(fiber.Map{
			"message": "Create Book successful",
		})
	})
	app.Put("/books/:id", func(c *fiber.Ctx) error {
		bookId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		bookUpdate := new(Book)
		if err := c.BodyParser(bookUpdate); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		bookUpdate.ID = uint(bookId)

		err = updateBook(db, bookUpdate)
		if err != nil {
			return c.SendStatus(fiber.ErrBadRequest.Code)
		}

		return c.JSON(fiber.Map{
			"message": "Update Book successful",
		})
	})
	app.Delete("/books/:id", func(c *fiber.Ctx) error {
		bookId, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		err = deleteBook(db, uint(bookId))
		if err != nil {
			return c.SendStatus(fiber.ErrBadRequest.Code)
		}

		return c.JSON(fiber.Map{
			"message": "Delete Book successful",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}
