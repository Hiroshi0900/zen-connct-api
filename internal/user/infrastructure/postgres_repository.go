package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"zen-connect/internal/user/domain"
)

// PostgresUserRepository implements UserRepository interface
type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		pool: pool,
	}
}

// Save saves a user to the database
func (r *PostgresUserRepository) Save(user *domain.User) error {
	ctx := context.Background()
	
	query := `
		INSERT INTO users (
			id, auth0_user_id, email, display_name, bio, profile_image_url,
			email_verified, created_at, verified_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			auth0_user_id = EXCLUDED.auth0_user_id,
			email = EXCLUDED.email,
			display_name = EXCLUDED.display_name,
			bio = EXCLUDED.bio,
			profile_image_url = EXCLUDED.profile_image_url,
			email_verified = EXCLUDED.email_verified,
			verified_at = EXCLUDED.verified_at,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.pool.Exec(ctx, query,
		user.ID(),
		user.Auth0UserID(),
		user.Email().String(),
		user.Profile().DisplayName(),
		user.Profile().Bio(),
		user.Profile().ProfileImageURL(),
		user.EmailVerified(),
		user.CreatedAt(),
		user.VerifiedAt(),
		user.UpdatedAt(),
	)

	return err
}

// FindByID finds a user by ID
func (r *PostgresUserRepository) FindByID(id string) (*domain.User, error) {
	ctx := context.Background()
	
	query := `
		SELECT id, auth0_user_id, email, display_name, bio, profile_image_url,
			   email_verified, created_at, verified_at, updated_at
		FROM users
		WHERE id = $1
	`

	var userID, auth0UserID, email string
	var displayName, bio, profileImageURL sql.NullString
	var emailVerified bool
	var createdAt, updatedAt time.Time
	var verifiedAt sql.NullTime

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&userID,
		&auth0UserID,
		&email,
		&displayName,
		&bio,
		&profileImageURL,
		&emailVerified,
		&createdAt,
		&verifiedAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return r.reconstructUser(
		userID,
		auth0UserID,
		email,
		displayName,
		bio,
		profileImageURL,
		emailVerified,
		createdAt,
		verifiedAt,
		updatedAt,
	)
}

// FindByEmail finds a user by email
func (r *PostgresUserRepository) FindByEmail(email *domain.Email) (*domain.User, error) {
	ctx := context.Background()
	
	query := `
		SELECT id, auth0_user_id, email, display_name, bio, profile_image_url,
			   email_verified, created_at, verified_at, updated_at
		FROM users
		WHERE email = $1
	`

	var id, auth0UserID, emailStr string
	var displayName, bio, profileImageURL sql.NullString
	var emailVerified bool
	var createdAt, updatedAt time.Time
	var verifiedAt sql.NullTime

	err := r.pool.QueryRow(ctx, query, email.String()).Scan(
		&id,
		&auth0UserID,
		&emailStr,
		&displayName,
		&bio,
		&profileImageURL,
		&emailVerified,
		&createdAt,
		&verifiedAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return r.reconstructUser(
		id,
		auth0UserID,
		emailStr,
		displayName,
		bio,
		profileImageURL,
		emailVerified,
		createdAt,
		verifiedAt,
		updatedAt,
	)
}

// FindByAuth0UserID finds a user by Auth0 user ID
func (r *PostgresUserRepository) FindByAuth0UserID(auth0UserID string) (*domain.User, error) {
	ctx := context.Background()
	
	query := `
		SELECT id, auth0_user_id, email, display_name, bio, profile_image_url,
			   email_verified, created_at, verified_at, updated_at
		FROM users
		WHERE auth0_user_id = $1
	`

	var id, emailStr string
	var displayName, bio, profileImageURL sql.NullString
	var emailVerified bool
	var createdAt, updatedAt time.Time
	var verifiedAt sql.NullTime

	err := r.pool.QueryRow(ctx, query, auth0UserID).Scan(
		&id,
		&auth0UserID,
		&emailStr,
		&displayName,
		&bio,
		&profileImageURL,
		&emailVerified,
		&createdAt,
		&verifiedAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return r.reconstructUser(
		id,
		auth0UserID,
		emailStr,
		displayName,
		bio,
		profileImageURL,
		emailVerified,
		createdAt,
		verifiedAt,
		updatedAt,
	)
}

// reconstructUser reconstructs a user from database values
func (r *PostgresUserRepository) reconstructUser(
	id string,
	auth0UserID string,
	email string,
	displayName sql.NullString,
	bio sql.NullString,
	profileImageURL sql.NullString,
	emailVerified bool,
	createdAt time.Time,
	verifiedAt sql.NullTime,
	updatedAt time.Time,
) (*domain.User, error) {
	var verifiedAtPtr *time.Time
	if verifiedAt.Valid {
		verifiedAtPtr = &verifiedAt.Time
	}

	return domain.FromSnapshot(
		id,
		auth0UserID,
		email,
		displayName.String,
		bio.String,
		profileImageURL.String,
		emailVerified,
		createdAt,
		verifiedAtPtr,
		updatedAt,
	)
}