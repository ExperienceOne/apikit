package tests

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/ExperienceOne/apikit/tests/api"
)

// Tests setting a custom prefix for api routes using ServerOpts.
func TestCustomPrefix(t *testing.T) {
	middlewares := []api.Middleware{
		{Handler: api.RouterPanicMiddleware()},
		{Handler: api.RouterPopulateContextMiddleware()},
	}

	testServerWrapper = api.NewVisAdminServer(&api.ServerOpts{
		Middleware:   middlewares,
		ErrorHandler: log.Println,
		Prefix:       "/my-custom-prefix",
	})

	testServerWrapper.SetGetUsersHandler(api.GetUsers)

	go testServerWrapper.Start(4568)

	defer testServerWrapper.Server.Stop()

	time.Sleep(1 * time.Second)

	opts := api.Opts{
		Hooks: api.DevHook(),
		Ctx:   nil,
	}

	client := api.NewVisAdminClient(new(http.Client), "http://localhost:4568/my-custom-prefix", opts)

	response, err := client.GetUsers(&api.GetUsersRequest{
		XAuth: "sessionID",
	})

	if err != nil {
		t.Fatalf("error GetUsers failed: %v", err)
	}

	if _, ok := response.(*api.GetUsers200Response); !ok {
		t.Fatalf("error GetUsers response is bad: %#v", err)
	}
}
