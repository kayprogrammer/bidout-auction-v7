package managers

import (
	"time"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
	"github.com/mitchellh/mapstructure"
)

type BaseManager struct {
	model interface{}
}

func NewBaseManager(model interface{}) *BaseManager {
	return &BaseManager{model: model}
}

func (m *BaseManager) GetAll(db *gorm.DB) ([]interface{}, error) {
	var result []interface{}
	if err := db.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (m *BaseManager) GetAllIDs(db *gorm.DB) ([]uuid.UUID, error) {
	var result []struct{ ID uuid.UUID }
	if err := db.Model(m.model).Pluck("id", &result).Error; err != nil {
		return nil, err
	}
	var ids []uuid.UUID
	for _, item := range result {
		ids = append(ids, item.ID)
	}
	return ids, nil
}

func (m *BaseManager) GetByID(db *gorm.DB, id uuid.UUID) (interface{}, error) {
	var result interface{}
	if err := db.First(result, id).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (m *BaseManager) Create(db *gorm.DB, objIn map[string]interface{}) (interface{}, error) {
	objIn["created_at"] = time.Now().UTC()
	objIn["updated_at"] = objIn["created_at"]
	obj := m.model
	if err := mapstructure.Decode(objIn, obj); err != nil {
		return nil, err
	}
	if err := db.Create(obj).Error; err != nil {
		return nil, err
	}
	return obj, nil
}

func (m *BaseManager) BulkCreate(db *gorm.DB, objIn []map[string]interface{}) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	for _, item := range objIn {
		item["created_at"] = time.Now().UTC()
		item["updated_at"] = item["created_at"]
		obj := m.model
		if err := mapstructure.Decode(item, obj); err != nil {
			return nil, err
		}
		if err := db.Create(obj).Error; err != nil {
			return nil, err
		}
		ids = append(ids, obj.ID)
	}
	return ids, nil
}

func (m *BaseManager) Update(db *gorm.DB, dbObj interface{}, objIn map[string]interface{}) (interface{}, error) {
	if dbObj == nil {
		return nil, nil
	}
	for attr, value := range objIn {
		if err := db.Model(dbObj).Update(attr, value).Error; err != nil {
			return nil, err
		}
	}
	dbObj.(BaseModel).UpdatedAt = time.Now().UTC()
	if err := db.Save(dbObj).Error; err != nil {
		return nil, err
	}
	return dbObj, nil
}

func (m *BaseManager) Delete(db *gorm.DB, dbObj interface{}) error {
	if dbObj != nil {
		if err := db.Delete(dbObj).Error; err != nil {
			return err
		}
	}
	return nil
}

func (m *BaseManager) DeleteByID(db *gorm.DB, id uuid.UUID) error {
	if err := db.Delete(m.model, id).Error; err != nil {
		return err
	}
	return nil
}

func (m *BaseManager) DeleteAll(db *gorm.DB) error {
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(m.model).Error; err != nil {
		return err
	}
	return nil
}
