package repository

import (
    "context"
    "database/sql"
    "fmt"
    "log"
)

type User struct {
    ID   int32
    Name string
}

type UserRepository interface {
    Create(ctx context.Context, name string) (int32, error)
    GetByID(ctx context.Context, id int32) (string, error)
    Delete(ctx context.Context, id int32) (string, error)
}

type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, name string) (int32, error) {
    query := `INSERT INTO users (name) VALUES ($1) RETURNING user_id`
    var userID int32
    err := r.db.QueryRowContext(ctx, query, name).Scan(&userID)
    if err != nil {
        log.Printf("Create user failed: %v", err)
        return 0, err
    }
    return userID, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id int32) (string, error) {
    query := `SELECT name FROM users WHERE user_id = $1`
    var name string
    err := r.db.QueryRowContext(ctx, query, id).Scan(&name)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("user not found")
        }
        log.Printf("Get user failed: %v", err)
        return "", err
    }
    return name, nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id int32) (string, error) {
    query := `DELETE FROM users WHERE user_id = $1`
    res, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        log.Printf("Delete user failed: %v", err)
        return "", err
    }

    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        return "", fmt.Errorf("no user found to delete")
    }

    return fmt.Sprintf("User with ID %d deleted successfully", id), nil
}