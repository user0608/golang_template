package repositories

import (
	"database/sql"
	"fmt"
	"goqrs/models"

	"github.com/ksaucedo002/answer/errores"
	"gorm.io/gorm"
)

type LoginRepository interface {
	HashPassword(tx *gorm.DB, password, hash string) (string, error)
	IncrementPasswordAttempt(tx *gorm.DB, username string) error
	ResetPasswordAttempt(tx *gorm.DB, username string) error
	FindAccount(tx *gorm.DB, username string) (*models.Account, error)
}
type login struct {
}

var _ LoginRepository = &login{}

func NewLoginRepository() *login {
	return &login{}
}

func (r *login) HashPassword(tx *gorm.DB, password, hash string) (string, error) {
	var currenthash sql.NullString
	const sql = "select crypt from crypt(?,?)"
	rs := tx.Raw(sql, password, hash)
	if rs.Error != nil {
		return "", errores.NewInternalDBf(rs.Error)
	}
	if rs := rs.Scan(&currenthash); rs.Error != nil {
		return "", errores.NewInternalDBf(rs.Error)
	}
	if !currenthash.Valid {
		return "", errores.NewBadRequestf(
			fmt.Errorf("login_repository:HashPassword error:%s", password),
			"no se pudo procesar la password",
		)
	}
	return currenthash.String, nil
}
func (r *login) IncrementPasswordAttempt(tx *gorm.DB, username string) error {
	account := models.Account{Username: username}
	rs := tx.Model(&account).Update("password_attempt", gorm.Expr("password_attempt + ?", 1))
	if rs.Error != nil {
		return errores.NewInternalDBf(rs.Error)
	}
	return nil
}
func (r *login) ResetPasswordAttempt(tx *gorm.DB, username string) error {
	account := models.Account{Username: username}
	rs := tx.Model(&account).Update("password_attempt", 0)
	if rs.Error != nil {
		return errores.NewInternalDBf(rs.Error)
	}
	return nil
}
func (r *login) FindAccount(tx *gorm.DB, username string) (*models.Account, error) {
	var account models.Account
	rs := tx.Find(&account, "username=?", username)
	if rs.Error != nil {
		return nil, errores.NewInternalDBf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return nil, errores.NewBadRequestf(nil, "usuario o contraseña invalidos")
	}
	return &account, nil
}
