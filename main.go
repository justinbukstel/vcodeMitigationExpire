package main

import (
	"flag"
	"net"
	"net/http"
	"os"
	"os/exec"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {

	var socketPath string
	flag.StringVar(&socketPath, "socket", "/run/guest-services/backend.sock", "Unix domain socket to listen on")
	flag.Parse()

	_ = os.RemoveAll(socketPath)

	logger.SetOutput(os.Stdout)

	logMiddleware := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}",` +
			`"method":"${method}","uri":"${uri}",` +
			`"status":${status},"error":"${error}"` +
			`}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output:           logger.Writer(),
	})

	logger.Infof("Starting listening on %s\n", socketPath)
	router := echo.New()
	router.HideBanner = true
	router.Use(logMiddleware)
	startURL := ""

	ln, err := listen(socketPath)
	if err != nil {
		logger.Fatal(err)
	}
	router.Listener = ln

	router.GET("/scan", hello)

	logger.Fatal(router.Start(startURL))
}

func listen(path string) (net.Listener, error) {
	return net.Listen("unix", path)
}

func hello(ctx echo.Context) error {
	cmd := exec.Command("bash", "-c", "echo 'Hello, world!'")

	// Execute the command and capture the output
	output, err := cmd.CombinedOutput()

	if err != nil {
		return ctx.JSON(http.StatusNotAcceptable, HTTPMessageBody{Message: "Error"})
	} else if output != nil {
		return ctx.JSON(http.StatusOK, HTTPMessageBody{Message: "Hello"})
	} else{
		return ctx.JSON(http.StatusOK, HTTPMessageBody{Message: "Hello"})
	}

	
}

type HTTPMessageBody struct {
	Message string
}

type ContainerScanResults struct {
	TargetID        string `json:"target_id"`
	Findings        Findings
	Vulnerabilities Vulnerabilities
	Secrets         Secrets
	Configs         Configs
	PolicyResults   []PolicyResult `json:"policy-results"`
	PolicyPassed    bool           `json:"policy-passed"`
}

type Findings struct {
	// Define the fields of the Findings struct here
}

type Vulnerabilities struct {
	// Define the fields of the Vulnerabilities struct here
}

type Secrets struct {
	// Define the fields of the Secrets struct here
}

type Configs struct {
	// Define the fields of the Configs struct here
}

type PolicyResult struct {
	// Define the fields of the PolicyResult struct here
}
