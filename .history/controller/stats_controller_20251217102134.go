package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "uas_pbe/service"
    "uas_pbe/utils"
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

    Statistik prestasi per mahasiswa

    KETERKAITAN SRS:
    - FR-011: Achievement Statistics
*/
func (c *StatsController) GetStudentReport(ctx *gin.Context) {

    // Ambil student_id dari URL
    studentID := ctx.Param("id")

    // Panggil service
    data, err := c.service.GetStudentStats(ctx.Request.Context(), studentID)
    if err != nil {
        utils.Error(ctx, http.StatusInternalServerError, err.Error())
        return
    }

    utils.SuccessData(ctx, data)
}
