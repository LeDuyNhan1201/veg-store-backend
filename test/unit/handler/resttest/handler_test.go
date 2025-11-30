package resttest

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"veg-store-backend/internal/infrastructure/core"
	"veg-store-backend/test/unit/injectiontest"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type HandlerTest[THandler any, TService any] struct {
	*core.Core
	*injectiontest.TestHTTPRouter
	MockInstance THandler
	MockService  TService
}

func NewHandlerTest[THandler any, TService any](
	handler THandler,
	service TService,
) *HandlerTest[THandler, TService] {
	return &HandlerTest[THandler, TService]{
		Core:           injectiontest.MockCore(),
		TestHTTPRouter: injectiontest.InitTestHTTPRouter(),
		MockInstance:   handler,
		MockService:    service,
	}
}

func (h *HandlerTest[THandler, TService]) AppURI(path string) string {
	return h.AppConfig.Server.ApiPrefix + h.AppConfig.Server.ApiVersion + path
}

func (h *HandlerTest[THandler, TService]) Get(t *testing.T, path string, headers ...map[string]string) *httptest.ResponseRecorder {
	return h.makeRequest(t, http.MethodGet, path, nil, mergeHeaders(headers...))
}

func (h *HandlerTest[THandler, TService]) Post(t *testing.T, path string, body any, headers ...map[string]string) *httptest.ResponseRecorder {
	return h.makeRequest(t, http.MethodPost, path, body, mergeHeaders(headers...))
}

func (h *HandlerTest[THandler, TService]) Put(t *testing.T, path string, body any, headers ...map[string]string) *httptest.ResponseRecorder {
	return h.makeRequest(t, http.MethodPut, path, body, mergeHeaders(headers...))
}

func (h *HandlerTest[THandler, TService]) Delete(t *testing.T, path string, headers ...map[string]string) *httptest.ResponseRecorder {
	return h.makeRequest(t, http.MethodDelete, path, nil, mergeHeaders(headers...))
}

func (h *HandlerTest[THandler, TService]) Patch(t *testing.T, path string, body any, headers ...map[string]string) *httptest.ResponseRecorder {
	return h.makeRequest(t, http.MethodPatch, path, body, mergeHeaders(headers...))
}

func (h *HandlerTest[THandler, TService]) UploadFile(
	test *testing.T,
	method string,
	path string,
	filePath string,
	fieldName string,
	extraFormFields map[string]string,
	headers ...map[string]string,
) *httptest.ResponseRecorder {
	// Open file base on 'filePath'
	file, err := os.Open(filePath)
	assert.NoError(test, err)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// Write file content to multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	assert.NoError(test, err)
	_, err = io.Copy(part, file)
	assert.NoError(test, err)

	// Add some extra fields to multipart form
	for k, v := range extraFormFields {
		err := writer.WriteField(k, v)
		if err != nil {
			panic(err)
		}
	}
	err = writer.Close()
	if err != nil {
		panic(err)
	}

	// Make HTTP request
	request, err := http.NewRequest(method, path, body)
	assert.NoError(test, err)

	// Set headers
	request.Header.Set("Content-Type", writer.FormDataContentType())
	for key, value := range mergeHeaders(headers...) {
		request.Header.Set(key, value)
	}

	// Handle HTTP response
	responseRecorder := httptest.NewRecorder()
	h.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code >= 400 {
		h.Logger.Info("HTTP Response Body Info",
			zap.String("method", method),
			zap.String("path", path),
			zap.Any("code", responseRecorder.Code),
			zap.String("body", responseRecorder.Body.String()),
		)
	}

	return responseRecorder
}

func (h *HandlerTest[THandler, TService]) DecodeResponse(
	test *testing.T,
	responseRecorder *httptest.ResponseRecorder,
	output any,
) {
	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType == "" ||
		contentType == "application/json" ||
		contentType == "application/json; charset=utf-8" {
		err := json.Unmarshal(responseRecorder.Body.Bytes(), output)
		assert.NoError(test, err, "failed to decode JSON response")

	} else {
		str, ok := output.(*string)
		assert.True(test, ok, "'output' must be *string for non-JSON response")
		*str = responseRecorder.Body.String()
	}
}

// ================================ //
// ====== PRIVATE FUNCTIONS ======= //
// ================================ //

func (h *HandlerTest[THandler, TService]) makeRequest(
	test *testing.T,
	method, path string,
	body any,
	headers map[string]string,
) *httptest.ResponseRecorder {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		assert.NoError(test, err)
		bodyReader = bytes.NewReader(jsonBody)
	}

	request, err := http.NewRequest(method, path, bodyReader)
	assert.NoError(test, err)
	request.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	responseRecorder := httptest.NewRecorder()
	h.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code >= 400 {
		h.Logger.Info("HTTP Response Body Info",
			zap.String("method", method),
			zap.String("path", path),
			zap.Any("code", responseRecorder.Code),
			zap.String("body", responseRecorder.Body.String()),
		)
	}
	return responseRecorder
}

func mergeHeaders(headerList ...map[string]string) map[string]string {
	merged := make(map[string]string)
	for _, header := range headerList {
		for key, value := range header {
			merged[key] = value
		}
	}
	return merged
}
