package managers

import (
	// "github.com/satori/go.uuid"
	"gorm.io/gorm"
	// "github.com/kayprogrammer/bidout-auction-v7/models"

)

type SiteDetailStruct struct {
	*BaseStruct
}

func SiteDetailManager(model interface{}) *SiteDetailStruct {
	return &SiteDetailStruct{
		BaseStruct: BaseManager(model),
	}
}

func (m *SiteDetailStruct) Get(db *gorm.DB) (interface{}) {
	var siteDetail interface{}
	if err := db.First(&siteDetail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a new SiteDetail record if not found
			siteDetail = m.Create(db, make(map[string]interface{}))
		}
	}
	return siteDetail
}
