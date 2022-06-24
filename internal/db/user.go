package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rinatkh/db_forum/internal/model/core"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *core.User) error
	GetUserByEmail(ctx context.Context, email string) (*core.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (*core.User, error)
	GetUsersByEmailOrNickname(ctx context.Context, email, nickname string) ([]*core.User, error)
	EditUser(ctx context.Context, user *core.User) (*core.User, error)
}

type userRepositoryImpl struct {
	dbConn *pgxpool.Pool
}

func (repo *userRepositoryImpl) CreateUser(ctx context.Context, user *core.User) error {
	_, err := repo.dbConn.Exec(ctx, "INSERT INTO Users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4);", user.Nickname, user.Fullname, user.About, user.Email)
	return err
}

func (repo *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*core.User, error) {
	user := &core.User{}
	err := repo.dbConn.QueryRow(ctx, "SELECT nickname, fullname, about, email FROM Users where email = $1;", email).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	return user, err
}

func (repo *userRepositoryImpl) GetUserByNickname(ctx context.Context, nickname string) (*core.User, error) {
	user := &core.User{}
	err := repo.dbConn.QueryRow(ctx, "SELECT nickname, fullname, about, email FROM Users where nickname = $1;", nickname).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	return user, err
}

func (repo *userRepositoryImpl) GetUsersByEmailOrNickname(ctx context.Context, email, nickname string) ([]*core.User, error) {
	rows, err := repo.dbConn.Query(ctx, "SELECT nickname, fullname, about, email FROM Users WHERE email = $1 OR nickname = $2;", email, nickname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*core.User
	for rows.Next() {
		u := &core.User{}
		if err := rows.Scan(&u.Nickname, &u.Fullname, &u.About, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (repo *userRepositoryImpl) EditUser(ctx context.Context, user *core.User) (*core.User, error) {
	updatedUser := &core.User{Nickname: user.Nickname}
	if err := repo.dbConn.QueryRow(ctx, "UPDATE Users SET fullname = COALESCE(NULLIF(TRIM($1), ''), fullname), about = COALESCE(NULLIF(TRIM($2), ''), about), email = COALESCE(NULLIF(TRIM($3), ''), email) WHERE nickname = $4 RETURNING fullname, about, email;", user.Fullname, user.About, user.Email, user.Nickname).Scan(&updatedUser.Fullname, &updatedUser.About, &updatedUser.Email); err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func NewUserRepository(dbConn *pgxpool.Pool) *userRepositoryImpl {
	return &userRepositoryImpl{dbConn: dbConn}
}
