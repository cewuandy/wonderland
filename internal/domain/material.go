package domain

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Material struct {
	Name string

	Quantity float32

	Tool string // Which tool or platform will use

	Time float32

	Dependencies []*Material
}

type MaterialHandler interface {
	GetProductionProcess(ctx *gin.Context)
}

type MaterialUseCase interface {
	GetProductionProcess(ctx context.Context, req *ProductionProcessReq) (string, error)
}
