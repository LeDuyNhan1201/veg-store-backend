package util

import (
	"log"
	"math/rand"
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

func GetTraceID(ginContext *gin.Context) string {
	if v, ok := ginContext.Get(TraceIDContextKey); ok {
		return v.(string)
	}
	return ""
}

func ParseDate(dayStr, timezone string) time.Time {
	format := "2006-01-02"
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		panic("Failed to load location: " + err.Error())
	}
	day, err := time.ParseInLocation(format, dayStr, loc)
	if err != nil {
		panic("Failed to load location: " + err.Error())
	}
	return day
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

func Map[T any, R any](in []T, f func(T) R) []R {
	out := make([]R, 0, len(in))
	for _, v := range in {
		out = append(out, f(v))
	}
	return out
}

func RandomDateTime() time.Time {
	from := time.Date(2000, 1, 1, 6, 0, 0, 0, time.UTC)
	to := time.Now()
	return RandomTimeBetween(from, to)
}

func RandomTimeBetween(from, to time.Time) time.Time {
	if from.After(to) {
		from, to = to, from
	}

	start := from.UnixNano()
	end := to.UnixNano()

	random := rand.Int63n(end-start) + start
	return time.Unix(0, random).In(from.Location())
}

func ParseDay(dayStr string) time.Time {
	format := "2006-01-02"
	day, err := time.Parse(format, dayStr)
	if err != nil {
		panic(err)
	}
	return day
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
