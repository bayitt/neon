package services

import (
	"errors"
	"fmt"
	"maps"
	"math/rand"
	"neon/dto"
	"neon/models"
	"neon/utilities"
	"net/http"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReviewService struct {
	DB *gorm.DB
}

func (rs *ReviewService) FindUnique(
	field string,
	value string,
	preload bool,
) (models.Review, error) {
	var review models.Review
	var result *gorm.DB

	if preload {
		result = rs.DB.Preload(clause.Associations).
			Where(fmt.Sprintf("%s = ?", field), value).
			First(&review)
	} else {
		result = rs.DB.Where(fmt.Sprintf("%s = ?", field), value).First(&review)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Review{}, utilities.ThrowError(
			http.StatusNotFound,
			"REVIEW_002",
			fmt.Sprintf("review with field %s and value %s does not exist", field, value),
		)
	}

	review.Excerpt = review.Content[0:200]
	return review, nil
}

func (rs *ReviewService) Find(
	offset uint,
	count uint,
	where map[string]uint,
) ([]models.Review, error) {
	var reviews []models.Review
	var result *gorm.DB

	if len(where) > 0 {
		mapKeys := slices.Collect(maps.Keys(where))
		mapValues := slices.Collect(maps.Values(where))
		result = rs.DB.
			Where(fmt.Sprintf("%s = ?", mapKeys[0]), mapValues[0]).
			Order("created_at desc").
			Offset(int(offset)).
			Limit(int(count)).
			Preload(clause.Associations).
			Find(&reviews)
	} else {
		result = rs.DB.
			Order("created_at desc").
			Offset(int(offset)).
			Limit(int(count)).
			Preload(clause.Associations).
			Find(&reviews)
	}

	if result.Error != nil {
		return []models.Review{}, utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			result.Error.Error(),
		)
	}

	parsedReviews := []models.Review{}
	for _, review := range reviews {
		review.Excerpt = review.Content[0:200]
		parsedReviews = append(parsedReviews, review)
	}
	return parsedReviews, nil
}

func (rs *ReviewService) FindCategoryReviews(
	category models.Category,
	offset uint,
	count uint,
) ([]models.Review, error) {
	var categoryReviews, reviewIds, reviews = []models.CategoryReview{},
		[]uint{},
		[]models.Review{}

	result := rs.DB.
		Where("category_id = ?", category.ID).
		Order("created_at desc").
		Offset(int(offset)).
		Limit(int(count)).
		Find(&categoryReviews)

	if result.Error != nil {
		return []models.Review{}, utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			result.Error.Error(),
		)
	}

	for _, categoryReview := range categoryReviews {
		reviewIds = append(reviewIds, categoryReview.ReviewID)
	}

	reviewResult := rs.DB.
		Where(reviewIds).
		Order("created_at desc").
		Preload(clause.Associations).
		Find(&reviews)
	if reviewResult.Error != nil {
		return []models.Review{}, utilities.ThrowError(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			reviewResult.Error.Error(),
		)
	}

	parsedReviews := []models.Review{}
	for _, review := range reviews {
		review.Excerpt = review.Content[0:200]
		parsedReviews = append(parsedReviews, review)
	}
	return parsedReviews, nil
}

func (rs *ReviewService) FindCategoriesReviews(categories []models.Category) ([]models.Review, error) {
	var categoryIds = []uint{}

	for _, category := range categories {
		categoryIds = append(categoryIds, category.ID)
	}

	var totalCategoriesReviews int64
	rs.DB.Model(models.CategoryReview{}).Where("category_id IN ?", categoryIds).Count(&totalCategoriesReviews)

	rand.Seed(time.Now().UnixNano())
	randomOffset := rand.Int63n(totalCategoriesReviews - 3)

	var categoryReviews []models.CategoryReview
	result := rs.DB.Offset(int(randomOffset)).Limit(3).Find(&categoryReviews)
}

func (rs *ReviewService) Count(where map[string]uint) uint {
	var totalReviews int64
	if len(where) > 0 {
		mapKeys := slices.Collect(maps.Keys(where))
		mapValues := slices.Collect(maps.Keys(where))
		rs.DB.Model(models.Review{}).
			Where(fmt.Sprintf("%s = ?", mapKeys[0]), mapValues[0]).
			Count(&totalReviews)
	} else {
		rs.DB.Model(models.Review{}).Count(&totalReviews)
	}

	return uint(totalReviews)
}

func (rs *ReviewService) CountCategoryReviews(category models.Category) uint {
	totalCategoryReviews := rs.DB.Model(&category).Association("Reviews").Count()
	return uint(totalCategoryReviews)
}

