package main

import (
	"net/http"

	"github.com/cewuandy/wonderland/internal/controller"
	wonderlandDo "github.com/cewuandy/wonderland/pkg/do"

	"github.com/samber/do"
)

func main() {
	var err error

	injector := wonderlandDo.Injector

	// Provide all instances
	wonderlandDo.ProvideThirdPartyElement(injector)
	wonderlandDo.ProvideRepository(injector)
	wonderlandDo.ProvideUseCase(injector)

	// Register handler
	controller.NewMaterialHandler(injector)

	err = do.MustInvoke[*http.Server](injector).ListenAndServe()
	if err != nil {
		panic(err)
	}

	err = injector.Shutdown()
	if err != nil {
		panic(err)
	}
}
