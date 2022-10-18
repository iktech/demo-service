package version

import (
	"encoding/json"
	"github.com/iktech/demo-service/core"
	"net/http"
)

type Version struct {
}

func NewHandler() Version {
	return Version{}
}

func (c Version) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := core.Version{
		Version: "__VERSION__",
	}

	body, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(500)
	}

	w.WriteHeader(200)
	_, err = w.Write(body)
}
