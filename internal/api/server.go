package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type Medication struct {
	ID           int
	Name         string
	Dosage       string
	Manufacturer string
	Price        float64
	ImagePath    string
}

var medications = []Medication{
	{1, "Парацетомол", "500 мг", "Производитель A", 500, "image/Парацетомол.png"},
	{2, "Ибупрофен", "200 мг", "Производитель B", 700, "image/Ибупрофен.png"},
	{3, "Аспирин", "100 мг", "Производитель C", 300, "image/Аспирин.png"},
	{4, "Амоксициллин", "100 мг", "Производитель D", 1200, "image/Амоксициллин.png"},
}

func StartServer() {
	log.Println("Server start up")
	r := gin.Default()

	r.LoadHTMLGlob("templates/*html")
	r.Static("/image", "./resources/image")
	r.Static("/css", "./templates/css")

	r.GET("/", loadMedications)
	r.GET("/:name", loadMedication)

	r.Run()

	log.Println("Server down")
}
func loadMedications(c *gin.Context) {
	medicationName := c.Query("medication_name")

	if medicationName == "" {
		c.HTML(http.StatusOK, "medications.html", gin.H{
			"medications": medications,
		})
		return
	}

	foundMedications := []Medication{}
	lowerMedicationName := strings.ToLower(medicationName)
	for i := range medications {
		if strings.Contains(strings.ToLower(medications[i].Name), lowerMedicationName) {
			foundMedications = append(foundMedications, medications[i])
		}
	}

	c.HTML(http.StatusOK, "medications.html", gin.H{
		"medications":    foundMedications,
		"medicationName": medicationName,
	})
}

func loadMedication(c *gin.Context) {
	name := c.Param("name")

	for i := range medications {
		if medications[i].Name == name {
			c.HTML(http.StatusOK, "medication.html", gin.H{
				"Name":         medications[i].Name,
				"Dosage":       medications[i].Dosage,
				"Manufacturer": medications[i].Manufacturer,
				"Price":        medications[i].Price,
				"ImagePath":    medications[i].ImagePath,
			})
			return
		}
	}
}
