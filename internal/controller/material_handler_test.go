package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cewuandy/wonderland/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type materialHandlerTestSuite struct {
	suite.Suite

	r        *gin.Engine
	recorder *httptest.ResponseRecorder
	req      domain.ProductionProcessReq
}

func TestMaterialHandler(t *testing.T) {
	suite.Run(t, &materialHandlerTestSuite{})
}

func (t *materialHandlerTestSuite) SetupTest() {
	t.r = gin.Default()
	t.recorder = httptest.NewRecorder()
	t.req = domain.ProductionProcessReq{
		Name:             "test",
		Quantity:         1,
		ClockEnable:      true,
		StandClockEnable: true,
		WindmillEnable:   true,
		ACEnable:         true,
	}
}

func (t *materialHandlerTestSuite) TestGetProductionProcess() {
	t.Run(
		"success", func() {
			raw, _ := json.Marshal(t.req)
			request, err := http.NewRequest(
				http.MethodGet,
				"/wonderland/v1/production/process",
				bytes.NewBuffer(raw),
			)
			t.Nil(err)

			t.r.ServeHTTP(t.recorder, request)

			t.Equal(http.StatusOK, t.recorder.Code)
			t.Contains(t.recorder.Body.String(), "test")
		},
	)
}
