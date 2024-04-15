package main

import (
	"github.com/cloudflare/workers-sdk-go/kv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	r := gin.Default()
	r.GET("/todos", getTodos)
	r.POST("/todos", createTodo)
	r.PUT("/todos/:id", updateTodo)
	r.DELETE("/todos/:id", deleteTodo)
	r.Run()
}

func getTodos(c *gin.Context) {
	// Récupérer les todos depuis Cloudflare Workers KV
	todos, err := kv.Get("todos")
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error getting todos",
		})
		return
	}
	c.JSON(200, todos)
}

func createTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.BindJSON(&newTodo); err != nil {
		return
	}
	newTodo.ID = uuid.New().String()

	// Enregistrer le nouveau todo dans Cloudflare Workers KV
	err := kv.Set("todos", newTodo)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error creating todo",
		})
		return
	}
	c.JSON(201, newTodo)
}

func updateTodo(c *gin.Context) {
	id := c.Param("id")
	var updatedTodo Todo
	if err := c.BindJSON(&updatedTodo); err != nil {
		return
	}

	// Mettre à jour le todo dans Cloudflare Workers KV
	err := kv.Set(id, updatedTodo)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error updating todo",
		})
		return
	}
	c.JSON(200, updatedTodo)
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")

	// Supprimer le todo dans Cloudflare Workers KV
	err := kv.Delete(id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error deleting todo",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Todo deleted",
	})
}
