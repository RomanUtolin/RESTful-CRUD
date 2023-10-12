package middleware_test

import (
	"context"
	"github.com/RomanUtolin/RESTful-CRUD/internall/http/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	server := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	ctx := server.NewContext(req, rec)
	middl := middleware.InitMiddleware(context.Background())

	handler := middl.CORS(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))

	err := handler(ctx)
	require.NoError(t, err)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
}
