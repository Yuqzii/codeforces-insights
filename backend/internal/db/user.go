package db

import "context"

func (db *db) UpsertUser(ctx context.Context, handle string, rating int) (user_id int, err error) {
	err = db.conn.QueryRow(ctx, `
		INSERT INTO users (handle, current_rating) VALUES ($1, $2)
		ON CONFLICT (handle) DO UPDATE
		SET current_rating = EXCLUDED.current_rating
		RETURNING id`,
		handle, rating).Scan(&user_id)
	return user_id, err
}
