package util_test

import (
	"banking/util"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"
)

type reader struct {
	bfr []byte
}

func (r *reader) Read(dst []byte) (int, error) {
	for i, b := range r.bfr {
		dst[i] = b
	}
	return len(r.bfr), nil
}

type reqBody struct {
	ID   int64           `json:"id"`
	Name string          `json:"name"`
	Nums []int           `json:"nums"`
	Map  map[string]bool `json:"map"`
}

func (reqBody) Valid() map[string]string {
	return map[string]string{}
}

func TestDecodeBody_ReturnCorrectBody(t *testing.T) {
	body := reqBody{
		1,
		"John",
		[]int{0, 1, 2, 3},
		map[string]bool{
			"1": true,
			"2": false,
		},
	}
	bodyBfr, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/", &reader{bodyBfr})
	w := httptest.NewRecorder()
	newBody, ok := util.DecodeBody[reqBody](w, req)
	if !ok {
		t.Error("DecodeBody().ok = false; want true")
	}

	if !reflect.DeepEqual(body, newBody) {
		t.Error("DeepEqual() = false; want true")
	}
}
