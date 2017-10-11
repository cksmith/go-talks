package main

import (
	"context"
	"fmt"
	"github.com/cksmith/go-talks/Golang-KW-201710/ble/blelib"
	"github.com/jawher/mow.cli"
	"os"
	"os/signal"
	"syscall"
)

// START1 OMIT
func main() {
	app := cli.App("ble", "Go Meetup BLE Demo")

	ctx, doneFunc := context.WithCancel(context.Background())

	waitForCancel := func() {
		// Wait for a signal (ctrl-c)
		stop := make(chan os.Signal)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			select {
			case signal := <-stop:
				fmt.Println("Got signal", signal.String())
				doneFunc()
			}
		}()

		<-ctx.Done()
	}

	// END1 OMIT
	// START2 OMIT
	app.Command("service", "Run the BLE service", func(cmd *cli.Cmd) {
		cmd.Spec = "[--mock]"

		var (
			mockParam = cmd.BoolOpt("mock", false, "Mock the GPIO interface")
		)

		cmd.Action = func() {
			var led blelib.Led
			var sensor blelib.TemperatureHumiditySensor
			if *mockParam {
				led = blelib.NewMockLedService()
				sensor = blelib.NewMockTemperatureHumidityService()
			} else {
				led = blelib.NewLedService()
				sensor = blelib.NewTemperatureHumidityService()
			}
			err := led.Open()
			if err != nil {
				panic("Cannot open GPIO")
			}
			defer led.Close()
			blelib.NewServer(ctx, led, sensor)
			waitForCancel()
		}
	})

	// END2 OMIT

	app.Run(os.Args)
}
