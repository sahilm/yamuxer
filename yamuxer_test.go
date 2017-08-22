package yamuxer_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"fmt"

	"github.com/sahilm/yamuxer"
)

// nolint: gocyclo
func TestYamuxer(t *testing.T) {
	t.Run("It should match the correct route", func(t *testing.T) {
		route := newTestRoute(regexp.MustCompile(`/users/(\d+)`))
		mux := yamuxer.New([]*yamuxer.Route{route.Route})

		server := httptest.NewServer(mux)
		defer server.Close()

		userid := "1234"
		r, err := http.Get(fmt.Sprintf("%v/users/%v", server.URL, userid))
		if err != nil {
			t.Fatal(err)
		}

		if r.StatusCode != 200 {
			t.Errorf("got status code: %v, want: 200", r.StatusCode)
		}

		got := route.matches[0]
		want := userid

		if got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("It should return 404 if there is no matching route", func(t *testing.T) {
		route := newTestRoute(regexp.MustCompile(`/$`))
		mux := yamuxer.New([]*yamuxer.Route{route.Route})

		server := httptest.NewServer(mux)
		defer server.Close()

		r, err := http.Get(fmt.Sprintf("%v/users", server.URL))
		if err != nil {
			t.Fatal(err)
		}

		if r.StatusCode != 404 {
			t.Errorf("got status code: %v, want: 404", r.StatusCode)
		}
	})

	t.Run("It selects the first match found", func(t *testing.T) {
		route1 := newTestRoute(regexp.MustCompile(`/users/(\d+)`))
		route2 := newTestRoute(regexp.MustCompile(`/admins/(\d+)`))
		route3 := newTestRoute(regexp.MustCompile(`/admins/(\w+)`))

		mux := yamuxer.New([]*yamuxer.Route{route1.Route, route2.Route, route3.Route})
		server := httptest.NewServer(mux)
		defer server.Close()

		adminName := "sahil"
		r, err := http.Get(fmt.Sprintf("%v/admins/%v", server.URL, adminName))
		if err != nil {
			t.Fatal(err)
		}

		if r.StatusCode != 200 {
			t.Errorf("got status code: %v, want: 200", r.StatusCode)
		}

		got := route3.matches[0]
		want := adminName

		if got != want {
			t.Errorf("got: %v, want: %v", got, want)
		}

		if len(route1.matches) > 0 {
			t.Errorf("route1 should not have matched")
		}

		if len(route2.matches) > 0 {
			t.Errorf("route2 should not have matched")
		}

	})
}

type testroute struct {
	*yamuxer.Route
	matches []string
}

func newTestRoute(pattern *regexp.Regexp) *testroute {
	var t *testroute
	t = &testroute{
		Route: &yamuxer.Route{
			Pattern: pattern,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				t.matches = r.Context().Value(yamuxer.ContextKey("matches")).([]string)
			}},
	}
	return t
}
