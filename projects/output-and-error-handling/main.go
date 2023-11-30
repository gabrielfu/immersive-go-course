package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetWeather() (string, error) {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080", nil)
	if err != nil {
		return "", fmt.Errorf("error: failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error: server connection terminated: %v", err)
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return "", fmt.Errorf("error: server returned status %d", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")

		// try to parse into integer
		seconds, err := strconv.Atoi(retryAfter)
		var duration time.Duration
		if err == nil {
			duration = time.Duration(seconds) * time.Second
		} else if date, err := http.ParseTime(retryAfter); err == nil {
			// try to parse as date
			duration = date.Sub(time.Now().UTC())
		} else {
			// if we get here, we couldn't parse the header
			return "", fmt.Errorf("error: server returned invalid Retry-After header: %s", retryAfter)
		}

		if duration > 5*time.Second {
			return "", fmt.Errorf("error: cannot retry within %s, giving up", duration)
		} else if duration > 1*time.Second {
			fmt.Fprintf(os.Stderr, "retrying after %s...\n", duration)
		}
		time.Sleep(duration)
		return GetWeather()
	}

	weather := make([]byte, resp.ContentLength)
	resp.Body.Read(weather)
	defer resp.Body.Close()
	return string(weather), nil
}

func main() {
	weather, err := GetWeather()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "The weather is: %s\n", weather)
}
