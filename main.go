package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello world")

	if os.Getenv("ENV") != "production" {
		// Load .env file if not in production
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
			return
		}
	}

	
	MONGODB_URI := os.Getenv("MONGODB_URI")
	fmt.Println(MONGODB_URI)
	clientOption := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		log.Fatal("Error connecting to MongoDB")
		return
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
		return
	}
	fmt.Println("Connected to MongoDB")

	defer client.Disconnect(context.Background())

	collection = client.Database("Cluster0").Collection("todos")

	app := fiber.New()
	
	if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	} else {
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept",
		}))
	}
	
	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Get("/api/todos/:id", getTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)
	
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}
	log.Fatal(app.Listen("0.0.0.0:" + PORT))
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}
	defer cursor.Close(context.Background())
	
	if err = cursor.All(context.Background(), &todos); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}
	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err,
		})
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Todo body cannot be empty"})
	}

	result, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)
	return c.JSON(todo)
}

func getTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	var todo Todo
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&todo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}
	return c.JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid ID"})
	}
	filter := bson.M{"_id": objectID}

	var todo Todo
	err = collection.FindOne(context.Background(), filter).Decode(&todo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err})
	}

	update := bson.M{"$set": bson.M{"completed": !todo.Completed}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	todo.Completed = !todo.Completed
	return c.Status(200).JSON(todo)

}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid ID"})
	}

	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"deletedCount": result.DeletedCount})
}
