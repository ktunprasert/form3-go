package repository_test

import (
	"os"
	"strings"
	"testing"

	"github.com/ktunprasert/form3-go/sdk/client"
	"github.com/ktunprasert/form3-go/sdk/domain"
	"github.com/ktunprasert/form3-go/sdk/repository"
)

func TestCreateAccount(t *testing.T) {
	accountRepo := setupClient(t)

	var ID string = "9d932617-294d-4065-a2cd-ba9de4af0988"
	var version int64 = 0
	var country string = "UK"
	account := &domain.Account{
		Attributes: &domain.AccountAttributes{
			Name:    []string{"James"},
			Country: &country,
		},
		ID:             ID,
		OrganisationID: "7f249544-8841-4a1e-9576-8c6b32dce797",
		Type:           "accounts",
		Version:        &version,
	}

	created, createErr := accountRepo.Create(account)
	if createErr != nil {
		t.Error("creation error", createErr)
	}
	defer func() {
		err := accountRepo.Delete(created.ID, *created.Version)
		if err != nil {
			t.Error("teardown error", err)
		}
	}()

	fetched, fetchErr := accountRepo.Fetch(ID)
	if fetchErr != nil {
		t.Error(fetchErr)
	}

	if fetched.ID != created.ID {
		t.Error("Unexpected error")
	}
}

func TestCreateAccount_ReturnsConflictOnDuplicateID(t *testing.T) {
	accountRepo := setupClient(t)
	createdAccount, teardown := setupAccount(t, accountRepo)
	defer teardown()

	_, err := accountRepo.Create(createdAccount)
	if err == nil {
		t.Error("Conflict not thrown")
	}

	if err.Error() != "[Conflict]: Account cannot be created as it violates a duplicate constraint" {
		t.Error("Unexpected error thrown")
	}
}

func TestCreateAccount_ReturnsBadRequestOnInvalidDomain(t *testing.T) {
	accountRepo := setupClient(t)

	emptyAccountObj := &domain.Account{}

	_, err := accountRepo.Create(emptyAccountObj)
	if err == nil {
		t.Error("Bad request not thrown")
	}

	if !strings.HasPrefix(err.Error(), "[Bad Request]") {
		t.Error("Unexpected error thrown")
	}
}

func TestFetchAccount(t *testing.T) {
	accountRepo := setupClient(t)
	createdAccount, teardown := setupAccount(t, accountRepo)
	defer teardown()

	fetched, err := accountRepo.Fetch(createdAccount.ID)
	if err != nil {
		t.Error("Unexpected error thrown")
	}

	if fetched.ID != createdAccount.ID {
		t.Error("Account doesn't match")
	}
}

func TestFetchAccount_ReturnsNotFoundOnNonExistentId(t *testing.T) {
	accountRepo := setupClient(t)

	unknownUUID := "95288fda-8f71-416a-b600-2d53e0a04688"
	_, err := accountRepo.Fetch(unknownUUID)
	if err == nil {
		t.Error("Not found not thrown")
	}

	if !strings.HasPrefix(err.Error(), "[Not Found]") {
		t.Error("Unexpected error thrown")
	}
}

func TestFetchAccount_ReturnsBadRequestOnInvalidUUID(t *testing.T) {
	accountRepo := setupClient(t)

	invalidIDs := []string{
		"1",
		"string_id",
	}

	for _, id := range invalidIDs {
		_, err := accountRepo.Fetch(id)
		if err == nil {
			t.Error("Not found not thrown")
		}

		if !strings.HasPrefix(err.Error(), "[Bad Request]") {
			t.Error("Unexpected error thrown")
		}
	}
}

func TestDeleteAccount(t *testing.T) {
	accountRepo := setupClient(t)
	createdAccount, _ := setupAccount(t, accountRepo)

	err := accountRepo.Delete(createdAccount.ID, *createdAccount.Version)
	if err != nil {
		t.Error("Unexpected error thrown")
	}
}

func TestDeleteAccount_ReturnsNotFoundOnNonExistentId(t *testing.T) {
	accountRepo := setupClient(t)

	unknownUUID := "95288fda-8f71-416a-b600-2d53e0a04688"
	err := accountRepo.Delete(unknownUUID, 0)
	if err == nil {
		t.Error("Not found not thrown")
	}

	if !strings.HasPrefix(err.Error(), "[Not Found]") {
		t.Error("Unexpected error thrown")
	}
}

func TestDeleteAccount_ReturnsBadRequestOnInvalidUUID(t *testing.T) {
	accountRepo := setupClient(t)

	invalidIDs := []string{
		"1",
		"string_id",
	}

	for _, id := range invalidIDs {
		err := accountRepo.Delete(id, 0)
		if err == nil {
			t.Error("Bad request not thrown")
		}

		if !strings.HasPrefix(err.Error(), "[Bad Request]") {
			t.Error("Unexpected error thrown")
		}
	}
}

func TestDeleteAccount_ReturnsConflictOnInvalidVersion(t *testing.T) {
	accountRepo := setupClient(t)
	createdAccount, teardown := setupAccount(t, accountRepo)
    defer teardown()

	invalidVersion := []int64{
        -1,
        10000,
	}

	for _, version := range invalidVersion {
		err := accountRepo.Delete(createdAccount.ID, version)
		if err == nil {
			t.Error("Conflict not thrown")
		}

		if !strings.HasPrefix(err.Error(), "[Conflict]") {
			t.Error("Unexpected error thrown")
		}
	}

}

func setupAccount(t *testing.T, repo repository.AccountRepositoryInterface) (*domain.Account, func()) {
	var ID string = "1237573b-185b-4f80-bb70-ea295975ca96"
	var version int64 = 0
	var country string = "UK"
	account := &domain.Account{
		Attributes: &domain.AccountAttributes{
			Name:    []string{"James"},
			Country: &country,
		},
		ID:             ID,
		OrganisationID: "7f249544-8841-4a1e-9576-8c6b32dce797",
		Type:           "accounts",
		Version:        &version,
	}

	created, err := repo.Create(account)
	if err != nil {
		t.Error("Setup: Error creating account")
	}

	teardown := func() {
		err := repo.Delete(created.ID, *created.Version)
		if err != nil {
			t.Error("Teardown: Error deleting account")
		}
	}

	return created, teardown
}

func setupClient(t *testing.T) repository.AccountRepositoryInterface {
	var HOST string = os.Getenv("HOST")
	testClient := client.New(HOST)
	accountRepo := repository.NewAccountRepo(testClient)

	return accountRepo
}
