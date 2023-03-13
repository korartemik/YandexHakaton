package bd

import (
	"context"

	"awesomeProject1/model"
	"github.com/yandex-cloud/ydb-go-sdk"
	"github.com/yandex-cloud/ydb-go-sdk/table"
)

func (r *repository) SaveAnswer(ctx context.Context, user *model.User) error {
	const query = `
DECLARE $id AS string;
DECLARE $answer AS string;
UPSERT INTO user(id, answer) VALUES ($id, $answer);
`
	return r.execute(ctx, func(ctx context.Context, s *table.Session, txc *table.TransactionControl) (*table.Transaction, error) {
		tx, _, err := s.Execute(ctx, txc, query, table.NewQueryParameters(
			table.ValueParam("$id", ydb.StringValue([]byte(user.ID))),
			table.ValueParam("$name", ydb.StringValue([]byte(user.Answer))),
		))
		return tx, err
	})
}

func (r *repository) GetAnswer(ctx context.Context, id string) (*model.User, error) {
	const query = `
DECLARE $id AS string;
SELECT id, name, yandex_avatar_id FROM user WHERE id = $id;
`
	var user *model.User
	err := r.execute(ctx, func(ctx context.Context, s *table.Session, txc *table.TransactionControl) (*table.Transaction, error) {
		tx, res, err := s.Execute(ctx, txc, query, table.NewQueryParameters(
			table.ValueParam("$id", ydb.StringValue([]byte(id))),
		))
		if err != nil {
			return nil, err
		}
		defer res.Close()
		if !res.NextSet() || !res.NextRow() {
			return tx, nil
		}
		user = &model.User{}
		return tx, readAnswer(res, user)
	})
	return user, err
}

func readAnswer(res *table.Result, u *model.User) error {
	er := entityReader("answers")

	if id, err := er.fieldString(res, "id"); err != nil {
		return err
	} else {
		u.ID = string(id)
	}
	if answer, err := er.fieldString(res, "answer"); err != nil {
		return err
	} else {
		u.Answer = string(answer)
	}
	return nil
}
