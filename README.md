# ПЗ №5 Подключение к PostgreSQL через database/sql. Выполнение простых запросов (INSERT, SELECT)
# Седова Мария Александровна, ЭФМО-01-25
# 21.10.2025

## Описание окружения: 
- go1.25.1 windows/amd64,
- PostgreSQL 17.6,
- ОС Windows 11.
  
## Скриншоты:
-	создание БД/таблицы в psql
<img width="836" height="418" alt="image" src="https://github.com/user-attachments/assets/af1ce564-f815-4dfd-b732-fca33f8923f8" />

- успешный вывод go run . (вставка и список задач);
<img width="870" height="263" alt="image" src="https://github.com/user-attachments/assets/8835bcf8-6255-4b88-a690-95e2373e49df" />

- SELECT * FROM tasks; в psql /не знаю что у меня с кодировкой, она не чинится
<img width="816" height="358" alt="image" src="https://github.com/user-attachments/assets/e46942a7-876d-4eb0-9098-5a593edd847c" />

- Доп задания
<img width="769" height="631" alt="image" src="https://github.com/user-attachments/assets/59d6e272-35b1-4db0-b89a-4aa5318c0cc3" />
<img width="613" height="260" alt="image" src="https://github.com/user-attachments/assets/a48c90c5-c399-4d62-853f-2559d5c2ef9e" />

## Код: db.go, repository.go, фрагменты main.go.

## 📄 **main.go**
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:password@localhost:5432/todo?sslmode=disable"
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("openDB error: %v", err)
	}
	defer db.Close()

	repo := NewRepo(db)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	titles := []string{"Сделать ПЗ №5", "Купить кофе", "Проверить отчёты"}
	for _, title := range titles {
		id, err := repo.CreateTask(ctx, title)
		if err != nil {
			log.Fatalf("CreateTask error: %v", err)
		}
		log.Printf("Inserted task id=%d (%s)", id, title)
	}

	fmt.Println("\n=== Все задачи ===")
	tasks, err := repo.ListTasks(ctx)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n", t.ID, t.Title, t.Done, t.CreatedAt.Format("2006-01-02 15:04"))
	}

	fmt.Println("\n=== Невыполненные задачи ===")
	undoneTasks, err := repo.ListDone(ctx, false)
	if err != nil {
		log.Fatalf("ListDone error: %v", err)
	}
	for _, t := range undoneTasks {
		fmt.Printf("#%d | %s\n", t.ID, t.Title)
	}

	fmt.Println("\n=== Поиск задачи с ID=1 ===")
	task, err := repo.FindByID(ctx, 1)
	if err != nil {
		log.Fatalf("FindByID error: %v", err)
	}
	fmt.Printf("Найдено: #%d | %s | done=%v\n", task.ID, task.Title, task.Done)

	fmt.Println("\n=== Массовая вставка ===")
	newTitles := []string{"Новая задача 1", "Новая задача 2", "Новая задача 3"}
	err = repo.CreateMany(ctx, newTitles)
	if err != nil {
		log.Fatalf("CreateMany error: %v", err)
	}
	fmt.Println("Добавлено 3 новые задачи")

	fmt.Println("\n=== Проверка пула соединений ===")
	fmt.Printf("MaxOpenConns: %d\n", db.Stats().MaxOpenConnections)
	fmt.Printf("OpenConnections: %d\n", db.Stats().OpenConnections)
	fmt.Printf("InUse: %d\n", db.Stats().InUse)
	fmt.Printf("Idle: %d\n", db.Stats().Idle)
}

package main

import (
	"context"
	"database/sql"
	"time"
)

type Task struct {
	ID        int
	Title     string
	Done      bool
	CreatedAt time.Time
}

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo { return &Repo{DB: db} }

func (r *Repo) CreateTask(ctx context.Context, title string) (int, error) {
	var id int
	const q = `INSERT INTO tasks (title) VALUES ($1) RETURNING id;`
	err := r.DB.QueryRowContext(ctx, q, title).Scan(&id)
	return id, err
}

func (r *Repo) ListTasks(ctx context.Context) ([]Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks ORDER BY id;`
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repo) ListDone(ctx context.Context, done bool) ([]Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks WHERE done = $1 ORDER BY id;`
	rows, err := r.DB.QueryContext(ctx, q, done)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repo) FindByID(ctx context.Context, id int) (*Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks WHERE id = $1;`
	var t Task
	err := r.DB.QueryRowContext(ctx, q, id).Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repo) CreateMany(ctx context.Context, titles []string) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const q = `INSERT INTO tasks (title) VALUES ($1);`
	for _, title := range titles {
		_, err := tx.ExecContext(ctx, q, title)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

1. Пул соединений sql.DB и зачем его настраивать?
Пул соединений хранит несколько подключений к базе данных готовыми к использованию. Это ускоряет работу, потому что не нужно каждый раз заново подключаться. 

2. Почему используем плейсхолдеры $1, $2?
Чтобы защититься от SQL-инъекций, когда злоумышленник может вставить вредоносный код. 

3. Чем Query, QueryRow и Exec отличаются?
Query - для нескольких строк результата
QueryRow - для одной строки
Exec - для команд без возврата данных (INSERT, UPDATE, DELETE)

4. Обоснование настроек пула:
SetMaxOpenConns(10) - максимум 10 одновременных подключений
SetMaxIdleConns(5) - 5 подключений в режиме ожидания
SetConnMaxLifetime(30 минут) - переподключение каждые 30 минут для надежности
