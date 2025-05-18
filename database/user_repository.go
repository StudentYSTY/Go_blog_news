package database

import (
	"myproject/models"
	"time"
)

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := DB.QueryRow(`
		SELECT id, username, password, email, is_blocked, created_at, updated_at 
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, 
		&user.IsBlocked, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := DB.QueryRow(`
		SELECT id, username, password, email, is_blocked, created_at, updated_at 
		FROM users WHERE username = $1
	`, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, 
		&user.IsBlocked, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) (int, error) {
	var id int
	err := DB.QueryRow(`
		INSERT INTO users (username, password, email, is_blocked, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, user.Username, user.Password, user.Email, user.IsBlocked, user.CreatedAt, user.UpdatedAt).Scan(&id)
	
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	_, err := DB.Exec(`
		UPDATE users
		SET username = $1, password = $2, email = $3, is_blocked = $4, updated_at = $5
		WHERE id = $6
	`, user.Username, user.Password, user.Email, user.IsBlocked, user.UpdatedAt, user.ID)
	return err
}

func DeleteUser(id int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	_, err = tx.Exec("DELETE FROM comments WHERE user_id = $1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM comments WHERE news_id IN (SELECT id FROM news WHERE author_id = $1)", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM news WHERE author_id = $1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetAllUsers() ([]models.User, error) {
	rows, err := DB.Query(`
		SELECT id, username, password, email, is_blocked, created_at, updated_at
		FROM users
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Password, &user.Email, 
			&user.IsBlocked, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	
	return users, nil
}

func BlockUser(id int) error {
	_, err := DB.Exec(`
		UPDATE users
		SET is_blocked = TRUE, updated_at = $1
		WHERE id = $2
	`, time.Now(), id)
	return err
}

func UnblockUser(id int) error {
	_, err := DB.Exec(`
		UPDATE users
		SET is_blocked = FALSE, updated_at = $1
		WHERE id = $2
	`, time.Now(), id)
	return err
} 