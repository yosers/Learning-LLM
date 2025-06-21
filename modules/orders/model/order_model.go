package order_model

type OrderQuery struct {
	Limit       int `form:"limit" binding:"min=1"`
	CurrentPage int `form:"page" binding:"min=1"`
}
