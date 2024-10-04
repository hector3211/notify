package service

import (
	"database/sql"
	"log"
	"server/models"
	"strconv"

	"github.com/hector3211/shogun"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (u UserService) CreateUser(user *models.User) *models.UserResponse {
	hashedPassword, err := u.HashPassword(user.Password)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	query := shogun.NewInsertBuilder().
		Insert("users").
		Columns("first_name", "last_name", "email", "password", "role").
		Values(user.FirstName, user.LastName, user.Email, hashedPassword, user.Role.String())

	result, err := u.db.Exec(query.Build())
	if err != nil {
		log.Fatal(err)
		return nil
	}

	userId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &models.UserResponse{
		ID:        int(userId),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      models.UserRole(user.Role),
	}
}

func (u UserService) DeleteUser(userId string) error {
	id, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}
	query := shogun.Delete("users").Where(shogun.Equal("id", id))

	_, err = u.db.Exec(query.Build())
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) ChangeUserRole(userId string, newRole string) error {
	id, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}

	stmt := shogun.Update("users").
		Set(shogun.Equal("role", newRole)).
		Where(shogun.Equal("id", id))

	_, err = u.db.Exec(stmt.Build())
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) GetAllUsers() []models.User {
	var users []models.User

	query := shogun.Select(
		"id",
		"first_name",
		"last_name",
		"email",
		"role",
		"created_at",
	).From("users").OrderBy("created_at").Desc()
	rows, err := u.db.Query(query.Build())
	if err != nil {
		return users
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			return users
		}
		users = append(users, user)
	}

	return users
}

func (u UserService) GetUserByEmail(logInEmail string) *models.UserResponse {
	var user models.UserResponse

	query := shogun.Select("id", "first_name", "last_name", "role").
		From("users").
		Where(shogun.Equal("email", logInEmail))
	// fmt.Printf("GetUserByEmail query: %s\n", query.String())
	//
	err := u.db.QueryRow(query.Build()).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Printf("error getting user existsence: %v", err)
		return nil
	}
	return &user
}

func (u UserService) GetUserByID(userId int) *models.UserResponse {
	var user models.UserResponse

	query := shogun.Select("id", "first_name", "last_name", "role").
		From("users").
		Where(shogun.Equal("id", userId))
		// fmt.Printf("GetUserByEmail query: %s\n", query.String())

	err := u.db.QueryRow(query.Build()).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Role)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Printf("error getting user existsence: %v", err)
		return nil
	}
	return &user
}

func (u UserService) GetUserHash(userId int) (string, error) {
	var hash string

	query := shogun.Select("password").
		From("users").
		Where(shogun.Equal("id", userId))

	err := u.db.QueryRow(query.Build()).Scan(&hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (u UserService) CheckIfUsersExists(userId int) bool {
	query := shogun.Select("id").From("users").Where(shogun.Equal("id", userId))
	var id int

	err := u.db.QueryRow(query.Build()).Scan(&id)
	return err == nil
}

func (u UserService) CheckEmailExists(userEmail string) bool {
	query := shogun.Select("email").From("users").Where(shogun.Equal("email", userEmail))

	var email string
	err := u.db.QueryRow(query.Build()).Scan(&email)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Printf("error checking email existsence: %v", err)
		return false
	}

	return err == nil
}

func (u UserService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (u UserService) CheckPasswordHash(password string, userId int) bool {
	hash, err := u.GetUserHash(userId)
	if err != nil {
		log.Printf("Error fetching user hash: %v", err)
		return false
	}
	// log.Printf("Hash from DB: %s", hash) // Log the hash for debugging
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("userID: %d Password comparison failed: %v\n", userId, err)
		return false
	}
	return true
}
