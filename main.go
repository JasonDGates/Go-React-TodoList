package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty`
	Completed bool `json:"completed"`
	Body string `json:"body"`
}

var collection *mongo.Collection

func main(){
	fmt.Println("Hello World")
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading envionrment variables")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI((MONGODB_URI))
	client,err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to MongoDB")

	collection = client.Database("golang_db").Collection("todos")
	todos := []Todo{}

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	// app.Patch("/api/todos/:id", updateTodo)
	// app.Delete("/api/todos/:id", deleteTodo)

	// Create Todo
	// app.Post("/api/todos", func(c *fiber.Ctx) error {
	// 	todo := &Todo{}

	// 	if err := c.BodyParser(todo); err != nil {
	// 		return err
	// 	}

	// 	if todo.Body == "" {
	// 		return c.Status(400).JSON(fiber.Map{"error":"Todo body is required"})
	// 	}

	// 	todo.ID = len(todos) + 1
	// 	todos = append(todos, *todo)

	// 	return c.Status(201).JSON(todo)
	// })

	// Complete todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error":"Todo not found"})
	})

	// Delete Todo
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error":"Todo not found"})
	})
	log.Fatal(app.Listen(":"+PORT))
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return c.JSON(todos)
}
func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	c.BodyParser(todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error":"Todo body is required"})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)

	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo);
}