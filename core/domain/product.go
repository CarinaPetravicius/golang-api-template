package domain

import "time"

// Product request product
type Product struct {
	Name        string `json:"name" validate:"required,not_blank,min=2,max=256"`
	Description string `json:"description" validate:"omitempty,min=2,max=256"`
	UnitType    string `json:"unitType" validate:"required,oneof=unit kilos grams liters box size"`
	Unit        string `json:"unit" validate:"required,not_blank,min=1,max=50"`
	Brand       string `json:"brand" validate:"required,not_blank,min=1,max=50"`
	Color       string `json:"color" validate:"required,not_blank,min=1,max=50"`
	Style       string `json:"style" validate:"required,not_blank,min=1,max=50"`
	Status      string `json:"status" validate:"required,oneof=available pending inactive"`
}

// ProductResponse product response
type ProductResponse struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	UnitType     string    `json:"unitType"`
	Unit         string    `json:"unit"`
	Brand        string    `json:"brand"`
	Color        string    `json:"color"`
	Style        string    `json:"style"`
	Status       string    `json:"status"`
	CreationDate time.Time `json:"creationDate"`
	UpdateDate   time.Time `json:"updateDate"`
	AuditUser    string    `json:"auditUser"`
}
