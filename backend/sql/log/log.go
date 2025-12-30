package log

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type LogEntry struct {
	ID        int64
	Level     string
	Message   string
	CreatedAt time.Time
}

func SaveLogs(ctx context.Context, db *sql.DB, entries []LogEntry) error {
	if entries == nil {
		return nil
	}

	stmt, err := db.PrepareContext(ctx,
		"INSERT INTO logs (level,message) values(?,?)",
	)
	if err != nil {
		return nil
	}
	defer stmt.Close()
	for _, v := range entries {
		_, err := stmt.ExecContext(ctx, v.Level, v.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func FetchLogsByLevels(ctx context.Context, db *sql.DB, levels []string) (map[string][]LogEntry, error) {

	if len(levels) == 0 {
		fmt.Printf("Пустой levels")
		return map[string][]LogEntry{}, nil
	}
	logs := make(map[string][]LogEntry)

	stmt, err := db.PrepareContext(ctx,
		"SELECT id,level,message,created_at from logs where level = ?",
	)
	if err != nil {
		return map[string][]LogEntry{}, nil
	}
	defer stmt.Close()

	for _, v := range levels {

		rows, err := stmt.QueryContext(ctx, v)
		if err != nil {
			return nil, err
		}
		for rows.Next() {

			var l LogEntry
			err = rows.Scan(&l.ID, &l.Level, &l.Message, &l.CreatedAt)
			if err != nil {
				return nil, err
			}
			logs[v] = append(logs[v], LogEntry{
				ID:        l.ID,
				Level:     l.Level,
				Message:   l.Message,
				CreatedAt: l.CreatedAt,
			})

		}
		rows.Close()
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return logs, nil

}

func DoubleChangeAge(ctx context.Context, db *sql.DB, log1 LogEntry, log2 LogEntry) error {

	err := withTx(ctx, db, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO logs (level,message) values(?,?)",
			log1.Level, log1.Message,
		)
		if err != nil {
			return fmt.Errorf("Error: %w", err)
		}
		_, err = tx.ExecContext(ctx,
			"INSERT INTO logs (level,message) values(?,?)",
			log2.Level, log2.Message,
		)
		return err
	})

	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	return nil

}

func withTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})

	if err != nil {
		return err
	}

	defer tx.Rollback()
	if err = fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
