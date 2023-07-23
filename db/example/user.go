package example

import (
	"context"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/chaos-io/chaos/db"
	"github.com/chaos-io/chaos/logs"
)

type User struct {
	Id        string
	Name      string
	CreatedAt time.Time
}

var userModel *UserModel
var userModelOnce sync.Once

type UserModel struct {
	DB *db.DB
}

func CreateUserModel() *UserModel {
	userModelOnce.Do(func() {
		userModel = NewUserModel()
	})

	return userModel
}

func NewUserModel() *UserModel {
	u := &UserModel{DB: InitDB()}
	if !u.DB.Config.DisableAutoMigrate || !d.Migrator().HasTable(&User{}) {
		if err := d.AutoMigrate(&User{}); err != nil {
			logs.Error("Init UserModel model err: ", err)
			panic(err)
		}
	}

	return u
}

func (a *UserModel) Insert(ctx context.Context, users ...*User) (int64, error) {
	usersLen := len(users)
	var executionResult *gorm.DB

	if usersLen == 0 {
		return 0, nil
	} else if usersLen == 1 {
		cb := users[0]
		if cb.CreatedAt.IsZero() {
			cb.CreatedAt = time.Now()
		}

		if len(cb.Id) == 0 {
			cb.Id = "1111"
			executionResult = a.DB.WithContext(ctx).Create(users[0])
		} else {
			executionResult = a.DB.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(users[0])
		}
	} else {
		executionResult = a.DB.WithContext(ctx).CreateInBatches(users, len(users))
	}

	return executionResult.RowsAffected, executionResult.Error
}

func (a *UserModel) Get(ctx context.Context, uid string) (*User, error) {
	ag := &User{}
	return ag, a.DB.WithContext(ctx).First(ag, "id = ?", uid).Error
}

func (a *UserModel) GetIds(ctx context.Context, condition ...string) ([]string, error) {
	var ids []string
	return ids, a.DB.WithContext(ctx).Model(&User{}).Pluck("id", &ids).Error
}

func (a *UserModel) BatchGet(ctx context.Context, ids []string) ([]*User, error) {
	var users []*User
	return users, a.DB.WithContext(ctx).Find(&users, ids).Error
}

func (a *UserModel) Query(ctx context.Context, name string) ([]*User, error) {
	var users []*User

	tx := a.DB.DB.WithContext(ctx)
	if len(name) > 0 {
		tx = tx.Where("name = ?", name)
	}

	return users, tx.Find(&users).Error
}

func (a *UserModel) Delete(ctx context.Context, uid string) (int64, error) {
	executionResult := a.DB.WithContext(ctx).Where("id = ?", uid).Delete(&User{})
	return executionResult.RowsAffected, executionResult.Error
}

func (a *UserModel) BatchDelete(ctx context.Context, ids ...string) (int64, error) {
	executionResult := a.DB.WithContext(ctx).Delete(&User{}, ids)
	return executionResult.RowsAffected, executionResult.Error
}

func (a *UserModel) Update(ctx context.Context, user *User) (int64, error) {
	executionResult := a.DB.WithContext(ctx).Updates(user)
	return executionResult.RowsAffected, executionResult.Error
}
