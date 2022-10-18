package customer

import (
	"encoding/json"
	"github.com/iktech/demo-service/data"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type Customer struct {
}

func NewHandler() Customer {
	return Customer{}
}

func (c Customer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	cc := data.CustomersList[id]
	if cc == nil {
		w.WriteHeader(404)
		return
	}

	body, err := json.Marshal(cc)
	if err != nil {
		w.WriteHeader(500)
	}

	w.WriteHeader(200)
	_, err = w.Write(body)
}
