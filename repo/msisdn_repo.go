package repo

import (
	"context"
	"github.com/google/uuid"
	"huspass/ent"
	"huspass/external"
	"huspass/model"
)

type MsisdnRepository interface {
	CreateMsisdn(ctx context.Context, msisdn *model.Msisdn) (*model.MsisdnDto, error)
	GetAllMsisdn(ctx context.Context) ([]*model.MsisdnDto, error)
	UpdateMsisdn(ctx context.Context, msisdn *model.Msisdn, id uuid.UUID) (*model.MsisdnDto, error)
	DeleteMsisdn(ctx context.Context, id uuid.UUID) error
	GetMsisdn(ctx context.Context, id uuid.UUID) (*model.MsisdnDto, error)
}

func (m *MsisdnRepositoryImpl) CreateMsisdn(ctx context.Context, msisdn *model.Msisdn) (*model.MsisdnDto, error) {
	id := uuid.New()
	provisioned := msisdn.Provisioned
	savedMsisdn, err := m.client.Msisdn.Create().
		SetID(id).
		SetPrimaryIdentity(msisdn.PrimaryIdentity).
		SetPrimaryIdentityType(msisdn.MsisdnType).
		SetProvisioned(provisioned).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return model.NewMsisdnDto(id, savedMsisdn.PrimaryIdentity, savedMsisdn.PrimaryIdentityType, savedMsisdn.Provisioned), nil
}

func (m *MsisdnRepositoryImpl) GetAllMsisdn(ctx context.Context) ([]*model.MsisdnDto, error) {
	allMsisdnFromRepo, err := m.client.Msisdn.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	var msisdnToReturn []*model.MsisdnDto
	for _, m := range allMsisdnFromRepo {
		msisdnDto := model.NewMsisdnDto(
			m.ID, m.PrimaryIdentity, m.PrimaryIdentityType, m.Provisioned,
		)
		msisdnToReturn = append(msisdnToReturn, msisdnDto)
	}
	return msisdnToReturn, nil
}

func (m *MsisdnRepositoryImpl) UpdateMsisdn(ctx context.Context, msisdn *model.Msisdn, id uuid.UUID) (*model.MsisdnDto, error) {
	updatedMsisdn, err := m.client.Msisdn.UpdateOneID(id).
		SetPrimaryIdentity(msisdn.PrimaryIdentity).
		SetPrimaryIdentityType(msisdn.MsisdnType).
		SetProvisioned(msisdn.Provisioned).Save(ctx)
	if err != nil {
		return nil, err
	}
	return model.NewMsisdnDto(
		id, updatedMsisdn.PrimaryIdentity, updatedMsisdn.PrimaryIdentityType, updatedMsisdn.Provisioned,
	), nil
}

// DeleteMsisdn should do a soft delete
func (m *MsisdnRepositoryImpl) DeleteMsisdn(ctx context.Context, id uuid.UUID) error {
	err := m.client.Msisdn.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *MsisdnRepositoryImpl) GetMsisdn(ctx context.Context, id uuid.UUID) (*model.MsisdnDto, error) {
	msisdnFromRepo, err := m.client.Msisdn.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.NewMsisdnDto(msisdnFromRepo.ID, msisdnFromRepo.PrimaryIdentity, msisdnFromRepo.PrimaryIdentityType, msisdnFromRepo.Provisioned), nil
}

// MsisdnRepositoryImpl Provides MsisdnRepository implementation to save entity in db
type MsisdnRepositoryImpl struct {
	client *ent.Client
}

func NewMsisdnRepositoryImpl(client *ent.Client) *MsisdnRepositoryImpl {
	return &MsisdnRepositoryImpl{client: client}
}

// DatabaseImplProvider provides implementation of interface MsisdnRepository
func DatabaseImplProvider(h *external.Host) (*MsisdnRepositoryImpl, error) {
	client, err := ent.Open("mysql", h.ConnString)
	if err != nil {
		return nil, err
	}
	return NewMsisdnRepositoryImpl(client), nil
}
