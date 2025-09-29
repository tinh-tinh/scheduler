package scheduler_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/scheduler"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func SchedulerApp() *core.App {
	type CounterSvc struct {
		counter int
	}

	service := func(module core.Module) core.Provider {
		return module.NewProvider(&CounterSvc{})
	}

	schedule := func(module core.Module) core.Provider {
		svc := core.Inject[CounterSvc](module)
		task := scheduler.NewTask(module)

		task.Cron("* * * * *", func() {
			fmt.Println(1)
			svc.counter++
		})

		return task
	}

	controller := func(module core.Module) core.Controller {
		svc := core.Inject[CounterSvc](module)
		ctrl := module.NewController("schedulers")

		ctrl.Get("", func(ctx core.Ctx) error {
			return ctx.JSON(svc.counter)
		})

		return ctrl
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports:     []core.Modules{scheduler.ForRoot()},
			Controllers: []core.Controllers{controller},
			Providers:   []core.Providers{service, schedule},
		})

		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("api")

	return app
}

func TestModule(t *testing.T) {
	serverReady := make(chan struct{})
	var testServer *httptest.Server

	go func() {
		app := SchedulerApp()
		testServer = httptest.NewServer(app.PrepareBeforeListen())

		time.Sleep(5 * time.Second)

		close(serverReady)
		select {}
	}()

	<-serverReady
	testClient := testServer.Client()

	time.Sleep(5 * time.Second)
	resp, err := testClient.Get(testServer.URL + "/api/schedulers")
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.Nil(t, err)

	fmt.Println(string(data))
}
