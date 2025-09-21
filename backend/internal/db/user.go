package db

import "context"

func (db *db) UpsertUser(ctx context.Context, handle string, rating int) error {
	_, err := db.conn.Exec(ctx, `
		INSERT INTO users (handle, current_rating) VALUES ($1, $2)
		ON CONFLICT (handle) DO UPDATE
		SET current_rating = EXCLUDED.current_rating
	`, handle, rating)
	return err
}
