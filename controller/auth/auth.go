package auth

import (
	"database/sql"	
    "encoding/json"
	"kuburan/database"
	"log"
	"net/http"
	"time"
	"strings"

	"kuburan/model/users"

	"github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

// func Registration(w http.ResponseWriter, r *http.Request) {
// 	var creds users.Users
// 	err := json.NewDecoder(r.Body).Decode(&creds)
//     if err != nil {
//         http.Error(w, "Invalid request payload", http.StatusBadRequest)
//         return
//     }

// 	var existingUser users.Users
// 	err = database.DB.QueryRow("SELECT username FROM users WHERE username = (?)", creds.Username).Scan(&existingUser.Username)

// 	if err != nil && err != sql.ErrNoRows {
//         http.Error(w, "Internal server error", http.StatusInternalServerError)
//         return
//     }

// 	// Hash password
//     hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
//     if err != nil {
//         http.Error(w, "Internal server error Password", http.StatusInternalServerError)
//         return
//     }

// 	_,err = database.DB.Exec("INSERT INTO users (username,email,password) VALUES (?,?,?)",creds.Username,creds.Email,hashedPassword)
// 	if err != nil {
//         http.Error(w, "Internal server error Insert", http.StatusInternalServerError)
//         return
//     }

// 	// Berikan respon sukses
//     w.Header().Set("Content-Type", "application/json")
//     response := map[string]interface{}{
//         "message": "User registered successfully",
//     }
//     err = json.NewEncoder(w).Encode(response)
//     if err != nil {
//         log.Printf("Error encoding response: %v", err)
//         http.Error(w, "Internal server error", http.StatusInternalServerError)
//     }
// }

func Registration(w http.ResponseWriter, r *http.Request) {
	var creds users.Users
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	var existingUser users.Users
	err = database.DB.QueryRow("SELECT username FROM users WHERE username = ?", creds.Username).Scan(&existingUser.Username)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if existingUser.Username != "" {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error during password hashing", http.StatusInternalServerError)
		return
	}

	// Insert new user into database
	_, err = database.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", creds.Username, creds.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Internal server error during user insertion", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "User registered successfully",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// func Login(w http.ResponseWriter, r *http.Request) {
// 	var creds users.Users
// 	err := json.NewDecoder(r.Body).Decode(&creds)
//     if err != nil {
//         http.Error(w, "Invalid request payload", http.StatusBadRequest)
//         return
//     }

// 	var user users.Users
//     err = database.DB.QueryRow("SELECT user_id, username, email, password FROM users WHERE username= (?)", creds.Username).Scan(&user.UserId, &user.Username, &user.Password, &user.Email)
//     if err != nil {
//         if err == sql.ErrNoRows {
//             http.Error(w, "User not found", http.StatusUnauthorized)
//             return
//         }
//         http.Error(w, "Internal server error", http.StatusInternalServerError)
//         return
//     }

// 	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
//     if err != nil {
//         http.Error(w, "Invalid password", http.StatusUnauthorized)
//         return
//     }

// 	expirationTime := time.Now().Add(120 * time.Minute)
//     claims := &Claims{
//         Username: creds.Username,
//         StandardClaims: jwt.StandardClaims{
//             ExpiresAt: expirationTime.Unix(),
//         },
//     }

//     token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//     tokenString, err := token.SignedString(jwtKey)
//     if err != nil {
//         http.Error(w, "Internal server error", http.StatusInternalServerError)
//         return
//     }

// 	w.Header().Set("Content-Type", "application/json")
//     response := map[string]interface{}{
//         "message": "Login successful",
//         "token":   tokenString,
//     }
//     err = json.NewEncoder(w).Encode(response)
//     if err != nil {
//         log.Printf("Error encoding response: %v", err)
//         http.Error(w, "Internal server error", http.StatusInternalServerError)
//     }
// }

func Login(w http.ResponseWriter, r *http.Request) {
	var creds users.Users
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch user from database
	var user users.Users
	err = database.DB.QueryRow("SELECT user_id, username, email, password FROM users WHERE username = ?", creds.Username).Scan(&user.UserId, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Compare the hashed password with the password provided by the user
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	expirationTime := time.Now().Add(120 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error during token generation", http.StatusInternalServerError)
		return
	}

	// Respond with the token
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		sttArr := strings.Split(bearerToken, " ")
		if len(sttArr) == 2 {
			isValid, _ := ValidateToken(sttArr[1])
			if isValid {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}
