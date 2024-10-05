package do

import (
	"github.com/samber/do"

	"github.com/cewuandy/wonderland/internal/usecase"
)

func ProvideUseCase(injector *do.Injector) {
	do.Provide(injector, usecase.NewMaterialUseCase)
}
