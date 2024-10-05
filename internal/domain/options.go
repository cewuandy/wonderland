package domain

type Options struct {
	Addr string `default:"0.0.0.0" usage:"[Server Mode] Listen address."`
	Port int    `default:"8081" usage:"[Server Mode] Listen port."`
}