func (rs *ReviewService) Create(crDto *dto.CreateReviewDto) (models.Review, error) {
	var review models.Review
	err := rs.DB.Transaction(func(tx *gorm.DB) error {
		slug := "/" + strings.Replace(
			crDto.Title,
			" ",
			"-",
			-1,
		) + "-" + utilities.GenerateRandomString(
			4,
		)
		var status uint

		if crDto.Status {
			status = 1
		} else {
			status = 0
		}

		review = models.Review{
			Uuid:    crDto.Uuid,
			Title:   crDto.Title,
			Slug:    slug,
			Author:  crDto.Author,
			Content: crDto.Content,
			Image:   crDto.Image,
			Status:  status,
		}

		if crDto.Series != nil {
			review.SeriesID = &crDto.Series.ID
		}
		result := tx.Create(&review)

		if result.Error != nil {
			return utilities.ThrowError(
				http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				result.Error.Error(),
			)
		}

		var reviewCategories []models.CategoryReview
		for i := 0; i < len(crDto.Categories); i++ {
			reviewCategories = append(
				reviewCategories,
				models.CategoryReview{CategoryID: crDto.Categories[i].ID, ReviewID: review.ID},
			)
		}

		associateResult := tx.Create(&reviewCategories)
		if associateResult.Error != nil {
			return utilities.ThrowError(
				http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				associateResult.Error.Error(),
			)
		}

		return nil
	})

	if err != nil {
		return models.Review{}, err
	}

	var categories = []*models.Category{}
	for i := 0; i < len(crDto.Categories); i++ {
		categories = append(categories, &(crDto.Categories[i]))
	}
	review.Categories = categories
	review.Series = crDto.Series

	return review, nil
}

func (rs *ReviewService) Update(
	review models.Review,
	urDto *dto.UpdateReviewDto,
) (models.Review, error) {
	updateStringField := func(initialValue string, newValuePointer *string) string {
		if newValuePointer != nil {
			return *newValuePointer
		}

		return initialValue
	}

	err := rs.DB.Transaction(func(tx *gorm.DB) error {
		if urDto.Title != nil && (review.Title != *urDto.Title) {
			review.Slug = "/" + strings.Replace(
				*urDto.Title,
				" ",
				"-",
				-1,
			) + "-" + utilities.GenerateRandomString(
				4,
			)
		}

		if urDto.Image != nil {
			review.Image = urDto.Image
		}

		if urDto.Series != nil {
			review.SeriesID = &(*urDto.Series).ID
		}

		if urDto.Status != nil {
			if *urDto.Status {
				review.Status = 1
			} else {
				review.Status = 0
			}
		}

		review.Title = updateStringField(review.Title, urDto.Title)
		review.Author = updateStringField(review.Author, urDto.Author)
		review.Content = updateStringField(review.Content, urDto.Content)

		result := tx.Save(&review)
		if result.Error != nil {
			return utilities.ThrowError(
				http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				result.Error.Error(),
			)
		}

		if urDto.Categories != nil {
			var reviewCategories []models.Category
			tx.Model(review).Association("Categories").Find(&reviewCategories)

			var isCategoriesChanged bool

			for _, category := range *urDto.Categories {
				index := slices.IndexFunc(
					reviewCategories,
					func(cat models.Category) bool { return cat.ID == category.ID },
				)

				if index == -1 {
					isCategoriesChanged = true
					break
				}
			}

			if isCategoriesChanged {
				tx.Model(&review).Association("Categories").Clear()
				var reviewCategories []models.CategoryReview
				for i := 0; i < len(*urDto.Categories); i++ {
					reviewCategories = append(
						reviewCategories,
						models.CategoryReview{
							CategoryID: (*urDto.Categories)[i].ID,
							ReviewID:   review.ID,
						},
					)
				}

				associateResult := tx.Create(&reviewCategories)
				if associateResult.Error != nil {
					return utilities.ThrowError(
						http.StatusInternalServerError,
						"INTERNAL_SERVER_ERROR",
						associateResult.Error.Error(),
					)
				}
			}
		}

		return nil
	})

	if err != nil {
		return models.Review{}, err
	}

	var series, categories = models.Series{}, []*models.Category{}
	rs.DB.Model(review).Association("Series").Find(&series)
	rs.DB.Model(review).Association("Categories").Find(&categories)

	if len(series.Name) > 0 {
		review.Series = &series
	}
	review.Categories = categories
	review.Excerpt = review.Content[0:200]
	return review, nil
}
