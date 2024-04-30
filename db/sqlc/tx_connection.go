package db

import (
	"context"
)

type AcceptConnectionTxParams struct {
	ID int32 `json:"id" binding:"required"`
}

func (store *SQLStore) AcceptConnectionTx(ctx context.Context, arg AcceptConnectionTxParams) (ConnectionInvitation, error) {
	var connInvitation ConnectionInvitation
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		conn, err := q.AcceptConnection(ctx, arg.ID)
		connInvitation = conn

		if err != nil {
			return err
		}

		err = q.AddConnection(ctx, AddConnectionParams{
			SenderID:    conn.SenderID,
			RecipientID: conn.RecipientID,
		})
		if err != nil {
			return err
		}
		return err
	})

	return connInvitation, err
}
