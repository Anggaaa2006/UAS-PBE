package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"uas_pbe/service"
	"uas_pbe/utils"
)

type AdminUserController struct {
	svc service.AdminUserService
}

func NewAdminUserController(s service.AdminUserService) *AdminUserController {
	return &AdminUserController{svc: s}
}

// GET /users
func (c *AdminUserController) List(ctx *gin.Context) {
	data, err := c.svc.List(ctx.Request.Context())
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessData(ctx, data)
}

// GET /users/:id
func (c *AdminUserController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	data, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, err.Error())
		return
	}
	utils.SuccessData(ctx, data)
}

// POST /users
func (c *AdminUserController) Create(ctx *gin.Context) {
	var req service.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, 400, "invalid body")
		return
	}

	if err := c.svc.Create(ctx.Request.Context(), req); err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}
	utils.Success(ctx, "user berhasil dibuat")
}

// PUT /users/:id
func (c *AdminUserController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req service.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, 400, "invalid body")
		return
	}

	if err := c.svc.Update(ctx.Request.Context(), id, req); err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}
	utils.Success(ctx, "user berhasil diupdate")
}

// DELETE /users/:id
func (c *AdminUserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.svc.Delete(ctx.Request.Context(), id); err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}
	utils.Success(ctx, "user berhasil dihapus")
}

// PUT /users/:id/role
func (c *AdminUserController) UpdateRole(ctx *gin.Context) {
	id := ctx.Param("id")

	var req struct {
		Role string `json:"role"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, 400, "invalid body")
		return
	}

	if err := c.svc.UpdateRole(ctx.Request.Context(), id, req.Role); err != nil {
		utils.Error(ctx, 400, err.Error())
		return
	}
	utils.Success(ctx, "role berhasil diubah")
}
