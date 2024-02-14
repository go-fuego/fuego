package controller

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-fuego/fuego"
)

type MyContext struct {
	fuego.BaseContext
}

// Controllers tests for /NewController
// Does not test the actual HTTP requests, only the controller logic.
// Beware of the Context, it is not exactly the same as the one used in the actual HTTP requests.

func TestNewControllerRessources_Routes(t *testing.T) {
	s := fuego.NewServer()

	rs := NewControllerRessources{
		// Dependency Injection
		NewControllerService: NewControllerServiceMock{
			getAllNewControllerLength: 2,
		},
	}

	rs.Routes(s)

	t.Run("GET /newController", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/newController", nil)

			newControllers, err := rs.getAllNewController(fuego.BaseContext{
				Req: r,
			})
			if err != nil {
				t.Errorf("error: " + err.Error())
			}
			if len(newControllers) != 2 {
				t.Errorf("error: newControllers is of wrong length, expected %d", 2)
			}
		})
	})

	t.Run("getNewController", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/newController/123", nil)
		w := httptest.NewRecorder()

		_, err := rs.getNewController(fuego.BaseContext{Req: r, Res: w})
		if err != nil {
			t.Errorf("error: " + err.Error())
		}
	})

	t.Run("POST /newController", func(t *testing.T) {
		body := `{"name": "hi"}`
		r := httptest.NewRequest("GET", "/newController/123", strings.NewReader(body))
		w := httptest.NewRecorder()

		_, err := rs.postNewController(&fuego.BaseContextWithBody[NewControllerCreate]{
			BaseContext: fuego.BaseContext{Req: r, Res: w},
		})
		if err != nil {
			t.Errorf("error: " + err.Error())
		}
	})
}
