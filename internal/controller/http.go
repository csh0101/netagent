package controller

import (
	"net/http"

	"github.com/csh0101/netagent.git/internal/protocol"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RunApiServer(controller *Controller) {
	e := echo.New()
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	registerHandler(e, controller)
	e.Logger.Fatal(e.Start(":8080"))
}

func registerHandler(e *echo.Echo, controller *Controller) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "hello world")
	})
	e.GET("/nodes", func(c echo.Context) error {
		nodes, err := controller.RetrieveAllNode()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, nodes)
	})
	e.POST("/nodes", func(c echo.Context) error {
		req := &SendDataMsgReq{}
		if err := c.Bind(req); err != nil {
			return err
		}
		target := controller.nodeManager.GetNodeDataTunnelByIp(req.IP)
		msg := protocol.DataMessage{
			RequestID: uuid.New(),
			ClientID:  1,
			TunnelID:  1,
			DstAddr:   "127.0.0.1:10001",
			SrcAddr:   "",
			Error:     "",
			Data:      []byte("Hello World"),
		}

		if target != nil {
			data, err := msg.Encode()
			if err != nil {
				return err
			}

			if _, err := target.Write(data); err != nil {
				return err
			}
		} else {
			return c.JSON(http.StatusOK, "target conn is not find")
		}

		return nil
	})
}

type SendDataMsgReq struct {
	IP string `json:"ip"`
}
