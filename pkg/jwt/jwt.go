package jwt

import (
	"errors"
	"log"
	"net/http"
	"os"
	"server/models"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt"
)

var Secret = os.Getenv("NOTIFY_SECRET")

var (
	ErrTokenExpired         = errors.New("token has expired")
	ErrTokenNotEligible     = errors.New("token is not eligible for refresh yet")
	ErrInvalidToken         = errors.New("invalid token")
	ErrCouldNotParseToken   = errors.New("couldn't parse claims")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)

type Token struct {
	ID   int
	Role models.UserRole
}

func NewJwtService() *Token {
	return &Token{}
}

func (t Token) Init(userId int, userRole string) (string, error) {
	// Create JWT claims
	claims := jwt.MapClaims{
		"id":   strconv.Itoa(userId),
		"role": userRole,
		"exp":  time.Now().Add(time.Hour * 72).Unix(), // Token expiration time
	}

	// Create a new JWT token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}

	generatedToken := tokenString
	return generatedToken, nil
}

func (t Token) Parse(tok *jwtauth.JWTAuth, cookie string) (*Token, error) {
	token, err := jwtauth.VerifyToken(tok, cookie)
	if err != nil {
		return nil, err
	}
	claims := token.PrivateClaims()

	claimsID, ok := claims["id"].(string)
	id, err := strconv.Atoi(claimsID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("invalid token payload, id")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid token payload, role")
	}

	return &Token{
		ID:   id,
		Role: models.RoleStrConv(role),
	}, nil
}

func (t Token) ParseFromHeader(header string) (*Token, error) {
	// Split the header to extract the token part
	tokenStr := header

	log.Println("header:" + tokenStr)
	// Parse the JWT token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		// Check if the error is a JWT validation error
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Extract claims from the token
	claimsID, ok := claims["id"].(string)
	if !ok {
		return nil, errors.New("invalid or missing 'id' claim")
	}

	id, err := strconv.Atoi(claimsID)
	if err != nil {
		return nil, err
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid or missing 'role' claim")
	}

	return &Token{
		ID:   id,
		Role: models.RoleStrConv(role),
	}, nil
}

func (t Token) ParseFromCookieString(cookie string) (*Token, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cookie, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		// Check if the error is a JWT validation error
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claimsID, ok := claims["id"].(string)
	id, err := strconv.Atoi(claimsID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("invalid or missing 'id' claim")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid or missing 'role' claim")
	}
	return &Token{
		ID:   id,
		Role: models.RoleStrConv(role),
	}, nil
}

func (t Token) RefreshJwtToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return "", ErrInvalidToken
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	if err != nil {
		return "", ErrCouldNotParseToken
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", ErrCouldNotParseToken
	}

	if time.Until(time.Unix(int64(exp/1000), 0)) > 6*time.Hour {
		return "", ErrTokenNotEligible
	}
	//
	// Create a new token using the existing claims (e.g., id, name, role)
	claimsID, ok := claims["id"].(string)
	userID, err := strconv.Atoi(claimsID)
	if err != nil || !ok {
		return "", err
	}

	userRole, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("invalid or missing 'role' claim")
	}

	// Create a new token
	newTokenString, err := t.Init(userID, models.RoleStrConv(userRole).String())
	if err != nil {
		return "", err
	}

	return newTokenString, nil
}

func (t Token) DeleteJwtCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
		Path:    "/",
	})
}
