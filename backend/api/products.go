package api

import (
	"fmt"
	"net/http"
	"strconv"

	"tg-robot-sim/services"
)

// handleProducts 处理产品列表请求
func (h *MiniAppApiService) handleProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("====获取产品列表 ==========")

	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	ctx := r.Context()

	// 获取查询参数
	productType := r.URL.Query().Get("type")
	country := r.URL.Query().Get("country")
	search := r.URL.Query().Get("search")
	limit := h.parseIntParam(r, "limit", 20)
	offset := h.parseIntParam(r, "offset", 0)

	// 构建筛选条件
	filters := services.ProductFilters{
		Type:    productType,
		Country: country,
		Search:  search,
		Limit:   limit,
		Offset:  offset,
	}

	// 获取产品列表
	products, err := h.productService.GetProducts(ctx, filters)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to get products", err.Error())
		return
	}

	// 获取总数
	total, _ := h.productService.CountProducts(ctx, filters)

	// 返回响应
	h.sendSuccess(w, map[string]interface{}{
		"products": products,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

// handleProductDetail 处理产品详情请求
func (h *MiniAppApiService) handleProductDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	ctx := r.Context()

	// 从 URL 路径提取产品 ID
	// /api/miniapp/products/123
	path := r.URL.Path
	idStr := path[len("/api/miniapp/products/"):]
	if idStr == "" {
		h.sendError(w, http.StatusBadRequest, "Product ID is required", "")
		return
	}

	productID, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	// 获取产品详情
	product, err := h.productService.GetProductByID(ctx, productID)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "Product not found", err.Error())
		return
	}

	h.sendSuccess(w, product)
}
