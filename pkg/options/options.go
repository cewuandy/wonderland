package options

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/cewuandy/wonderland/internal/domain"

	"github.com/iancoleman/strcase"
)

func LoadDefaultConfig(flagSet *flag.FlagSet, option *domain.Options) error {
	rType := reflect.TypeOf(*option)
	rElem := reflect.ValueOf(option).Elem()

	for i := 0; i < rType.NumField(); i++ {
		name := rType.Field(i).Name
		field := rElem.FieldByName(name)
		value, isDefaultExisted := rType.Field(i).Tag.Lookup("default")
		usage, _ := rType.Field(i).Tag.Lookup("usage")

		if !field.IsValid() || !field.CanSet() || !isDefaultExisted {
			continue
		}

		switch field.Kind() {
		case reflect.Int:
			v, _ := strconv.ParseInt(value, 10, 0)
			field.SetInt(v)
			flagSet.IntVar(
				(*int)(unsafe.Pointer(field.Addr().Pointer())), strcase.ToKebab(name), int(v),
				usage,
			)
		case reflect.String:
			field.SetString(value)
			flagSet.StringVar(
				(*string)(unsafe.Pointer(field.Addr().Pointer())), strcase.ToKebab(name), value,
				usage,
			)
		default:
			return fmt.Errorf("get unrecognized builtin type with key %s", name)
		}
	}

	return nil
}

func LoadCliFlagConfigs(flagSet *flag.FlagSet) error {
	var err error

	err = flagSet.Parse(os.Args[1:])

	if errors.Is(err, flag.ErrHelp) {
		if errors.Is(err, flag.ErrHelp) {
			_, _ = fmt.Fprintf(os.Stderr, "usage: %s [command line options]\n", os.Args[0])
			_, _ = fmt.Fprintf(os.Stderr, "Available command line options:\n")
			flagSet.PrintDefaults()
			os.Exit(0)
		}
		os.Exit(0)
	}

	if err != nil {
		return err
	}

	return nil
}
