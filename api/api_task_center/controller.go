package apitaskcenter

import (
	"HttpScheduleBE/api/helper"
	"HttpScheduleBE/config"
	domainTaskCenter "HttpScheduleBE/domain/domain_task_center"
	typesTaskCenter "HttpScheduleBE/api/types"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	projectService *domainTaskCenter.Service
	appConfig      *config.Config
}

func NewTaskCenterController(projectService *domainTaskCenter.Service, appConfig *config.Config) *Controller {
	return &Controller{
		projectService: projectService,
		appConfig:      appConfig,
	}
}

func (c *Controller) CreateTask(g *gin.Context) {
	var req typesTaskCenter.RequestTaskCenter
	if err := g.ShouldBindJSON(&req); err != nil {
		helper.RespondWithError(g, 400, err.Error())
		return
	}

	if err := c.projectService.CreateTask(&req); err != nil {
		helper.RespondWithError(g, 500, err.Error())
		return
	}
	helper.RespondWithSuccess(g, 200, "Task created successfully", &req)
}

func (c *Controller) UpdateTask(g *gin.Context) {
	taskId := g.Param("id")
	if taskId == "" {
		helper.RespondWithError(g, 400, "Task ID is required")
		return
	}
	var req typesTaskCenter.RequestTaskCenter
	if err := g.ShouldBindJSON(&req); err != nil {
		helper.RespondWithError(g, 400, err.Error())
		return
	}
	if err := c.projectService.UpdateTask(taskId, &req); err != nil {
		helper.RespondWithError(g, 500, err.Error())
		return
	}
	helper.RespondWithSuccess(g, 200, "Task updated successfully", &req)
}

func (c *Controller) DeleteTask(g *gin.Context) {
	id := g.Param("id")
	if err := c.projectService.DeleteTask(id); err != nil {
		helper.RespondWithError(g, 500, err.Error())
		return
	}
	helper.RespondWithSuccess(g, 200, "Task deleted successfully", nil)
}

func (c *Controller) GetAllTasks(g *gin.Context) {
	tasks, err := c.projectService.GetAllTasks()
	if err != nil {
		helper.RespondWithError(g, 500, err.Error())
		return
	}
	helper.RespondWithSuccess(g, 200, "Tasks retrieved successfully", tasks)
}