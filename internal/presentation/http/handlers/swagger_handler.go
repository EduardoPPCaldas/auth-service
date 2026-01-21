package handlers

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed swagger.json
var swaggerJSONData []byte

//go:embed swagger-ui.html
var swaggerUIHTML []byte

// SwaggerHandler handles Swagger/OpenAPI documentation
type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// GetSwaggerJSON serves the Swagger JSON specification
func (h *SwaggerHandler) GetSwaggerJSON(c echo.Context) error {
	var swagger map[string]any
	if err := json.Unmarshal(swaggerJSONData, &swagger); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to parse swagger specification",
		})
	}

	return c.JSON(http.StatusOK, swagger)
}

// GetSwaggerUI serves the Swagger UI HTML page
func (h *SwaggerHandler) GetSwaggerUI(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return c.Blob(http.StatusOK, echo.MIMETextHTMLCharsetUTF8, swaggerUIHTML)
}
