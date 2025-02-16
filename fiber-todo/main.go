package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var db *sqlx.DB
var redisClient *redis.Client
var ctx = context.Background()

func initDB() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the PostgreSQL connection URL from environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	var err error
	db, err = sqlx.Connect("postgres", databaseURL)
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

	// Query the last record

	query := "SELECT * FROM todos ORDER BY id DESC LIMIT 1"
	row := db.QueryRow(query)

	// Scan result into the Todo struct
	err = row.Scan(&todo.ID, &todo.Title, &todo.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows returned")
		} else {
			log.Fatal(err)
		}
	}

	print("hello ", todo)
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

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the PostgreSQL connection URL from environment variables
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL is required")
	}

	// Connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisURL,
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

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port if not specified
	}

	fmt.Printf("Starting server on port %s...\n", port)
	app.Listen(":" + port)
}
