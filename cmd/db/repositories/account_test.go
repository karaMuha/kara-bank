package db

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AccountTestSuite struct {
	suite.Suite
	ctx context.Context
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, &AccountTestSuite{
		ctx: context.Background(),
	})
}

func (suite *AccountTestSuite) AfterTest(suiteName string, testName string) {
	// clear accounts table after every test to avoid dependencies and side effects between tests
	_, err := testStore.ClearAccountsTable()
	require.NoError(suite.T(), err)

	_, err = testStore.ClearUsersTable()
	require.NoError(suite.T(), err)
}

func (suite *AccountTestSuite) TestCreateAccount() {
	registerUserParam := &RegisterUserParams{
		Email:          "Max@Mustermann.de",
		HashedPassword: "",
		FirstName:      "Max",
		LastName:       "Mustermann",
	}
	user := registerTestUser(suite.T(), registerUserParam)

	createAccountParam := CreateAccountParams{
		Owner:    user.Email,
		Balance:  100,
		Currency: "EUR",
	}

	account, err := testStore.CreateAccount(suite.ctx, &createAccountParam)

	require.NoError(suite.T(), err)
	require.NotEmpty(suite.T(), account)

	require.Equal(suite.T(), createAccountParam.Owner, account.Owner)
	require.Equal(suite.T(), createAccountParam.Balance, account.Balance)
	require.Equal(suite.T(), createAccountParam.Currency, account.Currency)

	require.NotZero(suite.T(), account.ID)
	require.NotZero(suite.T(), account.CreatedAt)
}

func (suite *AccountTestSuite) TestGetAccount() {
	registerUserParam := &RegisterUserParams{
		Email:          "Max@Mustermann.de",
		HashedPassword: "",
		FirstName:      "Max",
		LastName:       "Mustermann",
	}
	user := registerTestUser(suite.T(), registerUserParam)

	arg := CreateAccountParams{
		Owner:    user.Email,
		Balance:  100,
		Currency: "EUR",
	}

	account1 := createTestAccount(suite.T(), arg)

	account2, err := testStore.GetAccount(suite.ctx, account1.ID)

	require.NoError(suite.T(), err)
	require.NotEmpty(suite.T(), account2)

	require.Equal(suite.T(), account1.ID, account2.ID)
	require.Equal(suite.T(), account1.Owner, account2.Owner)
	require.Equal(suite.T(), account1.Balance, account2.Balance)
	require.Equal(suite.T(), account1.Currency, account2.Currency)
	require.WithinDuration(suite.T(), account1.CreatedAt, account2.CreatedAt, time.Second)
}

func (suite *AccountTestSuite) TestUpdateAccount() {
	registerUserParam := &RegisterUserParams{
		Email:          "Max@Mustermann.de",
		HashedPassword: "",
		FirstName:      "Max",
		LastName:       "Mustermann",
	}
	user := registerTestUser(suite.T(), registerUserParam)

	arg := CreateAccountParams{
		Owner:    user.Email,
		Balance:  100,
		Currency: "EUR",
	}

	account1 := createTestAccount(suite.T(), arg)

	arg2 := UpdateAccountParams{
		ID:      account1.ID,
		Balance: 200,
	}

	account2, err := testStore.UpdateAccount(suite.ctx, &arg2)

	require.NoError(suite.T(), err)
	require.NotEmpty(suite.T(), account2)

	require.Equal(suite.T(), account1.ID, account2.ID)
	require.Equal(suite.T(), account1.Owner, account2.Owner)
	require.Equal(suite.T(), arg2.Balance, account2.Balance)
	require.Equal(suite.T(), account1.Currency, account2.Currency)
	require.WithinDuration(suite.T(), account1.CreatedAt, account2.CreatedAt, time.Second)
}

func (suite *AccountTestSuite) TestDeleteAccount() {
	registerUserParam := &RegisterUserParams{
		Email:          "Max@Mustermann.de",
		HashedPassword: "",
		FirstName:      "Max",
		LastName:       "Mustermann",
	}
	user := registerTestUser(suite.T(), registerUserParam)

	arg := CreateAccountParams{
		Owner:    user.Email,
		Balance:  100,
		Currency: "EUR",
	}

	account1 := createTestAccount(suite.T(), arg)

	err := testStore.DeleteAccount(suite.ctx, account1.ID)
	require.NoError(suite.T(), err)

	account2, err := testStore.GetAccount(suite.ctx, account1.ID)
	require.Error(suite.T(), err)
	require.EqualError(suite.T(), err, pgx.ErrNoRows.Error())
	require.Empty(suite.T(), account2)
}

func (suite *AccountTestSuite) TestListAccounts() {
	for i := 0; i < 10; i++ {
		email := "Max" + strconv.Itoa(i) + "@Mustermann.de"
		registerUserParam := &RegisterUserParams{
			Email:          email,
			HashedPassword: "",
			FirstName:      "Max",
			LastName:       "Mustermann",
		}
		user := registerTestUser(suite.T(), registerUserParam)

		arg := CreateAccountParams{
			Owner:    user.Email,
			Balance:  100,
			Currency: "EUR",
		}

		_ = createTestAccount(suite.T(), arg)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testStore.ListAccounts(suite.ctx, &arg)
	require.NoError(suite.T(), err)
	require.Len(suite.T(), accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(suite.T(), account)
	}
}

func createTestAccount(t *testing.T, arg CreateAccountParams) *Account {
	account, err := testStore.CreateAccount(context.Background(), &arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
