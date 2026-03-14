// 迁移程序：清空当前数据库后按顺序执行 db/migrations/*.sql（001_schema.sql → 002_seed.sql）
// 用法：go run ./db/cmd/migrate [-dsn "postgres://..."]  或设置环境变量 BEEHIVE_POSTGRES_DSN
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/lib/pq"
)

func main() {
	dsn := flag.String("dsn", "", "PostgreSQL DSN (default: $BEEHIVE_POSTGRES_DSN)")
	flag.Parse()
	if *dsn == "" {
		*dsn = os.Getenv("BEEHIVE_POSTGRES_DSN")
	}
	if *dsn == "" {
		*dsn = "postgres://beehive:beehive@127.0.0.1:5432/beehive?sslmode=disable"
	}

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open db: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "ping db: %v\n", err)
		os.Exit(1)
	}

	// 清空 public schema
	fmt.Println("Dropping schema public CASCADE...")
	_, err = db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		fmt.Fprintf(os.Stderr, "drop schema: %v\n", err)
		os.Exit(1)
	}

	// 从 repo 根目录用 go run ./db/cmd/migrate 时用 db/migrations；在 db 目录下时用 migrations
	var migrationsDir string
	for _, d := range []string{"db/migrations", "migrations"} {
		if _, err := os.Stat(d); err == nil {
			migrationsDir = d
			break
		}
	}
	if migrationsDir == "" {
		fmt.Fprintf(os.Stderr, "migrations dir not found (run from repo root or db/)\n")
		os.Exit(1)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read migrations dir: %v\n", err)
		os.Exit(1)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sql" {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		path := filepath.Join(migrationsDir, name)
		fmt.Printf("Running %s ...\n", name)
		body, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
			os.Exit(1)
		}
		_, err = db.Exec(string(body))
		if err != nil {
			fmt.Fprintf(os.Stderr, "exec %s: %v\n", name, err)
			os.Exit(1)
		}
	}

	fmt.Println("Migration done.")
}
