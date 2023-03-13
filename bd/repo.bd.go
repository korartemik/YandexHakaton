package bd

import (
	"context"
	"fmt"

	"awesomeProject1/model"
	"github.com/yandex-cloud/ydb-go-sdk/table"
)

type Repository interface {
	SaveAnswer(ctx context.Context, user *model.User) error
	GetAnswer(ctx context.Context, name string) (*model.User, error)
}

var _ Repository = &repository{}

type repository struct {
}

func NewRepository() (Repository, error) {
	return &repository{}, nil
}

func missingField(entity, field string) error {
	return fmt.Errorf("missing required %s field: %s", entity, field)
}

type entityReader string

func (er entityReader) fieldString(r *table.Result, name string) ([]byte, error) {
	if !r.SeekItem(name) {
		return nil, missingField(string(er), name)
	}
	r.Unwrap()
	result := r.String()
	if r.Err() != nil {
		return nil, r.Err()
	}
	return result, nil
}

func (er entityReader) fieldUtf8(r *table.Result, name string) (string, error) {
	if !r.SeekItem(name) {
		return "", missingField(string(er), name)
	}
	r.Unwrap()
	result := r.UTF8()
	if r.Err() != nil {
		return "", r.Err()
	}
	return result, nil
}

func (er entityReader) fieldBool(r *table.Result, name string) (bool, error) {
	if !r.SeekItem(name) {
		return false, missingField(string(er), name)
	}
	r.Unwrap()
	result := r.Bool()
	if r.Err() != nil {
		return false, r.Err()
	}
	return result, nil
}
