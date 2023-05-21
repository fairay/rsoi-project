package repository

import (
	"log"
	"statistics/objects"

	"github.com/jinzhu/gorm"
)

type StatisticsRep interface {
	Create(*objects.RequestStat) error
}

type PGStatisticsRep struct {
	db *gorm.DB
}

func NewPGStatisticsRep(db *gorm.DB) *PGStatisticsRep {
	return &PGStatisticsRep{db}
}

func (rep *PGStatisticsRep) Create(statistics *objects.RequestStat) error {
	log.Println("writing statistics to db")
	return rep.db.Create(statistics).Error
}
