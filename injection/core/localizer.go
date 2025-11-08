package core

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
- Load all JSON locale files from the ./i18n directory.
- Log loaded messages for debugging.
- Provide a Localize function to retrieve localized messages by ID and language.

Example JSON structure:
{
  "hello_world": {
	"message": "Hello, World!"
  }
}
*/

type Localizer struct {
	Bundle *i18n.Bundle
}

func InitI18n() *Localizer {
	// Initialize i18n bundle
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Set path to i18n directory
	i18nPath := util.GetConfigPathFromGoMod("i18n")

	// Walk through i18n directory and load all .toml files
	loadI18nMessages(bundle, i18nPath)

	Logger.Info("All locale files loaded successfully.")
	return &Localizer{Bundle: bundle}
}

// T - Usage: Translator.T("message_id", params) to get a localized message.
func (localizer *Localizer) T(locale string, msgID string, params ...map[string]interface{}) string {
	return localizer.Localize(locale, msgID, params...)
}

func (localizer *Localizer) Localize(lang, msgID string, params ...map[string]interface{}) string {
	// Create a localizer for the specified language
	specificLocalizer := i18n.NewLocalizer(localizer.Bundle, lang)

	/*
		Create a map for template repository if provided Example template repository:
		params := map[string]interface{}{
			"Name": "John",
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
		Logger.Warn("Failed to localize message",
			zap.String("lang", lang),
			zap.String("ID", msgID),
			zap.Error(err),
		)
		return msgID // fallback
	}

	return msg
}

func loadI18nMessages(bundle *i18n.Bundle, absPath string) {
	err := filepath.WalkDir(absPath, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			Logger.Fatal("Error walking through 'i18n' directory", zap.Error(err))
		}

		if dirEntry.IsDir() || filepath.Ext(path) != ".toml" {
			return nil
		}

		locale := extractLocaleFromFilename(path)
		if locale == "" {
			Logger.Warn("Failed to determine locale from file", zap.String("file", path))
			return nil
		}

		messageFile, err := bundle.LoadMessageFile(path)
		if err != nil {
			Logger.Warn("Failed to load message file",
				zap.String("file", path),
				zap.Error(err),
			)
			return nil
		}

		if Configs.Mode != "prod" && Configs.Mode != "production" {
			logLoadedLocaleMessages(locale, path, messageFile)
		}

		return nil
	})

	if err != nil {
		Logger.Fatal("Error after walking through 'i18n' directory", zap.Error(err))
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

	Logger.Info("Loaded locale messages",
		zap.String("locale", locale),
		zap.String("file", path),
		zap.Int("message_count", len(messageFile.Messages)),
	)
	Logger.Debug("Message details", fields...) // log details only in debug mode
}
