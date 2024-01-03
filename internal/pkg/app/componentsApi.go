package app

import (
	"fmt"
	"net/http"

	"lab3/internal/app/ds"
	"lab3/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

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
	draftMedicine, err := app.repo.GetDraftMedicine(app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response := schemes.GetAllComponentsResponse{DraftMedicine: nil, Components: components}
	if draftMedicine != nil {
		response.DraftMedicine = &schemes.MedicineShort{UUID: draftMedicine.UUID}
		components, err := app.repo.GetMedicineProduction(draftMedicine.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		response.DraftMedicine.ComponentCount = len(components)
	}
	c.JSON(http.StatusOK, response)
}

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

	c.Status(http.StatusOK)
}

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
	medicine, err = app.repo.GetDraftMedicine(app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if medicine == nil {
		medicine, err = app.repo.CreateDraftMedicine(app.getCustomer())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err = app.repo.AddToMedicine(medicine.UUID, request.ComponentId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var components []ds.Component
	components, err = app.repo.GetMedicineProduction(medicine.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllComponentsResponse{Components: components})
}