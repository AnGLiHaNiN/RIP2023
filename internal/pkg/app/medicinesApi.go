package app

import (
	"fmt"
	"net/http"

	_ "R_I_P_labs/docs"
	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

// @Summary		Получить все лекарства
// @Tags		Лекарства
// @Description	Возвращает все доступные лекарства с опциональной фильтрацией по Названию
// @Produce		json
// @Param		name query string false "Название для фильтрации"
// @Success		200 {object} schemes.GetAllMedicinesResponse
// @Router		/api/medicines [get]
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
	var draftComponent *ds.Component = nil
	if userId, exists := c.Get("userId"); exists {
		draftComponent, err = app.repo.GetDraftComponent(userId.(string))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	response := schemes.GetAllMedicinesResponse{DraftComponent: nil, Medicines: medicines}
	if draftComponent != nil {
		response.DraftComponent = &schemes.ComponentShort{UUID: draftComponent.UUID}
		medicinesCount, err := app.repo.CountMedicines(draftComponent.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		response.DraftComponent.MedicineCount = int(medicinesCount)
	}
	c.JSON(http.StatusOK, response)
}

// @Summary		Получить одно лекарство
// @Tags		Лекарства
// @Description	Возвращает более подробную информацию об одном лекарстве
// @Produce		json
// @Param		medicine_id path string true "id лекарства"
// @Success		200 {object} ds.Medicine
// @Router		/api/medicines/{medicine_id} [get]
func (app *Application) GetMedicine(c *gin.Context) {
	var request schemes.MedicineRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	medicine, err := app.repo.GetMedicineByID(request.MedicineId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}
	c.JSON(http.StatusOK, medicine)
}

// @Summary		Удалить лекарство
// @Tags		Лекарства
// @Description	Удаляет лекраство по id
// @Param		medicine_id path string true "id лекарства"
// @Success		200
// @Router		/api/medicines/{medicine_id} [delete]
func (app *Application) DeleteMedicine(c *gin.Context) {
	var request schemes.MedicineRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	medicine, err := app.repo.GetMedicineByID(request.MedicineId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}
	medicine.ImageURL = nil
	medicine.IsDeleted = true
	if err := app.repo.SaveMedicine(medicine); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Добавить лекарство
// @Tags		Лекарства
// @Description	Добавить новое лекарство
// @Accept		mpfd
// @Param     	image formData file false "Изображение лекарства"
// @Param     	name formData string true "Название" format:"string" maxLength:100
// @Param     	manufacturer formData string true "Производитель" format:"string" maxLength:100
// @Param     	amount formData int true "Количество" format:"int"
// @Param     	dosage formData string true "Дозировка" format:"string" maxLength:100
// @Success		200
// @Router		/api/medicines/ [post]
func (app *Application) AddMedicine(c *gin.Context) {
	var request schemes.AddMedicineRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	medicine := ds.Medicine(request.Medicine)
	if err := app.repo.AddMedicine(&medicine); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if request.Image != nil {
		imageURL, err := app.uploadImage(c, request.Image, medicine.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		medicine.ImageURL = imageURL
	}
	if err := app.repo.SaveMedicine(&medicine); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary		Изменить лекарство
// @Tags		Лекарства
// @Description	Изменить данные полей о лекарстве
// @Accept		mpfd
// @Produce		json
// @Param		medicine_id path string true "Идентификатор лекарства" format:"uuid"
// @Param		name formData string false "Название" format:"string" maxLength:100
// @Param		manufacturer formData string false "Производитель" format:"string" maxLength:100
// @Param		amount formData int false "Количество" format:"int"
// @Param		image formData file false "Изображение лекарства"
// @Param		dosage formData string false "Дозировка" format:"string" maxLength:100
// @Router		/api/medicines/{medicine_id} [put]
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

	medicine, err := app.repo.GetMedicineByID(request.MedicineId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}

	if request.Name != nil {
		medicine.Name = *request.Name
	}
	if request.Image != nil {
		if medicine.ImageURL != nil {
			if err := app.deleteImage(c, medicine.UUID); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		imageURL, err := app.uploadImage(c, request.Image, medicine.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		medicine.ImageURL = imageURL
	}
	if request.Manufacturer != nil {
		medicine.Manufacturer = *request.Manufacturer
	}
	if request.Amount != nil {
		medicine.Amount = *request.Amount
	}
	if request.Dosage != nil {
		medicine.Dosage = *request.Dosage
	}

	if err := app.repo.SaveMedicine(medicine); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, medicine)
}

// @Summary		Добавить в компонент
// @Tags		Лекарства
// @Description	Добавить выбранное лекарство в черновик копмпонента
// @Produce		json
// @Param		medicine_id path string true "id лекарства"
// @Success		200 {object} schemes.AddToComponentResp
// @Router		/api/medicines/{medicine_id}/add_to_component [post]
func (app *Application) AddToComponent(c *gin.Context) {
	var request schemes.AddToComponentRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var err error

	medicine, err := app.repo.GetMedicineByID(request.MedicineId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}

	var component *ds.Component
	userId := getUserId(c)
	component, err = app.repo.GetDraftComponent(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		component, err = app.repo.CreateDraftComponent(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err = app.repo.AddToComponent(component.UUID, request.MedicineId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	medicinesCount, err := app.repo.CountMedicines(component.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AddToComponentResp{MedicinesCount: medicinesCount})
}