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
		dsn = "postgres://postgres:4471716@localhost:5432/todo?sslmode=disable"
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