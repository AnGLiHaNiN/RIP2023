package app

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"R_I_P_labs/internal/app/config"
	"R_I_P_labs/internal/app/dsn"
	"R_I_P_labs/internal/app/redis"
	"R_I_P_labs/internal/app/repository"
	"R_I_P_labs/internal/app/role"


	"github.com/swaggo/files"      
	"github.com/swaggo/gin-swagger" 
	_ "R_I_P_labs/docs"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
	redisClient *redis.Client
}

func (app *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.Use(ErrorHandler())

	// Услуги - лекарства
	api := r.Group("/api")
	{
		components := api.Group("/components")
		{
			components.GET("/", app.WithAuthCheck(role.NotAuthorized, role.Customer, role.Moderator), app.GetAllComponents)                     // Список с поиском
			components.GET("/:component_id", app.WithAuthCheck(role.NotAuthorized, role.Customer, role.Moderator), app.GetComponent)            // Одна услуга
			components.DELETE("/:component_id", app.WithAuthCheck(role.Moderator), app.DeleteComponent)                         				// Удаление
			components.PUT("/:component_id", app.WithAuthCheck(role.Moderator), app.ChangeComponent)                            				// Изменение
			components.POST("/", app.WithAuthCheck(role.Moderator), app.AddComponent)                                           				// Добавление
			components.POST("/:component_id/add_to_medicine", app.WithAuthCheck(role.Customer,role.Moderator), app.AddToMedicine) 					// Добавление в заявку
		}

		// Заявки - компоненты
		medicines := api.Group("/medicines")
		{
			medicines.GET("/", app.WithAuthCheck(role.Customer, role.Moderator), app.GetAllMedicines)                                         				  // Список (отфильтровать по дате формирования и статусу)
			medicines.GET("/:medicine_id",app.WithAuthCheck(role.Customer, role.Moderator),  app.GetMedicine)                             				  // Одна заявка
			medicines.PUT("/:medicine_id/update", app.WithAuthCheck(role.Customer, role.Moderator), app.UpdateMedicine)                                	  // Изменение (добавление транспорта)
			medicines.DELETE("/:medicine_id", app.WithAuthCheck(role.Customer,role.Moderator), app.DeleteMedicine)                                      				  // Удаление
			medicines.DELETE("/:medicine_id/delete_component/:component_id", app.WithAuthCheck(role.Customer, role.Moderator), app.DeleteFromMedicine) 	  // Изменеие (удаление услуг)
			medicines.PUT("/user_confirm", app.WithAuthCheck(role.Customer, role.Moderator), app.UserConfirm)                                    				  // Сформировать создателем
			medicines.PUT("/:medicine_id/moderator_confirm", app.WithAuthCheck(role.Moderator), app.ModeratorConfirm)                         				  // Завершить отклонить модератором
		}

		// Пользователи (авторизация)
		user := api.Group("/user")
		{
			user.POST("/sign_up", app.Register)
			user.POST("/login", app.Login)
			user.POST("/logout", app.Logout)
		}

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		r.Run(fmt.Sprintf("%s:%d", app.config.ServiceHost, app.config.ServicePort))

		log.Println("Server down")
	}
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

	app.minioClient, err = minio.New(app.config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	app.redisClient, err = redis.New(app.config.Redis)
	if err != nil {
		return nil, err
	}

	return &app, nil
}