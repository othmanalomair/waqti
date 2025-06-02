package services

import (
	"encoding/json"
	"time"
	"waqti/internal/models"

	"github.com/google/uuid"
)

type CreatorService struct {
	creators []models.Creator
}

func NewCreatorService() *CreatorService {
	// Generate fixed UUIDs for demo data consistency
	creatorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	creators := []models.Creator{
		{
			ID:        creatorID,
			Name:      "Ahmed Al-Kuwaiti",
			NameAr:    "أحمد الكويتي",
			Username:  "ahmed",
			Email:     "ahmed@example.com",
			Avatar:    "/static/avatars/ahmed.jpg",
			Plan:      "Free",
			PlanAr:    "مجاني",
			IsActive:  true,
			CreatedAt: time.Now().AddDate(0, -2, 0),
			UpdatedAt: time.Now(),
		},
	}

	return &CreatorService{
		creators: creators,
	}
}

func (s *CreatorService) GetCreatorByID(id uuid.UUID) (*models.Creator, error) {
	for _, creator := range s.creators {
		if creator.ID == id {
			return &creator, nil
		}
	}
	return nil, nil
}

func (s *CreatorService) GetCreatorByUsername(username string) (*models.Creator, error) {
	for _, creator := range s.creators {
		if creator.Username == username {
			return &creator, nil
		}
	}
	return nil, nil
}

func (s *CreatorService) ToJSON(creator *models.Creator) string {
	data, _ := json.Marshal(creator)
	return string(data)
}

func (s *CreatorService) GetCreator(id uuid.UUID) *models.Creator {
	creator, _ := s.GetCreatorByID(id)
	if creator != nil {
		return creator
	}

	// Return dummy data if not found
	return &models.Creator{
		ID:        id,
		Name:      "John Doe",
		NameAr:    "جون دو",
		Username:  "johndoe",
		Email:     "john@example.com",
		Avatar:    "/static/images/avatar.jpg",
		Plan:      "Pro",
		PlanAr:    "احترافي",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
