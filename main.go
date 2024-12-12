package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
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

	r := gin.Default()

	if os.Getenv("ENV") != "production" {
		r.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
			c.Next()
		})
	}

	r.GET("/api/todos", getTodos)
	r.POST("/api/todos", createTodo)
	r.GET("/api/todos/:id", getTodo)
	r.PATCH("/api/todos/:id", updateTodo)
	r.DELETE("/api/todos/:id", deleteTodo)

	if os.Getenv("ENV") == "production" {
		fmt.Println("Running in production mode")
		r.Static("/", "./client/dist")
		r.StaticFS("/assets", http.Dir("./client/dist/assets"))
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "80"
	}
	log.Fatal(r.Run(":" + PORT))
}

func getTodos(c *gin.Context) {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context) {
	var todo Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	if todo.Body == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Todo body cannot be empty"})
		return
	}

	result, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusOK, todo)
}

func getTodo(c *gin.Context) {
	id := c.Param("id")
	var todo Todo
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func updateTodo(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	filter := bson.M{"_id": objectID}

	var todo Todo
	err = collection.FindOne(context.Background(), filter).Decode(&todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	update := bson.M{"$set": bson.M{"completed": !todo.Completed}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	todo.Completed = !todo.Completed
	c.JSON(http.StatusOK, todo)
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deletedCount": result.DeletedCount})
}
