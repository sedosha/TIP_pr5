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

5.	Краткие ответы:
	Что такое пул соединений *sql.DB и зачем его настраивать?
	Почему используем плейсхолдеры $1, $2?
	Чем Query, QueryRow и Exec отличаются?
6.	Обоснование транзакций и настроек пула.
