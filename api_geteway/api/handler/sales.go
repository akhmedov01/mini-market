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

// CreateSale godoc
// @Router       /sales [POST]
// @Summary      Create Sale
// @Description  Create Sale
// @Tags         SALE
// @Accept       json
// @Produce      json
// @Param        data  body      models.CreateSale  true  "Sale data"
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) CreateSale(c *gin.Context) {

	var sale models.CreateSale
	err := c.ShouldBindJSON(&sale)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	resp, err := h.strg.Sales().CreateSale(c.Request.Context(), sale)
	if err != nil {
		fmt.Println("error Sale Create:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	cashier, err := h.strg.Staff().GetStaff(c.Request.Context(), models.IdRequest{ID: sale.CashierId})

	if err != nil {

		fmt.Println("error from import cashier:", err.Error())

		return
	}

	tarifCashier, err := h.strg.StaffTarif().GetStaffTarif(c.Request.Context(), models.IdRequest{ID: cashier.TarifId})

	if err != nil {

		fmt.Println("error from import tariff:", err.Error())

		return
	}

	amount := 0.0

	if tarifCashier.Type == "Fixed" {

		if sale.PaymentType == "Cash" {
			amount = tarifCashier.AmountForCash
		} else {
			amount = tarifCashier.AmountForCard
		}

		//h.strg.Staff().ChangeBalance(cashierId, amount)

	} else if tarifCashier.Type == "Percent" {

		if sale.PaymentType == "Cash" {
			amount = (sale.Price * tarifCashier.AmountForCash) / 100
		} else {
			amount = (sale.Price * tarifCashier.AmountForCard) / 100
		}

		//h.strg.Staff().ChangeBalance(cashierId, amount)

	}

	reqUpB := models.UpdateBalanceRequest{
		SaleId: resp,
		Cashier: models.StaffType{
			StaffId: sale.CashierId,
			Amount:  amount,
		},
		TransactionType: "Topup",
		SourceType:      "Sales",
		Text:            "plus",
	}

	if sale.ShopAssistentId != "" {

		shopAssistant, err := h.strg.Staff().GetStaff(c.Request.Context(), models.IdRequest{ID: sale.ShopAssistentId})

		if err != nil {

			fmt.Println("error from import shopAssistent:", err.Error())

			return
		}

		tarifAssistent, err := h.strg.StaffTarif().GetStaffTarif(c.Request.Context(), models.IdRequest{ID: shopAssistant.TarifId})

		if err != nil {

			fmt.Println("error from import tariff:", err.Error())

			return
		}

		amount := 0.0

		if tarifAssistent.Type == "Fixed" {

			if sale.PaymentType == "Cash" {
				amount = tarifAssistent.AmountForCash
			} else {
				amount = tarifAssistent.AmountForCard
			}

			//h.strg.Staff().ChangeBalance(shopAssistentId, amount)

		} else if tarifAssistent.Type == "Percent" {

			if sale.PaymentType == "Cash" {
				amount = (sale.Price * tarifAssistent.AmountForCash) / 100
			} else {
				amount = (sale.Price * tarifAssistent.AmountForCard) / 100
			}

			//h.strg.Staff().ChangeBalance(shopAssistentId, amount)
		}

		reqUpB.ShopAssistent.StaffId = sale.ShopAssistentId
		reqUpB.ShopAssistent.Amount = amount
	}

	h.strg.Staff().UpdateBalance(c.Request.Context(), reqUpB)

	c.JSON(http.StatusCreated, response.CreateResponse{Id: resp})

}

// @Router       /sales/{id} [put]
// @Summary      update sale
// @Description  api for update sales
// @Tags         SALE
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of sale"
// @Param        sale    body     models.CreateSale  true  "data of sale"
// @Success      200  {string}   string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) UpdateSale(c *gin.Context) {
	var sale models.CreateSale
	err := c.ShouldBindJSON(&sale)
	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		c.JSON(http.StatusBadRequest, "invalid body")
		return
	}
	id := c.Param("id")

	resp, err := h.strg.Sales().UpdateSale(c.Request.Context(), id, sale)
	if err != nil {
		fmt.Println("error Sale Update:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)

	err = h.red.Cache().Delete(c.Request.Context(), id)

	if err != nil {
		fmt.Println("Error while delete sale in cache")
	}
}

// @Router       /sales/{id} [GET]
// @Summary      Get By Id
// @Description  get sale by ID
// @Tags         SALE
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Sale ID" format(uuid)
// @Success      200  {object}  models.Sale
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetSale(c *gin.Context) {

	id := c.Param("id")

	var sale = &models.Sale{}

	response, err := h.red.Cache().Get(c.Request.Context(), id, sale)

	if err != nil {
		fmt.Println("Error while geting sale in cache")
	}

	if response {
		c.JSON(http.StatusOK, sale)
		return
	}

	sale, err = h.strg.Sales().GetSale(c.Request.Context(), models.IdRequest{ID: id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, "internal server error")
		fmt.Println("error Sale Get:", err.Error())
		return
	}

	c.JSON(http.StatusOK, sale)

	err = h.red.Cache().Create(c.Request.Context(), id, sale, 0)

	if err != nil {
		fmt.Println("Error while Create sale in cache")
	}

}

// @Router       /sales/{id} [DELETE]
// @Summary      Delete By Id
// @Description  delete sale by Id
// @Tags         SALE
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Sale ID" format(uuid)
// @Success      200  {string}  string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteSale(c *gin.Context) {

	id := c.Param("id")

	if !helper.IsValidUUID(id) {

		h.log.Error("error Sale Delete:", logger.Error(errors.New("invalid id")))
		c.JSON(http.StatusBadRequest, "invalid id")
		return

	}
	resp, err := h.strg.Sales().DeleteSale(c.Request.Context(), models.IdRequest{ID: id})
	if err != nil {
		h.log.Error("error Sale Delete:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, resp)

	err = h.red.Cache().Delete(c.Request.Context(), id)

	if err != nil {
		fmt.Println("Error while delete sale in cache")
	}

}

// @Router       /sales [get]
// @Summary      List Sales
// @Description  get Sale
// @Tags         SALE
// @Accept       json
// @Produce      json
// @Param        limit    query     integer  true  "limit for response"  Default(10)
// @Param        page    query     integer  true  "page of req"  Default(1)
// @Param        client_name    query     string  false  "filter by client_name"
// @Success      200  {array}   models.Sale
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) GetAllSale(c *gin.Context) {

	h.log.Info("request GetAllSale")

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

	resp, err := h.strg.Sales().GetAllSale(c.Request.Context(), models.GetAllSaleRequest{
		Page:       page,
		Limit:      limit,
		ClientName: c.Query("client_name"),
	})
	if err != nil {
		h.log.Error("error Sale GetAll:", logger.Error(err))
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	h.log.Warn("response to GetAllSale")
	c.JSON(http.StatusOK, resp)

}

/* func (h *Handler) BranchTotal() {
	resp, err := h.strg.Sales().BranchTotal()

	if err != nil {
		fmt.Println("Error geting BranchTotal")
		return
	}

	branches, _ := h.strg.Branch().GetAllBranch(models.GetAllBranchRequest{
		Page:   1,
		Limit:  100,
		Search: "",
	})

	branchesMap := make(map[string]models.Branch)

	for _, v := range branches.Branches {
		branchesMap[v.ID] = v
	}

	for i, v := range resp {
		fmt.Println(branchesMap[i].Name, v.Count, v.Sum)
	}
}

func (h *Handler) GetSalesInDay() {
	resp, err := h.strg.Sales().GetSalesInDay()

	if err != nil {
		fmt.Println("Error geting GetSalesInDay")
		return
	}

	branches, _ := h.strg.Branch().GetAllBranch(models.GetAllBranchRequest{
		Page:   1,
		Limit:  100,
		Search: "",
	})

	branchesMap := make(map[string]models.Branch)

	for _, v := range branches.Branches {
		branchesMap[v.ID] = v
	}

	for branchId, bValue := range resp {
		max := 0.0
		fmt.Println(branchesMap[branchId].Name)
		for _, v := range bValue {
			if max < v {
				max = v
			}
		}
		fmt.Print(max)
	}
}
*/
/* func (h *Handler) CancelSale(id string) {
	resp, err := h.strg.Sales().CancelSale(id)

	if err != nil {
		fmt.Println("error from CreateSale:", err.Error())
		return
	}

	cashier, err := h.strg.Staff().GetStaff(models.IdRequest{ID: resp.CashierId})

	if err != nil {

		fmt.Println("error from import cashier:", err.Error())

		return
	}

	tarifCashier, err := h.strg.StaffTarif().GetStaffTarif(models.IdRequest{ID: cashier.TarifId})

	if err != nil {

		fmt.Println("error from import tariff:", err.Error())

		return
	}

	amount := 0.0

	if tarifCashier.Type == config.Fixed {

		if resp.PaymentType == config.Cash {
			amount = tarifCashier.AmountForCash
		} else {
			amount = tarifCashier.AmountForCard
		}

		h.strg.Staff().ChangeBalance(resp.CashierId, -amount)

	} else if tarifCashier.Type == config.Percent {

		if resp.PaymentType == config.Cash {
			amount = (resp.Price * tarifCashier.AmountForCash) / 100
		} else {
			amount = (resp.Price * tarifCashier.AmountForCard) / 100
		}

		h.strg.Staff().ChangeBalance(resp.CashierId, -amount)

	}

	h.strg.Transaction().CreateStaffTransaction(models.CreateStaffTransaction{
		SaleID:          resp.ID,
		StaffID:         resp.CashierId,
		TransactionType: 1,
		SourceType:      1,
		Amount:          amount,
		Text:            "minus",
	})

	if resp.ShopAssistentId != "" {

		shopAssistant, err := h.strg.Staff().GetStaff(models.IdRequest{ID: resp.ShopAssistentId})

		if err != nil {

			fmt.Println("error from import shopAssistent:", err.Error())

			return
		}

		tarifAssistent, err := h.strg.StaffTarif().GetStaffTarif(models.IdRequest{ID: shopAssistant.TarifId})

		if err != nil {

			fmt.Println("error from import tariff:", err.Error())

			return
		}

		amount := 0.0

		if tarifAssistent.Type == config.Fixed {

			if resp.PaymentType == config.Cash {
				amount = -tarifAssistent.AmountForCash
			} else {
				amount = -tarifAssistent.AmountForCard
			}

			h.strg.Staff().ChangeBalance(resp.ShopAssistentId, amount)

		} else if tarifAssistent.Type == config.Percent {

			if resp.PaymentType == config.Cash {
				amount = -(resp.Price * tarifAssistent.AmountForCash)
			} else {
				amount = -(resp.Price * tarifAssistent.AmountForCard)
			}

			h.strg.Staff().ChangeBalance(resp.ShopAssistentId, amount)
		}

		h.strg.Transaction().CreateStaffTransaction(models.CreateStaffTransaction{
			SaleID:          resp.ID,
			StaffID:         resp.ShopAssistentId,
			TransactionType: 1,
			SourceType:      1,
			Amount:          amount,
			Text:            "minus",
		})

	}

} */

/* resp, err := h.strg.Sales().CreateSale(models.CreateSale{
	BranchId:        branchId,
	ShopAssistentId: shopAssistentId,
	CashierId:       cashierId,
	PaymentType:     paymentType,
	Status:          status,
	ClientName:      clientName,
	Price:           price,
})

if err != nil {
	fmt.Println("error from CreateSale:", err.Error())
	return
}

cashier, err := h.strg.Staff().GetStaff(models.IdRequest{ID: cashierId})

if err != nil {

	fmt.Println("error from import cashier:", err.Error())

	return
}

tarifCashier, err := h.strg.StaffTarif().GetStaffTarif(models.IdRequest{ID: cashier.TarifId})

if err != nil {

	fmt.Println("error from import tariff:", err.Error())

	return
}

amount := 0.0

if tarifCashier.Type == config.Fixed {

	if paymentType == config.Cash {
		amount = tarifCashier.AmountForCash
	} else {
		amount = tarifCashier.AmountForCard
	}

	h.strg.Staff().ChangeBalance(cashierId, amount)

} else if tarifCashier.Type == config.Percent {

	if paymentType == config.Cash {
		amount = price * tarifCashier.AmountForCash
	} else {
		amount = price * tarifCashier.AmountForCard
	}

	h.strg.Staff().ChangeBalance(cashierId, amount)

}

h.strg.Transaction().CreateStaffTransaction(models.CreateStaffTransaction{
	SaleID:          resp,
	StaffID:         cashierId,
	TransactionType: 2,
	SourceType:      1,
	Amount:          amount,
	Text:            "plus",
})

if shopAssistentId != "" {

	shopAssistant, err := h.strg.Staff().GetStaff(models.IdRequest{ID: shopAssistentId})

	if err != nil {

		fmt.Println("error from import shopAssistent:", err.Error())

		return
	}

	tarifAssistent, err := h.strg.StaffTarif().GetStaffTarif(models.IdRequest{ID: shopAssistant.TarifId})

	if err != nil {

		fmt.Println("error from import tariff:", err.Error())

		return
	}

	amount := 0.0

	if tarifAssistent.Type == config.Fixed {

		if paymentType == config.Cash {
			amount = tarifAssistent.AmountForCash
		} else {
			amount = tarifAssistent.AmountForCard
		}

		h.strg.Staff().ChangeBalance(shopAssistentId, amount)

	} else if tarifAssistent.Type == config.Percent {

		if paymentType == config.Cash {
			amount = price * tarifAssistent.AmountForCash
		} else {
			amount = price * tarifAssistent.AmountForCard
		}

		h.strg.Staff().ChangeBalance(shopAssistentId, amount)
	}

	h.strg.Transaction().CreateStaffTransaction(models.CreateStaffTransaction{
		SaleID:          resp,
		StaffID:         shopAssistentId,
		TransactionType: 2,
		SourceType:      1,
		Amount:          amount,
		Text:            "plus",
	})

}

*/
