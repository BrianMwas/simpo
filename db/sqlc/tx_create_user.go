package db

import "context"

// CreateUserTxParams transaction request body with a callback which will be returned and used to handle
// an action only when the CreateUserTX passes
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// CreateUserTxResult a transaction response to return a user.
type CreateUserTxResult struct {
	User
}

// CreateUserTX creates a user with a rollback if the process fails
func (store *SQLStore) CreateUserTX(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})
	return result, err
}
