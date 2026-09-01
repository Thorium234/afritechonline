package database

import (
	"context"
	"database/sql"
	"fmt"
)

// Tx wraps a database transaction with a callback pattern for clean transaction handling.
type Tx struct {
	*sql.Tx
}

// WithTx executes the given callback within a transaction.
// If the callback returns an error, the transaction is rolled back.
// If the callback succeeds, the transaction is committed.
// This pattern ensures transactions are always properly cleaned up.
func WithTx(db *sql.DB, ctx context.Context, callback func(tx *Tx) error) error {
	// Start transaction
	sqlTx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	
	tx := &Tx{sqlTx}
	
	// Execute callback
	if err := callback(tx); err != nil {
		// Rollback on error
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("callback failed: %w, rollback failed: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}
	
	// Commit on success
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	
	return nil
}

// WithTxDefault is a convenience function that uses context.Background() and *sql.DB.
// Use this only when you don't have a specific context to pass.
func WithTxDefault(db *sql.DB, callback func(tx *Tx) error) error {
	return WithTx(db, context.Background(), callback)
}

// Isolation level constants for transaction control (optional)
var (
	ReadCommitted   = sql.LevelReadCommitted
	ReadUncommitted = sql.LevelReadUncommitted
	Repeatable      = sql.LevelRepeatableRead
	Serializable    = sql.LevelSerializable
)

// WithTxIsolation executes a callback within a transaction with a specific isolation level.
func WithTxIsolation(db *sql.DB, ctx context.Context, isolationLevel sql.IsolationLevel, callback func(tx *Tx) error) error {
	sqlTx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: isolationLevel,
	})
	if err != nil {
		return fmt.Errorf("begin transaction with isolation level: %w", err)
	}
	
	tx := &Tx{sqlTx}
	
	if err := callback(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("callback failed: %w, rollback failed: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}
	
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	
	return nil
}
