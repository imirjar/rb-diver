package http

import (
	"log"
	"net/http"
)

type client struct {
	*http.Client
}

func (c *client) POST(req *http.Request) (string, error) {

	resp, err := c.Do(req)
	if err != nil {
		log.Println("ERROR:", err)
		return resp.Status, err
	}
	defer resp.Body.Close()
	return resp.Status, nil

}
