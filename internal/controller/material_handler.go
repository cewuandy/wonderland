package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"

	"github.com/cewuandy/wonderland/internal/domain"
	"github.com/cewuandy/wonderland/pkg/gin/routes"
)

type MaterialHandler struct {
	materialUseCase domain.MaterialUseCase
}

func (m *MaterialHandler) GetProductionProcess(ctx *gin.Context) {
	var (
		req    *domain.ProductionProcessReq
		result string
		err    error
	)
	req = &domain.ProductionProcessReq{}

	err = ctx.ShouldBindQuery(req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	result, err = m.materialUseCase.GetProductionProcess(ctx, req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Data(http.StatusOK, "text/plain; charset-utf-8", []byte(result))
}

func NewMaterialHandler(i *do.Injector) {
	handlers := &MaterialHandler{materialUseCase: do.MustInvoke[domain.MaterialUseCase](i)}
	r := do.MustInvoke[*gin.Engine](i)
	routes.RegisterMaterialRoutes(r, handlers)
}
