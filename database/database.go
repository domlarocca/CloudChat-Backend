package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
)

var DB *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the MySQL database
	DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	log.Println("Connected to the database")
}

func RegisterUser(email, birthday, username, password string) error {
	if usernameExists(username) {
		return errors.New("Username already in use")
	}
	if emailExists(email) {
		return errors.New("Email already in use")
	}

	if err := validatePassword(password); err != nil {
		return err
	}

	tx, err := DB.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	userQuery := "INSERT INTO users (email, birthday, username) VALUES (?, ?, ?)"
	result, err := tx.Exec(userQuery, email, birthday, username)
	if err != nil {
		log.Println("Error registering user:", err)
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error obtaining user ID:", err)
		return err
	}

	credentialsQuery := "INSERT INTO user_credentials (userID, password) VALUES (?, ?)"
	_, err = tx.Exec(credentialsQuery, userID, password)
	if err != nil {
		log.Println("Error registering user credentials:", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return err
	}
	log.Println("User registered successfully")
	return nil
}

func LoginUser(username, password string) error {
	// Check if the user exists & grab their ID
	userID, err := getUserIDByUsername(username)
	if err != nil {
		return err
	}

	if userID == 0 {
		return errors.New("user not found")
	}

	// Validate the password
	if err := checkUserPassword(userID, password); err != nil {
		return err
	}

	log.Println("User logged in successfully")
	return nil
}

func getUserIDByUsername(username string) (int64, error) {
	var userID int64
	err := DB.QueryRow("SELECT userID FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // User not found
		}
		log.Println("Error querying user ID:", err)
		return 0, err
	}
	return userID, nil
}

func checkUserPassword(userID int64, password string) error {
	var storedPassword string
	err := DB.QueryRow("SELECT password FROM user_credentials WHERE userID = ?", userID).Scan(&storedPassword)
	if err != nil {
		log.Println("Error querying user password:", err)
		return err
	}

	// Validate the password (you may want to use a more secure password validation)
	if storedPassword != password {
		return errors.New("Incorrect password")
	}

	return nil
}

func usernameExists(username string) bool {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		log.Println("Error checking username existence:", err)
		return true
	}
	return count > 0
}

func emailExists(email string) bool {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		log.Println("Error checking email existence:", err)
		return true
	}
	return count > 0
}

func validatePassword(password string) error {
	// minimum 8 characters, at least one number, and one special character
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}
	if matched, _ := regexp.MatchString(`\d`, password); !matched {
		return errors.New("Password must contain at least one numeric character")
	}
	if matched, _ := regexp.MatchString(`[!@#$%^&*(),.?":{}|<>]`, password); !matched {
		return errors.New("Password must contain at least one special character")
	}
	return nil
}
