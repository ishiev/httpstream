package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Config конфигурация сервиса
type Config struct {
	addr    string
	isDebug bool
	dataDir string
}

var config Config

func main() {
	log.Printf("HTTP Stream Store Service -- простой сервис хранения потоков данных c RESTful интерфейсом.\n")

	// Флаги и их разбор
	flag.StringVar(&config.addr, "addr", "0.0.0.0:8080", "Service listennig address & port")
	flag.StringVar(&config.dataDir, "path", "DATA", "Streams storage path")
	flag.BoolVar(&config.isDebug, "d", false, "Debug mode")
	flag.Parse()

	// Создание роутера
	if config.isDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	//
	// HTTP API
	//

	// Проверка сервиса на работоспособность /ping
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "pong",
		})
	})

	// API /v1
	v1 := router.Group("/v1")
	{
		// Сохранить новый поток и вернуть его идентификатор
		v1.POST("/streams", ginSaveStream)

		// Сохранить новый поток и вернуть его идентификатор (PUT)
		v1.PUT("/streams", ginSaveStream)

		// Вернуть поток по идентификатору
		v1.GET("/streams/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.File(GetStreamPath(id))
		})

		// Удалить поток по идентификатору
		v1.DELETE("/streams/:id", func(c *gin.Context) {
			id := c.Param("id")
			err := DeleteStream(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"status": "error",
					"error":  err,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status": "deleted",
					"id":     id,
				})
			}
		})
	}

	// Обработчик для возврата 404 в формате JSON
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"error":  http.StatusNotFound,
		})
	})

	//
	// Запуск в работу
	//

	// Печать конфигурации
	log.Printf("Current config:\n")
	log.Printf(" - Streams path:      %s\n", config.dataDir)
	log.Printf(" - Listening address: %s\n", config.addr)

	// Cоздание пути хранения (если нужно)
	err := os.MkdirAll(config.dataDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Fatal: error creating %s: %v", config.dataDir, err)
	}
	// Запуск сервера
	log.Fatal(router.Run(config.addr))
}

// Обработчик сохранения HTTP запроса в поток
func ginSaveStream(c *gin.Context) {
	id, err := SaveStream(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err,
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"status": "created",
			"id":     id,
		})
	}
}
