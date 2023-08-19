package initials

import (
	"time"
	"log"
	"math/rand"
	"errors"
	"os"
	"path/filepath"

	"github.com/gosimple/slug"


	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"github.com/shopspring/decimal"
)
var	truth = true
func createSuperUser(db *gorm.DB) models.User {
	user := models.User{
		FirstName:       "Test",
		LastName:        "Admin",
		Email:           "testadmin@email.com",
		Password:        "testadmin",
		IsSuperuser:     &truth,
		IsStaff:         &truth,
		IsEmailVerified: &truth,
	}
	db.Where(models.User{Email: user.Email}).FirstOrCreate(&user)
	return user
}

func createAutioneer(db *gorm.DB) models.User {
	user := models.User{
		FirstName:       "Test",
		LastName:        "Auctioneer",
		Email:           "testauctioneer@email.com",
		Password:        "testauctioneer",
		IsEmailVerified: &truth,
	}
	db.Where(models.User{Email: user.Email}).FirstOrCreate(&user)
	return user
}

func createReviewer(db *gorm.DB) models.User {
	user := models.User{
		FirstName:       "Test",
		LastName:        "Reviewer",
		Email:           "testreviewer@email.com",
		Password:        "testreviewer",
		IsEmailVerified: &truth,
	}
	db.Where(models.User{Email: user.Email}).FirstOrCreate(&user)
	return user
}

func createReviews(db *gorm.DB, reviewerId uuid.UUID) {
	review := &models.Review{
		ReviewerId: reviewerId,
		Show:       true,
		Text:       "Maecenas vitae porttitor neque, ac porttitor nunc. Duis venenatis lacinia libero. Nam nec augue ut nunc vulputate tincidunt at suscipit nunc.",
	}
	reviews := []*models.Review{review, review, review}

	// Create Reviews if none in the database
	result := db.Take(&models.Review{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		db.Create(&reviews)
	}
}

func createCategories(db *gorm.DB) []models.Category {
	technology := "technology"
	accessories := "accessories"
	automobile := "automobile"
	fashion := "fashion"
	categories := []models.Category{
		{Name: "Technology", Slug: &technology},
		{Name: "Accessories", Slug: &accessories},
		{Name: "Automobile", Slug: &automobile},
		{Name: "Fashion", Slug: &fashion},
	}

	// Create Categories if none in the database
	existingCategories := []models.Category{}
	db.Find(&existingCategories)
	if len(existingCategories) == 0 {
		db.Create(&categories)
		return categories
	}
	return existingCategories
}

func createListingImages(db *gorm.DB) []models.File {
	image := models.File{
		ResourceType: "image/png",
	}
	images := []models.File{image, image, image, image, image, image}
	db.Create(images)
	return images
}

func randomSelect(list []models.Category) (models.Category, error) {
	// Create a custom random source seeded with the current time
	source := rand.NewSource(time.Now().UnixNano())

	// Create a new random number generator from the custom source
	rng := rand.New(source)

	// Generate a random index between 0 and the length of the list - 1
	randomIndex := rng.Intn(len(list))

	// Return the randomly selected person from the list
	return list[randomIndex], nil
}

func uploadListingImages(images []models.File) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd: ", err)
	}

	// Create a directory path for the "images" subdirectory
	imagesDir := filepath.Join(currentDir, "/initials/images")

	fileList := []string{}
	err = filepath.Walk(imagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal("File Path", err)
	}
	
	// Upload each file with the corresponding public ID
	for i, filePath := range fileList {
		if i >= len(images) {
			log.Println("Not enough public IDs specified for all files.")
			break
		}
		publicID := images[i].ID.String()
		file, err := os.Open(filePath)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
		utils.UploadImage(file, publicID, "listings")
	}
}

func createListings(db *gorm.DB, auctioneerId uuid.UUID, categories []models.Category) {
	// Create Listings if none in the database
	result := db.Take(&models.Listing{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		images := createListingImages(db)
		listings := []models.Listing{}
		listingsMapping := ListingsMapping() 

		for i := range listingsMapping {
			slugStr := slug.Make(listingsMapping[i].Name)
			randomCategory, err := randomSelect(categories)
			if err != nil {
				log.Println("Error:", err)
				return
			}

			currentDateTime := time.Now()

			// Add several days to the current date and time
			daysToAdd := 7 + i
			severalDaysLater := currentDateTime.AddDate(0, 0, daysToAdd)

			listing := models.Listing{
				AuctioneerId: auctioneerId, 
				Name: listingsMapping[i].Name, 
				Slug: &slugStr, 
				Desc: "Korem ipsum dolor amet, consectetur adipiscing elit. Maece nas in pulvinar neque. Nulla finibus lobortis pulvinar. Donec a consectetur nulla.", 
				CategoryId: &randomCategory.ID, 
				Price: decimal.NewFromInt(int64((i+1) * 1000)).Round(2),
				ClosingDate: severalDaysLater, 
				ImageId: images[i].ID,
			}
			listings = append(listings, listing)
		}
		
		db.Create(&listings)

		// Upload Images
		log.Println("Uploading Images to cloudinary....")
		uploadListingImages(images)
		log.Println("Images uploaded....")

	}
}

func CreateInitialData(db *gorm.DB) {
	log.Println("Creating Initial Data....")
	createSuperUser(db)
	auctioneer := createAutioneer(db)
	reviewer := createReviewer(db)
	createReviews(db, reviewer.ID)
	categories := createCategories(db)
	createListings(db, auctioneer.ID, categories)
	log.Println("Initial Data Created....")
}
