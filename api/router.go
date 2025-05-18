package api

import (
	apiTaskCenter "HttpScheduleBE/api/api_task_center"
	apiExecutionCenter "HttpScheduleBE/api/api_execution_center"
	domainTaskCenter "HttpScheduleBE/domain/domain_task_center"
	domainExecutionCenter "HttpScheduleBE/domain/domain_execution_center"
	"HttpScheduleBE/utils/database"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, dbs database.Databases) {
	registryTaskCenterHandler(r, dbs)
	registryExecutionCenterHandler(r, dbs)
}

func registryTaskCenterHandler(r *gin.Engine, dbs database.Databases) {
	taskCenterService := domainTaskCenter.NewTaskCenterService(*dbs.TaskCenterRepository)
	taskCenterController := apiTaskCenter.NewTaskCenterController(taskCenterService, database.AppConfig)
	taskCenterGroup := r.Group("/task")
	taskCenterGroup.POST("/", taskCenterController.CreateTask)
	taskCenterGroup.PUT("/:id", taskCenterController.UpdateTask)
	taskCenterGroup.DELETE("/:id", taskCenterController.DeleteTask)
	taskCenterGroup.GET("/", taskCenterController.GetAllTasks)
}

func registryExecutionCenterHandler(r *gin.Engine, dbs database.Databases) {
	executionCenterService := domainExecutionCenter.NewExecutionCenterService(*dbs.ExecutionCenterRepository)
	executionCenterController := apiExecutionCenter.NewExecutionCenterController(executionCenterService, database.AppConfig)
	executionCenterGroup := r.Group("/execution")
	executionCenterGroup.GET("/", executionCenterController.GetAllExecution)
	// executionCenterGroup.POST("/start/:id", executionCenterController.StartExecution)
	// executionCenterGroup.POST("/stop/:id", executionCenterController.StopExecution)
	// executionCenterGroup.GET("/status/:id", executionCenterController.GetExecutionStatus)
}
