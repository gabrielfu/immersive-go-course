package main

import "github.com/gin-gonic/gin"

type Request struct {
	Schedule string `json:"schedule"`
	Command  string `json:"command"`
}

type JobHandler func(schedule, command string) error

func ServeAPI(jobHandler JobHandler, port string) {
	router := gin.Default()
	router.POST("/jobs", func(ctx *gin.Context) {
		var req Request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := jobHandler(req.Schedule, req.Command); err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"msg": "ok"})
	})
	router.Run(":" + port)
}
