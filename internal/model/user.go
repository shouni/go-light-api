package model

// User: ユーザーデータ構造体（DBとのやり取り、リクエストボディのデコードに使用）
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserResponse: APIレスポンス構造体
type UserResponse struct {
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"` // データがない場合は省略
}
