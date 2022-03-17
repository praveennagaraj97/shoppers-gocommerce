package middlewares

import (
	conf "github.com/praveennagaraj97/shoppers-gocommerce/config"
)

type Middlewares struct {
	conf *conf.GlobalConfiguration
}

func (m *Middlewares) Initialize(cfg *conf.GlobalConfiguration) {
	m.conf = cfg
}
