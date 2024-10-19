package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"server/models"
	"server/pkg/jwt"
	"server/utils"
	"strings"
	"time"
)

var protectedURLs = []*regexp.Regexp{
	regexp.MustCompile("^/admin$"),
	regexp.MustCompile("^/admin/account$"),
	regexp.MustCompile("^/admin/jobs$"),
	regexp.MustCompile("^/admin/jobs/new$"),
	regexp.MustCompile("^/admin/jobs/edit$"),
	regexp.MustCompile("^/admin/users$"),
	regexp.MustCompile("^/admin/users/new$"),
	regexp.MustCompile("^/admin/users/edit$"),
	regexp.MustCompile("^/profile$"),
}

type UserContext struct {
	Role models.UserRole
	ID   int
}

type UserContextKey string

var UserKey UserContextKey = "user"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAuth(r) {
			next.ServeHTTP(w, r)
			return
		}
		userCtx := GetUserCtxFromCookie(w, r)
		if userCtx == nil {
			http.Error(w, "Token parsing error in middlware", http.StatusUnauthorized)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			log.Println("Error retrieving JWT cookie")
			return
		}

		log.Printf("Setting user context: %+v", userCtx)
		c := context.WithValue(r.Context(), UserKey, userCtx)
		next.ServeHTTP(w, r.WithContext(c))
	})
}

func GetUserCtx(r *http.Request) *UserContext {
	if user, ok := r.Context().Value(UserKey).(*UserContext); ok && user != nil {
		return user
	}
	return nil
}

func GetUserCtxFromCookie(w http.ResponseWriter, r *http.Request) *UserContext {
	// Retrieve JWT cookie
	jwtCookie, err := r.Cookie("jwt")
	if err != nil {
		log.Printf("no JWT cookie: %v", err)
		return nil
	}

	// Parse JWT token
	user, err := jwt.NewJwtService().ParseFromCookieString(jwtCookie.Value)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Println("Access token expired, attemping to refresh")
			refreshedTokenStr, err := jwt.NewJwtService().RefreshJwtToken(r)
			if err != nil {
				if errors.Is(err, jwt.ErrTokenNotEligible) {
					log.Printf("token not eligible for refresh, continuing with the current token: %v", err)
					return createUserContext(&jwt.Token{ID: user.ID, Role: user.Role})
				}
				log.Printf("failed to refresh token: %v", err)
				return nil
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "jwt",
				Value:    refreshedTokenStr,
				Expires:  time.Now().Add(72 * time.Hour),
				HttpOnly: true,
				Secure:   utils.IsProduction(),
				SameSite: http.SameSiteStrictMode,
				Path:     "/",
			})

			user, err = jwt.NewJwtService().ParseFromCookieString(refreshedTokenStr)
			if err != nil {
				log.Printf("Error parsing new access token: %v", err)
				return nil
			}
		} else {
			log.Printf("JWT token parsing error: %v", err)
			return nil
		}
	}

	// Construct UserContext
	userCtx := createUserContext(&jwt.Token{ID: user.ID, Role: user.Role})
	return userCtx
}

func GetUserCtxFromHeader(w http.ResponseWriter, r *http.Request) (*UserContext, bool) {
	authHeader := r.Header.Get("Authorization")
	log.Printf("Authorization header: %s", authHeader)
	if !strings.HasPrefix(authHeader, "Bearer ") {
		log.Println("Authorization header format is incorrect")
		return nil, false
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := jwt.NewJwtService().ParseFromHeader(tokenStr) // Make sure to create a method that accepts the token string
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Println("Access token expired, attempting to refresh")
			refreshedTokenStr, err := jwt.NewJwtService().RefreshJwtToken(r) // Ensure this handles refreshing
			if err != nil {
				if errors.Is(err, jwt.ErrTokenNotEligible) {
					log.Printf("token not eligible for refresh, continuing with the current token: %v", err)
					return createUserContext(&jwt.Token{ID: user.ID, Role: user.Role}), true
				}
				log.Printf("failed to refresh token: %v", err)
				return nil, false
			}

			// Set the new token in the Authorization header for future requests
			w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", refreshedTokenStr))

			user, err = jwt.NewJwtService().ParseFromHeader(refreshedTokenStr)
			if err != nil {
				log.Printf("Error parsing new access token: %v", err)
				return nil, false
			}
		} else {
			log.Printf("JWT token parsing error: %v", err)
			return nil, false
		}
	}

	// Construct UserContext
	userCtx := createUserContext(&jwt.Token{ID: user.ID, Role: user.Role})
	return userCtx, true
}

func isAuth(r *http.Request) bool {
	originalURL := strings.ToLower(r.URL.String())

	for _, pattern := range protectedURLs {
		if pattern.MatchString(originalURL) {
			_, err := r.Cookie("jwt")
			if err != nil {
				log.Printf("JWT cookie missing or invalid for protected URL: %s", originalURL)
				return false
			}
			return true
		}
	}
	return false
}

func createUserContext(user *jwt.Token) *UserContext {
	return &UserContext{
		ID:   user.ID,
		Role: user.Role,
		// Role: models.ADMIN,
	}
}
