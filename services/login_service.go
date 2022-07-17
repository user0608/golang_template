package services

import (
	"context"
	"goqrs/database"
	"goqrs/models"
	"goqrs/repositories"
	"log"

	"github.com/ksaucedo002/answer/errores"
)

type LoginService interface {
	Login(ctx context.Context, username, password string) (*models.Account, error)
}
type login struct {
	r repositories.LoginRepository
}

var _ LoginService = &login{}

func NewLoginService(r repositories.LoginRepository) *login {
	return &login{r: r}
}

func (s *login) Login(ctx context.Context, username, password string) (*models.Account, error) {
	if username == "" || password == "" {
		return nil, errores.NewBadRequestf(nil, "usuario o contraseña no encontrados")
	}
	tx := database.Conn(ctx)
	account, err := s.r.FindAccount(tx, username)
	if err != nil {
		return nil, err
	}
	if account.PasswordAttempt > 8 {
		return nil, errores.NewBadRequestf(nil, "cuenta suspendida, por contraseña invalida")
	}
	passwordHash, err := s.r.HashPassword(tx, password, account.Password)
	if err != nil {
		return nil, err
	}
	if account.Password != passwordHash {
		if err := s.r.IncrementPasswordAttempt(tx, username); err != nil {
			log.Println(err)
		}
		return nil, errores.NewBadRequestf(nil, "usuario o contraseña invalidos ")
	}
	if err := s.r.ResetPasswordAttempt(tx, username); err != nil {
		log.Println(err)
	}
	return account, nil
}
