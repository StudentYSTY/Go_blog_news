package database

import (
	"myproject/models"
	"time"
)

func GetCommentByID(id int) (*models.Comment, error) {
	var comment models.Comment
	err := DB.QueryRow(`
		SELECT id, news_id, user_id, username, content, created_at, updated_at 
		FROM comments WHERE id = $1
	`, id).Scan(
		&comment.ID, &comment.NewsID, &comment.UserID, &comment.Username, 
		&comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func CreateComment(comment *models.Comment) (int, error) {
	var id int
	err := DB.QueryRow(`
		INSERT INTO comments (news_id, user_id, username, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, comment.NewsID, comment.UserID, comment.Username, comment.Content, comment.CreatedAt, comment.UpdatedAt).Scan(&id)
	
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateComment(comment *models.Comment) error {
	comment.UpdatedAt = time.Now()
	_, err := DB.Exec(`
		UPDATE comments
		SET content = $1, updated_at = $2
		WHERE id = $3
	`, comment.Content, comment.UpdatedAt, comment.ID)
	return err
}

func DeleteComment(id int) error {
	_, err := DB.Exec("DELETE FROM comments WHERE id = $1", id)
	return err
}

func GetCommentsByNewsID(newsID int) ([]models.Comment, error) {
	rows, err := DB.Query(`
		SELECT id, news_id, user_id, username, content, created_at, updated_at
		FROM comments
		WHERE news_id = $1
		ORDER BY created_at ASC
	`, newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(
			&c.ID, &c.NewsID, &c.UserID, &c.Username, 
			&c.Content, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	
	return comments, nil
}

func GetCommentsByUserID(userID int) ([]models.Comment, error) {
	rows, err := DB.Query(`
		SELECT id, news_id, user_id, username, content, created_at, updated_at
		FROM comments
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(
			&c.ID, &c.NewsID, &c.UserID, &c.Username, 
			&c.Content, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	
	return comments, nil
}

func DeleteCommentsByNewsID(newsID int) error {
	_, err := DB.Exec("DELETE FROM comments WHERE news_id = $1", newsID)
	return err
} 