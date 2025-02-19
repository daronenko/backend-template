package repo

import (
	"context"

	"emperror.dev/errors"
	"github.com/daronenko/backend-template/internal/models"
	"github.com/daronenko/backend-template/internal/pkg/pgerrs"
	"github.com/daronenko/backend-template/pkg/logger"
	"github.com/daronenko/backend-template/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
)

// Repo for users
type User struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewUser(db *sqlx.DB, logger logger.Logger) Repo {
	return &User{db, logger}
}

// Create new user
func (r *User) Create(ctx context.Context, user *models.User) (*models.User, error) {
	createdUser := &models.User{}

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
func (r *User) Update(ctx context.Context, user *models.User) (*models.User, error) {
	updatedUser := &models.User{}

	err := r.db.QueryRowxContext(ctx, updateUserQuery,
		user.Username,
		user.Email,
		user.Role,
		user.Avatar,
		user.ID,
	).StructScan(updatedUser)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
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
func (r *User) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	foundUser := &models.User{}

	err := r.db.GetContext(ctx, foundUser, getUserQuery, userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "repo.User.GetByID.GetContext")
	}

	return foundUser, nil
}

// Get user by email
func (u *User) GetByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	foundUser := &models.User{}

	err := u.db.GetContext(ctx, foundUser, getUserByEmailQuery, user.Email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "repo.User.GetByEmail.GetContext")
	}

	return foundUser, nil
}

// Find users by name with pagination
func (r *User) FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error) {
	var totalUsers int
	if err := r.db.GetContext(ctx, &totalUsers, getTotalUsersQuery, name); err != nil {
		return nil, errors.Wrap(err, "repo.User.FindByName.GetContext")
	}

	if totalUsers == 0 {
		return &models.UsersList{
			TotalCount: totalUsers,
			TotalPages: utils.GetTotalPages(totalUsers, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalUsers, query.GetSize()),
			Users:      make([]*models.User, 0),
		}, nil
	}

	rows, err := r.db.QueryxContext(ctx, findUsersByNameQuery, name, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "repo.User.FindByName.QueryxContext")
	}
	defer rows.Close()

	var users = make([]*models.User, 0, query.GetSize())
	for rows.Next() {
		var user models.User
		if err = rows.StructScan(&user); err != nil {
			return nil, errors.Wrap(err, "repo.User.FindByName.StructScan")
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "repo.User.FindByName.rows.Err")
	}

	return &models.UsersList{
		TotalCount: totalUsers,
		TotalPages: utils.GetTotalPages(totalUsers, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalUsers, query.GetSize()),
		Users:      users,
	}, nil
}

// Get users with pagination
func (r *User) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	var totalUsers int
	if err := r.db.GetContext(ctx, &totalUsers, getTotalUsersQuery); err != nil {
		return nil, errors.Wrap(err, "repo.User.GetUsers.GetContext.totalCount")
	}

	if totalUsers == 0 {
		return &models.UsersList{
			TotalCount: totalUsers,
			TotalPages: utils.GetTotalPages(totalUsers, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalUsers, pq.GetSize()),
			Users:      make([]*models.User, 0),
		}, nil
	}

	var users = make([]*models.User, 0, pq.GetSize())
	if err := r.db.SelectContext(
		ctx,
		&users,
		getTotalUsersQuery,
		pq.GetOrderBy(),
		pq.GetOffset(),
		pq.GetLimit(),
	); err != nil {
		return nil, errors.Wrap(err, "repo.User.GetUsers.SelectContext")
	}

	return &models.UsersList{
		TotalCount: totalUsers,
		TotalPages: utils.GetTotalPages(totalUsers, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalUsers, pq.GetSize()),
		Users:      users,
	}, nil
}
