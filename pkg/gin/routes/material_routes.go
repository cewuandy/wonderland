package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/cewuandy/wonderland/internal/domain"
)

func RegisterMaterialRoutes(r *gin.Engine, handler domain.MaterialHandler) {
	group := r.Group(wonderland)
	routes := []Route{
		{
			Name:    "GetProductionProcess",
			Group:   v1,
			Pattern: material,
			Method:  http.MethodGet,
			Handler: handler.GetProductionProcess,
		},
	}

	for i := 0; i < len(routes); i++ {
		routes[i].registerURL(group)
	}
}
