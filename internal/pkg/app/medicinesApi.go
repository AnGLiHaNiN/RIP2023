package app

import (
	"fmt"
	"net/http"
	"time"

	"lab3/internal/app/ds"
	"lab3/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

func (app *Application) GetAllMedicines(c *gin.Context) {
	var request schemes.GetAllMedicinesRequst
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	medicines, err := app.repo.GetAllMedicines(request.FormationDateStart, request.FormationDateEnd, request.Status)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	outputMedicines := make([]schemes.MedicineOutput, len(medicines))
	for i, medicine := range medicines {
		outputMedicines[i] = schemes.ConvertMedicine(&medicine)
	}
	c.JSON(http.StatusOK, schemes.AllMedicinesResponse{Medicines: outputMedicines})
}

func (app *Application) GetMedicine(c *gin.Context) {
	var request schemes.MedicineRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	medicine, err := app.repo.GetMedicineById(request.MedicineId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}

	components, err := app.repo.GetMedicineProduction(request.MedicineId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.MedicineResponse{Medicine: schemes.ConvertMedicine(medicine), Components: components})
}

func (app *Application) UpdateMedicine(c *gin.Context) {
	var request schemes.UpdateMedicineRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	medicine, err := app.repo.GetMedicineById(request.URI.MedicineId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}
	medicine.MedicineName = request.MedicineName
	if app.repo.SaveMedicine(medicine); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.UpdateMedicineResponse{Medicine:schemes.ConvertMedicine(medicine)})
}

func (app *Application) DeleteMedicine(c *gin.Context) {
	var request schemes.MedicineRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	medicine, err := app.repo.GetMedicineById(request.MedicineId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("увдомление не найдено"))
		return
	}
	medicine.Status = ds.DELETED

	if err := app.repo.SaveMedicine(medicine); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) DeleteFromMedicine(c *gin.Context) {
	var request schemes.DeleteFromMedicineRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	medicine, err := app.repo.GetMedicineById(request.MedicineId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}
	if medicine.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя редактировать лекарство со статусом: %s", medicine.Status))
		return
	}

	if err := app.repo.DeleteFromMedicine(request.MedicineId, request.ComponentId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	components, err := app.repo.GetMedicineProduction(request.MedicineId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllComponentsResponse{Components: components})
}

func (app *Application) UserConfirm(c *gin.Context) {
	var request schemes.UserConfirmRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	medicine, err := app.repo.GetMedicineById(request.URI.MedicineId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}
	if medicine.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя сформировать лекарство со статусом %s", medicine.Status))
		return
	}
	if request.Confirm {
		medicine.Status = ds.FORMED
		now := time.Now()
		medicine.FormationDate = &now
	} else {
		medicine.Status = ds.DELETED
	}

	if err := app.repo.SaveMedicine(medicine); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) ModeratorConfirm(c *gin.Context) {
	var request schemes.ModeratorConfirmRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	

	medicine, err := app.repo.GetMedicineById(request.URI.MedicineId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}
	if medicine.Status != ds.FORMED {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%s\" на \"%s\"", medicine.Status,  ds.FORMED))
		return
	}
	if request.Confirm {
		medicine.Status = ds.COMPELTED
		now := time.Now()
		medicine.CompletionDate = &now
	
	} else {
		medicine.Status = ds.REJECTED
	}
	medicine.ModeratorId = app.getModerator()
	
	if err := app.repo.SaveMedicine(medicine); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
