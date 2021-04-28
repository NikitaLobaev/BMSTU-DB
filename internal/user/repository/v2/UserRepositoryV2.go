package v2

import (
	"database/sql"
	"github.com/NikitaLobaev/BMSTU-DB/internal/models"
	"github.com/labstack/gommon/log"
)

type UserRepositoryV2 struct {
	dbConnection *sql.DB
}

func NewUserRepositoryV2(dbConnection *sql.DB) *UserRepositoryV2 {
	return &UserRepositoryV2{
		dbConnection: dbConnection,
	}
}

func (userRepositoryV2 *UserRepositoryV2) Insert(user *models.User) (*models.User, error) {
	const query = "INSERT INTO user_ (nickname, fullname, about, email) VALUES ($1, $2, $3, $4) RETURNING nickname, fullname, about, email"
	if err := userRepositoryV2.dbConnection.QueryRow(query, user.Nickname, user.FullName, user.About, user.Email).
		Scan(&user.Nickname, &user.FullName, &user.About, &user.Email); err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepositoryV2 *UserRepositoryV2) SelectByNicknameOrEmail(nickname string, email string) (*models.Users, error) {
	const query = "SELECT nickname, fullname, about, email FROM user_ WHERE nickname = $1 OR email = $2"
	rows, err := userRepositoryV2.dbConnection.Query(query, nickname, email)
	defer func() {
		if err := rows.Close(); err != nil {
			log.Print(err)
		}
	}()
	if err != nil {
		return nil, err
	}

	users := make(models.Users, 0)
	for rows.Next() {
		user := new(models.User)
		if err := rows.Scan(&user.Nickname, &user.FullName, &user.About, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

func (userRepositoryV2 *UserRepositoryV2) SelectByNickname(nickname string) (*models.User, error) {
	const query = "SELECT nickname, fullname, about, email FROM user_ WHERE nickname = $1"
	user := new(models.User)
	if err := userRepositoryV2.dbConnection.QueryRow(query, nickname).Scan(&user.Nickname, &user.FullName, &user.About,
		&user.Email); err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepositoryV2 *UserRepositoryV2) SelectByEmail(email string) (*models.User, error) {
	const query = "SELECT nickname, fullname, about, email FROM user_ WHERE email = $1"
	user := new(models.User)
	if err := userRepositoryV2.dbConnection.QueryRow(query, email).Scan(&user.Nickname, &user.FullName, &user.About,
		&user.Email); err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepositoryV2 *UserRepositoryV2) Update(nickname string, userUpdate *models.UserUpdate) (*models.User, error) {
	const query = "UPDATE user_ SET fullname = $2, about = $3, email = $4 WHERE nickname = $1 RETURNING nickname, fullname, about, email"
	user := new(models.User)
	if err := userRepositoryV2.dbConnection.QueryRow(query, nickname, userUpdate.FullName, userUpdate.About,
		userUpdate.Email).Scan(&user.Nickname, &user.FullName, &user.About, &user.Email); err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepositoryV2 *UserRepositoryV2) SelectUsersByForumSlug(forumSlug string, userParams *models.UserParams) (*models.Users, error) {
	const query1 = "SELECT user_.nickname, user_.fullname, user_.about, user_.email FROM forum_user JOIN user_ ON forum_user.user_nickname = user_.nickname WHERE forum_slug = $1"
	const queryNickname1 = " AND user_.nickname "
	const queryLess = "<"
	const queryMore = ">"
	const queryNickname2 = " $2"
	const queryOrderBy = " ORDER BY user_.nickname"
	const queryDesc = " DESC"
	const queryLimit1 = " LIMIT $2"
	const queryLimit2 = " LIMIT $3"

	query := query1
	var rows *sql.Rows
	var err error
	if userParams.IsSinceSet() {
		query += queryNickname1
		if userParams.IsDescSet() && userParams.Desc {
			query += queryLess + queryNickname2 + queryOrderBy + queryDesc
		} else {
			query += queryMore + queryNickname2 + queryOrderBy
		}
		if userParams.IsLimitSet() {
			query += queryLimit2
			rows, err = userRepositoryV2.dbConnection.Query(query, forumSlug, userParams.Since, userParams.Limit)
		} else {
			rows, err = userRepositoryV2.dbConnection.Query(query, forumSlug, userParams.Since)
		}
	} else {
		query += queryOrderBy
		if userParams.IsDescSet() && userParams.Desc {
			query += queryDesc
		}
		if userParams.IsLimitSet() {
			query += queryLimit1
			rows, err = userRepositoryV2.dbConnection.Query(query, forumSlug, userParams.Limit)
		} else {
			rows, err = userRepositoryV2.dbConnection.Query(query, forumSlug)
		}
	}

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error(err)
		}
	}()

	users := make(models.Users, 0)
	for rows.Next() {
		user := new(models.User)
		if err := rows.Scan(&user.Nickname, &user.FullName, &user.About, &user.Email); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return &users, nil
}
