package models

import (
	"time"

	"gorm.io/gorm"
)

type Emprunt struct {
	IDEmprunt          uint       `gorm:"primaryKey;column:id_emprunt"`
	UtilisateurID      uint       `gorm:"not null"`
	LivreID            uint       `gorm:"not null"`
	DateEmprunt        time.Time  `gorm:"not null"`
	DateRetourPrevu    time.Time  `gorm:"not null"`
	DateRetourEffectif *time.Time `gorm:"default:null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type Penalite struct {
	IDPenalite    uint           `gorm:"primaryKey;column:id_penalite"`
	UtilisateurID uint           `gorm:"not null;column:utilisateur_id"`
	EmpruntID     uint           `gorm:"not null;column:emprunt_id"`
	Montant       float64        `gorm:"not null"`
	DatePaiement  *time.Time     `gorm:"column:date_paiement"`
	Emprunt       Emprunt        `gorm:"foreignKey:EmpruntID;references:IDEmprunt;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index;column:deleted_at"`
}
