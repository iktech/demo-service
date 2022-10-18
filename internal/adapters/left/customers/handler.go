package customers

import (
	"encoding/json"
	"github.com/iktech/demo-service/core"
	"github.com/iktech/demo-service/data"
	"net/http"
)

type Customers struct {
}

func NewHandler() Customers {
	return Customers{}
}

func (c Customers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := make([]*core.Customer, 0)
	for _, v := range data.CustomersList {
		resp = append(resp, v)
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
	}

	w.WriteHeader(200)
	_, err = w.Write(body)
}
