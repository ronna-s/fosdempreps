package fosdem2026

import (
	"io"
	"net/http"
	"runtime"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
)

func startServer(t *testing.T, server *http.Server) {
	go func() {
		t.Log("spinning up the server...")
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				t.Errorf("server error: %v", err)
			} else {
				t.Log("server closed")
			}
		}
	}()
}

func runCheckServer(t *testing.T, url string) {
	// run the server
	// wait for service to become responsive
	// make some requests to the server here
	t.Log("making the request")
	resp, err := http.Get("http://127.0.0.1:0/")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	t.Log("request made")
	defer resp.Body.Close()
	t.Log("before reading body")
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "hello there!", string(b))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

// Just test a resgular server
func TestNormalServerWithSleep(t *testing.T) {
	server := &http.Server{Addr: ":0", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello there!"))
	})}
	go startServer(t, server)
	defer server.Close()
	time.Sleep(100 * time.Millisecond) // wait for server goroutine to start
	runCheckServer(t, server.Addr)
}

func TestNormalServerWithGosched(t *testing.T) {
	server := &http.Server{Addr: ":0", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello there!"))
	})}
	go startServer(t, server)
	defer server.Close()
	runtime.Gosched() // yield to allow server goroutine to start
	runCheckServer(t, server.Addr)
}

func TestNormalServerWithSynctestGosched(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		server := &http.Server{Addr: ":0", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello there!"))
		})}
		go startServer(t, server)
		defer server.Close()
		runtime.Gosched() // yield to allow server goroutine to start
		runCheckServer(t, server.Addr)
	})
}

// Will fail -  run with timeout
func TestNormalServerWithSynctestSleep(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		server := &http.Server{Addr: ":0", Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello there!"))
		})}
		go startServer(t, server)
		defer server.Close()
		time.Sleep(time.Millisecond * 100)
		runCheckServer(t, server.Addr)
	})
}
