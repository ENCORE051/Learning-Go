// main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "gouser:GoPass123@tcp(127.0.0.1:3306)/go_learning?parseTime=true&charset=utf8mb4"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// sql.Open 创建连接池，PingContext 才会实际检查连接。
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("连接 MySQL 失败：", err)
	}

	fmt.Println("连接 MySQL 成功！")

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(100) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// 用 ? 传入参数，避免通过拼接字符串构造 SQL。
	result, err := db.ExecContext(ctx,
		"INSERT INTO users (name) VALUES (?)", "小明")
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	var name string
	err = db.QueryRowContext(ctx,
		"SELECT name FROM users WHERE id = ?", id,
	).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("查询成功：id=%d，name=%s\n", id, name)
}
