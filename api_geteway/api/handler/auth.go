package handler

import (
	"fmt"
	"main/api/response"
	"main/config"
	staff_service "main/genproto/staff-server"
	"main/models"
	"main/packages/helper"
	"main/packages/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Router       /login [post]
// @Summary      create person
// @Description  api for create persons
// @Tags         persons
// @Accept       json
// @Produce      json
// @Param        person    body     models.LoginReq  true  "data of person"
// @Success      200  {object}  models.LoginRes
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) Login(c *gin.Context) {

	var req models.LoginReq

	err := c.ShouldBindJSON(&req)

	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		res := response.ErrorResp{Code: "BAD REQUEST", Message: "invalid fields in body"}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	hashPass, err := helper.GeneratePasswordHash(req.Password)

	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		res := response.ErrorResp{Code: "INVALID Password", Message: "invalid password"}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	resp, err := h.grpcClient.StaffService().GetByUsername(c.Request.Context(), &staff_service.RequestByUsername{
		Login: req.Loging,
	})

	if err != nil {
		fmt.Println("error Staff GetByLoging:", err.Error())
		res := response.ErrorResp{Code: "INTERNAL ERROR", Message: "internal server error"}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	err = helper.ComparePasswords([]byte(hashPass), []byte(resp.Password))

	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		res := response.ErrorResp{Code: "INVALID Password", Message: "invalid password"}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	m := make(map[string]interface{})
	m["user_id"] = resp.Id
	m["branch_id"] = resp.BranchId
	token, err := helper.GenerateJWT(m, config.TokenExpireTime, config.JWTSecretKey)

	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, models.LoginRes{Token: token})
}
