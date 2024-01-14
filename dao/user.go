package dao

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/samber/lo"
)

func AdminGetAllUsers() ([]*models.UserEntity, error) {
	db := models.GetDBManager().GetDefaultDB()
	users := make([]*models.User, 0)
	err := db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(users,
		func(u *models.User, _ int) *models.UserEntity {
			return u.UserEntity
		}), nil
}

func AdminCountUsers() (int64, error) {
	db := models.GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.User{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetUserByUserID(userID int) (*models.UserEntity, error) {
	if userID == 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	db := models.GetDBManager().GetDefaultDB()
	u := &models.User{}
	err := db.Where(&models.User{
		UserEntity: &models.UserEntity{
			UserID: userID,
		},
	}).First(u).Error
	if err != nil {
		return nil, err
	}
	return u.UserEntity, nil
}

func UpdateUser(userInfo models.UserInfo, user *models.UserEntity) error {
	db := models.GetDBManager().GetDefaultDB()
	user.UserID = userInfo.GetUserID()
	return db.Model(&models.User{}).Where(
		&models.User{
			UserEntity: &models.UserEntity{
				UserID: userInfo.GetUserID(),
			},
		},
	).Save(&models.User{
		UserEntity: user,
	}).Error
}

func GetUserByUserName(userName string) (*models.UserEntity, error) {
	if userName == "" {
		return nil, fmt.Errorf("invalid user name")
	}
	db := models.GetDBManager().GetDefaultDB()
	u := &models.User{}
	err := db.Where(&models.User{
		UserEntity: &models.UserEntity{
			UserName: userName,
		},
	}).First(u).Error
	if err != nil {
		return nil, err
	}
	return u.UserEntity, nil
}

func CheckUserPassword(userNameOrEmail, password string) (bool, models.UserInfo, error) {
	var user models.User
	db := models.GetDBManager().GetDefaultDB()

	if err := db.Where(&models.User{
		UserEntity: &models.UserEntity{
			UserName: userNameOrEmail,
		},
	}).Or(&models.User{
		UserEntity: &models.UserEntity{
			Email: userNameOrEmail,
		},
	}).First(&user).Error; err != nil {
		return false, nil, err
	}
	return utils.CheckPasswordHash(password, user.Password), user, nil
}

func CheckUserNameAndEmail(userName, email string) error {
	var user models.User
	db := models.GetDBManager().GetDefaultDB()

	if err := db.Where(&models.User{
		UserEntity: &models.UserEntity{
			UserName: userName,
		},
	}).Or(&models.User{
		UserEntity: &models.UserEntity{
			Email: email,
		},
	}).First(&user).Error; err != nil {
		return err
	}
	return nil
}

func CreateUser(user *models.UserEntity) error {
	u := &models.User{
		UserEntity: user,
	}
	db := models.GetDBManager().GetDefaultDB()
	return db.Create(u).Error
}
