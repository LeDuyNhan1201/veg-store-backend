package localizer

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"veg-store-backend/util"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

/*
This file handles internationalization (i18n) using go-i18n.
Logic:
- Initialize an i18n bundle with English as the default language.
- Load all .toml locale files from the ./i18n directory.
- Log loaded messages for debugging.
- Provide a Localize function to retrieve localized messages by Id and language.

Example .toml structure:
[Hello]
one = "Hello {{.DBName}}"
other = "Hello everyone"
*/

type Localizer struct {
	Bundle *i18n.Bundle
}

func Init(mode string) *Localizer {
	// Initialize i18n bundle
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Set path to i18n directory
	i18nPath := util.GetConfigPathFromGoMod("i18n")

	// Walk through i18n directory and load all .toml files
	loadI18nMessages(mode, bundle, i18nPath)

	zap.L().Info("All locale files loaded successfully.")
	return &Localizer{Bundle: bundle}
}

// T - Usage: Localizer.T("message_id", params) to get a localized message.
func (l *Localizer) T(locale string, msgID string, params ...map[string]interface{}) string {
	return l.Localize(locale, msgID, params...)
}

func (l *Localizer) Localize(lang, msgID string, params ...map[string]interface{}) string {
	// Create a l for the specified language
	specificLocalizer := i18n.NewLocalizer(l.Bundle, lang)

	/*
		Create a map for template repository if provided Example template repository:
		params := map[string]interface{}{
			"DBName": "John",
			"Age": "30",
		}
	*/

	// Get TemplateData if existed
	var templateData map[string]interface{}
	pluralCount := 1

	if len(params) > 0 {
		templateData = params[0]

		// Check if "count" key exists and is numeric
		if val, ok := templateData["Count"]; ok {
			switch v := val.(type) {
			case int:
				pluralCount = v
			case int32:
				pluralCount = int(v)
			case int64:
				pluralCount = int(v)
			case float64:
				pluralCount = int(v)
			}
		}
	}

	// Build localize config
	config := &i18n.LocalizeConfig{
		MessageID:    msgID,
		TemplateData: templateData,
		PluralCount:  pluralCount,
	}

	// Localize message
	msg, err := specificLocalizer.Localize(config)
	if err != nil {
		zap.L().Warn("Failed to localize message",
			zap.String("lang", lang),
			zap.String("Id", msgID),
			zap.Error(err),
		)
		return msgID // fallback
	}

	return msg
}

func loadI18nMessages(mode string, bundle *i18n.Bundle, absPath string) {
	err := filepath.WalkDir(absPath, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			zap.L().Fatal("Error walking through 'i18n' directory", zap.Error(err))
		}

		if dirEntry.IsDir() || filepath.Ext(path) != ".toml" {
			return nil
		}

		locale := extractLocaleFromFilename(path)
		if locale == "" {
			zap.L().Warn("Failed to determine locale from file", zap.String("file", path))
			return nil
		}

		messageFile, err := bundle.LoadMessageFile(path)
		if err != nil {
			zap.L().Warn("Failed to load message file",
				zap.String("file", path),
				zap.Error(err),
			)
			return nil
		}

		if mode != "prod" && mode != "production" {
			logLoadedLocaleMessages(locale, path, messageFile)
		}

		return nil
	})

	if err != nil {
		zap.L().Fatal("Error after walking through 'i18n' directory", zap.Error(err))
	}
}

func extractLocaleFromFilename(path string) string {
	base := filepath.Base(path)       // e.g. "active.vi.toml"
	parts := strings.Split(base, ".") // ["active", "vi", "toml"]

	if len(parts) >= 3 {
		return parts[len(parts)-2] // "vi"
	}
	return ""
}

func logLoadedLocaleMessages(locale, path string, messageFile *i18n.MessageFile) {
	var fields []zap.Field

	for _, message := range messageFile.Messages {
		fields = append(fields,
			zap.String(fmt.Sprintf("%s.%s.One", locale, message.ID), message.One),
			zap.String(fmt.Sprintf("%s.%s.Other", locale, message.ID), message.Other),
		)
	}

	zap.L().Info("Loaded locale messages",
		zap.String("locale", locale),
		zap.String("file", path),
		zap.Int("message_count", len(messageFile.Messages)),
	)
	zap.L().Debug("Message details", fields...) // log details only in debug mode
}
