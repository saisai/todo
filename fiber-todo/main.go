package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var db *sqlx.DB
var redisClient *redis.Client
var ctx = context.Background()

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

	// Try fetching from Redis cache
	todosJSON, err := redisClient.Get(ctx, "todos").Result()
	fmt.Println("err", err)
	if err == nil {
		var todos []Todo
		json.Unmarshal([]byte(todosJSON), &todos)
		fmt.Println("📌 Data served from Redis cache")
		return c.JSON(todos)
	}

	var todos []Todo
	err = db.Select(&todos, "SELECT * FROM todos")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Store in Redis for caching
	todosBytes, _ := json.Marshal(todos)
	redisClient.Set(ctx, "todos", string(todosBytes), 10*time.Minute)

	fmt.Println("📌 Data fetched from PostgreSQL and cached")
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

	// Clear Redis cache (force refresh)
	log.Println("clear from cache in createTodo")
	redisClient.Del(ctx, "todos")

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

	// Clear Redis cache (force refresh)
	log.Println("clear from cache in updateTodo")
	redisClient.Del(ctx, "todos")

	return c.JSON(todo)
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Clear Redis cache (force refresh)
	log.Println("clear from cache in deleteTodo")
	redisClient.Del(ctx, "todos")

	return c.SendStatus(204) // 204 No Content (successful delete)
}

func main() {
	initDB()
	app := fiber.New()
	// app.Use(cors.New())

	// Connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: "192.168.2.39:6379",
	})
	fmt.Println("✅ Connected to Redis")

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
