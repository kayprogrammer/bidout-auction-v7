package managers

import (
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
	"github.com/kayprogrammer/bidout-auction-v7/models"

)

type SiteDetailManager struct {
	*BaseManager
}

func NewSiteDetailManager(model interface{}) *SiteDetailManager {
	return &SiteDetailManager{
		BaseManager: NewBaseManager(model),
	}
}

func (m *SiteDetailManager) Get(db *gorm.DB) (interface{}, error) {
	var siteDetail interface{}
	if err := db.First(&siteDetail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new SiteDetail record if not found
			siteDetail = m.Create(db, make(map[string]interface{}))
		} else {
			return nil, err
		}
	}
	return siteDetail, nil
}
