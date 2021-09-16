package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/nazeemnato/employee-go/db"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {

	client := db.NewClient()

	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	fields := []string{"fullName", "username", "email", "password"}
	errors, totalError := bodyChecker(&data, fields)

	if totalError > 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Fields missig  errors",
			"errors":  errors,
			"total":   totalError,
		})
	}

	// check username already taken
	usernameExist, _ := client.User.FindUnique(
		db.User.Username.Equals(data["username"]),
	).Exec(ctx)

	if usernameExist != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Username already taken",
		})
	}

	// check email already taken
	emailExist, _ := client.User.FindUnique(
		db.User.Email.Equals(data["email"]),
	).Exec(ctx)

	if emailExist != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Email already taken",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)

	user, err := client.User.CreateOne(
		db.User.FullName.Set(data["fullName"]),
		db.User.Username.Set(data["username"]),
		db.User.Email.Set(data["email"]),
		db.User.Password.Set(string(password)),
	).Exec(ctx)

	if err != nil {
		fmt.Println(err)
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Something went wrong",
		})
	}

	payload := jwt.StandardClaims{
		Subject:   user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("sce"))

	if err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Something went wrong",
		})
	}

	cookie := fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"status": true,
		"token":  token,
	})
}

func Login(c *fiber.Ctx) error {

	client := db.NewClient()

	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	fields := []string{"username", "password"}
	errors, totalError := bodyChecker(&data, fields)

	if totalError > 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Fields missig  errors",
			"errors":  errors,
			"total":   totalError,
		})
	}

	// check username exist
	user, _ := client.User.FindUnique(
		db.User.Username.Equals(data["username"]),
	).Exec(ctx)

	if user == nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid username",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		c.Status(403)
		return c.JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	payload := jwt.StandardClaims{
		Subject:   user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("sce"))

	if err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Something went wrong",
		})
	}

	cookie := fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"status": true,
		"token":  token,
	})
}

func User(c *fiber.Ctx) error {
	client := db.NewClient()

	id := c.GetRespHeader("x-userId")

	fmt.Println("id from request", id)

	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	
	user, _ := client.User.FindUnique(
		db.User.ID.Equals(id),
	).Exec(ctx)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		HTTPOnly: true,
		Expires:  time.Now().Add(-time.Hour),
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "Bye",
	})
}
