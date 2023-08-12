package models

import (
	"fmt"
	"log"
	"time"

	"github.com/gosimple/slug"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/kayprogrammer/bidout-auction-v7/utils"
)

// CATEGORY
type Category struct {
	BaseModel
	Name				string				`json:"name" gorm:"not null" example:"Category"`
	Slug				*string				`json:"slug" gorm:"not null;unique" example:"category_slug"`
}

// Function to retrieve a category by slug
func getCategoryBySlug(db *gorm.DB, slug *string) Category {
	var category Category
	result := db.Where("slug = ?", slug).First(&category)
	err := result.Error 
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return Category{}
		}
		log.Println(err)
	}
	return category
}

func (c *Category) BeforeSave(tx *gorm.DB) (err error) {
	// Check if the Name field has changed
	if c.ID != uuid.Nil {
		var oldCategory Category
		if result := tx.First(&oldCategory, "id = ?", c.ID); result.Error == nil {
			// Compare the old Name with the new Name
			if oldCategory.Name != c.Name {
				// Generate new slug based on the updated Name
				createdSlug := slug.Make(c.Name)
				c.Slug = &createdSlug
			}
		}
	}

	// Generate unique slug
	if c.Slug == nil {
		createdSlug := slug.Make(c.Name)
		c.Slug = &createdSlug
	}

	for {
		slugExists := getCategoryBySlug(tx, c.Slug)
		if slugExists.ID == c.ID || slugExists.ID == uuid.Nil {
			// Unique slug found, break the loop
			break
		}

		// Slug exists, generate a new random string
		randomStr := utils.GetRandomString(4)
		newSlug := fmt.Sprintf("%s-%s", *c.Slug, randomStr)
		c.Slug = &newSlug
	}	
	return
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
	CategoryObj			*Category			`json:"-" gorm:"foreignKey:CategoryId;constraint:OnDelete:SET NULL;null"`
	Category			string				`json:"category" gorm:"-"`

	Active				bool				`json:"active" gorm:"default:true"`
	Price				decimal.Decimal		`json:"price" gorm:"default:0"`
	HighestBid			decimal.Decimal		`json:"highest_bid" gorm:"-"`
	BidsCount			int					`json:"bids_count" gorm:"-"`
	ClosingDate			time.Time			`json:"closing_date" gorm:"not null"`

	ImageId				uuid.UUID			`json:"-" gorm:"not null"`
	ImageObj			File				`json:"-" gorm:"foreignKey:ImageId;constraint:OnDelete:SET NULL;null;"`
	Image				string				`json:"image" gorm:"-"`

	Watchlist			bool				`json:"watchlist" gorm:"-"`
	TimeLeftSecs		int64				`json:"time_left_seconds" gorm:"-"`

	Bids				[]Bid				`json:"-"`
}

// Function to retrieve a listing by slug
func getListingBySlug(db *gorm.DB, slug *string) Listing {
	var listing Listing
	result := db.Where("slug = ?", slug).First(&listing)
	err := result.Error 
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return Listing{}
		}
		log.Println(err)
	}
	return listing
}

