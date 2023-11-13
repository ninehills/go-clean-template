package integration_test

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/Eun/go-hit" //nolint: revive
)

const (
	// Attempts connection.
	host       = "app:8080"
	healthPath = "http://" + host + "/healthz"
	attempts   = 20

	// HTTP REST.
	basePath = "http://" + host + "/v1"
)

func TestMain(m *testing.M) {
	if err := healthCheck(attempts); err != nil {
		log.Fatalf("Integration tests: host %s is not available: %s", host, err)
	}

	log.Printf("Integration tests: host %s is available", host)

	code := m.Run()
	os.Exit(code)
}

func healthCheck(attempts int) error {
	var err error

	for attempts > 0 {
		err = Do(Get(healthPath), Expect().Status().Equal(http.StatusOK))
		if err == nil {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return err
}

// HTTP GET: /.
func TestHTTPCreateUser(t *testing.T) {
	t.Parallel()

	body := `{
		"username": "user",
		"email": "user@example.com",
		"description": "aaa",
		"password": "pass@123",
		"confirmPassword": "pass@123"
	}`
	Test(t,
		Description("Create User Success"),
		Post(basePath+"/users"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".username").Equal("user"),
	)
}
