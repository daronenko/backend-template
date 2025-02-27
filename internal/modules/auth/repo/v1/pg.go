package repo

import (
	"context"
	"database/sql"

	"emperror.dev/errors"
	"github.com/daronenko/backend-template/internal/model/v1"
	"github.com/daronenko/backend-template/internal/util"
	"github.com/daronenko/backend-template/pkg/pgerrs"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repo for users
type User struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) Repo {
	return &User{db}
}

// Create new user
func (r *User) Create(ctx context.Context, user *model.User) (*model.User, error) {
	createdUser := &model.User{}

	err := r.db.QueryRowxContext(ctx, createUserQuery,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.Avatar,
	).StructScan(createdUser)
	if err != nil {
		if pgerrs.Is(err, pgerrs.UniqueViolation) {
			return nil, ErrUserExists
		}
		return nil, errors.Wrap(err, "repo.User.Create.StructScan")
	}

	return createdUser, nil
}

// Update existing user
func (r *User) Update(ctx context.Context, user *model.User) (*model.User, error) {
	updatedUser := &model.User{}

	err := r.db.QueryRowxContext(ctx, updateUserQuery,
		user.Username,
		user.Email,
		user.Role,
		user.Avatar,
		user.ID,
	).StructScan(updatedUser)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		case pgerrs.Is(err, pgerrs.UniqueViolation):
			return nil, ErrUserExists
		default:
			return nil, errors.Wrap(err, "repo.User.Update.StructScan")
		}
	}

	return updatedUser, nil
}

// Delete existing user
func (u *User) Delete(ctx context.Context, userID uuid.UUID) error {
	result, err := u.db.ExecContext(ctx, deleteUserQuery, userID)
	if err != nil {
		return errors.Wrap(err, "repo.User.Delete.ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "repo.User.Delete.RowsAffected")
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Get user by id
func (r *User) GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	foundUser := &model.User{}

	err := r.db.GetContext(ctx, foundUser, getUserQuery, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "repo.User.GetByID.GetContext")
	}

	return foundUser, nil
}

// Get user by email
func (u *User) GetByEmail(ctx context.Context, user *model.User) (*model.User, error) {
	foundUser := &model.User{}

	err := u.db.GetContext(ctx, foundUser, getUserByEmailQuery, user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "repo.User.GetByEmail.GetContext")
	}

	return foundUser, nil
}

// Find users by name with pagination
func (r *User) FindByName(ctx context.Context, name string, query *util.PaginationQuery) (*model.UsersList, error) {
	var totalUsers int
	if err := r.db.GetContext(ctx, &totalUsers, getTotalUsersQuery, name); err != nil {
		return nil, errors.Wrap(err, "repo.User.FindByName.GetContext")
	}

	if totalUsers == 0 {
		return &model.UsersList{
			TotalCount: totalUsers,
			TotalPages: util.GetTotalPages(totalUsers, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    util.GetHasMore(query.GetPage(), totalUsers, query.GetSize()),
			Users:      make([]*model.User, 0),
		}, nil
	}

	rows, err := r.db.QueryxContext(ctx, findUsersByNameQuery, name, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "repo.User.FindByName.QueryxContext")
	}
	defer rows.Close()

	users := make([]*model.User, 0, query.GetSize())
	for rows.Next() {
		var user model.User
		if err = rows.StructScan(&user); err != nil {
			return nil, errors.Wrap(err, "repo.User.FindByName.StructScan")
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "repo.User.FindByName.rows.Err")
	}

	return &model.UsersList{
		TotalCount: totalUsers,
		TotalPages: util.GetTotalPages(totalUsers, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    util.GetHasMore(query.GetPage(), totalUsers, query.GetSize()),
		Users:      users,
	}, nil
}

// Get users with pagination
func (r *User) GetUsers(ctx context.Context, pq *util.PaginationQuery) (*model.UsersList, error) {
	var totalUsers int
	if err := r.db.GetContext(ctx, &totalUsers, getTotalUsersQuery); err != nil {
		return nil, errors.Wrap(err, "repo.User.GetUsers.GetContext.totalCount")
	}

	if totalUsers == 0 {
		return &model.UsersList{
			TotalCount: totalUsers,
			TotalPages: util.GetTotalPages(totalUsers, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    util.GetHasMore(pq.GetPage(), totalUsers, pq.GetSize()),
			Users:      make([]*model.User, 0),
		}, nil
	}

	users := make([]*model.User, 0, pq.GetSize())
	if err := r.db.SelectContext(
		ctx,
		&users,
		getUsers,
		pq.GetOrderBy(),
		pq.GetOffset(),
		pq.GetLimit(),
	); err != nil {
		return nil, errors.Wrap(err, "repo.User.GetUsers.SelectContext")
	}

	return &model.UsersList{
		TotalCount: totalUsers,
		TotalPages: util.GetTotalPages(totalUsers, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    util.GetHasMore(pq.GetPage(), totalUsers, pq.GetSize()),
		Users:      users,
	}, nil
}
