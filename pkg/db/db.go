package db

import (
    "database/sql"
    "fmt"
    "os"

    _ "modernc.org/sqlite"
)

const (
    // Создаю SQL-схему (таблицу и индекс)
    schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    date    CHAR(8) NOT NULL DEFAULT '',
    title   VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT    NOT NULL DEFAULT '',
    repeat  VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`
)

var db *sql.DB

// Открывает или создаёт файл базы и, если он новый, выполняет schema
func Init(dbFile string) error {
    _, err := os.Stat(dbFile)
    install := os.IsNotExist(err)

    // Открываю базу (файл будет создан автоматически, если его нет)
    datasource := fmt.Sprintf("file:%s?_foreign_keys=1", dbFile)
    conn, err := sql.Open("sqlite", datasource)
    if err != nil {
        return fmt.Errorf("db open: %w", err)
    }

    if install {
        if _, err := conn.Exec(schema); err != nil {
            conn.Close()
            return fmt.Errorf("db init schema: %w", err)
        }
    }

    db = conn
    return nil
}

func Close() error {
    if db != nil {
        return db.Close()
    }
    return nil
}
