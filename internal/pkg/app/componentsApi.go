package app

import (
	"fmt"
	"net/http"
	"time"

	"lab3/internal/app/ds"
	"lab3/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

func (app *Application) GetAllComponents(c *gin.Context) {
	var request schemes.GetAllComponentsRequst
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	components, err := app.repo.GetAllComponents(request.FormationDateStart, request.FormationDateEnd, request.Status)
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

func (app *Application) GetComponent(c *gin.Context) {
	var request schemes.ComponentRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	component, err := app.repo.GetComponentById(request.ComponentId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("компонент не найден"))
		return
	}

	medicines, err := app.repo.GetComponentContent(request.ComponentId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.ComponentResponse{Component: schemes.ConvertComponent(component), Medicines: medicines})
}

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
	component, err := app.repo.GetComponentById(request.URI.ComponentId, app.getCustomer())
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

func (app *Application) DeleteComponent(c *gin.Context) {
	var request schemes.ComponentRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	component, err := app.repo.GetComponentById(request.ComponentId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("компонент не найден"))
		return
	}
	component.Status = ds.DELETED

	if err := app.repo.SaveComponent(component); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) DeleteFromComponent(c *gin.Context) {
	var request schemes.DeleteFromComponentRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	component, err := app.repo.GetComponentById(request.ComponentId, app.getCustomer())
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

	medicines, err := app.repo.GetComponentContent(request.ComponentId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllMedicinesResponse{Medicines: medicines})
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

	component, err := app.repo.GetComponentById(request.URI.ComponentId, app.getCustomer())
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
	if request.Confirm {
		component.Status = ds.FORMED
		now := time.Now()
		component.FormationDate = &now
	} else {
		component.Status = ds.DELETED
	}

	if err := app.repo.SaveComponent(component); err != nil {
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

	

	component, err := app.repo.GetComponentById(request.URI.ComponentId, app.getCustomer())
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
	component.ModeratorId = app.getModerator()
	
	if err := app.repo.SaveComponent(component); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
