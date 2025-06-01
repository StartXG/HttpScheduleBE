package api

import (
	ExecutionHandler "HttpScheduleBE/services/execution/handler"
	ExecutionService "HttpScheduleBE/services/execution/service"
	TaskHandler "HttpScheduleBE/services/task/handler"
	TaskService "HttpScheduleBE/services/task/service"
	"HttpScheduleBE/utils/database"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, dbs database.Databases) {
	registryTaskCenterHandler(r, dbs)
	registryExecutionCenterHandler(r, dbs)
}

func registryTaskCenterHandler(r *gin.Engine, dbs database.Databases) {
	taskCenterService := TaskService.NewTaskCenterService(*dbs.TaskCenterRepository)
	taskCenterController := TaskHandler.NewTaskCenterController(taskCenterService, database.AppConfig)
	taskCenterGroup := r.Group("/task")
	taskCenterGroup.POST("/", taskCenterController.CreateTask)
	taskCenterGroup.PUT("/:id", taskCenterController.UpdateTask)
	taskCenterGroup.DELETE("/:id", taskCenterController.DeleteTask)
	taskCenterGroup.GET("/", taskCenterController.GetAllTasks)
}

func registryExecutionCenterHandler(r *gin.Engine, dbs database.Databases) {
	executionCenterService := ExecutionService.NewExecutionCenterService(*dbs.ExecutionCenterRepository)
	executionCenterController := ExecutionHandler.NewExecutionCenterController(executionCenterService, database.AppConfig)
	executionCenterGroup := r.Group("/execution")
	executionCenterGroup.GET("/", executionCenterController.GetAllExecution)
	executionCenterGroup.GET("/running", executionCenterController.GetExecutingTask)
	// executionCenterGroup.POST("/start/:id", executionCenterController.StartExecution)
	// executionCenterGroup.POST("/stop/:id", executionCenterController.StopExecution)
	// executionCenterGroup.GET("/status/:id", executionCenterController.GetExecutionStatus)
}
