package services

import (
	"encoding/json"
	"time"
	"waqti/internal/models"
)

type CreatorService struct {
	// In real app, this would be database connection
	creators []models.Creator
}

func NewCreatorService() *CreatorService {
	// Dummy data - replace with database later
	creators := []models.Creator{
		{
			ID:        1,
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

func (s *CreatorService) GetCreatorByID(id int) (*models.Creator, error) {
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

// Add the missing GetCreator method
func (s *CreatorService) GetCreator(id int) *models.Creator {
	// Return dummy data for now
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
