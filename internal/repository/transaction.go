package repository

import (
	"gorm.io/gorm"

	"github.com/hacker4257/pet_charity/internal/database"
)

type GormTransactionManager struct{}

func NewTransactionManager() *GormTransactionManager {
	return &GormTransactionManager{}
}

func (m *GormTransactionManager) Transaction(fn func(tx TransactionContext) error) error {
	return database.DB.Transaction(func(gormTx *gorm.DB) error {
		txCtx := &gormTxContext{tx: gormTx}
		return fn(txCtx)
	})
}

// gormTxContext 实现 TransactionContext
type gormTxContext struct {
	tx *gorm.DB
}

func (c *gormTxContext) AdoptionRepo() AdoptionRepository {
	return &AdoptionRepo{db: c.tx}
}

func (c *gormTxContext) PetRepo() PetRepository {
	return &PetRepo{db: c.tx}
}

func (c *gormTxContext) DonationRepo() DonationRepository {
	return &DonationRepo{db: c.tx}
}
