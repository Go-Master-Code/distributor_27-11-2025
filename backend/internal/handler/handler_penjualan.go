package handler

import (
	"api-distributor/helper"
	"api-distributor/internal/dto"
	"api-distributor/internal/service"

	"github.com/gin-gonic/gin"
)

type handlerPenjualan struct {
	service service.ServicePenjualan
}

func NewHandlerPenjualan(service service.ServicePenjualan) *handlerPenjualan {
	return &handlerPenjualan{service}
}

// func (h *handlerPenjualan) GetAllPenjualan(c *gin.Context) {
// 	penjualan, err := h.service.GetAllPenjualan()
// 	if err != nil {
// 		helper.ErrorDataNotFound(c)
// 		return
// 	}

// 	helper.StatusSuksesGetData(c, penjualan)
// }

func (h *handlerPenjualan) CreatePenjualan(c *gin.Context) {
	// parsing json body
	var penjualan dto.CreatePenjualanRequest
	err := c.ShouldBindJSON(&penjualan)

	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	newPenjualan, err := h.service.CreatePenjualan(penjualan)
	if err != nil {
		helper.ErrorCreateData(c, err)
		return
	}

	helper.StatusSuksesCreateData(c, newPenjualan)
}
