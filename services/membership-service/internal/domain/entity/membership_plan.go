package entity

import (
	"time"

	"github.com/google/uuid"
)

type MembershipPlan struct {
	ID             uuid.UUID
	Name           string
	DurationMonths int
	Price          float64
	Description    string
	IsDeleted      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewMembershipPlan(
	name string,
	durationMonths int,
	price float64,
	description string,
) *MembershipPlan {
	now := time.Now().UTC()

	return &MembershipPlan{
		ID:             uuid.New(),
		Name:           name,
		DurationMonths: durationMonths,
		Price:          price,
		Description:    description,
		IsDeleted:      false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (p *MembershipPlan) Update(
	name string,
	durationMonths int,
	price float64,
	description string,
) {
	p.Name = name
	p.DurationMonths = durationMonths
	p.Price = price
	p.Description = description
	p.UpdatedAt = time.Now().UTC()
}

func (p *MembershipPlan) SoftDelete() {
	p.IsDeleted = true
	p.UpdatedAt = time.Now().UTC()
}

func (p MembershipPlan) IsValidDuration() bool {
	return p.DurationMonths == 1 ||
		p.DurationMonths == 6 ||
		p.DurationMonths == 12
}
