package models

import (
	"fmt"
	"time"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
	"github.com/shopspring/decimal"
	"github.com/gosimple/slug"

	"github.com/kayprogrammer/bidout-auction-v7/utils"
)

// CATEGORY
type Category struct {
	BaseModel
	Name				string				`json:"name" gorm:"not null"`
	Slug				*string				`json:"slug" gorm:"not null;unique"`
}

// Function to retrieve a category by slug
func getCategoryBySlug(db *gorm.DB, slug *string) (*Category, error) {
	var category Category
	err := db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (c *Category) BeforeSave(tx *gorm.DB) error {
	// Generate unique slug
	createdSlug := slug.Make(c.Name)
	updatedSlug := c.Slug
	slug := updatedSlug
	if updatedSlug == nil {
		slug = &createdSlug
	}
	for {
		slugExists, err := getCategoryBySlug(tx, slug)
		if err != nil {
			return err
		}
		if slugExists.ID == c.ID || slugExists == nil {
			// Unique slug found, break the loop
			break
		}

		// Slug exists, generate a new random string
		randomStr := utils.GetRandomString(4)
		newSlug := fmt.Sprintf("%s-%s", createdSlug, randomStr)
		slug = &newSlug
	}
	c.Slug = slug
	return nil
}

// ---------------------------------------------------------------------------------

// LISTING
type Listing struct {
	BaseModel

	AuctioneerId		uuid.UUID			`json:"-" gorm:"not null"`
	AuctioneerObj		User				`json:"-" gorm:"foreignKey:AuctioneerId;constraint:OnDelete:CASCADE;not null"`
	Auctioneer			ShortUserData		`json:"auctioneer" gorm:"-"`
	
	Name 				string 				`json:"name" gorm:"type:varchar(70);not null"`
	Slug				*string				`json:"slug" gorm:"not null;unique"`
	Desc 				string 				`json:"desc" gorm:"not null"`

	CategoryId			*uuid.UUID			`json:"-" gorm:"null"`
	CategoryObj			*Category			`json:"-" gorm:"foreignKey:CategoryId;constraint:OnDelete:SET NULL;unique;null"`
	Category			*string				`json:"category" gorm:"-"`

	Active				bool				`json:"active" gorm:"default:true"`
	Price				decimal.Decimal		`json:"price" gorm:"default:0"`
	HighestBid			decimal.Decimal		`json:"highest_bid" gorm:"default:0.00"`
	BidsCount			uint64				`json:"bids_count" gorm:"default:0"`
	ClosingDate			time.Time			`json:"closing_date" gorm:"not null"`

	ImageId				uuid.UUID			`json:"-" gorm:"not null"`
	ImageObj			File				`gorm:"foreignKey:ImageId;constraint:OnDelete:SET NULL;null;"`
	Image				string				`json:"image" gorm:"-"`
}

// Function to retrieve a category by slug
func getListingBySlug(db *gorm.DB, slug *string) (*Listing, error) {
	var listing Listing
	err := db.Where("slug = ?", slug).First(&listing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &listing, nil
}

func (listing *Listing) BeforeSave(tx *gorm.DB) (err error) {
    listing.Price = listing.Price.Round(2)
    listing.HighestBid = listing.HighestBid.Round(2)

	// Generate unique slug
	createdSlug := slug.Make(listing.Name)
	updatedSlug := listing.Slug
	slug := updatedSlug
	if updatedSlug == nil {
		slug = &createdSlug
	}
	for {
		slugExists, err := getListingBySlug(tx, slug)
		if err != nil {
			return err
		}
		if slugExists.ID == listing.ID || slugExists == nil {
			// Unique slug found, break the loop
			break
		}

		// Slug exists, generate a new random string
		randomStr := utils.GetRandomString(4)
		newSlug := fmt.Sprintf("%s-%s", createdSlug, randomStr)
		slug = &newSlug
	}
	listing.Slug = slug
	return nil
}

func (listing Listing) TimeLeftSeconds() int64 {
	closingDate := listing.ClosingDate.UTC()
	currentDate := time.Now().UTC()
	remainingTimeInSeconds := int64(closingDate.Sub(currentDate).Seconds())
	return remainingTimeInSeconds
}

func (listing Listing) TimeLeft() int64 {
	if !listing.Active {
		return 0
	}
	return listing.TimeLeftSeconds()
}

func (listing Listing) Init(db *gorm.DB) Listing {
	user := User{}
	db.Find(&user,"id = ?", listing.AuctioneerId)
	name := user.FullName()
	listing.Auctioneer.Name = name

	avatarId := user.AvatarId
	if avatarId != nil {
		avatar := File{}
		db.Find(&avatar,"id = ?", avatarId)
		url := utils.GenerateFileUrl(avatarId.String(), "avatars", avatar.ResourceType)
		listing.Auctioneer.Avatar = &url
	}

	listing.Price = listing.Price.Round(2)
	listing.HighestBid = listing.HighestBid.Round(2)
	return listing
}

// ---------------------------------------------------------------

// BID
type Bid struct {
	BaseModel
	UserId				uuid.UUID			`json:"-" gorm:"not null;uniqueIndex:idx_user_id_listing_id"`
	UserObj				User				`gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;not null;"`
	User				ShortUserData		`json:"user" gorm:"-"`

	ListingId			uuid.UUID			`json:"-" gorm:"not null;uniqueIndex:idx_user_id_listing_id,uniqueIndex:idx_listing_id_amount"`
	Listing				Listing				`json:"-" gorm:"foreignKey:ListingId;constraint:OnDelete:CASCADE;not null;"`
	Amount				decimal.Decimal		`json:"amount" gorm:"not null;uniqueIndex:idx_listing_id_amount"`
}

func (bid *Bid) BeforeSave(tx *gorm.DB) (err error) {
    bid.Amount = bid.Amount.Round(2)
    return
}

func (bid Bid) Init(db *gorm.DB) Bid {
	user := User{}
	db.Find(&user,"id = ?", bid.UserId)
	name := user.FullName()
	bid.User.Name = name

	avatarId := user.AvatarId
	if avatarId != nil {
		avatar := File{}
		db.Find(&avatar,"id = ?", avatarId)
		url := utils.GenerateFileUrl(avatarId.String(), "avatars", avatar.ResourceType)
		bid.User.Avatar = &url
	}

	bid.Amount = bid.Amount.Round(2)
	return bid
}

// -------------------------------------------------------------------------

// WATCHLIST
type Watchlist struct {
	BaseModel
	UserId				*uuid.UUID			`json:"-" gorm:"uniqueIndex:idx_user_id_listing_id,where:user_id IS NOT NULL"`
	User				*User				`json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;null;"`

	ListingId			uuid.UUID			`json:"-" gorm:"not null;uniqueIndex:idx_user_id_listing_id,uniqueIndex:idx_listing_id_session_key"`
	Listing				Listing				`json:"-" gorm:"foreignKey:ListingId;constraint:OnDelete:CASCADE;not null;"`

	SessionKey			uuid.UUID			`json:"-" gorm:"uniqueIndex:idx_listing_id_session_key,where:session_key IS NOT NULL"`
}

// --------------------------------------------------------------