package repository

import (
	"../../models"
	"database/sql"
	"log"
)

type PostRepository struct {
	dbConnection *sql.DB
}

func NewPostRepository(dbConnection *sql.DB) *PostRepository {
	return &PostRepository{
		dbConnection: dbConnection,
	}
}

func (postRepository *PostRepository) Insert(user *models.User) error {
	const query = "INSERT INTO profile (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)"
	_, err := postRepository.dbConnection.Exec(query, user.Nickname, user.About, user.Email, user.FullName)
	return err
}

func (postRepository *PostRepository) SelectByNicknameOrEmail(nickname string, email string) ([]*models.User, error) {
	const query = "SELECT nickname, fullname, about, email FROM profile WHERE nickname = $1 OR email = $2"
	rows, err := postRepository.dbConnection.Query(query, nickname, email)
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

func (postRepository *PostRepository) SelectByNickname(nickname string) (*models.User, error) {
	const query = "SELECT nickname, fullname, about, email FROM profile WHERE nickname = $1"
	user := new(models.User)
	if err := postRepository.dbConnection.QueryRow(query, nickname).Scan(&user.Nickname, &user.FullName, &user.About,
		&user.Email); err != nil {
		return nil, err
	}
	return user, nil
}

func (postRepository *PostRepository) Update(nickname string, userUpdate *models.UserUpdate) (*models.User, error) {
	const query = "UPDATE profile SET fullname = $2, about = $3, email = $4 WHERE nickname = $1 RETURNING nickname, fullname, about, email"
	user := new(models.User)
	if err := postRepository.dbConnection.QueryRow(query, nickname, userUpdate.FullName, userUpdate.About,
		userUpdate.Email).Scan(&user.Nickname, &user.FullName, &user.About, &user.Email); err != nil {
		return nil, err
	}
	return user, nil
}
