package repository

import (
	"database/sql"
	"fmt" // エラーラップのために追加
	"go-light-api/internal/model"
)

// SQL文を定数として定義
const createTableSQL = `
CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT
);`

const insertUserSQL = "INSERT INTO users(id, name, email) VALUES(?, ?, ?)"
const selectUserByIDSQL = "SELECT id, name, email FROM users WHERE id = ?"

// ------------------------------------
// インターフェース定義
// ------------------------------------

// UserRepository: DB操作の抽象化
type UserRepository interface {
	InitTable() error
	Create(user *model.User) error
	FindByID(id string) (*model.User, error)
}

// userRepository: UserRepositoryを実装する構造体 (unexported)
type userRepository struct {
	db *sql.DB // フィールド名を小文字の 'db' に変更 (unexported)
}

// NewUserRepository: UserRepositoryのインスタンスを作成するコンストラクタ
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// ------------------------------------
// メソッド実装
// ------------------------------------

// InitTable: テーブルが存在しない場合に作成
func (r *userRepository) InitTable() error {
	_, err := r.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize users table: %w", err) // エラーをラップ
	}
	return nil
}

// Create: ユーザーをデータベースに挿入
func (r *userRepository) Create(user *model.User) error {
	// プリペアドステートメントでSQLインジェクション対策
	_, err := r.db.Exec(insertUserSQL, user.ID, user.Name, user.Email)
	if err != nil {
		return fmt.Errorf("failed to create user %s: %w", user.ID, err) // エラーをラップ
	}
	return nil
}

// FindByID: IDに基づいてユーザーを検索
func (r *userRepository) FindByID(id string) (*model.User, error) {
	u := &model.User{}

	row := r.db.QueryRow(selectUserByIDSQL, id)

	err := row.Scan(&u.ID, &u.Name, &u.Email)

	if err == sql.ErrNoRows {
		return nil, nil // データが見つからない場合はnil, nilを返す
	} else if err != nil {
		// エラーにコンテキストを追加し、元のエラーをラップ
		return nil, fmt.Errorf("failed to find user by ID '%s': %w", id, err)
	}

	return u, nil
}
