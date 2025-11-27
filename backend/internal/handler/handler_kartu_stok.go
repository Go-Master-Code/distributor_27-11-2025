package handler

import (
	"api-distributor/helper"
	"api-distributor/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// struct implementasi
type handlerKartuStok struct {
	service service.ServiceKartuStok
}

// constructor
func NewHandlerKartuStok(service service.ServiceKartuStok) *handlerKartuStok {
	return &handlerKartuStok{service}
}

// func (h *handlerKartuStok) GetAllKartuStok(c *gin.Context) {
// 	barang, err := h.service.GetAllKartuStok()
// 	if err != nil {
// 		helper.ErrorFetchDataFromDB(c, err)
// 		return
// 	}

// 	helper.StatusSuksesGetData(c, barang)
// }

// func (h *handlerKartuStok) CreateKartuStok(c *gin.Context) {
// 	// parsing request body
// 	var barang dto.CreateKartuStokRequest
// 	err := c.ShouldBindJSON(&barang)
// 	if err != nil {
// 		helper.ErrorParsingRequestBody(c, err)
// 		return
// 	}

// 	newKartuStok, err := h.service.CreateKartuStok(barang)
// 	if err != nil {
// 		helper.ErrorCreateData(c, err)
// 		return
// 	}

// 	helper.StatusSuksesCreateData(c, newKartuStok)
// }

func (h *handlerKartuStok) GetAllKartuStok(c *gin.Context) {
	kartuStok, err := h.service.GetAllKartuStok()
	if err != nil {
		helper.ErrorDataNotFound(c)
		return
	}

	helper.StatusSuksesGetData(c, kartuStok)
}

func (h *handlerKartuStok) GetKartuStokByIdBarang(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}

	barang, err := h.service.GetKartuStokByIdBarang(idInt)
	if err != nil {
		helper.ErrorDataNotFound(c)
		return
	}

	helper.StatusSuksesGetData(c, barang)
}

// func (h *handlerKartuStok) GetKartuStok(c *gin.Context) {
// 	id := c.Query("id")
// 	nama := c.Query("nama")

// 	if id == "" && nama == "" {
// 		h.GetAllKartuStok(c)
// 	}

// 	if id != "" {
// 		// cari data by ID
// 		idInt, err := strconv.Atoi(id)
// 		if err != nil {
// 			helper.ErrorParsingAtoi(c, err)
// 			return
// 		}

// 		merk, err := h.service.GetKartuStokById(idInt)
// 		if err != nil {
// 			helper.ErrorDataNotFound(c)
// 			return
// 		}

// 		helper.StatusSuksesGetData(c, merk)
// 	}

// 	if nama != "" {
// 		merk, err := h.service.SearchKartuStok(nama)
// 		if err != nil {
// 			helper.ErrorDataNotFound(c)
// 			return
// 		}

// 		helper.StatusSuksesGetData(c, merk)
// 	}
// }

// func (h *handlerKartuStok) UpdateKartuStok(c *gin.Context) {
// 	// parsing json
// 	var req dto.UpdateKartuStokRequest

// 	err := c.ShouldBindJSON(&req)
// 	if err != nil {
// 		helper.ErrorParsingRequestBody(c, err)
// 		return
// 	}

// 	// ambil param id
// 	id := c.Param("id")
// 	idInt, err := strconv.Atoi(id)
// 	if err != nil {
// 		helper.ErrorParsingAtoi(c, err)
// 		return
// 	}

// 	barangDTO, err := h.service.UpdateKartuStok(idInt, req)

// 	if err != nil {
// 		helper.ErrorUpdateData(c, err)
// 		return
// 	}

// 	helper.StatusSuksesUpdateData(c, barangDTO)
// }
