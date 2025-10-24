package repository

import (
	"database/sql"
	"go-light-api/internal/model" // modelパッケージをインポート
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
// インターフェース定義 (Goの慣習に従い 'I' を削除)
// ------------------------------------

// UserRepository: DB操作の抽象化
type UserRepository interface {
	InitTable() error
	Create(user *model.User) error
	FindByID(id string) (*model.User, error)
}

// userRepository: UserRepositoryを実装する構造体 (unexportedに変更)
type userRepository struct {
	DB *sql.DB
}

// NewUserRepository: UserRepositoryのインスタンスを作成するコンストラクタ
func NewUserRepository(db *sql.DB) UserRepository { // 戻り値の型をUserRepositoryに変更
	return &userRepository{DB: db} // 構造体インスタンスも変更
}

// ------------------------------------
// メソッド実装
// ------------------------------------

// InitTable: テーブルが存在しない場合に作成
func (r *userRepository) InitTable() error {
	_, err := r.DB.Exec(createTableSQL)
	return err
}

// Create: ユーザーをデータベースに挿入
func (r *userRepository) Create(user *model.User) error {
	// プリペアドステートメントでSQLインジェクション対策
	_, err := r.DB.Exec(insertUserSQL, user.ID, user.Name, user.Email)
	return err
}

// FindByID: IDに基づいてユーザーを検索
func (r *userRepository) FindByID(id string) (*model.User, error) {
	u := &model.User{}

	row := r.DB.QueryRow(selectUserByIDSQL, id)

	// 取得したデータをGoの構造体にマッピング
	err := row.Scan(&u.ID, &u.Name, &u.Email)

	if err == sql.ErrNoRows {
		return nil, nil // データが見つからない場合はnil, nilを返す
	} else if err != nil {
		// 問題点2の修正: リポジトリ層でのログ出力を削除し、呼び出し元にエラーを返す
		return nil, err // DBエラーをそのまま呼び出し元に返す
	}

	return u, nil
}
