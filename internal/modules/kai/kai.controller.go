package kai

import (
	"net/http"

	helper "kai-back/internal/shared/helpers"
	"kai-back/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type KaiController struct {
	service ServicePort
}

func NewKaiController(service ServicePort) *KaiController {
	return &KaiController{service: service}
}

func (ctrl *KaiController) GetKaiDashboard(
	c *gin.Context,
) {

	userID, err := helper.ValidateUserUUID(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp, err := ctrl.service.GetKaiDashboard(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(
		http.StatusOK,
		response.Success(
			resp,
		),
	)
}
