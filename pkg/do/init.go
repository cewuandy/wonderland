package do

import "github.com/samber/do"

var Injector *do.Injector

func init() {
	Injector = do.New()
}
