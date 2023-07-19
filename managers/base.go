package managers

import (
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type BaseStruct struct {
	model interface{}
}

func BaseManager(model interface{}) *BaseStruct {
	return &BaseStruct{model: model}
}

func (m *BaseStruct) GetAll(db *gorm.DB) ([]interface{}) {
	var result []interface{}
	db.Find(&result)
	return result
}

func (m *BaseStruct) GetAllIDs(db *gorm.DB) ([]uuid.UUID) {
	var result []struct{ ID uuid.UUID }
	db.Model(m.model).Pluck("id", &result)
	var ids []uuid.UUID
	for _, item := range result {
		ids = append(ids, item.ID)
	}
	return ids
}

func (m *BaseStruct) GetByID(db *gorm.DB, id uuid.UUID) (interface{}) {
	var result interface{}
	db.First(result, id)
	return result
}

func (m *BaseStruct) Create(db *gorm.DB, objIn interface{}) (interface{}) {
	db.Create(&objIn)
	return objIn
}

func (m *BaseStruct) BulkCreate(db *gorm.DB, objIn []map[string]interface{}) (bool) {
	db.Create(objIn)
	return true
}

func (m *BaseStruct) Update(db *gorm.DB, dbObj interface{}, objIn map[string]interface{}) (interface{}) {
	for attr, value := range objIn {
		db.Model(dbObj).Update(attr, value)
	}
	db.Save(dbObj)
	return dbObj
}

func (m *BaseStruct) Delete(db *gorm.DB, dbObj interface{}) error {
	if dbObj != nil {
		if err := db.Delete(dbObj).Error; err != nil {
			return err
		}
	}
	return nil
}

func (m *BaseStruct) DeleteByID(db *gorm.DB, id uuid.UUID) error {
	if err := db.Delete(m.model, id).Error; err != nil {
		return err
	}
	return nil
}

func (m *BaseStruct) DeleteAll(db *gorm.DB) error {
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(m.model).Error; err != nil {
		return err
	}
	return nil
}
