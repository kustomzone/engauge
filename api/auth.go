package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/EngaugeAI/engauge/db"
	"github.com/EngaugeAI/engauge/types"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Login --
func Login(c echo.Context) error {
	var adminCreds types.AdminCreds
	err := json.NewDecoder(c.Request().Body).Decode(&adminCreds)
	if err != nil {
		return echo.ErrBadRequest
	}

	token, err := genToken(adminCreds.Username)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(15 * time.Minute),
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":       token,
		"tokenExpiry": time.Duration(15 * time.Minute).Seconds(),
	})
}

// Logout --
func Logout(c echo.Context) error {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrCookieNotFound
	}
	cookie.Value = ""
	cookie.Expires = time.Time{}
	cookie.MaxAge = -1
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":       "",
		"tokenExpiry": 0,
	})
}

// RefreshToken --
func RefreshToken(c echo.Context) error {
	cookie, err := c.Request().Cookie("token")
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrCookieNotFound
	}
	tString := cookie.Value

	token, err := jwt.Parse(tString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(db.GlobalSettings.JWTSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var username string
		un, ok := claims["username"]
		if ok {
			username = un.(string)
		}
		if int(claims["sub"].(float64)) == 1 {
			t, err := genToken(username)
			if err != nil {
				c.Logger().Error(err)
				return echo.ErrInternalServerError
			}

			c.SetCookie(&http.Cookie{
				Name:    "token",
				Value:   t,
				Expires: time.Now().Add(15 * time.Minute),
			})

			return c.JSON(http.StatusOK, map[string]interface{}{
				"token":       t,
				"tokenExpiry": time.Duration(15 * time.Minute).Seconds(),
			})
		}

		return echo.ErrUnauthorized
	}

	return err
}

func genToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		claims = jwt.MapClaims{}
	}
	claims["username"] = username
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	claims["sub"] = float64(1)

	t, err := token.SignedString([]byte(db.GlobalSettings.JWTSecret))
	if err != nil {
		return t, err
	}

	return t, nil
}
