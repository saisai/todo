package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func initDB() {
	var err error
	db, err = sqlx.Connect("postgres", "postgres://postgres:postgres@192.168.2.39:5432/todo_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

type Todo struct {
	ID        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Completed bool   `json:"completed" db:"completed"`
}

func getTodos(c *fiber.Ctx) error {

	var todos []Todo
	err := db.Select(&todos, "SELECT * FROM todos")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(todos)

}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := db.Exec("INSERT INTO todos (title, completed) VALUES ($1, $2)", todo.Title, todo.Completed)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	_, err := db.Exec("UPDATE todos SET title=$1, completed=$2 WHERE id=$3", todo.Title, todo.Completed, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(todo)
}

// func deleteTodo(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	_, err := db.Exec("DELETE FROM todos WHERE id=$1", id)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
// 	}
// 	return c.SendStatus(204)
// }

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("Received request to /debug")
	fmt.Println("id delete ", id)
	_, err := db.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204) // 204 No Content (successful delete)
}

func main() {
	initDB()
	app := fiber.New()
	// app.Use(cors.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	app.Get("/todos", getTodos)
	app.Post("/todos", createTodo)
	app.Put("/todos/:id", updateTodo)
	app.Delete("/todos/:id", deleteTodo)

	fmt.Println("Server running on http://localhost:3001")
	log.Fatal(app.Listen(":3001"))
}
