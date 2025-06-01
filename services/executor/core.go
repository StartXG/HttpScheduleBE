package executor

import (
	"HttpScheduleBE/entity"
	ExecutionRepo "HttpScheduleBE/services/execution/repo"
	"HttpScheduleBE/services/execution/types"
	TaskRepo "HttpScheduleBE/services/task/repo"
	"encoding/json"
	"strconv"
	"time"

	"fmt"
	"io"
	"net/http"
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
func StartScheduler() {
	fmt.Println("Starting scheduler...")
	c.Start()
	fmt.Println("Scheduler started.")
}

// StopScheduler 停止调度器
func StopScheduler() {
	fmt.Println("Stopping scheduler...")
	c.Stop()
	fmt.Println("Scheduler stopped.")
}

func StartExecutionAutomation(isAuto bool, taskRepo *TaskRepo.Repository, execRepo *ExecutionRepo.Repository) {
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
			fmt.Println("Task Name:", task)
			// 创建一个新的 TaskExecution 实例
			te := &TaskExecution{
				Name:     task.TaskName,
				Schedule: task.TaskCron,
				TaskID:   task.ID,
				Job: func() {
					// 这里可以调用实际的 HTTP 请求逻辑
					// 例如使用 http_task 包中的方法
					headers := make(map[string]string)
					err := json.Unmarshal([]byte(task.TaskHeader), &headers)
					if err != nil {
						return
					}
					executeHttpTask(
						task.ID,
						task.TaskUrl,
						task.TaskMethod,
						headers,
						task.TaskBody,
						execRepo,
					)
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

		StartScheduler()
	} else {
		fmt.Println("Automation is disabled.")
	}
}

func executeHttpTask(
	taskID uint,
	url string,
	method string,
	header map[string]string,
	body string,
	execRepo *ExecutionRepo.Repository,
) {
	var ErrLog string
	startTime := time.Now().Format("2006-01-02 15:04:05")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		fmt.Println(err)
		return
	}
	for key, value := range header {
		req.Header.Add(key, value)
	}
	if value, ok := header["Content-Type"]; !ok || value != "application/json" {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	statusCode := 0
	var responseBody []byte
	if err == nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println("Error closing response body:", err)
			}
		}(res.Body)
		statusCode = res.StatusCode
		if statusCode != 200 {
			ErrLog = fmt.Sprintf("Error: %s", res.Status)
		} else {
			ErrLog = ""
		}
		responseBody, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	endTime := time.Now().Format("2006-01-02 15:04:05")

	// 记录到表中
	if err != nil {
		ErrLog = fmt.Sprintf("Error: %s", err.Error())
	}
	execRecord := &entity.ExecutionCenter{
		TaskID:    taskID,
		Status:    strconv.Itoa(statusCode),
		StartTime: startTime,
		EndTime:   endTime,
		ErrorLog:  ErrLog,
	}

	if execRepo != nil {
		if saveErr := execRepo.CreateExecution(execRecord); saveErr != nil {
			fmt.Println("Failed to save execution record:", saveErr)
		}
	}

	fmt.Println(string(responseBody))
}
