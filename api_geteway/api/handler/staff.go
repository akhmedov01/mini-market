package handler

import (
	"errors"
	"fmt"
	"main/api/response"
	staff_service "main/genproto/staff-server"
	"main/models"
	"main/packages/helper"
	"main/packages/logger"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

// @Router       /staffs [POST]
// @Summary      Create Staff
// @Description  Create Stsff
// @Tags         STAFF
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateStaff  true  "staff data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateStaff(c *gin.Context) {

	var staff models.CreateStaff
	err := c.ShouldBindJSON(&staff)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.grpcClient.StaffService().Create(c.Request.Context(), &staff_service.CreateStaff{
		BranchId:  staff.BranchId,
		TarifId:   staff.TarifId,
		TypeStaff: staff.TypeStaffs,
		Name:      staff.Name,
		Balance:   float32(staff.Balance),
		Age:       int64(staff.Age),
		BirthDate: staff.Birthday_Date,
		Login:     staff.Loging,
		Password:  staff.Password,
	})
	if err != nil {
		fmt.Println("error Staff Create:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusCreated, response.CreateResponse{Id: resp.GetId()})

}

// @Router       /staffs/{id} [put]
// @Summary      Update Staff
// @Description  api for update staffs
// @Tags         STAFF
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of staff"
// @Param        staff    body     models.CreateStaff  true  "data of staff"
// @Success      200  {string}   string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdateStaff(c *gin.Context) {

	var staff models.CreateStaff
	err := c.ShouldBindJSON(&staff)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}
	id := c.Param("id")

	resp, err := h.grpcClient.StaffService().Update(c.Request.Context(), &staff_service.Staff{
		Id:        id,
		BranchId:  staff.BranchId,
		TarifId:   staff.TarifId,
		TypeStaff: staff.TypeStaffs,
		Name:      staff.Name,
		Balance:   float32(staff.Balance),
		Age:       int64(staff.Age),
		BirthDate: staff.Birthday_Date,
	})
	if err != nil {
		fmt.Println("error Staff Update:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)

	err = h.red.Cache().Delete(c.Request.Context(), id)

	if err != nil {
		fmt.Println("Error while delete staff in cache")
	}

}

// @Router       /staffs/{id} [GET]
// @Summary      Get By Id
// @Description  get staff by ID
// @Tags         STAFF
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Staff ID" format(uuid)
// @Success      200  {object}  models.Staff
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetStaff(c *gin.Context) {

	id := c.Param("id")

	var resp = &models.Staff{}

	response, err := h.red.Cache().Get(c.Request.Context(), id, resp)

	if err != nil {
		fmt.Println("Error while geting staff in cache")
	}

	if response {
		c.JSON(http.StatusOK, resp)
		return
	}

	respService, err := h.grpcClient.StaffService().Get(c.Request.Context(), &staff_service.IdRequest{
		Id: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
		fmt.Println("error Sale Get:", err.Error())
		return
	}

	c.JSON(http.StatusOK, respService)

	err = h.red.Cache().Create(c.Request.Context(), id, respService, 0)

	if err != nil {
		fmt.Println("Error while Create staff in cache")
	}

}

// @Router       /staffs/{id} [DELETE]
// @Summary      Delete By Id
// @Description  delete staff by Id
// @Tags         STAFF
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Staff ID" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteStaff(c *gin.Context) {

	id := c.Param("id")

	if !helper.IsValidUUID(id) {

		h.log.Error("error Staff Delete:", logger.Error(errors.New("invalid id")))
		c.JSON(http.StatusBadRequest, "invalid id")
		return

	}
	resp, err := h.grpcClient.StaffService().Delete(c.Request.Context(), &staff_service.IdRequest{Id: id})
	if err != nil {
		h.log.Error("error Staff Delete:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)

	err = h.red.Cache().Delete(c.Request.Context(), id)

	if err != nil {
		fmt.Println("Error while delete staff in cache")
	}

}

// @Router       /staffs [get]
// @Summary      List Staffs
// @Description  get Staff
// @Tags         STAFF
// @Accept       json
// @Produce      json
// @Param        limit    query     integer  true  "limit for response"  Default(10)
// @Param        page    query     integer  true  "page of req"  Default(1)
// @Param        name    query     string  false  "filter by name"
// @Success      200  {array}   models.Staff
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllStaff(c *gin.Context) {

	h.log.Info("request GetAllStaff")

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

	resp, err := h.grpcClient.StaffService().GetAll(c.Request.Context(), &staff_service.GetAllStaffRequest{
		Page:  int64(page),
		Limit: int64(limit),
		Name:  c.Query("name"),
	})

	if err != nil {
		h.log.Error("error Staff GetAll:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	h.log.Warn("response to GetAllStaff")
	c.JSON(http.StatusOK, resp)

}

// @Router       /change-password/{id} [PUT]
// @Summary      Change Password
// @Description  change staff's password
// @Tags         PASSWORD
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Staff ID" format(uuid)
// @Param        password    body     models.ChangePassword  true  "data of password"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) ChangePassword(c *gin.Context) {

	var pass models.ChangePassword
	err := c.ShouldBindJSON(&pass)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	id := c.Param("id")

	respGetStaff, err := h.grpcClient.StaffService().Get(c.Request.Context(), &staff_service.IdRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
		fmt.Println("error Sale Get:", err.Error())
		return
	}

	err = helper.ComparePasswords([]byte(respGetStaff.Password), []byte(pass.OldPassword))

	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		res := response.ErrorResp{Code: "INVALID Password", Message: "invalid password"}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	newHashPas, err := helper.GeneratePasswordHash(pass.NewPassword)

	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		res := response.ErrorResp{Code: "INVALID Password", Message: "invalid password"}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	resp, err := h.grpcClient.StaffService().ChangePassword(c.Request.Context(), &staff_service.RequestByPassword{
		Id:       id,
		Password: string(newHashPas),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
		fmt.Println("error Sale Get:", err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)

}

/* func (h *Handler) FindErnedSum(dateFrom, dateTo string) {

	resp, err := h.strg.Transaction().FindErnedSum(dateFrom, dateTo)

	if err != nil {
		fmt.Println("error from GetAllStaff:", err.Error())
		return
	}

	staffMap, err := h.strg.Staff().GetMapOfStaffs()

	if err != nil {
		fmt.Println("error from GetAllStaff:", err.Error())
		return
	}

	if err != nil {
		fmt.Println("error from GetAllStaff:", err.Error())
		return
	}

	branchMap := make(map[string]models.Branch)

	branch, err := h.strg.Branch().GetAllBranch(models.GetAllBranchRequest{
		Page:   1,
		Limit:  100,
		Search: "",
	})

	if err != nil {
		fmt.Println("error from GetAllStaff:", err.Error())
		return
	}

	for _, v := range branch.Branches {
		branchMap[v.ID] = v
	}

	for i, v := range resp {

		fmt.Println(staffMap[i].Name, branchMap[staffMap[i].BranchId].Name, v)

	}

} */
