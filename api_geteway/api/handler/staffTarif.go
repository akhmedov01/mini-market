package handler

import (
	"errors"
	"fmt"
	"main/api/response"
	staff_service "main/genproto/staff-server"
	"main/models"
	"main/packages/helper"
	"main/packages/logger"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Router       /tarifs [POST]
// @Summary      Create StaffTarif
// @Description  Create StsffTarif
// @Tags         TARIF
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateStaffTarif  true  "tarif data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateStaffTarif(c *gin.Context) {

	var tarif models.CreateStaffTarif
	err := c.ShouldBindJSON(&tarif)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.grpcClient.TarifService().Create(c.Request.Context(), &staff_service.CreateTarif{
		Name:          tarif.Name,
		Type:          tarif.Type,
		AmountForCash: float32(tarif.AmountForCash),
		AmountForCard: float32(tarif.AmountForCard),
		FoundedAt:     tarif.FoundedAt,
	})
	if err != nil {
		fmt.Println("error Tarif Create:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, response.CreateResponse{Id: resp.GetId()})

}

// @Router       /tarifs/{id} [put]
// @Summary      Update StaffTarif
// @Description  api for update tarifs
// @Tags         TARIF
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of staff tarif"
// @Param        tarif    body     models.CreateStaffTarif  true  "data of tarif"
// @Success      200  {string}   string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdateStaffTarif(c *gin.Context) {

	var tarif models.CreateStaffTarif

	err := c.ShouldBindJSON(&tarif)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}
	id := c.Param("id")

	resp, err := h.grpcClient.TarifService().Update(c.Request.Context(), &staff_service.Tarif{
		Id:            id,
		Name:          tarif.Name,
		Type:          tarif.Type,
		AmountForCash: float32(tarif.AmountForCash),
		AmountForCard: float32(tarif.AmountForCard),
		FoundedAt:     tarif.FoundedAt,
	})
	if err != nil {
		fmt.Println("error Tarif Update:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)

	err = h.red.Cache().Delete(c.Request.Context(), id)

	if err != nil {
		fmt.Println("Error while delete staffTarif in cache")
	}
}

// @Router       /tarifs/{id} [GET]
// @Summary      Get By Id
// @Description  get staff tarif by ID
// @Tags         TARIF
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tarif ID" format(uuid)
// @Success      200  {object}  models.StaffTarif
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetStaffTarif(c *gin.Context) {

	id := c.Param("id")

	var resp = &models.StaffTarif{}

	response, err := h.red.Cache().Get(c.Request.Context(), id, resp)

	if err != nil {
		fmt.Println("Error while geting stafTarif in cache")
	}

	if response {
		c.JSON(http.StatusOK, resp)
		return
	}

	respService, err := h.grpcClient.TarifService().Get(c.Request.Context(), &staff_service.IdRequest{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
		fmt.Println("error Tarif Get:", err.Error())
		return
	}

	c.JSON(http.StatusOK, respService)

	err = h.red.Cache().Create(c.Request.Context(), id, respService, 0)

	if err != nil {
		fmt.Println("Error while Create tarif in cache")
	}

}

// @Router       /tarifs/{id} [DELETE]
// @Summary      Delete By Id
// @Description  delete tarif by Id
// @Tags         TARIF
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tarif ID" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteStaffTarif(c *gin.Context) {

	id := c.Param("id")

	if !helper.IsValidUUID(id) {

		h.log.Error("error Tarif Delete:", logger.Error(errors.New("invalid id")))
		c.JSON(http.StatusBadRequest, "invalid id")
		return

	}
	resp, err := h.grpcClient.TarifService().Delete(c.Request.Context(), &staff_service.IdRequest{
		Id: id,
	})

	if err != nil {
		h.log.Error("error Tarif Delete:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)

	err = h.red.Cache().Delete(c.Request.Context(), id)

	if err != nil {
		fmt.Println("Error while delete tarif in cache")
	}

}

// @Router       /tarifs [get]
// @Summary      List Tarifs
// @Description  get Tarif
// @Tags         TARIF
// @Accept       json
// @Produce      json
// @Param        limit    query     integer  true  "limit for response"  Default(10)
// @Param        page    query     integer  true  "page of req"  Default(1)
// @Param        name    query     string  false  "filter by name"
// @Success      200  {array}   models.StaffTarif
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllStaffTarif(c *gin.Context) {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		h.log.Error("error get page:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid page param")
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		h.log.Error("error get limit:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid page param")
		return
	}

	resp, err := h.grpcClient.TarifService().GetAll(c.Request.Context(), &staff_service.GetAllTarifRequest{
		Page:  int64(page),
		Limit: int64(limit),
		Name:  c.Query("name"),
	})
	if err != nil {
		h.log.Error("error Tarif GetAll:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)

}
