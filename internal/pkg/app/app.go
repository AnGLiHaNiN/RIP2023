package app

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"lab3/internal/app/config"
	"lab3/internal/app/dsn"
	"lab3/internal/app/repository"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
	// dsn string
}

func (app *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.Use(ErrorHandler())

	// Услуги - получатели
	r.GET("/api/components", app.GetAllComponents)                                     // Список с поиском
	r.GET("/api/components/:component_id", app.GetComponent)                           // Одна услуга
	r.DELETE("/api/components/:component_id", app.DeleteComponent)              // Удаление
	r.PUT("/api/components/:component_id", app.ChangeComponent)                 // Изменение
	r.POST("/api/components", app.AddComponent)                                    // Добавление
	r.POST("/api/components/:component_id/add_to_medicine", app.AddToMedicine) // Добавление в заявку

	// Заявки - уведомления
	r.GET("/api/medicines", app.GetAllMedicines)                                                       // Список (отфильтровать по дате формирования и статусу)
	r.GET("/api/medicines/:medicine_id", app.GetMedicine)                                          // Одна заявка
	r.PUT("/api/medicines/:medicine_id/update", app.UpdateMedicine)                                // Изменение (добавление транспорта)
	r.DELETE("/api/medicines/:medicine_id", app.DeleteMedicine)                             //Удаление
	r.DELETE("/api/medicines/:medicine_id/delete_component/:component_id", app.DeleteFromMedicine) // Изменеие (удаление услуг)
	r.PUT("/api/medicines/:medicine_id/user_confirm", app.UserConfirm)                                 // Сформировать создателем
	r.PUT("/api/medicines/:medicine_id/moderator_confirm", app.ModeratorConfirm)                        // Завершить отклонить модератором

	r.Static("/image", "./resources")
	r.Static("/css", "./static/css")
	r.Run("localhost:7000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	app.minioClient, err = minio.New(app.config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			log.Println(err.Err)
		}
		lastError := c.Errors.Last()
		if lastError != nil {
			switch c.Writer.Status() {
			case http.StatusBadRequest:
				c.JSON(-1, gin.H{"error": "wrong request"})
			case http.StatusNotFound:
				c.JSON(-1, gin.H{"error": lastError.Error()})
			case http.StatusMethodNotAllowed:
				c.JSON(-1, gin.H{"error": lastError.Error()})
			default:
				c.Status(-1)
			}
		}
	}
}
