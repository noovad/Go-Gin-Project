package data

type TagRequest struct {
	Name string `validate:"required,min=4,max=200" json:"name"`
}

type TagResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}