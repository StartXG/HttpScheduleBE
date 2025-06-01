package executor

import (
	httpinvoke "HttpScheduleBE/pkgs/http-invoke"
	ExecutionRepo "HttpScheduleBE/services/execution/repo"
	"HttpScheduleBE/services/execution/types"
	TaskRepo "HttpScheduleBE/services/task/repo"
	"encoding/json"
	"strconv"
	"time"

	"fmt"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
)

var (
	executionList = make(map[cron.EntryID]*TaskExecution) // 使用 map 存储任务
	executionLock sync.Mutex                              // 用于并发安全
	c             = cron.New()                            // 初始化 cron 实例
)

type TaskExecution struct {
	ID       cron.EntryID
	TaskID   uint
	Name     string
	Schedule string
	Job      func()
	Status   string
}

type ExecuteResultForRecord struct {
	TaskID    uint      `json:"task_id"`
	Result    string    `json:"result"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	ErrMsg   string    `json:"error_msg,omitempty"`
}

// AddTask 添加任务
func AddTask(execution *TaskExecution) (string, error) {
	executionLock.Lock()
	defer executionLock.Unlock()

	executionId, err := c.AddFunc(execution.Schedule, execution.Job)
	if err != nil {
		return "", err
	}
	execution.ID = executionId
	executionList[executionId] = execution
	execution.Status = "pending" // 初始化状态为 pending
	return strconv.Itoa(int(executionId)), nil
}

// DeleteTask 删除任务
func DeleteTask(executionId cron.EntryID) {
	executionLock.Lock()
	defer executionLock.Unlock()

	c.Remove(executionId)
	delete(executionList, executionId)
}

// UpdateTask 更新任务
func UpdateTask(executionId cron.EntryID, newTask *TaskExecution) (string, error) {
	DeleteTask(executionId) // 删除旧任务
	return AddTask(newTask)
}

// GetTask 获取任务
func GetTask(executionId cron.EntryID) (*TaskExecution, bool) {
	executionLock.Lock()
	defer executionLock.Unlock()

	execution, exists := executionList[executionId]
	return execution, exists
}

// GetAllTasks 获取所有任务
func GetAllTasks() []*TaskExecution {
	executionLock.Lock()
	defer executionLock.Unlock()

	executions := make([]*TaskExecution, 0, len(executionList))
	for _, execution := range executionList {
		fmt.Println("Task Name:", execution)
		executions = append(executions, execution)
	}
	return executions
}

func GetAllExecutingTasks() []*types.ResponseExecutingTask {
	executionLock.Lock()
	defer executionLock.Unlock()

	var executingTasks []*types.ResponseExecutingTask
	for _, execution := range executionList {
		executingTasks = append(executingTasks, &types.ResponseExecutingTask{
			TaskID: execution.TaskID,
			Status: execution.Status,
			Name:   execution.Name,
		})
	}
	return executingTasks
}

// StartScheduler 启动调度器
func startScheduler() {
	fmt.Println("Starting scheduler...")
	c.Start()
	fmt.Println("Scheduler started.")
}

// StopScheduler 停止调度器
// func stopScheduler() {
// 	fmt.Println("Stopping scheduler...")
// 	c.Stop()
// 	fmt.Println("Scheduler stopped.")
// }

func StartExecutionAutomation(
	isAuto bool, 
	taskRepo *TaskRepo.Repository, 
	execRepo *ExecutionRepo.Repository,
	executeResult chan<- ExecuteResultForRecord,
) {
	if isAuto {
		taskRepos, err := taskRepo.GetAllTasks()
		if err != nil {
			fmt.Println("Failed to get tasks:", err)
			return
		}
		for _, task := range taskRepos {
			if !task.IsTaskEnabled {
				continue
			}
			t := task
			fmt.Println("Task detail:", t)
			// 创建一个新的 TaskExecution 实例
			te := &TaskExecution{
				Name:     t.TaskName,
				Schedule: t.TaskCron,
				TaskID:   t.ID,
				Job: func() {
					StartTime := time.Now()
					headers := make(map[string]string)
					if t.TaskHeader != "" && t.TaskHeader != "{}" && t.TaskHeader != "null" {
						err := json.Unmarshal([]byte(t.TaskHeader), &headers)
						if err != nil {
							fmt.Println("Failed to unmarshal task headers:", err)
							return
						}
					}
					httpinvoke := &httpinvoke.HttpInvoke{
						URL:     t.TaskUrl,
						Method:  t.TaskMethod,
						Headers: headers,
						Body:    strings.NewReader(t.TaskBody),
					}
					if response, err := httpinvoke.Invoke(); err != nil {
						executeResult <- ExecuteResultForRecord{
							TaskID:  t.ID,
							Result:  string(response),
							ErrMsg:  err.Error(),
							StartTime: StartTime,
							EndTime:   time.Now(),
						}
					}else{
						// 任务执行成功
						executeResult <- ExecuteResultForRecord{
							TaskID:    t.ID,
							Result:    string(response),
							StartTime: StartTime,
							EndTime:   time.Now(),
						}	
					}
					
				},
			}
			// 添加任务到调度器
			id, err := AddTask(te)
			fmt.Println("Task ID:", id)
			if err != nil {
				fmt.Println("Failed to add task:", err)
			} else {
				fmt.Println("Task added with ID:", id)
			}
		}

		startScheduler()
	} else {
		fmt.Println("Automation is disabled.")
	}
}