package cherrygw

import (
	"fmt"
	"net/http"
)

func (c *Gateway) GetServerStatusByID(id string) (bool, error) {
	server, response, err := c.client.Server.List(id)
	if err != nil {
		c.logger.Error("error getting server by ID:", err, map[string]interface{}{"id": id})
		return false, err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to get server by ID, status code: %d", response.StatusCode)
		c.logger.Error("error getting server by ID:", err, map[string]interface{}{"id": id})
		return false, err
	}

	if server.State == "active" {
		return true, nil
	}
	return false, nil
}
