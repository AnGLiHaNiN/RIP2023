package app

import (
	"fmt"
	"net/http"

	_ "R_I_P_labs/docs"
	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

// @Summary		Получить все компоненты
// @Tags		Компоненты
// @Description	Возвращает все доступные компоненты с опциональной фильтрацией по Названию
// @Produce		json
// @Param		name query string false "Название для фильтрации"
// @Success		200 {object} schemes.GetAllComponentsResponse
// @Router		/api/components [get]
func (app *Application) GetAllComponents(c *gin.Context) {
	var request schemes.GetAllComponentsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	components, err := app.repo.GetComponentByName(request.Name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var draftMedicine *ds.Medicine = nil
	if userId, exists := c.Get("userId"); exists {
		draftMedicine, err = app.repo.GetDraftMedicine(userId.(string))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	response := schemes.GetAllComponentsResponse{DraftMedicine: nil, Components: components}
	if draftMedicine != nil {
		response.DraftMedicine = &schemes.MedicineShort{UUID: draftMedicine.UUID}
		componentsCount, err := app.repo.CountComponents(draftMedicine.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		response.DraftMedicine.ComponentCount = int(componentsCount)
	}
	c.JSON(http.StatusOK, response)
}

// @Summary		Получить один компонент
// @Tags		Компоненты
// @Description	Возвращает более подробную информацию об одном компоненте
// @Produce		json
// @Param		component_id path string true "id компонента"
// @Success		200 {object} ds.Component
// @Router		/api/components/{component_id} [get]
func (app *Application) GetComponent(c *gin.Context) {
	var request schemes.ComponentRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	component, err := app.repo.GetComponentByID(request.ComponentId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("компонент не найден"))
		return
	}
	c.JSON(http.StatusOK, component)
}

// @Summary		Удалить компонент
// @Tags		Компоненты
// @Description	Удаляет лекраство по id
// @Param		component_id path string true "id компонента"
// @Success		200
// @Router		/api/components/{component_id} [delete]
func (app *Application) DeleteComponent(c *gin.Context) {
	var request schemes.ComponentRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	component, err := app.repo.GetComponentByID(request.ComponentId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("компонент не найден"))
		return
	}
	component.ImageURL = nil
	component.IsDeleted = true
	if err := app.repo.SaveComponent(component); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Добавить компонент
// @Tags		Компоненты
// @Description	Добавить новый компонент
// @Accept		mpfd
// @Param     	image formData file false "Изображение компонента"
// @Param     	name formData string true "Название" format:"string" maxLength:100
// @Param     	world_name formData string true "Всемирное наименование" format:"string" maxLength:100
// @Param     	amount formData int true "Количество" format:"int"
// @Param     	properties formData string true "Свойства" format:"string" maxLength:100
// @Success		200
// @Router		/api/components/ [post]
func (app *Application) AddComponent(c *gin.Context) {
	var request schemes.AddComponentRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	component := ds.Component(request.Component)
	if err := app.repo.AddComponent(&component); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if request.Image != nil {
		imageURL, err := app.uploadImage(c, request.Image, component.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		component.ImageURL = imageURL
	}
	if err := app.repo.SaveComponent(&component); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary		Изменить компонент
// @Tags		Компоненты
// @Description	Изменить данные полей о компоненте
// @Accept		mpfd
// @Produce		json
// @Param		component_id path string true "Идентификатор компонента" format:"uuid"
// @Param		name formData string false "Название" format:"string" maxLength:100
// @Param		world_name formData string false "Всемирное наименование" format:"string" maxLength:100
// @Param		amount formData int false "Количество" format:"int"
// @Param		image formData file false "Изображение компоненты"
// @Param		properties formData string false "Свойства" format:"string" maxLength:100
// @Router		/api/components/{component_id} [put]
func (app *Application) ChangeComponent(c *gin.Context) {
	var request schemes.ChangeComponentRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	component, err := app.repo.GetComponentByID(request.ComponentId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("компонент не найден"))
		return
	}

	if request.Name != nil {
		component.Name = *request.Name
	}
	if request.Image != nil {
		if component.ImageURL != nil {
			if err := app.deleteImage(c, component.UUID); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		imageURL, err := app.uploadImage(c, request.Image, component.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		component.ImageURL = imageURL
	}
	if request.WorldName != nil {
		component.WorldName = *request.WorldName
	}
	if request.Amount != nil {
		component.Amount = *request.Amount
	}
	if request.Properties != nil {
		component.Properties = *request.Properties
	}

	if err := app.repo.SaveComponent(component); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, component)
}

// @Summary		Добавить в компонент
// @Tags		Компоненты
// @Description	Добавить выбранный компонент в черновик лекарства
// @Produce		json
// @Param		component_id path string true "id компонента"
// @Success		200 {object} schemes.AddToMedicineResp
// @Router		/api/components/{component_id}/add_to_medicine [post]
func (app *Application) AddToMedicine(c *gin.Context) {
	var request schemes.AddToMedicineRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var err error

	component, err := app.repo.GetComponentByID(request.ComponentId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if component == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("компонент не найден"))
		return
	}

	var medicine *ds.Medicine
	userId := getUserId(c)
	medicine, err = app.repo.GetDraftMedicine(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		medicine, err = app.repo.CreateDraftMedicine(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err = app.repo.AddToMedicine(medicine.UUID, request.ComponentId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	componentsCount, err := app.repo.CountComponents(medicine.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AddToMedicineResp{ComponentsCount: componentsCount})
}