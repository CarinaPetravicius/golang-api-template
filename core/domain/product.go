package domain

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
