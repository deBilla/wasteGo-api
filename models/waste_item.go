package models

type WasteItem struct {
	ID       string `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}
