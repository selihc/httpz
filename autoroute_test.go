package httpz

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"selihc.com/glaive/testassist"
)

func TestAutoroute(t *testing.T) {
	mapReturnFn := func(ctx context.Context, input struct {
		Name string
	}) map[string]string {
		if input.Name == "test" {
			return map[string]string{"test": "booo"}
		}
		return map[string]string{"test": "awooo"}
	}

	cases := []struct {
		Name         string
		Fn           interface{}
		Body         io.Reader
		ExpectStatus int
		ExpectRes    []byte
	}{
		{
			"only-error",
			func(ctx context.Context, input struct {
				Name string
			}) error {
				if input.Name == "test" {
					return nil
				}
				return errors.New("illegal")
			},
			strings.NewReader(`{"Name": "nah"}`),
			http.StatusInternalServerError,
			[]byte(`{"error": "illegal"}`),
		},
		{
			"only-struct",
			mapReturnFn,
			strings.NewReader(`{"Name": "nah"}`),
			http.StatusOK,
			[]byte(`{"test": "awooo"}`),
		},
		{
			"only-struct",
			mapReturnFn,
			strings.NewReader(`{"Name": "test"}`),
			http.StatusOK,
			[]byte(`{"test": "booo"}`),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			ar, err := NewAutoroute(testassist.TestLog(t), NewJSONDecoder(), &JSONEncoder{}, c.Fn)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", c.Body)
			r.Header.Set("Content-Type", "application/json")

			ar.ServeHTTP(w, r)

			if w.Code != c.ExpectStatus {
				t.Errorf("expected %d got %d", c.ExpectStatus, w.Code)
			}

			if !testassist.JSONEquals(t, w.Body.Bytes(), c.ExpectRes) {
				t.Errorf("json not equals: %s", testassist.JSONDiffMessage(t, w.Body.Bytes(), c.ExpectRes))
			}
		})
	}
}
