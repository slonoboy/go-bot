package database

import (
	"gorm.io/gorm"
)

type DatabaseHandler struct {
	db *gorm.DB
}

func (dh *DatabaseHandler) CreateUser(user User) error {
	dh.db.Create(&user)
	if dh.db.Error != nil {
		return dh.db.Error
	}
	return nil
}

func (dh *DatabaseHandler) FindUserByTGID(tgid int64) (User, error) {
	var user User
	dh.db.Where("tg_id = ?", tgid).First(&user)
	if dh.db.Error != nil {
		return User{}, dh.db.Error
	}
	return user, nil
}

func (dh *DatabaseHandler) AddConversion(conversion Conversion) error {
	dh.db.Create(&conversion)
	if dh.db.Error != nil {
		return dh.db.Error
	}
	return nil
}

func (dh *DatabaseHandler) UserMostFrequentCryptos(limit int) ([]ConversionCount, error) {
	var counts []ConversionCount
	dh.db.Model(&Conversion{}).
		Select("crypto, COUNT(*) as count").
		Group("crypto").
		Order("count DESC").
		Limit(limit).
		Scan(&counts)

	if dh.db.Error != nil {
		return []ConversionCount{}, dh.db.Error
	}
	return counts, nil
}

func (dh *DatabaseHandler) UserMostFrequentCurrencies(limit int) ([]ConversionCount, error) {
	var counts []ConversionCount
	dh.db.Model(&Conversion{}).
		Select("currency, COUNT(*) as count").
		Group("currency").
		Order("count DESC").
		Limit(limit).
		Scan(&counts)

	if dh.db.Error != nil {
		return []ConversionCount{}, dh.db.Error
	}
	return counts, nil
}

func (dh *DatabaseHandler) UserFirstConversion(tgid int64) (Conversion, error) {
	user, err := dh.FindUserByTGID(tgid)
	if err != nil {
		return Conversion{}, err
	}

	var conversion Conversion
	dh.db.Where("user_id = ?", user.ID).Order("created_at").First(&conversion)
	if dh.db.Error != nil {
		return Conversion{}, err
	}

	return conversion, nil
}

func (dh *DatabaseHandler) UserConversionsCount(tgid int64) (int, error) {
	user, err := dh.FindUserByTGID(tgid)
	if err != nil {
		return 0, err
	}

	var conversionCount ConversionCount
	dh.db.Model(&Conversion{}).
		Select("COUNT(*) as count").
		Where("user_id = ?", user.ID).
		Scan(&conversionCount)

	if dh.db.Error != nil {
		return 0, err
	}

	return conversionCount.Count, nil
}

func NewDatabaseHandler(db *gorm.DB) (*DatabaseHandler, error) {
	err := db.AutoMigrate(&User{}, &Conversion{})
	if err != nil {
		return &DatabaseHandler{}, err
	}

	return &DatabaseHandler{
		db: db,
	}, nil
}
