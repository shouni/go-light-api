package repository

import (
	"database/sql"
	"errors" // errorsパッケージをインポート
	"fmt"
	"github.com/mattn/go-sqlite3" // sqlite3ドライバをインポート
	"go-light-api/internal/model"
)

// SQL文を定数として定義 (変更なし)
const createTableSQL = `
CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT
);`

const insertUserSQL = "INSERT INTO users(id, name, email) VALUES(?, ?, ?)"
const selectUserByIDSQL = "SELECT id, name, email FROM users WHERE id = ?"

// ------------------------------------
// カスタムエラー定義
// ------------------------------------

// ErrDuplicateEntry: ID重複エラーを示すSentinel Error
var ErrDuplicateEntry = errors.New("duplicate entry")

// ------------------------------------
// インターフェース定義 (変更なし)
// ------------------------------------

type UserRepository interface {
	InitTable() error
	Create(user *model.User) error
	FindByID(id string) (*model.User, error)
}

// userRepository (変更なし)
type userRepository struct {
	db *sql.DB
}

// NewUserRepository (変更なし)
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// ------------------------------------
// メソッド実装 (Create を修正)
// ------------------------------------

// InitTable (変更なし)
func (r *userRepository) InitTable() error {
	_, err := r.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize users table: %w", err)
	}
	return nil
}

// Create: ユーザーをデータベースに挿入 (ID重複チェックを追加)
func (r *userRepository) Create(user *model.User) error {
	_, err := r.db.Exec(insertUserSQL, user.ID, user.Name, user.Email)
	if err != nil {
		// SQLiteのエラー型に変換可能かチェック
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
			// ID重複の場合はカスタムエラーをラップして返す
			return fmt.Errorf("%w: user with ID %s already exists", ErrDuplicateEntry, user.ID)
		}

		// その他のDBエラー
		return fmt.Errorf("failed to create user %s: %w", user.ID, err)
	}
	return nil
}

// FindByID (変更なし)
func (r *userRepository) FindByID(id string) (*model.User, error) {
	u := &model.User{}

	row := r.db.QueryRow(selectUserByIDSQL, id)

	err := row.Scan(&u.ID, &u.Name, &u.Email)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to find user by ID '%s': %w", id, err)
	}

	return u, nil
}
