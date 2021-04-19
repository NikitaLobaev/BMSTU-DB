package repository

import (
	"../../models"
	"database/sql"
	"log"
)

type UserRepository struct {
	dbConnection *sql.DB
}

func NewUserRepository(dbConnection *sql.DB) *UserRepository {
	return &UserRepository{
		dbConnection: dbConnection,
	}
}

func (userRepository *UserRepository) Insert(user *models.User) (*models.User, error) {
	const query = "INSERT INTO profile (nickname, fullname, about, email) VALUES ($1, $2, $3, $4) RETURNING nickname, fullname, about, email"
	if err := userRepository.dbConnection.QueryRow(query, user.Nickname, user.About, user.Email, user.FullName).
		Scan(&user.Nickname, &user.FullName, &user.About, &user.Email); err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepository *UserRepository) SelectByNicknameOrEmail(nickname string, email string) ([]*models.User, error) {
	const query = "SELECT nickname, fullname, about, email FROM profile WHERE nickname = $1 OR email = $2"
	rows, err := userRepository.dbConnection.Query(query, nickname, email)
	defer func() {
		if err := rows.Close(); err != nil {
			log.Print(err)
		}
	}()
	if err != nil {
		return nil, err
	}

	var users []*models.User
	for rows.Next() {
		user := new(models.User)
		if err := rows.Scan(&user.Nickname, &user.FullName, &user.About, &user.Email); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (userRepository *UserRepository) SelectByNickname(nickname string) (*models.User, error) {
	const query = "SELECT nickname, fullname, about, email FROM profile WHERE nickname = $1"
	user := new(models.User)
	if err := userRepository.dbConnection.QueryRow(query, nickname).Scan(&user.Nickname, &user.FullName, &user.About,
		&user.Email); err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepository *UserRepository) Update(nickname string, userUpdate *models.UserUpdate) (*models.User, error) {
	const query = "UPDATE profile SET fullname = $2, about = $3, email = $4 WHERE nickname = $1 RETURNING nickname, fullname, about, email"
	user := new(models.User)
	if err := userRepository.dbConnection.QueryRow(query, nickname, userUpdate.FullName, userUpdate.About,
		userUpdate.Email).Scan(&user.Nickname, &user.FullName, &user.About, &user.Email); err != nil {
		return nil, err
	}
	return user, nil
}
