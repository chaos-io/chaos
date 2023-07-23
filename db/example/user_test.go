package example

import (
	"context"
	"reflect"
	"testing"

	"github.com/chaos-io/chaos/db"
)

func TestCreateUserModel(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"CreateUserModel"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("CreateUserModel panic, %v", r)
				}
			}()

			if user := CreateUserModel(); user == nil {
				t.Errorf("CreateUserModel, the user is nil")
			}
		})
	}
}

func TestNewUserModel(t *testing.T) {
	tests := []struct {
		name string
		want *UserModel
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserModel_BatchDelete(t *testing.T) {
	type fields struct {
		DB *db.DB
	}
	type args struct {
		ctx context.Context
		ids []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserModel{
				DB: tt.fields.DB,
			}
			got, err := a.BatchDelete(tt.args.ctx, tt.args.ids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BatchDelete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserModel_BatchGet(t *testing.T) {
	type fields struct {
		DB *db.DB
	}
	type args struct {
		ctx context.Context
		ids []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserModel{
				DB: tt.fields.DB,
			}
			got, err := a.BatchGet(tt.args.ctx, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BatchGet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserModel_Delete(t *testing.T) {
	type fields struct {
		DB *db.DB
	}
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserModel{
				DB: tt.fields.DB,
			}
			got, err := a.Delete(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserModel_Get(t *testing.T) {
	type fields struct {
		DB *db.DB
	}
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserModel{
				DB: tt.fields.DB,
			}
			got, err := a.Get(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserModel_GetIds(t *testing.T) {
	type fields struct {
		DB *db.DB
	}
	type args struct {
		ctx       context.Context
		condition []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserModel{
				DB: tt.fields.DB,
			}
			got, err := a.GetIds(tt.args.ctx, tt.args.condition...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIds() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserModel_Insert(t *testing.T) {
	type fields struct {
		DB *db.DB
	}
	type args struct {
		ctx   context.Context
		users []*User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserModel{
				DB: tt.fields.DB,
			}
			got, err := a.Insert(tt.args.ctx, tt.args.users...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Insert() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserModel_Query(t *testing.T) {
	type fields struct {
		DB *db.DB
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserModel{
				DB: tt.fields.DB,
			}
			got, err := a.Query(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserModel_Update(t *testing.T) {
	type fields struct {
		DB *db.DB
	}
	type args struct {
		ctx  context.Context
		user *User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &UserModel{
				DB: tt.fields.DB,
			}
			got, err := a.Update(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}
