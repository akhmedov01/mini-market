package handler

/* package handler

import (
	"errors"
	"fmt"
	"main/api/response"
	"main/models"
	"main/packages/helper"
	"main/packages/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Router       /transactions [POST]
// @Summary      Create Transaction
// @Description  Create Transaction
// @Tags         TRANSACTION
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateStaffTransaction  true  "transaction data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateStaffTransaction(c *gin.Context) {

	var transaction models.CreateStaffTransaction
	err := c.ShouldBindJSON(&transaction)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.strg.Transaction().CreateStaffTransaction(c.Request.Context(), transaction)
	if err != nil {
		fmt.Println("error Transaction Create:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, response.CreateResponse{Id: resp})

}

// @Router       /transactions/{id} [put]
// @Summary      Update Transaction
// @Description  api for update staffs transaction
// @Tags         TRANSACTION
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of transaction"
// @Param        staff    body     models.CreateStaffTransaction  true  "data of transaction"
// @Success      200  {string}   string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdateStaffTransaction(c *gin.Context) {
	var transaction models.CreateStaffTransaction
	err := c.ShouldBindJSON(&transaction)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}
	id := c.Param("id")

	resp, err := h.strg.Transaction().UpdateStaffTransaction(c.Request.Context(), id, transaction)
	if err != nil {
		fmt.Println("error Transaction Update:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)

	err = h.red.Cache().Delete(c.Request.Context(), id)

	if err != nil {
		fmt.Println("Error while delete transaction in cache")
	}
}

// @Router       /transactions/{id} [GET]
// @Summary      Get By Id
// @Description  get transaction by ID
// @Tags         TRANSACTION
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction ID" format(uuid)
// @Success      200  {object}  models.StaffTransaction
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetStaffTransaction(c *gin.Context) {

	id := c.Param("id")

	var resp = &models.StaffTransaction{}

	response, err := h.red.Cache().Get(c.Request.Context(), id, resp)

	if err != nil {
		fmt.Println("Error while geting transaction in cache")
	}

	if response {
		c.JSON(http.StatusOK, resp)
		return
	}

	resp, err = h.strg.Transaction().GetStaffTransaction(c.Request.Context(), models.IdRequest{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
		fmt.Println("error Transaction Get:", err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)

	err = h.red.Cache().Create(c.Request.Context(), id, resp, 0)

	if err != nil {
		fmt.Println("Error while Create transaction in cache")
	}

}

// @Router       /transactions/{id} [DELETE]
// @Summary      Delete By Id
// @Description  delete transaction by Id
// @Tags         TRANSACTION
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Transaction ID" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteStaffTransaction(c *gin.Context) {

	id := c.Param("id")

	if !helper.IsValidUUID(id) {

		h.log.Error("error Transaction Delete:", logger.Error(errors.New("invalid id")))
		c.JSON(http.StatusBadRequest, "invalid id")
		return

	}
	resp, err := h.strg.Transaction().DeleteStaffTransaction(c.Request.Context(), models.IdRequest{ID: id})
	if err != nil {
		h.log.Error("error Transaction Delete:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)

	err = h.red.Cache().Delete(c.Request.Context(), id)

	if err != nil {
		fmt.Println("Error while delete transaction in cache")
	}

}

// @Router       /transactions [get]
// @Summary      List Transaction
// @Description  get Transaction
// @Tags         TRANSACTION
// @Accept       json
// @Produce      json
// @Param        limit    query     integer  true  "limit for response"  Default(10)
// @Param        page    query     integer  true  "page of req"  Default(1)
// @Param        type    query     string  false  "filter by type"
// @Success      200  {array}   models.StaffTransaction
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllStaffTransaction(c *gin.Context) {

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

	resp, err := h.strg.Transaction().GetAllStaffTransaction(c.Request.Context(), models.GetAllStaffTransactionRequest{
		Page:            page,
		Limit:           limit,
		TransactionType: c.Query("transaction_type"),
	})
	if err != nil {
		h.log.Error("error Transaction GetAll:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)

}
*/
