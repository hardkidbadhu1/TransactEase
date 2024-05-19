package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ValidateHeaderSuite struct {
	suite.Suite
	router *gin.Engine
}

func TestValidateHeaderSuite(t *testing.T) {
	suite.Run(t, new(ValidateHeaderSuite))
}

func (suite *ValidateHeaderSuite) SetupTest() {
	suite.router = gin.Default()
	suite.router.Use(ValidateHeader())
	suite.router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Passed")
	})
}

func (suite *ValidateHeaderSuite) TestValidateHeader_WithWrongHeader() {
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Requested-With", "WrongValue")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *ValidateHeaderSuite) TestValidateHeader_WithCorrectHeader() {
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}
