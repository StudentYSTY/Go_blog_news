package database

import (
	"myproject/models"
	"time"
)

func GetNewsByID(id int) (*models.News, error) {
	var news models.News
	err := DB.QueryRow(`
		SELECT id, title, content, author_id, author, created_at, updated_at 
		FROM news WHERE id = $1
	`, id).Scan(
		&news.ID, &news.Title, &news.Content, &news.AuthorID, 
		&news.Author, &news.CreatedAt, &news.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &news, nil
}

func CreateNews(news *models.News) (int, error) {
	var id int
	err := DB.QueryRow(`
		INSERT INTO news (title, content, author_id, author, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, news.Title, news.Content, news.AuthorID, news.Author, news.CreatedAt, news.UpdatedAt).Scan(&id)
	
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateNews(news *models.News) error {
	news.UpdatedAt = time.Now()
	_, err := DB.Exec(`
		UPDATE news
		SET title = $1, content = $2, updated_at = $3
		WHERE id = $4
	`, news.Title, news.Content, news.UpdatedAt, news.ID)
	return err
}

func DeleteNews(id int) error {
	_, err := DB.Exec("DELETE FROM news WHERE id = $1", id)
	return err
}

func GetAllNews() ([]models.News, error) {
	rows, err := DB.Query(`
		SELECT id, title, content, author_id, author, created_at, updated_at
		FROM news
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var news []models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(
			&n.ID, &n.Title, &n.Content, &n.AuthorID, 
			&n.Author, &n.CreatedAt, &n.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		news = append(news, n)
	}
	
	return news, nil
}

func GetNewsByAuthorID(authorID int) ([]models.News, error) {
	rows, err := DB.Query(`
		SELECT id, title, content, author_id, author, created_at, updated_at
		FROM news
		WHERE author_id = $1
		ORDER BY created_at DESC
	`, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var news []models.News
	for rows.Next() {
		var n models.News
		err := rows.Scan(
			&n.ID, &n.Title, &n.Content, &n.AuthorID, 
			&n.Author, &n.CreatedAt, &n.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		news = append(news, n)
	}
	
	return news, nil
} 