package handler

import (
	"errors"
	"fmt"
	"main/api/response"
	branch_service "main/genproto/branch-server"
	"main/models"
	"main/packages/helper"
	"main/packages/logger"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateBranch godoc
// @Router       /branches [POST]
// @Summary      Create Branch
// @Description  Create Branch
// @Tags         BRANCH
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateBranch  true  "branch data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateBranch(c *gin.Context) {

	var branch models.CreateBranch
	err := c.ShouldBindJSON(&branch)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.grpcClient.BranchService().Create(c.Request.Context(), &branch_service.CreateBranch{
		Name:    branch.Name,
		Address: branch.Address,
	})

	if err != nil {
		fmt.Println("error Branch Create:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusCreated, response.CreateResponse{Id: resp.GetId()})
}

// @Router       /branches/{id} [put]
// @Summary      update branch
// @Description  api for update persons
// @Tags         BRANCH
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of branch"
// @Param        branch    body     models.CreateBranch  true  "data of branch"
// @Success      200  {string}   string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdateBranch(c *gin.Context) {

	var branch models.CreateBranch
	err := c.ShouldBindJSON(&branch)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}
	id := c.Param("id")

	resp, err := h.grpcClient.BranchService().Update(c.Request.Context(), &branch_service.Branch{
		Id:      id,
		Name:    branch.Name,
		Address: branch.Address,
	})

	if err != nil {
		fmt.Println("error Branch Update:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)

	err = h.red.Cache().Delete(c.Request.Context(), id)

	if err != nil {
		fmt.Println("Error while delete branch in cache")
	}

}

// GetBranch godoc
// @Router       /branches/{id} [GET]
// @Summary      GET BY ID
// @Description  get branch by ID
// @Tags         BRANCH
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Branch ID" format(uuid)
// @Success      200  {object}  models.Branch
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetBranch(c *gin.Context) {

	id := c.Param("id")

	var resp = &models.Branch{}

	response, err := h.red.Cache().Get(c.Request.Context(), id, resp)

	if err != nil {
		fmt.Println("Error while geting branch in cache")
	}

	if response {
		c.JSON(http.StatusOK, resp)
		return
	}

	respService, err := h.grpcClient.BranchService().Get(c.Request.Context(), &branch_service.IdReqRes{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
		fmt.Println("error Branch Get:", err.Error())
		return
	}

	c.JSON(http.StatusOK, respService)

	err = h.red.Cache().Create(c.Request.Context(), id, respService, 0)

	if err != nil {
		fmt.Println("Error while Create branch in cache")
	}

}

// @Security ApiKeyAuth
// @Router       /branches [get]
// @Summary      List Branches
// @Description  get Branch
// @Tags         BRANCH
// @Accept       json
// @Produce      json
// @Param        limit    query     integer  true  "limit for response"  Default(10)
// @Param        page    query     integer  true  "page of req"  Default(1)
// @Param        name    query     string  false  "filter by name"
// @Success      200  {array}   models.Branch
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllBranch(c *gin.Context) {

	h.log.Info("request GetAllBranch")

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

	resp, err := h.grpcClient.BranchService().GetAll(c.Request.Context(), &branch_service.GetAllBranchRequest{
		Page:  int64(page),
		Limit: int64(limit),
		Name:  c.Query("name"),
	})
	if err != nil {
		h.log.Error("error Branch GetAll:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	h.log.Warn("response to GetAllBranch")
	c.JSON(http.StatusOK, resp)
}

// Delete godoc
// @Router       /branches/{id} [DELETE]
// @Summary      DELETE BY ID
// @Description  delete branch by ID
// @Tags         BRANCH
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Branch ID" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteBranch(c *gin.Context) {

	id := c.Param("id")

	if !helper.IsValidUUID(id) {

		h.log.Error("error Branch GetAll:", logger.Error(errors.New("invalid id")))
		c.JSON(http.StatusBadRequest, "invalid id")
		return

	}

	resp, err := h.grpcClient.BranchService().Delete(c.Request.Context(), &branch_service.IdReqRes{
		Id: id,
	})

	if err != nil {
		h.log.Error("error Branch GetAll:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, resp)

	err = h.red.Cache().Delete(c.Request.Context(), id)

	if err != nil {
		fmt.Println("Error while delete branch in cache")
	}
}
