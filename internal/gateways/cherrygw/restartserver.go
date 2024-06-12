package cherrygw

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cherryservers/cherrygo"
)

const (
	intervalStatus = 5 * time.Second
	timeoutStatus  = 5 * time.Minute
)

func (c *Gateway) RestartServerByID(id string) error {
	if err := c.powerOffServer(id); err != nil {
		return err
	}

	if err := c.waitForServerStatus(id, false, intervalStatus, timeoutStatus); err != nil {
		return err
	}

	if err := c.powerOnServer(id); err != nil {
		return err
	}

	if err := c.waitForServerStatus(id, true, intervalStatus, timeoutStatus); err != nil {
		return err
	}

	return nil
}

func (c *Gateway) powerOffServer(id string) error {
	_, response, err := c.client.Server.PowerOff(id)
	return c.handleResponse(response, err, "error turning off server by ID", "failed to turn server off by ID", id)
}

func (c *Gateway) powerOnServer(id string) error {
	_, response, err := c.client.Server.PowerOn(id)
	return c.handleResponse(response, err, "error turning on server by ID", "failed to turn server on by ID", id)
}

func (c *Gateway) handleResponse(response *cherrygo.Response, err error, logMsg, errMsg string, id string) error {
	if err != nil {
		c.logger.Error(logMsg, err, map[string]interface{}{"id": id})
		return err
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s, status code: %d", errMsg, response.StatusCode)
		c.logger.Error(logMsg, err, map[string]interface{}{"id": id})
		return err
	}

	return nil
}

func (c *Gateway) waitForServerStatus(id string, desiredStatus bool, interval, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	ticker := time.NewTicker(interval)
	defer timer.Stop()
	defer ticker.Stop()

	for {
		select {
		case <-timer.C:
			err := fmt.Errorf("timeout waiting for server %s to reach desired status %t", id, desiredStatus)
			c.logger.Error("error waiting for server status", err, map[string]interface{}{"id": id})
			return err
		case <-ticker.C:
			status, err := c.GetServerStatusByID(id)
			if err != nil {
				c.logger.Error("error getting server status by ID", err, map[string]interface{}{"id": id})
				return err
			}

			if status == desiredStatus {
				return nil
			}
		}
	}
}
