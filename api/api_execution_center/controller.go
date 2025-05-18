package apiexecutioncenter

import (
	"HttpScheduleBE/api/helper"
	"HttpScheduleBE/config"
	domainExecutionCenter "HttpScheduleBE/domain/domain_execution_center"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	projectService *domainExecutionCenter.Service
	appConfig      *config.Config
}

func NewExecutionCenterController(projectService *domainExecutionCenter.Service, appConfig *config.Config) *Controller {
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