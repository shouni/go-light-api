package repository

import (
	"database/sql"
	"go-light-api/internal/model" // modelパッケージをインポート
	"log"
)

// SQL文を定数として定義 (責務の分離)
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

// UserRepositoryI: DB操作の抽象化
type UserRepositoryI interface {
	InitTable() error
	Create(user *model.User) error
	FindByID(id string) (*model.User, error)
}

// UserRepository: UserRepositoryIを実装する構造体
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository: UserRepositoryのインスタンスを作成するコンストラクタ
func NewUserRepository(db *sql.DB) UserRepositoryI {
	return &UserRepository{DB: db}
}

// ------------------------------------
// メソッド実装
// ------------------------------------

// InitTable: テーブルが存在しない場合に作成
func (r *UserRepository) InitTable() error {
	_, err := r.DB.Exec(createTableSQL)
	return err
}

// Create: ユーザーをデータベースに挿入
func (r *UserRepository) Create(user *model.User) error {
	// プリペアドステートメントでSQLインジェクション対策
	_, err := r.DB.Exec(insertUserSQL, user.ID, user.Name, user.Email)
	return err
}

// FindByID: IDに基づいてユーザーを検索
func (r *UserRepository) FindByID(id string) (*model.User, error) {
	u := &model.User{}

	row := r.DB.QueryRow(selectUserByIDSQL, id)

	// 取得したデータをGoの構造体にマッピング
	err := row.Scan(&u.ID, &u.Name, &u.Email)

	if err == sql.ErrNoRows {
		return nil, nil // データが見つからない場合はnil, nilを返す
	} else if err != nil {
		log.Printf("Repository Error querying user %s: %v", id, err)
		return nil, err // その他のDBエラー
	}

	return u, nil
}
