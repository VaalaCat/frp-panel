package dao

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/samber/lo"
)

type UserQuery interface {
	AdminGetAllUsers() ([]*models.UserEntity, error)
	AdminCountUsers() (int64, error)
	GetUserByUserID(userID int) (*models.UserEntity, error)
	GetUserByUserName(userName string) (*models.UserEntity, error)
	CheckUserPassword(userNameOrEmail, password string) (bool, models.UserInfo, error)
	CheckUserNameAndEmail(userName, email string) error
}

type UserMutation interface {
	UpdateUser(userInfo models.UserInfo, user *models.UserEntity) error
	AdminUpdateUser(userInfo models.UserInfo, user *models.UserEntity) error
	CreateUser(user *models.UserEntity) error
}

type userQuery struct{ *queryImpl }
type userMutation struct{ *mutationImpl }

func newUserQuery(base *queryImpl) UserQuery          { return &userQuery{base} }
func newUserMutation(base *mutationImpl) UserMutation { return &userMutation{base} }

func (q *userQuery) AdminGetAllUsers() ([]*models.UserEntity, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
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

func (q *userQuery) AdminCountUsers() (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.User{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q *userQuery) GetUserByUserID(userID int) (*models.UserEntity, error) {
	if userID == 0 {
		return nil, fmt.Errorf("invalid user id")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
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

func (m *userMutation) UpdateUser(userInfo models.UserInfo, user *models.UserEntity) error {
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
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

func (m *userMutation) AdminUpdateUser(userInfo models.UserInfo, user *models.UserEntity) error {
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
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

func (q *userQuery) GetUserByUserName(userName string) (*models.UserEntity, error) {
	if userName == "" {
		return nil, fmt.Errorf("invalid user name")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
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

func (q *userQuery) CheckUserPassword(userNameOrEmail, password string) (bool, models.UserInfo, error) {
	var user models.User
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()

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

func (q *userQuery) CheckUserNameAndEmail(userName, email string) error {
	var user models.User
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()

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

func (m *userMutation) CreateUser(user *models.UserEntity) error {
	u := &models.User{
		UserEntity: user,
	}
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Create(u).Error
}
