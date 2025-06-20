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

type IInvitationService interface {
	GetAllInvitations(params queryparams.ListParams) (*queryparams.PaginatedResult, error)
	GetInvitationByID(ctx context.Context, id uint) (*models.Invitation, error)
	CreateInvitation(ctx context.Context, invitation *models.Invitation) error
	UpdateInvitation(ctx context.Context, id uint, invitation *models.Invitation) error
	DeleteInvitation(ctx context.Context, id uint) error
	GetParticipantsByInvitationID(invitationID uint) ([]models.InvitationParticipant, error)
	UpdateParticipant(id uint, participant *models.InvitationParticipant) error
	DeleteParticipant(id uint) error
}

type InvitationService struct {
	repo repositories.IInvitationRepository
}

func NewInvitationService() IInvitationService {
	return &InvitationService{repo: repositories.NewInvitationRepository()}
}

func (s *InvitationService) GetAllInvitations(params queryparams.ListParams) (*queryparams.PaginatedResult, error) {
	invitations, totalCount, err := s.repo.GetAllInvitations(params)
	if err != nil {
		logconfig.Log.Error("Davetiyeler alınamadı", zap.Error(err))
		return nil, errors.New("Davetiyeler getirilirken bir hata oluştu")
	}
	result := &queryparams.PaginatedResult{
		Data: invitations,
		Meta: queryparams.PaginationMeta{
			CurrentPage: params.Page,
			PerPage:     params.PerPage,
			TotalItems:  totalCount,
			TotalPages:  queryparams.CalculateTotalPages(totalCount, params.PerPage),
		},
	}
	return result, nil
}

func (s *InvitationService) GetInvitationByID(ctx context.Context, id uint) (*models.Invitation, error) {
	return s.repo.GetInvitationByID(id)
}

func (s *InvitationService) CreateInvitation(ctx context.Context, invitation *models.Invitation) error {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok || db == nil {
		db = databaseconfig.GetDB()
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.CreateInvitation(ctx, invitation); err != nil {
			return err
		}
		if invitation.InvitationDetail != nil {
			invitation.InvitationDetail.InvitationID = invitation.ID
			if err := tx.Create(invitation.InvitationDetail).Error; err != nil {
				return err
			}
		}
		for i := range invitation.Participants {
			invitation.Participants[i].InvitationID = invitation.ID
			if err := tx.Create(&invitation.Participants[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *InvitationService) UpdateInvitation(ctx context.Context, id uint, invitation *models.Invitation) error {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok || db == nil {
		db = databaseconfig.GetDB()
	}
	return db.Transaction(func(tx *gorm.DB) error {
		updateData := map[string]interface{}{
			"invitation_key": invitation.InvitationKey,
			"user_id":        invitation.UserID,
			"category_id":    invitation.CategoryID,
			"template":       invitation.Template,
			"type":           invitation.Type,
			"title":          invitation.Title,
			"image":          invitation.Image,
			"description":    invitation.Description,
			"venue":          invitation.Venue,
			"address":        invitation.Address,
			"location":       invitation.Location,
			"link":           invitation.Link,
			"telephone":      invitation.Telephone,
			"note":           invitation.Note,
			"date":           invitation.Date,
			"time":           invitation.Time,
			"is_confirmed":   invitation.IsConfirmed,
			"is_participant": invitation.IsParticipant,
		}
		if err := s.repo.UpdateInvitation(ctx, id, updateData, 0); err != nil {
			return err
		}
		tx.Where("invitation_id = ?", id).Delete(&models.InvitationDetail{})
		tx.Where("invitation_id = ?", id).Delete(&models.InvitationParticipant{})
		if invitation.InvitationDetail != nil {
			invitation.InvitationDetail.InvitationID = id
			if err := tx.Create(invitation.InvitationDetail).Error; err != nil {
				return err
			}
		}
		for i := range invitation.Participants {
			invitation.Participants[i].InvitationID = id
			if err := tx.Create(&invitation.Participants[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *InvitationService) DeleteInvitation(ctx context.Context, id uint) error {
	db := databaseconfig.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		tx.Where("invitation_id = ?", id).Delete(&models.InvitationDetail{})
		tx.Where("invitation_id = ?", id).Delete(&models.InvitationParticipant{})
		return s.repo.DeleteInvitation(ctx, id)
	})
}

func (s *InvitationService) GetParticipantsByInvitationID(invitationID uint) ([]models.InvitationParticipant, error) {
	return s.repo.GetParticipantsByInvitationID(invitationID)
}

func (s *InvitationService) UpdateParticipant(id uint, participant *models.InvitationParticipant) error {
	return s.repo.UpdateParticipant(id, participant)
}

func (s *InvitationService) DeleteParticipant(id uint) error {
	return s.repo.DeleteParticipant(id)
}
