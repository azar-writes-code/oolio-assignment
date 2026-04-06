package products

type Image struct {
	Thumbnail string `json:"thumbnail" binding:"required,url"`
	Mobile    string `json:"mobile" binding:"required,url"`
	Tablet    string `json:"tablet" binding:"required,url"`
	Desktop   string `json:"desktop" binding:"required,url"`
}

type Product struct {
	ID          int32    `json:"id,omitempty"`
	Name        string   `json:"name" binding:"required,min=3,max=100"`
	Description *string  `json:"description" binding:"omitempty,min=10"`
	Price       float64  `json:"price" binding:"required,gte=0"`
	Category    []string `json:"category" binding:"required,dive"`
	Image       Image    `json:"image" binding:"omitempty"`
	Stock       int32    `json:"stock" binding:"required,gte=0"`
}
