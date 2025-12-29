package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/singchia/liaison/pkg/utils"
)

func main() {
	var password string
	var email string
	var random bool
	var length int
	var hashOnly bool
	var createUser bool
	var status string

	flag.StringVar(&password, "password", "", "Password to hash (if not provided, will prompt for input)")
	flag.StringVar(&email, "email", "", "User email address (required for database operations)")
	flag.BoolVar(&random, "random", false, "Generate a random password")
	flag.IntVar(&length, "length", 16, "Length of random password (default: 16)")
	flag.BoolVar(&hashOnly, "hash-only", false, "Only output the hash, don't update database")
	flag.BoolVar(&createUser, "create", false, "Create new user if not exists")
	flag.StringVar(&status, "status", "active", "User status (active/inactive), default: active")
	flag.Parse()

	// 如果生成随机密码
	if random {
		randomPassword, err := utils.GenerateRandomPassword(length)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating random password: %v\n", err)
			os.Exit(1)
		}
		password = randomPassword
		fmt.Printf("Generated random password: %s\n", randomPassword)
		fmt.Println()
	}

	// 如果没有提供密码，提示输入
	if password == "" {
		fmt.Print("Enter password to hash: ")
		var err error
		_, err = fmt.Scanln(&password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
			os.Exit(1)
		}
	}

	// 生成密码哈希
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error hashing password: %v\n", err)
		os.Exit(1)
	}

	// 如果只需要输出哈希
	if hashOnly {
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("Password Hash (for database):")
		fmt.Println(hashedPassword)
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println()
		fmt.Println("Usage in SQL:")
		fmt.Printf("UPDATE users SET password = '%s' WHERE email = 'your-email@example.com';\n", hashedPassword)
		fmt.Println()
		fmt.Println("Or use in Go code:")
		fmt.Printf("hashedPassword := \"%s\"\n", hashedPassword)
		return
	}

	// 需要数据库操作，检查 email
	if email == "" {
		fmt.Fprintf(os.Stderr, "Error: email is required for database operations\n")
		fmt.Fprintf(os.Stderr, "Use -email flag or -hash-only to only generate hash\n")
		os.Exit(1)
	}

	// 获取数据库路径
	dbPath := getDBPath()

	// 打开数据库连接
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 检查数据库连接
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Printf("🔐 Liaison Password Generator\n")
	fmt.Printf("Database: %s\n", dbPath)
	fmt.Printf("Email: %s\n", email)
	fmt.Println(strings.Repeat("=", 50))

	// 检查用户是否存在
	var userID int
	var existingStatus string
	checkQuery := "SELECT id, status FROM users WHERE email = ?"
	err = db.QueryRow(checkQuery, email).Scan(&userID, &existingStatus)
	userExists := err == nil

	if userExists {
		// 更新现有用户
		fmt.Printf("✅ User found (ID: %d)\n", userID)
		fmt.Printf("📝 Updating password...\n")

		updateQuery := "UPDATE users SET password = ? WHERE email = ?"
		result, err := db.Exec(updateQuery, hashedPassword, email)
		if err != nil {
			log.Fatalf("Failed to update password: %v", err)
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			fmt.Printf("✅ Password updated successfully!\n")
		} else {
			fmt.Printf("⚠️  No rows updated\n")
		}
	} else {
		// 创建新用户
		if !createUser {
			fmt.Printf("❌ User with email %s not found\n", email)
			fmt.Printf("💡 Use -create flag to create a new user\n")
			os.Exit(1)
		}

		fmt.Printf("📝 Creating new user...\n")
		insertQuery := "INSERT INTO users (email, password, status) VALUES (?, ?, ?)"
		result, err := db.Exec(insertQuery, email, hashedPassword, status)
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}

		userID64, _ := result.LastInsertId()
		fmt.Printf("✅ User created successfully! (ID: %d)\n", userID64)
	}

	fmt.Println()
	fmt.Println("🔑 Password Hash:")
	fmt.Println(hashedPassword)
	fmt.Println()
	fmt.Println("✅ Operation completed!")
}

// getDBPath 获取数据库路径
func getDBPath() string {
	// 检查常见的数据库路径
	possiblePaths := []string{
		"/opt/liaison/data/liaison.db",
		"./etc/liaison.db",
		"./liaison.db",
		"./data/liaison.db",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 如果都找不到，使用默认路径
	return "/opt/liaison/data/liaison.db"
}
