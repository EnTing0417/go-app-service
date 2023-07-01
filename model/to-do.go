package model

type ToDoRequestBody struct {
	Title       string    `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

type ToDoPatchRequestBody struct {
	Completed   *bool      `json:"completed" bson:"completed"`
}

type ToDoDeleteRequestBody struct {
	ID		[]string	 `json:"id" bson:"id"`
}

type ToDoViewModel struct {
	ID          string      `json:"_id" bson:"_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	CreatedAt   string `json:"created_at" bson:"created_at"`
	UpdatedAt   string`json:"updated_at" bson:"updated_at"`
	DeletedAt	string`json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
	Completed   bool      `json:"completed" bson:"completed"`
	UserID		string	  `json:"user_id" bson:"user_id"`
}