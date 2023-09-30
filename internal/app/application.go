package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"pharmaBlend/internal/app/dsn"
	"pharmaBlend/internal/app/repository"
	"strconv"
)

// Config represents your application's configuration.
type Config struct {
	LocalHost string
	Port      string
}

// Application represents your application.
type Application struct {
	Config   *Config
	Router   *gin.Engine
	Database *gorm.DB
	Repo     *repository.Repository // Add the repository instance to the struct.
}

// New creates and initializes a new Application instance.
func New(config *Config) (*Application, error) {
	log.Println(dsn.FromEnv())
	// Initialize the Gin router.
	router := gin.Default()

	// Initialize the database connection.
	//log.Println(11111111111)
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()))
	if err != nil {
		return nil, err
	}
	repo := repository.New(db)

	// Check the database connection.

	// Return a new Application instance.
	return &Application{
		Config:   config,
		Router:   router,
		Database: db,
		Repo:     repo, // Assign the repository instance.
	}, nil
}

func (app *Application) StartServer() {
	log.Println("Server start up")
	r := gin.Default()

	r.LoadHTMLGlob("templates/*html")
	r.Static("/image", "./resources/image")
	r.Static("/css", "./templates/css")

	r.GET("/ping", func(c *gin.Context) {
		id := c.Query("id")

		if id != "" {
			log.Printf("id received %s\n", id)
			intID, err := strconv.Atoi(id)
			if err != nil {
				log.Printf("can't convert id %v", err)
				c.Error(err)
				return
			}

			// Use the repository to get a product by ID.
			ingredient, err := app.Repo.GetByID(intID)
			if err != nil {
				log.Printf("can't get product by id %v", err)
				c.Error(err)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"product_price": ingredient.Price,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		medicationName := c.Query("medication_name")

		if medicationName == "" {
			// Получите список всех ингредиентов из базы данных с помощью вашего репозитория.
			medications, err := app.Repo.GetAllIngredients()
			if err != nil {
				// Обработка ошибки
				c.String(http.StatusInternalServerError, "Internal Server Error")
				return
			}

			c.HTML(http.StatusOK, "medications.html", gin.H{
				"medicationName": medicationName,
				"medications":    medications,
			})
			return
		}

		foundMedications, err := app.Repo.SearchIngredientsByTitle(medicationName)
		if err != nil {
			// Обработка ошибки
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		c.HTML(http.StatusOK, "medications.html", gin.H{
			"medications":    foundMedications,
			"medicationName": medicationName,
		})
	})

	r.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		realID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.Abort()
		}

		// Получите информацию об ингредиенте из базы данных по его имени с помощью вашего репозитория.
		medication, err := app.Repo.GetByID(int(realID))
		if err != nil {
			// Обработка ошибки
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		log.Println(medication)
		c.HTML(http.StatusOK, "medication.html", medication)
	})

	r.POST("/:id", func(c *gin.Context) {
		id := c.Param("id")
		realID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.Abort()
		}

		// Получите информацию об ингредиенте из базы данных по его имени с помощью вашего репозитория.
		err = app.Repo.DeleteIngredientsByID(int(realID))
		if err != nil {
			c.Abort()
		}
		c.Redirect(http.StatusFound, "/")
	})

	r.Run(":" + app.Config.Port)

	log.Println("Server down")
}
