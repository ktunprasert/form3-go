package repository

import (
	"fmt"

	"github.com/ktunprasert/form3-go/sdk/client"
	"github.com/ktunprasert/form3-go/sdk/domain"
)

const (
	ORG_ACCOUNT_CREATE = "/v1/organisation/accounts"
	ORG_ACCOUNT_FETCH  = "/v1/organisation/accounts/%s"
	ORG_ACCOUNT_DELETE = "/v1/organisation/accounts/%s?version=%d"
)

type AccountRepository struct {
	Client client.ClientInterface
}

type AccountRepositoryInterface interface {
	Create(account *domain.Account) (*domain.Account, error)
	Fetch(id string) (*domain.Account, error)
	Delete(id string, version int64) error
}

func NewAccountRepo(client client.ClientInterface) AccountRepositoryInterface {
    return &AccountRepository{
        Client: client,
    }
}

func (r *AccountRepository) Create(account *domain.Account) (*domain.Account, error) {
	path := ORG_ACCOUNT_CREATE

	requestBody := domain.Request[domain.Account]{
		Data: account,
	}

	var resp domain.Response[domain.Account]
	err := r.Client.Create(path, requestBody, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (r *AccountRepository) Fetch(id string) (*domain.Account, error) {
	path := fmt.Sprintf(ORG_ACCOUNT_FETCH, id)

	var resp domain.Response[domain.Account]
	err := r.Client.Fetch(path, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (r *AccountRepository) Delete(id string, version int64) error {
	path := fmt.Sprintf(ORG_ACCOUNT_DELETE, id, version)

	err := r.Client.Delete(path, nil)
	if err != nil {
		return err
	}

	return nil
}
