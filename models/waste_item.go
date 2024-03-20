package models

type WasteItem struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
	UserID   string `json:"user_id"`
	ImgURL   string `json:"img_url"`
}
