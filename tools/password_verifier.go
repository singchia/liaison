package main

import (
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/argon2"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./password-verifier <email> <password>")
		fmt.Println("Example: ./password-verifier default@liaison.local mypassword")
		os.Exit(1)
	}

	email := os.Args[1]
	password := os.Args[2]

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

	fmt.Printf("🔐 Liaison Password Verifier\n")
	fmt.Printf("Database: %s\n", dbPath)
	fmt.Printf("Email: %s\n", email)
	fmt.Println(strings.Repeat("=", 50))

	// 获取用户的密码哈希
	var hashedPassword string
	query := "SELECT password FROM users WHERE email = ?"
	err = db.QueryRow(query, email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("❌ User with email %s not found\n", email)
		} else {
			fmt.Printf("❌ Error querying user: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("✅ User found!\n")
	fmt.Printf("🔑 Stored hash: %s\n", hashedPassword)
	fmt.Printf("🔤 Input password: %s\n", password)
	fmt.Println()

	// 验证密码
	isValid, err := verifyPassword(password, hashedPassword)
	if err != nil {
		fmt.Printf("❌ Error verifying password: %v\n", err)
		os.Exit(1)
	}

	if isValid {
		fmt.Printf("✅ Password is CORRECT! 🎉\n")
	} else {
		fmt.Printf("❌ Password is INCORRECT! 🚫\n")
	}

	// 显示一些常见的默认密码供参考
	fmt.Println()
	fmt.Println("💡 Common default passwords to try:")
	fmt.Println("   - default123")
	fmt.Println("   - password")
	fmt.Println("   - admin")
	fmt.Println("   - 123456")
	fmt.Println("   - liaison")
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

// verifyPassword 验证密码
func verifyPassword(password, hashedPassword string) (bool, error) {
	// 解析Argon2id哈希
	// 格式: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false, fmt.Errorf("invalid argon2id hash format")
	}

	// 解析参数
	var version int
	var memory, time, parallelism uint32
	var salt, hash []byte

	// 解析版本
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, fmt.Errorf("invalid version: %v", err)
	}

	// 解析内存、时间、并行度
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &parallelism); err != nil {
		return false, fmt.Errorf("invalid parameters: %v", err)
	}

	// 解码盐值和哈希
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("invalid salt: %v", err)
	}

	hash, err = base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("invalid hash: %v", err)
	}

	// 计算输入密码的哈希
	computedHash := argon2.IDKey([]byte(password), salt, time, memory, uint8(parallelism), uint32(len(hash)))

	// 比较哈希值
	return subtle.ConstantTimeCompare(hash, computedHash) == 1, nil
}
