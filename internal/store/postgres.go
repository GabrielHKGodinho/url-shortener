package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore conecta no banco e cria a tabela se não existir
func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Verifica se a conexão está viva
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Criação da tabela (Migration simplificada)
	query := `
	CREATE TABLE IF NOT EXISTS links (
		id SERIAL PRIMARY KEY,
		original_url TEXT NOT NULL,
		short_code TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		clicks INT DEFAULT 0
	);`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("falha ao criar tabela: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Save(l Link) error {
	query := `INSERT INTO links (original_url, short_code) VALUES ($1, $2)`
	_, err := s.db.Exec(query, l.OriginalURL, l.ShortCode)
	return err
}

func (s *PostgresStore) Get(code string) (Link, error) {
	query := `SELECT id, original_url, short_code, created_at FROM links WHERE short_code = $1`

	row := s.db.QueryRow(query, code)

	var l Link
	// Scan copia as colunas do SQL para os campos da Struct
	// A ordem DEVE ser a mesma do SELECT
	// created_at pode vir como string ou time.Time, o driver converte
	var createdAt string
	if err := row.Scan(&l.ID, &l.OriginalURL, &l.ShortCode, &createdAt); err != nil {
		return Link{}, err
	}
	l.CreatedAt = createdAt // Simplificação por enquanto

	return l, nil
}

func (s *PostgresStore) IncrementClick(code string) error {
	query := `UPDATE links set clicks = clicks + 1 where short_code = $1`

	if _, error := s.db.Exec(query, code); error != nil {
		return error
	}

	return nil
}
