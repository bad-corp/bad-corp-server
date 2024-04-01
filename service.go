package main

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func loginService(userId string) (string, int64, error) {
	token, expiresAt, err := makeToken(userId)
	if err != nil {
		return "", 0, err
	}

	return token, expiresAt.UnixMilli(), nil
}

func calcAndSave(db *gorm.DB) {
	m := calc(db)
	saveScore(db, &m)
}

func saveScore(db *gorm.DB, m *map[uint64]float64) {
	for anyCorpId, score := range *m {
		var corpScore CorpScore
		corpScore.CorpId = anyCorpId
		corpScore.Score = uint64(int(score * 10000))
		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "corp_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"score"}),
		}).Create(&corpScore)
	}
}

func calc(db *gorm.DB) map[uint64]float64 {
	var corpComments []CorpComment
	db.Find(&corpComments)
	return calcScore(&corpComments)
}

func score(score, addScore uint8, scoreCount uint32) (uint8, uint32) {
	var t1 = scoreCount + 1
	var t2 = uint32(score) * scoreCount
	var t3 = t2 + uint32(addScore)
	return uint8(t3 / t1), scoreCount + 1
}
