package commands

import (
	"HttpScheduleBE/domain/executor"
	"HttpScheduleBE/api"
	"HttpScheduleBE/config"
	"HttpScheduleBE/utils/database"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func HttpScheduleBeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "http-schedule",
		Short: "Schedule a command to run at a specific time",
		Long:  `Schedule a command to run at a specific time`,
		Run: func(cmd *cobra.Command, args []string) {
			cfgFile := "./etc/config.yaml"
			cfg, err := config.GetConfigValues(cfgFile)
			if err != nil {
				fmt.Printf("Failed to load config file: %v\n", err)
				return
			}
			dbs := *database.CreateDBs(cfg)
			r := gin.Default()
			r.Use(cors.New(cors.Config{
				AllowOrigins:     []string{"*"}, // 允许所有域名，生产环境建议修改
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour, // 预检请求的缓存时间
			}))
			api.RegisterRoutes(r, dbs)
			executor.StartExecutionAutomation(
				cfg.ExecuteAutomatic,
				dbs.TaskCenterRepository,
				dbs.ExecutionCenterRepository,
			)
			if err := r.Run(":8080"); err != nil {
				fmt.Printf("Failed to start server: %v\n", err)
			}
			fmt.Println("Server started on :8080")
		},
	}
	// cmd.AddCommand(SeedCmd()) // 添加 seed 命令
	return cmd
}
