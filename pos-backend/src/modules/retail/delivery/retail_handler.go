package delivery

import (
	"pos-backend/src/modules/retail/repository"
)

type RetailHandler struct {
	Repo repository.RetailRepository
}

func NewRetailHandler(repo repository.RetailRepository) *RetailHandler {
	return &RetailHandler{Repo: repo}
}