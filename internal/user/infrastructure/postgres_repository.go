package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
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
func (r *PostgresUserRepository) Save(user domain.User) error {
	ctx := context.Background()
	
	if activeUser, ok := domain.IsActive(user); ok {
		return r.saveActiveUser(ctx, activeUser)
	}
	
	if provisionalUser, ok := domain.IsProvisional(user); ok {
		return r.saveProvisionalUser(ctx, provisionalUser)
	}
	
	return fmt.Errorf("unknown user type")
}

// saveActiveUser saves an active user
func (r *PostgresUserRepository) saveActiveUser(ctx context.Context, user *domain.ActiveUser) error {
	query := `
		INSERT INTO users (id, auth0_user_id, email, display_name, bio, profile_image_url, email_verified, created_at, verified_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (auth0_user_id) DO UPDATE SET
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
		user.Email().Value(),
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

// saveProvisionalUser saves a provisional user
func (r *PostgresUserRepository) saveProvisionalUser(ctx context.Context, user *domain.ProvisionalUser) error {
	query := `
		INSERT INTO users (id, auth0_user_id, email, created_at)
		VALUES ($1, '', $2, $3)
		ON CONFLICT (email) DO NOTHING
	`
	
	_, err := r.pool.Exec(ctx, query,
		user.ID(),
		user.Email().Value(),
		user.CreatedAt(),
	)
	
	return err
}

// FindByEmail finds a user by email
func (r *PostgresUserRepository) FindByEmail(email *domain.Email) (domain.User, error) {
	ctx := context.Background()
	
	query := `
		SELECT id, auth0_user_id, email, display_name, bio, profile_image_url, 
			   email_verified, created_at, verified_at, updated_at
		FROM users 
		WHERE email = $1
	`
	
	row := r.pool.QueryRow(ctx, query, email.Value())
	
	var id, auth0UserID, emailStr string
	var displayName, bio, profileImageURL sql.NullString
	var emailVerified bool
	var createdAt, updatedAt time.Time
	var verifiedAt sql.NullTime
	
	err := row.Scan(&id, &auth0UserID, &emailStr, &displayName, &bio, &profileImageURL,
		&emailVerified, &createdAt, &verifiedAt, &updatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	
	userEmail, err := domain.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create email: %w", err)
	}
	
	var verifiedAtPtr *time.Time
	if verifiedAt.Valid {
		verifiedAtPtr = &verifiedAt.Time
	}
	
	return domain.FromSnapshot(
		id,
		auth0UserID,
		userEmail,
		displayName.String,
		bio.String,
		profileImageURL.String,
		emailVerified,
		createdAt,
		verifiedAtPtr,
		updatedAt,
	), nil
}

// FindByID finds a user by ID
func (r *PostgresUserRepository) FindByID(id string) (domain.User, error) {
	ctx := context.Background()
	
	query := `
		SELECT id, auth0_user_id, email, display_name, bio, profile_image_url, 
			   email_verified, created_at, verified_at, updated_at
		FROM users 
		WHERE id = $1
	`
	
	row := r.pool.QueryRow(ctx, query, id)
	
	var userID, auth0UserID, emailStr string
	var displayName, bio, profileImageURL sql.NullString
	var emailVerified bool
	var createdAt, updatedAt time.Time
	var verifiedAt sql.NullTime
	
	err := row.Scan(&userID, &auth0UserID, &emailStr, &displayName, &bio, &profileImageURL,
		&emailVerified, &createdAt, &verifiedAt, &updatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	
	userEmail, err := domain.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create email: %w", err)
	}
	
	var verifiedAtPtr *time.Time
	if verifiedAt.Valid {
		verifiedAtPtr = &verifiedAt.Time
	}
	
	return domain.FromSnapshot(
		userID,
		auth0UserID,
		userEmail,
		displayName.String,
		bio.String,
		profileImageURL.String,
		emailVerified,
		createdAt,
		verifiedAtPtr,
		updatedAt,
	), nil
}

// FindByAuth0UserID finds a user by Auth0 user ID
func (r *PostgresUserRepository) FindByAuth0UserID(auth0UserID string) (domain.User, error) {
	ctx := context.Background()
	
	query := `
		SELECT id, auth0_user_id, email, display_name, bio, profile_image_url, 
			   email_verified, created_at, verified_at, updated_at
		FROM users 
		WHERE auth0_user_id = $1
	`
	
	row := r.pool.QueryRow(ctx, query, auth0UserID)
	
	var id, auth0ID, emailStr string
	var displayName, bio, profileImageURL sql.NullString
	var emailVerified bool
	var createdAt, updatedAt time.Time
	var verifiedAt sql.NullTime
	
	err := row.Scan(&id, &auth0ID, &emailStr, &displayName, &bio, &profileImageURL,
		&emailVerified, &createdAt, &verifiedAt, &updatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to find user by Auth0 user ID: %w", err)
	}
	
	userEmail, err := domain.NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create email: %w", err)
	}
	
	var verifiedAtPtr *time.Time
	if verifiedAt.Valid {
		verifiedAtPtr = &verifiedAt.Time
	}
	
	return domain.FromSnapshot(
		id,
		auth0ID,
		userEmail,
		displayName.String,
		bio.String,
		profileImageURL.String,
		emailVerified,
		createdAt,
		verifiedAtPtr,
		updatedAt,
	), nil
}