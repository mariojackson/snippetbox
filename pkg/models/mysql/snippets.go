package mysql

import (
    "database/sql"
    "jackson.software/snippetbox/pkg/models"
)

type SnippetRepository struct {
    DB *sql.DB
}

func (m *SnippetRepository) Insert(title, content, expires string) (int, error) {
    stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

    result, err := m.DB.Exec(stmt, title, content, expires)
    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return int(id), nil
}

func (m *SnippetRepository) Get(id int) (*models.Snippet, error) {
    return nil, nil
}

func (m *SnippetRepository) Latest() ([]*models.Snippet, error) {
    return nil, nil
}
