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

// @Summary		Получить все компоненты
// @Tags		Компоненты
// @Description	Возвращает все компоненты с фильтрацией по статусу и дате формирования
// @Produce		json
// @Param		status query string false "статус компонента"
// @Param		formation_date_start query string false "начальная дата формирования"
// @Param		formation_date_end query string false "конечная дата формирвания"
// @Success		200 {object} schemes.AllComponentsResponse
// @Router		/api/components [get]
func (app *Application) GetAllComponents(c *gin.Context) {
	var request schemes.GetAllComponentsRequst
	var err error
	if err = c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	fmt.Println(userId, userRole)
	var components []ds.Component
	if userRole == role.Customer {
		components, err = app.repo.GetAllComponents(&userId, request.FormationDateStart, request.FormationDateEnd, request.Status)
	} else {
		components, err = app.repo.GetAllComponents(nil, request.FormationDateStart, request.FormationDateEnd, request.Status)
	}
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	outputComponents := make([]schemes.ComponentOutput, len(components))
	for i, component := range components {
		outputComponents[i] = schemes.ConvertComponent(&component)
	}
	c.JSON(http.StatusOK, schemes.AllComponentsResponse{Components: outputComponents})
}

// @Summary		Получить один компонент
// @Tags		Компоненты
// @Description	Возвращает подробную информацию о компоненте и его названии
// @Produce		json
// @Param		component_id path string true "id компонента"
// @Success		200 {object} schemes.ComponentResponse
// @Router		/api/components/{component_id} [get]
func (app *Application) GetComponent(c *gin.Context) {
	var request schemes.ComponentRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	component, err := app.repo.GetComponentById(request.ComponentId, userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("компонент не найден"))
		return
	}

	medicines, err := app.repo.GetMedicineProduction(request.ComponentId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.ComponentResponse{Component: schemes.ConvertComponent(component), Medicines: medicines})
}

type SwaggerUpdateComponentRequest struct {
	ComponentType string `json:"component_type"`
}

// @Summary		Указать название компонента
// @Tags		Компоненты
// @Description	Позволяет изменить название компонента и возвращает обновлённые данные
// @Access		json
// @Produce		json
// @Param		component_id path string true "id компонента"
// @Param		component_name body SwaggerUpdateComponentRequest true "Название компонента"
// @Success		200 {object} schemes.UpdateComponentResponse
// @Router		/api/components/{component_id} [put]
func (app *Application) UpdateComponent(c *gin.Context) {
	var request schemes.UpdateComponentRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userId := getUserId(c)
	component, err := app.repo.GetComponentById(request.URI.ComponentId, userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("компонент не найден"))
		return
	}
	component.ComponentName = request.ComponentName
	if app.repo.SaveComponent(component); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.UpdateComponentResponse{Component:schemes.ConvertComponent(component)})
}

// @Summary		Удалить компонент
// @Tags		Компоненты
// @Description	Удаляет компонент по id
// @Param		component_id path string true "id компонента"
// @Success		200
// @Router		/api/components/{component_id} [delete]
func (app *Application) DeleteComponent(c *gin.Context) {
	var request schemes.ComponentRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	component, err := app.repo.GetComponentById(request.ComponentId,userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("компонент не найден"))
		return
	}

	userROle := getUserRole(c)
	if userROle == role.Customer && component.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("компонент уже сформирован"))
		return
	}
	component.Status = ds.DELETED

	if err := app.repo.SaveComponent(component); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Удалить лекарство из компонента
// @Tags		Компоненты
// @Description	Удалить лекарство из компонента
// @Produce		json
// @Param		component_id path string true "id компонента"
// @Param		medicine_id path string true "id лекарства"
// @Success		200 {object} schemes.AllMedicinesResponse
// @Router		/api/components/{component_id}/delete_medicine/{medicine_id} [delete]
func (app *Application) DeleteFromComponent(c *gin.Context) {
	var request schemes.DeleteFromComponentRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userId := getUserId(c)
	component, err := app.repo.GetComponentById(request.ComponentId,userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("компонент не найден"))
		return
	}
	if component.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя редактировать компонент со статусом: %s", component.Status))
		return
	}

	if err := app.repo.DeleteFromComponent(request.ComponentId, request.MedicineId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	medicines, err := app.repo.GetMedicineProduction(request.ComponentId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllMedicinesResponse{Medicines: medicines})
}

// @Summary		Сформировать компонент
// @Tags		Компоненты
// @Description	Сформировать или удалить компонент пользователем
// @Success		200
// @Router		/api/components/user_confirm [put]
func (app *Application) UserConfirm(c *gin.Context) {
	userId := getUserId(c)
	component, err := app.repo.GetDraftComponent(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("компонент не найден"))
		return
	}
	if component.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя сформировать компонент со статусом %s", component.Status))
		return
	}
	component.Status = ds.FORMED
	now := time.Now()
	component.FormationDate = &now

	if err := app.repo.SaveComponent(component); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Подтвердить компонент
// @Tags		Компоненты
// @Description	Подтвердить или отменить компонент модератором
// @Param		component_id path string true "id компонента"
// @Param		confirm body boolean true "подтвердить"
// @Success		200
// @Router		/api/components/{component_id}/moderator_confirm [put]
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
	component, err := app.repo.GetComponentById(request.URI.ComponentId,userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("компонент не найден"))
		return
	}
	if component.Status != ds.FORMED {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%s\" на \"%s\"", component.Status,  ds.FORMED))
		return
	}


	if request.Confirm {
		component.Status = ds.COMPELTED
		now := time.Now()
		component.CompletionDate = &now
	
	} else {
		component.Status = ds.REJECTED
	}
	component.ModeratorId = &userId
	
	if err := app.repo.SaveComponent(component); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
