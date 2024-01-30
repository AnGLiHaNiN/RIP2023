package app

import (
	"fmt"
	"net/http"
	"time"

	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/role"
	"R_I_P_labs/internal/app/schemes"

	"github.com/gin-gonic/gin"
)

// @Summary		Получить все лекарства
// @Tags		Лекарства
// @Description	Возвращает все лекарства с фильтрацией по статусу и дате формирования
// @Produce		json
// @Param		status query string false "статус лекарствоа"
// @Param		formation_date_start query string false "начальная дата формирования"
// @Param		formation_date_end query string false "конечная дата формирвания"
// @Success		200 {object} schemes.AllMedicinesResponse
// @Router		/api/medicines [get]
func (app *Application) GetAllMedicines(c *gin.Context) {
	var request schemes.GetAllMedicinesRequst
	var err error
	if err = c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	fmt.Println(userId, userRole)
	var medicines []ds.Medicine
	if userRole == role.Customer {
		medicines, err = app.repo.GetAllMedicines(&userId, request.FormationDateStart, request.FormationDateEnd, request.Status)
	} else {
		medicines, err = app.repo.GetAllMedicines(nil, request.FormationDateStart, request.FormationDateEnd, request.Status)
	}
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

// @Summary		Получить одно лекарство
// @Tags		Лекарства
// @Description	Возвращает подробную информацию о лекарстве и его названии
// @Produce		json
// @Param		medicine_id path string true "id лекарствоа"
// @Success		200 {object} schemes.MedicineResponse
// @Router		/api/medicines/{medicine_id} [get]
func (app *Application) GetMedicine(c *gin.Context) {
	var request schemes.MedicineRequest
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	var medicine *ds.Medicine
	if userRole == role.Moderator {
		medicine, err = app.repo.GetMedicineById(request.MedicineId, nil)
	} else {
		medicine, err = app.repo.GetMedicineById(request.MedicineId, &userId)
	}
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

type SwaggerUpdateMedicineRequest struct {
	Name string `json:"name"`
}

// @Summary		Указать название лекарства
// @Tags		Лекарства
// @Description	Позволяет изменить название лекарства
// @Access		json
// @Param		name body SwaggerUpdateMedicineRequest true "Название лекарствоа"
// @Success		200
// @Router		/api/medicines [put]
func (app *Application) UpdateMedicine(c *gin.Context) {
	var request schemes.UpdateMedicineRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userId := getUserId(c)
	medicine, err := app.repo.GetDraftMedicine(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}
	medicine.Name = &request.Name
	if app.repo.SaveMedicine(medicine); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Удалить лекарство
// @Tags		Лекарства
// @Description	Удаляет лекарство по id
// @Success		200
// @Router		/api/medicines [delete]
func (app *Application) DeleteMedicine(c *gin.Context) {

	userId := getUserId(c)
	var medicine *ds.Medicine
	medicine, err := app.repo.GetDraftMedicine(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}

	medicine.Status = ds.DELETED

	if err := app.repo.SaveMedicine(medicine); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Удалить компонент из лекарства
// @Tags		Лекарства
// @Description	Удалить компонент из лекарства
// @Param		component_id path string true "id компонента"
// @Success		200
// @Router		/api/medicines/delete_component/{component_id} [delete]
func (app *Application) DeleteFromMedicine(c *gin.Context) {
	var request schemes.DeleteFromMedicineRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userId := getUserId(c)
	medicine, err := app.repo.GetDraftMedicine(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}

	if err := app.repo.DeleteFromMedicine(medicine.UUID, request.ComponentId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Сформировать лекарство
// @Tags		Лекарства
// @Description	Сформировать или удалить лекарство пользователем
// @Success		200
// @Router		/api/medicines/user_confirm [put]
func (app *Application) UserConfirm(c *gin.Context) {
	userId := getUserId(c)
	medicine, err := app.repo.GetDraftMedicine(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найден"))
		return
	}

	if err := verificationRequest(medicine.UUID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(`проверка ведомства в данный момент невозможна: {%s}`, err))
		return
	}

	verificationStatus := ds.VerificationStarted
	medicine.VerificationStatus = &verificationStatus
	medicine.Status = ds.FORMED
	now := time.Now()
	medicine.FormationDate = &now

	if err := app.repo.SaveMedicine(medicine); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Подтвердить лекарство
// @Tags		Лекарства
// @Description	Подтвердить или отменить лекарство модератором
// @Param		medicine_id path string true "id лекарства"
// @Param		confirm body boolean true "подтвердить"
// @Success		200
// @Router		/api/medicines/{medicine_id}/moderator_confirm [put]
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

	userId := getUserId(c)
	medicine, err := app.repo.GetMedicineById(request.URI.MedicineId, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}
	if medicine.Status != ds.FORMED {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%s\" на \"%s\"", medicine.Status, ds.FORMED))
		return
	}

	if *request.Confirm {
		medicine.Status = ds.COMPELTED
	} else {
		medicine.Status = ds.REJECTED
	}
	now := time.Now()
	medicine.CompletionDate = &now

	moderator, err := app.repo.GetUserById(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	medicine.ModeratorId = &userId
	medicine.Moderator = moderator

	if err := app.repo.SaveMedicine(medicine); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) Verification(c *gin.Context) {
	var request schemes.VerificationReq
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.Token != app.config.Token {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	transportation, err := app.repo.GetMedicineById(request.URI.MedicineId, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if transportation == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}

	var deliveryStatus string
	if *request.VerificationStatus {
		deliveryStatus = ds.VerificationCompleted
	} else {
		deliveryStatus = ds.VerificationFailed
	}
	transportation.VerificationStatus = &deliveryStatus

	if err := app.repo.SaveMedicine(transportation); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Изменить вес компонента
// @Tags		Лекарства
// @Description	Позволяет изменить вес компонента в лекарстве
// @Access		json
// @Param		component_id path string true "id компонента"
// @Param		count body schemes.ChangeCountReq true "Количество"
// @Success		200
// @Router		/api/medicines/change_count/{component_id} [put]
func (app *Application) ChangeCount(c *gin.Context) {
	var request schemes.ChangeCountReq
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userId := getUserId(c)
	medicine, err := app.repo.GetDraftMedicine(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("лекарство не найдено"))
		return
	}

	if app.repo.ChangeCount(medicine.UUID, request.URI.ComponentId, request.Count); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}