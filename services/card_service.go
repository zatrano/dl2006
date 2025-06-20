package services

import (
	"context"
	"errors"

	"davet.link/configs/databaseconfig"
	"davet.link/configs/logconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"
	"davet.link/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ICardService interface {
	GetAllCards(params queryparams.ListParams) (*queryparams.PaginatedResult, error)
	GetCardByID(ctx context.Context, id uint) (*models.Card, error)
	CreateCard(ctx context.Context, card *models.Card) error
	UpdateCard(ctx context.Context, id uint, card *models.Card) error
	DeleteCard(ctx context.Context, id uint) error
	GetCardsByUserID(userID uint) ([]models.Card, int64, error)
}

type CardService struct {
	repo repositories.ICardRepository
}

func NewCardService() ICardService {
	return &CardService{repo: repositories.NewCardRepository()}
}

func (s *CardService) GetAllCards(params queryparams.ListParams) (*queryparams.PaginatedResult, error) {
	cards, totalCount, err := s.repo.GetAllCards(params)
	if err != nil {
		logconfig.Log.Error("Kartlar alınamadı", zap.Error(err))
		return nil, errors.New("Kartlar getirilirken bir hata oluştu")
	}

	result := &queryparams.PaginatedResult{
		Data: cards,
		Meta: queryparams.PaginationMeta{
			CurrentPage: params.Page,
			PerPage:     params.PerPage,
			TotalItems:  totalCount,
			TotalPages:  queryparams.CalculateTotalPages(totalCount, params.PerPage),
		},
	}
	return result, nil
}

func (s *CardService) GetCardByID(ctx context.Context, id uint) (*models.Card, error) {
	return s.repo.GetCardByID(id)
}

func (s *CardService) CreateCard(ctx context.Context, card *models.Card) error {
	// Card ve ilişkili junction tabloları transaction ile ekle
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok || db == nil {
		db = databaseconfig.GetDB()
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.CreateCard(ctx, card); err != nil {
			return err
		}
		for _, cb := range card.CardBanks {
			if err := tx.Create(&cb).Error; err != nil {
				return err
			}
		}
		for _, csm := range card.CardSocialMedia {
			if err := tx.Create(&csm).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *CardService) UpdateCard(ctx context.Context, id uint, card *models.Card) error {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok || db == nil {
		db = databaseconfig.GetDB()
	}
	return db.Transaction(func(tx *gorm.DB) error {
		updateData := map[string]interface{}{
			"name":      card.Name,
			"slug":      card.Slug,
			"user_id":   card.UserID,
			"photo":     card.Photo,
			"telephone": card.Telephone,
			"email":     card.Email,
			"location":  card.Location,
			"website":   card.Website,
			"is_active": card.IsActive,
		}
		if err := s.repo.UpdateCard(ctx, id, updateData, 0); err != nil {
			return err
		}
		tx.Where("card_id = ?", id).Delete(&models.CardBank{})
		tx.Where("card_id = ?", id).Delete(&models.CardSocialMedia{})
		for _, cb := range card.CardBanks {
			if err := tx.Create(&cb).Error; err != nil {
				return err
			}
		}
		for _, csm := range card.CardSocialMedia {
			if err := tx.Create(&csm).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *CardService) DeleteCard(ctx context.Context, id uint) error {
	return s.repo.DeleteCard(ctx, id)
}

func (s *CardService) GetCardsByUserID(userID uint) ([]models.Card, int64, error) {
	params := queryparams.ListParams{
		Page:    1,
		PerPage: 1,
	}
	cards, totalCount, err := s.repo.GetAllCardsByUserID(userID, params)
	return cards, totalCount, err
}
