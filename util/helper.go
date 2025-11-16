package util

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetConfigPathFromGoMod(folderName string) string {
	return filepath.Join(findGoModuleRoot(), folderName)
}

func GetLocale(ginContext *gin.Context) string {
	if v, ok := ginContext.Get(LocaleContextKey); ok {
		return v.(string)
	}
	return "en"
}

func GetTraceId(ginContext *gin.Context) string {
	if v, ok := ginContext.Get(TraceIDContextKey); ok {
		return v.(string)
	}
	return ""
}

// ParseDuration support “d” (day), “h”, “m”, “s”
func ParseDuration(s string) (time.Duration, error) {
	if len(s) == 0 {
		return 0, nil
	}

	// If it has 'd' (day) then convert manually
	if len(s) > 1 && s[len(s)-1] == 'd' {
		daysString := s[:len(s)-1]
		days, err := time.ParseDuration(daysString + "h") // trick: parse like hours
		if err == nil {
			return days * 24, nil
		}
	}

	return time.ParseDuration(s)
}

func HashPassword(password string) string {
	// bcrypt.DefaultCost is 10, good balance between security and performance
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to hash password: " + err.Error())
	}
	return string(bytes)
}

func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		log.Println("Password mismatch:", err)
	}
	return err == nil
}

func findGoModuleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("Failed to find current directory" + err.Error())
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	panic("failed to find go.mod")
}
