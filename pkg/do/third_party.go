package do

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/cewuandy/wonderland/internal/domain"
	"github.com/cewuandy/wonderland/pkg/options"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func ProvideThirdPartyElement(injector *do.Injector) {
	r := gin.New()
	// TODO: should assign a real ip:port, that is workaround now
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	r.Use(cors.New(corsConfig))

	materialFilePath := "assets/material"

	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.Usage = func() {}
	o := &domain.Options{}
	_ = options.LoadDefaultConfig(flagSet, o)
	_ = options.LoadCliFlagConfigs(flagSet)

	do.ProvideNamedValue(injector, "material_path", materialFilePath)

	do.Provide(injector, func(i *do.Injector) (*gin.Engine, error) { return r, nil })
	do.Provide(
		injector, func(injector *do.Injector) (*http.Server, error) {
			return &http.Server{
				Addr:    fmt.Sprintf("%s:%d", o.Addr, o.Port),
				Handler: r,
			}, nil
		},
	)
}
