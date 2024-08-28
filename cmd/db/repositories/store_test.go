package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxTransferTestSuite struct {
	suite.Suite
	ctx context.Context
}

func TestTxTransferSuite(t *testing.T) {
	suite.Run(t, &TxTransferTestSuite{
		ctx: context.Background(),
	})
}

func (suite *TxTransferTestSuite) AfterTest(suiteName string, testName string) {
	_, err := testStore.ClearTransfersTable()
	require.NoError(suite.T(), err)

	_, err = testStore.ClearEntriesTable()
	require.NoError(suite.T(), err)

	_, err = testStore.ClearAccountsTable()
	require.NoError(suite.T(), err)

	_, err = testStore.ClearUsersTable()
	require.NoError(suite.T(), err)
}

func (suite *TxTransferTestSuite) TestTransferTx() {
	registerUserParam1 := RegisterUserParams{
		Email:          "Max@Mustermann.de",
		HashedPassword: "",
		FirstName:      "Max",
		LastName:       "Mustermann",
	}
	user1 := registerTestUser(suite.T(), &registerUserParam1)

	createAccountParam1 := CreateAccountParams{
		Owner:    user1.Email,
		Balance:  100,
		Currency: "EUR",
	}

	registerUserParam2 := RegisterUserParams{
		Email:          "Tom@Mustermann.de",
		HashedPassword: "",
		FirstName:      "Tom",
		LastName:       "Mustermann",
	}
	user2 := registerTestUser(suite.T(), &registerUserParam2)

	account1 := createTestAccount(suite.T(), createAccountParam1)

	createAccountParam2 := CreateAccountParams{
		Owner:    user2.Email,
		Balance:  200,
		Currency: "EUR",
	}
	account2 := createTestAccount(suite.T(), createAccountParam2)

	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(suite.T(), err)

		result := <-results
		require.NotEmpty(suite.T(), result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(suite.T(), transfer)
		require.Equal(suite.T(), account1.ID, transfer.FromAccountID)
		require.Equal(suite.T(), account2.ID, transfer.ToAccountID)
		require.Equal(suite.T(), amount, transfer.Amount)
		require.NotZero(suite.T(), transfer.ID)
		require.NotZero(suite.T(), transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(suite.T(), err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(suite.T(), fromEntry)
		require.Equal(suite.T(), account1.ID, fromEntry.AccountID)
		require.Equal(suite.T(), -amount, fromEntry.Amount)
		require.NotZero(suite.T(), fromEntry.ID)
		require.NotZero(suite.T(), fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(suite.T(), err)

		toEntry := result.ToEntry
		require.NotEmpty(suite.T(), toEntry)
		require.Equal(suite.T(), account2.ID, toEntry.AccountID)
		require.Equal(suite.T(), amount, toEntry.Amount)
		require.NotZero(suite.T(), toEntry.ID)
		require.NotZero(suite.T(), toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(suite.T(), err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(suite.T(), fromAccount)
		require.Equal(suite.T(), account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(suite.T(), toAccount)
		require.Equal(suite.T(), account2.ID, toAccount.ID)

		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(suite.T(), diff1, diff2)
		require.True(suite.T(), diff1 > 0)
		require.True(suite.T(), diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(suite.T(), k >= 1 && k <= n)
		require.NotContains(suite.T(), existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(suite.T(), err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(suite.T(), err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(suite.T(), account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(suite.T(), account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func (suite *TxTransferTestSuite) TestTransferTxDeadlock() {
	registerUserParam1 := RegisterUserParams{
		Email:          "Max@Mustermann.de",
		HashedPassword: "",
		FirstName:      "Max",
		LastName:       "Mustermann",
	}
	user1 := registerTestUser(suite.T(), &registerUserParam1)

	createAccountParam1 := CreateAccountParams{
		Owner:    user1.Email,
		Balance:  100,
		Currency: "EUR",
	}

	registerUserParam2 := RegisterUserParams{
		Email:          "Tom@Mustermann.de",
		HashedPassword: "",
		FirstName:      "Tom",
		LastName:       "Mustermann",
	}
	user2 := registerTestUser(suite.T(), &registerUserParam2)

	account1 := createTestAccount(suite.T(), createAccountParam1)

	createAccountParam2 := CreateAccountParams{
		Owner:    user2.Email,
		Balance:  200,
		Currency: "EUR",
	}
	account2 := createTestAccount(suite.T(), createAccountParam2)

	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(suite.T(), err)
	}

	// check the final updated balance
	updatedAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(suite.T(), err)

	updatedAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(suite.T(), err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(suite.T(), account1.Balance, updatedAccount1.Balance)
	require.Equal(suite.T(), account2.Balance, updatedAccount2.Balance)
}
