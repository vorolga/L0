package handlers

import (
	"L0/internal"
	"L0/static"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *zap.SugaredLogger
	usecase internal.Usecase
}

func NewHandler(u internal.Usecase, logger *zap.SugaredLogger) *Handler {
	return &Handler{usecase: u, logger: logger}
}

func (h *Handler) Register(router *echo.Echo) {
	router.GET("/api/getOrder", h.GetOrder())
}

func (h *Handler) ReceiveOrder(msg *stan.Msg) {
	err := msg.Ack()
	if err != nil {
		h.logger.Error(err)
		return
	}

	order := internal.Order{}
	err = easyjson.Unmarshal(msg.Data, &order)
	if err != nil {
		h.logger.Error(err)
		return
	}

	err = h.usecase.SaveOrder(context.Background(), &order)
	if err != nil {
		h.logger.Error(err)
	}
}

func (h *Handler) GetOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid := c.QueryParam("uid")
		if uid != "" {
			order, err := h.usecase.GetOrder(context.Background(), uid)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}
			if order == nil {
				return c.HTML(http.StatusNotFound, static.GenerateNotFound())
			}
			return c.HTML(http.StatusOK, static.GeneratePage(order))
		}

		return c.NoContent(http.StatusBadRequest)
	}
}
