package schemes

import (
	"R_I_P_labs/internal/app/ds"
	"time"
)

type AllMedicinesResponse struct {
	Medicines []ds.Medicine `json:"medicines"`
}

type ComponentShort struct {
	UUID           string `json:"uuid"`
	MedicineCount int    `json:"medicine_count"`
}

type GetAllMedicinesResponse struct {
	DraftComponent *ComponentShort         `json:"draft_component"`
	Medicines            []ds.Medicine `json:"medicines"`
}

type AllComponentsResponse struct {
	Components []ComponentOutput `json:"components"`
}

type ComponentResponse struct {
	Component ComponentOutput `json:"component"`
	Medicines    []ds.Medicine  `json:"medicines"`
}

type UpdateComponentResponse struct {
	Component ComponentOutput  `json:"components"`
}

type ComponentOutput struct {
	UUID           string  `json:"uuid"`
	Status         string  `json:"status"`
	CreationDate   string  `json:"creation_date"`
	FormationDate  *string `json:"formation_date"`
	CompletionDate *string `json:"completion_date"`
	Moderator      *string `json:"moderator"`
	Customer       string  `json:"customer"`
	ComponentName      string  `json:"component_name"`
}

func ConvertComponent(component *ds.Component) ComponentOutput {
	output := ComponentOutput{
		UUID:         component.UUID,
		Status:       component.Status,
		CreationDate: component.CreationDate.Format("2006-01-02 15:04:05"),
		ComponentName:    component.ComponentName,
		Customer:     component.Customer.Login,
	}

	if component.FormationDate != nil {
		formationDate := component.FormationDate.Format("2006-01-02 15:04:05")
		output.FormationDate = &formationDate
	}

	if component.CompletionDate != nil {
		completionDate := component.CompletionDate.Format("2006-01-02 15:04:05")
		output.CompletionDate = &completionDate
	}

	if component.Moderator != nil {
		output.Moderator = &component.Moderator.Login
	}

	return output
}

type AddToComponentResp struct {
	MedicinesCount int64 `json:"medicine_count"`
}

type LoginResp struct {
	ExpiresIn   time.Duration `json:"expires_in"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
}

type SwaggerLoginResp struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type RegisterResp struct {
	Ok bool `json:"ok"`
}