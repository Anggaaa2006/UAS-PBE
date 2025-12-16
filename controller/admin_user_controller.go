// path: controller/admin_user_controller.go
package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "uas_pbe/internal/service" // or "uas_pbe/service" depending path
    "uas_pbe/utils"
)

/*
    AdminUserController
    - POST /admin/users
    - GET /admin/users
    - GET /admin/users/:id
    - PUT /admin/users/:id
    - DELETE /admin/users/:id
    - PUT /admin/users/:id/role
*/
type AdminUserController struct {
    svc *service.AdminUserService
}

func NewAdminUserController(s *service.AdminUserService) *AdminUserController {
    return &AdminUserController{svc: s}
}

// Create user
func (c *AdminUserController) Create(ctx *gin.Context) {
    var req struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
        Role     string `json:"role"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        utils.Error(ctx, http.StatusBadRequest, "invalid request body")
        return
    }
    u, err := c.svc.Create(ctx, req.Name, req.Email, req.Password, req.Role)
    if err != nil {
        utils.Error(ctx, http.StatusBadRequest, err.Error())
        return
    }
    utils.SuccessData(ctx, u)
}

// List users
func (c *AdminUserController) List(ctx *gin.Context) {
    data, err := c.svc.List(ctx)
    if err != nil {
        utils.Error(ctx, http.StatusInternalServerError, err.Error())
        return
    }
    utils.SuccessData(ctx, data)
}

// Get by id
func (c *AdminUserController) GetByID(ctx *gin.Context) {
    id := ctx.Param("id")
    u, err := c.svc.GetByID(ctx, id)
    if err != nil {
        utils.Error(ctx, http.StatusInternalServerError, err.Error())
        return
    }
    if u == nil {
        utils.Error(ctx, http.StatusNotFound, "user tidak ditemukan")
        return
    }
    utils.SuccessData(ctx, u)
}

// Update user
func (c *AdminUserController) Update(ctx *gin.Context) {
    id := ctx.Param("id")
    var req struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
        Role     string `json:"role"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        utils.Error(ctx, http.StatusBadRequest, "invalid request body")
        return
    }
    if err := c.svc.Update(ctx, id, req.Name, req.Email, req.Password, req.Role); err != nil {
        utils.Error(ctx, http.StatusBadRequest, err.Error())
        return
    }
    utils.Success(ctx, "user updated")
}

// Delete user
func (c *AdminUserController) Delete(ctx *gin.Context) {
    id := ctx.Param("id")
    if err := c.svc.Delete(ctx, id); err != nil {
        utils.Error(ctx, http.StatusInternalServerError, err.Error())
        return
    }
    utils.Success(ctx, "user deleted")
}

// Update role
func (c *AdminUserController) UpdateRole(ctx *gin.Context) {
    id := ctx.Param("id")
    var req struct {
        Role string `json:"role"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        utils.Error(ctx, http.StatusBadRequest, "invalid request body")
        return
    }
    if err := c.svc.UpdateRole(ctx, id, req.Role); err != nil {
        utils.Error(ctx, http.StatusInternalServerError, err.Error())
        return
    }
    utils.Success(ctx, "role updated")
}
