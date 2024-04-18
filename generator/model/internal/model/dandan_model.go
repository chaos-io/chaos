package model

import (
	"context"
	"sync"

	"github.com/chaos-io/chaos/db"
	"github.com/chaos-io/chaos/generator/model/internal"
	"github.com/chaos-io/chaos/logs"
)

var (
	dandanModel     *DandanModel
	dandanModelOnce sync.Once
)

type DandanModel struct {
	DB *db.DB
}

func GetDandanModel() *DandanModel {
	dandanModelOnce.Do(func() {
		dandanModel = NewDandanModel()
	})

	return dandanModel
}

func NewDandanModel() *DandanModel {
	m := &DandanModel{DB: initDB()}
	if !m.DB.Config.DisableAutoMigrate || !d.Migrator().HasTable(&internal.Dandan{}) {
		if err := d.AutoMigrate(&internal.Dandan{}); err != nil {
			logs.Error("Init DandanModel model err: ", err)
			panic(err)
		}
	}

	return m
}

func (m *DandanModel) Create(ctx context.Context, dandan *internal.Dandan) (int64, error) {
	result := m.DB.WithContext(ctx).Create(dandan)
	return result.RowsAffected, result.Error
}

func (m *DandanModel) Get(ctx context.Context, id string) (*internal.Dandan, error) {
	dandan := &internal.Dandan{}
	return dandan, m.DB.WithContext(ctx).First(dandan, "id = ?", id).Error
}

func (m *DandanModel) Delete(ctx context.Context, id string) (int64, error) {
	result := m.DB.WithContext(ctx).Where("id = ?", id).Delete(&internal.Dandan{})
	return result.RowsAffected, result.Error
}

func (m *DandanModel) Update(ctx context.Context, dandan *internal.Dandan) (int64, error) {
	result := m.DB.WithContext(ctx).Updates(dandan)
	return result.RowsAffected, result.Error
}

func (m *DandanModel) List(ctx context.Context, filter string, condition ...string) ([]*internal.Dandan, error) {
	var dandan []*internal.Dandan

	tx := m.DB.WithContext(ctx)
	// todo add condition

	return dandan, tx.Find(&dandan).Error
}

func (m *DandanModel) BatchCreate(ctx context.Context, dandan ...*internal.Dandan) (int64, error) {
	result := m.DB.WithContext(ctx).CreateInBatches(dandan, len(dandan))
	return result.RowsAffected, result.Error
}

func (m *DandanModel) BatchGet(ctx context.Context, ids ...string) ([]*internal.Dandan, error) {
	var dandan []*internal.Dandan
	return dandan, m.DB.WithContext(ctx).Find(&dandan, "id = ?", ids).Error
}

func (m *DandanModel) BatchDelete(ctx context.Context, ids ...string) (int64, error) {
	result := m.DB.WithContext(ctx).Where("id = ?", ids).Delete(&internal.Dandan{})
	return result.RowsAffected, result.Error
}

func (m *DandanModel) BatchUpdate(ctx context.Context, dandan []*internal.Dandan) (int64, error) {
	result := m.DB.WithContext(ctx).Updates(dandan)
	return result.RowsAffected, result.Error
}
