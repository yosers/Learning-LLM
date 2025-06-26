package service

type ListShopRequest struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"limit"`
}

type ShopsResponse struct {
	ID            int32   `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	LogoUrl       string  `json:"logo_url"`
	WebsiteUrl    string  `json:"website_url"`
	Email         string  `json:"email"`
	WhatsappPhone string  `json:"whatsapp_phone"`
	Address       string  `json:"address"`
	City          string  `json:"city"`
	State         string  `json:"state"`
	IsActive      bool    `json:"is_active"`
	Latitude      float32 `json:"latitude"`
	Longitude     float32 `json:"longitude"`
	ZipCode       string  `json:"zip_code"`
	Country       string  `json:"country"`
}

type ListShopsResponse struct {
	Shops      []ShopsResponse `json:"shops"`
	Total      int32           `json:"total_items"`
	Page       int32           `json:"current_page"`
	PageSize   int32           `json:"page_size"`
	TotalPages int32           `json:"total_pages"`
}

type ShopsRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	LogoUrl       string  `json:"logo_url"`
	WebsiteUrl    string  `json:"website_url"`
	Email         string  `json:"email"`
	WhatsappPhone string  `json:"whatsapp_phone"`
	Address       string  `json:"address"`
	City          string  `json:"city"`
	State         string  `json:"state"`
	IsActive      bool    `json:"is_active"`
	Latitude      float32 `json:"latitude"`
	Longitude     float32 `json:"longitude"`
	ZipCode       string  `json:"zip_code"`
	Country       string  `json:"country"`
}
