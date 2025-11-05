package util

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
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
