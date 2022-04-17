package httptest_demo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_helloHandler(t *testing.T) {
	tests := []struct {
		name   string
		param  string
		expect string
	}{
		{"base case", `{"name": "lvsoso"}`, "hello lvsoso"},
		{"bad case", "", "we need a name"},
	}

	r := SetupRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				"POST",
				"/hello",
				strings.NewReader(tt.param),
			)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var resp map[string]string
			err := json.Unmarshal([]byte(w.Body.String()), &resp)
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, resp["msg"])
		})
	}
}
