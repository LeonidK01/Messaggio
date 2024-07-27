package delivery

import (
	"net/http"

	"github.com/LeonidK01/Messaggio/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type messageGinDelivery struct {
	msgUC model.MessageUsecase
}

func HandleMessageGinDelivery(gr gin.IRoutes, muc model.MessageUsecase) {
	d := &messageGinDelivery{
		msgUC: muc,
	}

	gr.POST("/send", d.Send)
}

// TODO: валидация
type SendRequest struct {
	CreatedBy string `json:"created_by"`
	From      string `json:"from"`
	To        string `json:"to"`
	Text      string `json:"text"`
}

// TODO: обробатывать корректно ошибки со стороны клиента
func (d *messageGinDelivery) Send(c *gin.Context) {
	req := &SendRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdBy, err := uuid.Parse(req.CreatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	from, err := uuid.Parse(req.From)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	to, err := uuid.Parse(req.To)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	msg := &model.Message{
		CreatedBy: createdBy,
		From:      from,
		To:        to,
		Text:      req.Text,
	}

	if err := d.msgUC.Send(c, msg); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
