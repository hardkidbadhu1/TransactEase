//go:generate mockgen -source=account_service.go -destination=mocks/account_service.go -package=mocks

package service

import (
	"github.com/gin-gonic/gin"
	"transact-api/model/dto/request"
	"transact-api/model/dto/response"
	"transact-api/repository"
)

type AccountService interface {
	InsertAccount(ctx *gin.Context, req request.AccountCreateRequest) error
	GetAccount(ctx *gin.Context, documentNumber string) (*response.AccountResponse, error)
}

type accountService struct {
	repo repository.AccountRepository
}

func (a accountService) InsertAccount(ctx *gin.Context, req request.AccountCreateRequest) error {
	return a.repo.InsertAccount(ctx, req.ToEntity())
}

func (a accountService) GetAccount(ctx *gin.Context, documentNumber string) (*response.AccountResponse, error) {
	account, err := a.repo.FindAccountByDocumentNumber(ctx, documentNumber)
	if err != nil {
		return nil, err
	}

	return account.ToResponse(), nil
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{
		repo: repo,
	}
}
