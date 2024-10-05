package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const wonderland = "wonderland"

const v1 = "v1"

const (
	material = "material"
)

type Route struct {
	Name    string
	Group   string
	Pattern string
	Method  string
	Handler func(ctx *gin.Context)
}

func (r *Route) registerURL(group *gin.RouterGroup) {
	url := group.Group(r.Group)

	switch r.Method {
	case http.MethodGet:
		url.GET(r.Pattern, r.Handler)
	case http.MethodPost:
		url.POST(r.Pattern, r.Handler)
	case http.MethodPut:
		url.PUT(r.Pattern, r.Handler)
	case http.MethodDelete:
		url.DELETE(r.Pattern, r.Handler)
	}
}