func (listing *Listing) BeforeSave(tx *gorm.DB) (err error) {
    listing.Price = listing.Price.Round(2)
    listing.HighestBid = listing.HighestBid.Round(2)

	// Check if the Name field has changed
	if listing.ID != uuid.Nil {
		var oldListing Listing
		if result := tx.First(&oldListing, "id = ?", listing.ID); result.Error == nil {
			// Compare the old Name with the new Name
			if oldListing.Name != listing.Name {
				// Generate new slug based on the updated Name
				createdSlug := slug.Make(listing.Name)
				listing.Slug = &createdSlug
			}
		}
	}

	// Generate unique slug
	if listing.Slug == nil {
		createdSlug := slug.Make(listing.Name)
		listing.Slug = &createdSlug
	}

	for {
		slugExists := getListingBySlug(tx, listing.Slug)
		if slugExists.ID == listing.ID || slugExists.ID == uuid.Nil {
			// Unique slug found, break the loop
			break
		}

		// Slug exists, generate a new random string
		randomStr := utils.GetRandomString(4)
		newSlug := fmt.Sprintf("%s-%s", *listing.Slug, randomStr)
		listing.Slug = &newSlug
	}
	return
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

func (listing Listing) GetHighestBid() decimal.Decimal {
	bids := listing.Bids
	bidsLength := len(bids)
	highestAmount := decimal.NewFromFloat(0.00)
	if bidsLength > 0 {
		highestAmount = bids[0].Amount
		for _, bid := range bids {
			if bid.Amount.GreaterThan(highestAmount) {
				highestAmount = bid.Amount
			}
		}
	}
	return highestAmount 
}

func (listing Listing) Init(db *gorm.DB) Listing {
	listing.Auctioneer.Name = listing.AuctioneerObj.FullName()

	avatarId := listing.AuctioneerObj.AvatarId
	if avatarId != nil {
		avatar := File{}
		db.Find(&avatar,"id = ?", avatarId)
		url := utils.GenerateFileUrl(avatarId.String(), "avatars", avatar.ResourceType)
		listing.Auctioneer.Avatar = &url
	}

	// Get Listing Image
	imageId := listing.ImageId
	image := File{}
	db.Find(&image,"id = ?", imageId)
	url := utils.GenerateFileUrl(imageId.String(), "listings", image.ResourceType)
	listing.Image = url

	listing.Price = listing.Price.Round(2)
	listing.HighestBid = listing.HighestBid.Round(2)
	if listing.CategoryId != nil {
		listing.Category = listing.CategoryObj.Name
	} else {
		listing.Category = "Other"
	}
	if listing.Active && (listing.TimeLeftSeconds() > 0) {
		listing.Active = true
	} else {
		listing.Active = false
	}
	listing.ClosingDate = listing.ClosingDate.UTC()
	listing.TimeLeftSecs = listing.TimeLeftSeconds()

	listing.BidsCount = len(listing.Bids)
	listing.HighestBid = listing.GetHighestBid()
	return listing
}

func (listing Listing) GetImageUploadData(db *gorm.DB) utils.SignatureFormat {
	imageId := listing.ImageId
	uploadData := utils.GenerateFileSignature(imageId.String(), "listings")
	return uploadData
}
// ---------------------------------------------------------------

// BID
type Bid struct {
	BaseModel
	UserId				uuid.UUID			`json:"-" gorm:"not null;index:,unique,composite:user_id_listing_id"`
	UserObj				User				`json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;not null;"`
	User				ShortUserData		`json:"user" gorm:"-"`

	ListingId			uuid.UUID			`json:"-" gorm:"not null;index:,unique,composite:user_id_listing_id;index:,unique,composite:listing_id_amount"`
	Listing				Listing				`json:"-" gorm:"foreignKey:ListingId;constraint:OnDelete:CASCADE;not null;"`
	Amount				decimal.Decimal		`json:"amount" gorm:"not null;index:,unique,composite:listing_id_amount"`
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
	UserId				*uuid.UUID			`json:"-" gorm:"null;index:,unique,composite:user_id_listing_id"`
	User				*User				`json:"-" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;null;"`

	ListingId			uuid.UUID			`json:"-" gorm:"not null;index:,unique,composite:user_id_listing_id;index:,unique,composite:listing_id_guest_user_id"`
	Listing				Listing				`json:"-" gorm:"foreignKey:ListingId;constraint:OnDelete:CASCADE;not null;"`

	GuestUserId			*uuid.UUID			`json:"-" gorm:"null;index:,unique,composite:listing_id_guest_user_id"`
	GuestUser			*GuestUser			`json:"-" gorm:"foreignKey:GuestUserId;constraint:OnDelete:CASCADE;null;"`
}

// --------------------------------------------------------------