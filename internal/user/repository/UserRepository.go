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

func (userRepository *UserRepository) Insert(user *models.User) error {
	const query = "INSERT INTO profile (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)"
	_, err := userRepository.dbConnection.Exec(query, user.Nickname, user.About, user.Email, user.FullName)
	return err
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
		&user.Email); err == sql.ErrNoRows {
		return nil, err
	}
	return user, nil
}
