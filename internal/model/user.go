package model

// User: ユーザーデータ構造体
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserResponse: APIレスポンス構造体 (ユーザー情報を含む)
type UserResponse struct {
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
}

// HealthCheckResponse: ヘルスチェック用のレスポンス構造体
type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
