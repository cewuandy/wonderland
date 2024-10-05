package do

import (
	"github.com/samber/do"

	"github.com/cewuandy/wonderland/internal/repository/files"
)

func ProvideRepository(injector *do.Injector) {
	do.Provide(injector, files.NewFileRepo)
}
