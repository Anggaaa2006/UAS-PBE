package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "uas_pbe/service"
)

/*
    StatsController
    Controller statistik prestasi
*/
type StatsController struct {
    service service.StatsService
}

func NewStatsController(s service.StatsService) *StatsController {
    return &StatsController{service: s}
}

/*
    GET /stats/achievements
*/
func (c *StatsController) GetAchievementStats(ctx *gin.Context) {

    data, err := c.service.GetAchievementStats(ctx.Request.Context())
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "message": err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": data,
    })
}

/*
    GET /reports/student/:id
    Statistik prestasi mahasiswa

    FR-011
*/
func (c *StatsController) GetStudentReport(ctx *gin.Context) {
    studentID := ctx.Param("id")

    data, err := c.service.GetStudentStats(ctx, studentID)
    if err != nil {
        utils.Error(ctx, 500, err.Error())
        return
    }

    utils.SuccessData(ctx, data)
}
