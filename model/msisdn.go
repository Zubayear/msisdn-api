package model

import "github.com/google/uuid"

type PrimaryIdentityType string

func (p PrimaryIdentityType) Values() (kinds []string) {
	for _, s := range []PrimaryIdentityType{Prepaid, Postpaid} {
		kinds = append(kinds, string(s))
	}
	return
}

const (
	Prepaid  PrimaryIdentityType = "Prepaid"
	Postpaid PrimaryIdentityType = "Postpaid"
)

// Msisdn API request from client
type Msisdn struct {
	PrimaryIdentity string              `json:"primaryIdentity" binding:"required"`
	MsisdnType      PrimaryIdentityType `json:"msisdnType" binding:"required"`
	Provisioned     bool                `json:"provisioned"`
}

func NewMsisdn(primaryIdentity string, msisdnType PrimaryIdentityType, provisioned bool) *Msisdn {
	return &Msisdn{PrimaryIdentity: primaryIdentity, MsisdnType: msisdnType, Provisioned: provisioned}
}

type MsisdnDto struct {
	Id              uuid.UUID           `json:"id"`
	PrimaryIdentity string              `json:"primaryIdentity"`
	MsisdnType      PrimaryIdentityType `json:"msisdnType"`
	Provisioned     bool                `json:"provisioned"`
}

func NewMsisdnDto(id uuid.UUID, primaryIdentity string, msisdnType PrimaryIdentityType, provisioned bool) *MsisdnDto {
	return &MsisdnDto{Id: id, PrimaryIdentity: primaryIdentity, MsisdnType: msisdnType, Provisioned: provisioned}
}
