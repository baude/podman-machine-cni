package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
)

type Expose struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

type Unexpose struct {
	Local string `json:"local"`
}

// getPrimaryIP extracts the host's IP address from an environment
// variable. It is an error if that IP is blank
func getPrimaryIP() (net.IP, error) {
	hostIP := os.Getenv("PODMAN_MACHINE_HOST")
	if len(hostIP) < 1 {
		return nil, errors.New("invalid PODMAN_MACHINE_HOST environment variable")
	}
	addr := net.ParseIP(hostIP)
	if addr == nil {
		return nil, errors.New("unable to parse PODMAN_HOST_MACHINE IP address")
	}
	return addr, nil
}

func postRequest(ctx context.Context, url *url.URL, body interface{}) error {
	var buf io.ReadWriter
	client := &http.Client{}
	buf = new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), buf)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("something went wrong with the request")
	}
	return nil
}
