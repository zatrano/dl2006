package repositories

import (
	"context"

	"davet.link/configs/databaseconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"

	"gorm.io/gorm"
)

type ICardRepository interface {
	GetAllCards(params queryparams.ListParams) ([]models.Card, int64, error)
	GetCardByID(id uint) (*models.Card, error)
	CreateCard(ctx context.Context, card *models.Card) error
	BulkCreateCards(ctx context.Context, cards []models.Card) error
	UpdateCard(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error
	BulkUpdateCards(ctx context.Context, condition map[string]interface{}, data map[string]interface{}, updatedBy uint) error
	DeleteCard(ctx context.Context, id uint) error
	BulkDeleteCards(ctx context.Context, condition map[string]interface{}) error
	GetCardCount() (int64, error)
	GetAllCardsByUserID(userID uint, params queryparams.ListParams) ([]models.Card, int64, error)
}

type CardRepository struct {
	base IBaseRepository[models.Card]
	db   *gorm.DB
}

func NewCardRepository() ICardRepository {
	base := NewBaseRepository[models.Card](databaseconfig.GetDB())
	base.SetAllowedSortColumns([]string{"id", "name", "slug", "created_at"})
	// İlişkili tabloları preload et
	base.SetPreloads(
		"User",
		"Banks",
		"SocialMedia",
		"CardBanks",
		"CardSocialMedia",
	)
	return &CardRepository{base: base, db: databaseconfig.GetDB()}
}

func (r *CardRepository) GetAllCards(params queryparams.ListParams) ([]models.Card, int64, error) {
	return r.base.GetAll(params)
}

func (r *CardRepository) GetCardByID(id uint) (*models.Card, error) {
	return r.base.GetByID(id)
}

func (r *CardRepository) CreateCard(ctx context.Context, card *models.Card) error {
	return r.base.Create(ctx, card)
}

func (r *CardRepository) BulkCreateCards(ctx context.Context, cards []models.Card) error {
	return r.base.BulkCreate(ctx, cards)
}

func (r *CardRepository) UpdateCard(ctx context.Context, id uint, data map[string]interface{}, updatedBy uint) error {
	return r.base.Update(ctx, id, data, updatedBy)
}

func (r *CardRepository) BulkUpdateCards(ctx context.Context, condition map[string]interface{}, data map[string]interface{}, updatedBy uint) error {
	return r.base.BulkUpdate(ctx, condition, data, updatedBy)
}

func (r *CardRepository) DeleteCard(ctx context.Context, id uint) error {
	return r.base.Delete(ctx, id)
}

func (r *CardRepository) BulkDeleteCards(ctx context.Context, condition map[string]interface{}) error {
	return r.base.BulkDelete(ctx, condition)
}

func (r *CardRepository) GetCardCount() (int64, error) {
	return r.base.GetCount()
}

func (r *CardRepository) GetAllCardsByUserID(userID uint, params queryparams.ListParams) ([]models.Card, int64, error) {
	var cards []models.Card
	db := r.db.Where("user_id = ?", userID)
	var totalCount int64
	db.Model(&models.Card{}).Count(&totalCount)
	db = db.Limit(params.PerPage).Offset((params.Page - 1) * params.PerPage)
	db = db.Order(params.SortBy + " " + params.OrderBy)
	db.Find(&cards)
	return cards, totalCount, db.Error
}

var _ ICardRepository = (*CardRepository)(nil)
var _ IBaseRepository[models.Card] = (*BaseRepository[models.Card])(nil)
