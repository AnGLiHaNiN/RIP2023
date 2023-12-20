package app

import (
	"fmt"
	"net/http"

	"lab3/internal/app/ds"
	"lab3/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

func (app *Application) GetAllMedicines(c *gin.Context) {
	var request schemes.GetAllMedicinesRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	medicines, err := app.repo.GetMedicineByName(request.Name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	draftComponent, err := app.repo.GetDraftComponent(app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response := schemes.GetAllMedicinesResponse{DraftComponent: nil, Medicines: medicines}
	if draftComponent != nil {
		response.DraftComponent = &schemes.ComponentShort{UUID: draftComponent.UUID}
		containers, err := app.repo.GetComponentContent(draftComponent.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		response.DraftComponent.MedicineCount = len(containers)
	}
	c.JSON(http.StatusOK, response)
}

func (app *Application) GetMedicine(c *gin.Context) {
	var request schemes.MedicineRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recipient, err := app.repo.GetMedicineByID(request.MedicineId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if recipient == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекраство не найдено"))
		return
	}
	c.JSON(http.StatusOK, recipient)
}

func (app *Application) DeleteMedicine(c *gin.Context) {
	var request schemes.MedicineRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recipient, err := app.repo.GetMedicineByID(request.MedicineId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if recipient == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекраство не найдено"))
		return
	}
	recipient.ImageURL = nil
	recipient.IsDeleted = true
	if err := app.repo.SaveMedicine(recipient); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (app *Application) AddMedicine(c *gin.Context) {
	var request schemes.AddMedicineRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recipient := ds.Medicine(request.Medicine)
	if err := app.repo.AddMedicine(&recipient); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if request.Image != nil {
		imageURL, err := app.uploadImage(c, request.Image, recipient.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		recipient.ImageURL = imageURL
	}
	if err := app.repo.SaveMedicine(&recipient); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (app *Application) ChangeMedicine(c *gin.Context) {
	var request schemes.ChangeMedicineRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recipient, err := app.repo.GetMedicineByID(request.MedicineId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if recipient == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекраство не найдено"))
		return
	}

	if request.Name != nil {
		recipient.Name = *request.Name
	}
	if request.Image != nil {
		if recipient.ImageURL != nil {
			if err := app.deleteImage(c, recipient.UUID); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		imageURL, err := app.uploadImage(c, request.Image, recipient.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		recipient.ImageURL = imageURL
	}
	if request.Dosage != nil {
		recipient.Dosage = *request.Dosage
	}
	if request.Amount != nil {
		recipient.Amount = *request.Amount
	}
	if request.Manufacturer != nil {
		recipient.Manufacturer = *request.Manufacturer
	}

	if err := app.repo.SaveMedicine(recipient); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, recipient)
}

func (app *Application) AddToComponent(c *gin.Context) {
	var request schemes.AddToComponentRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var err error

	recipient, err := app.repo.GetMedicineByID(request.MedicineId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if recipient == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекраство не найдено"))
		return
	}

	var component *ds.Component
	component, err = app.repo.GetDraftComponent(app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		component, err = app.repo.CreateDraftComponent(app.getCustomer())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err = app.repo.AddToComponent(component.UUID, request.MedicineId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var medicines []ds.Medicine
	medicines, err = app.repo.GetComponentContent(component.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllMedicinesResponse{Medicines: medicines})
}