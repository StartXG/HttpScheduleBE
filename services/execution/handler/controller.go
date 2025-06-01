package handler

import (
	"HttpScheduleBE/api/helper"
	"HttpScheduleBE/config"
	"HttpScheduleBE/services/execution/service"
	"HttpScheduleBE/services/executor"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	projectService *service.Service
	appConfig      *config.Config
}

func NewExecutionCenterController(projectService *service.Service, appConfig *config.Config) *Controller {
	return &Controller{
		projectService: projectService,
		appConfig:      appConfig,
	}
}

func (c *Controller) GetAllExecution(g *gin.Context) {
	executions, err := c.projectService.GetAllExecution()
	if err != nil {
		helper.RespondWithError(g, 500, err.Error())
		return
	}
	helper.RespondWithSuccess(g, 200, "Executions retrieved successfully", executions)
}

func (c *Controller) GetExecutingTask(g *gin.Context) {
	executions := executor.GetAllExecutingTasks()
	if executions == nil {
		helper.RespondWithError(g, 500, "No executing tasks found")
		return
	}
	helper.RespondWithSuccess(g, 200, "Executing tasks retrieved successfully", executions)
}
