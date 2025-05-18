package executor

import (
	domainExecutionCenter "HttpScheduleBE/domain/domain_execution_center"
	domaintaskcenter "HttpScheduleBE/domain/domain_task_center"
	"HttpScheduleBE/domain/entity"
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
	executionLock sync.Mutex                             // 用于并发安全
	c        = cron.New()                           // 初始化 cron 实例
)

type TaskExecution struct {
	ID       cron.EntryID
	Name     string
	Schedule string
	Job      func()
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
	return strconv.Itoa(int(executionId)), nil
}

// DeleteTask 删除任务
func DeleteTask(executionId cron.EntryID) {
	executionLock.Lock()
	defer executionLock.Unlock()

	c.Remove(executionId)
	delete(executionList, executionId)
}

// 更新任务
func UpdateTask(executionId cron.EntryID, newTask *TaskExecution) (string,error) {
	DeleteTask(executionId) // 删除旧任务
	return AddTask(newTask)
}

// 获取任务
func GetTask(executionId cron.EntryID) (*TaskExecution, bool) {
	executionLock.Lock()
	defer executionLock.Unlock()

	execution, exists := executionList[executionId]
	return execution, exists
}

// 获取所有任务
func GetAllTasks() []*TaskExecution {
	executionLock.Lock()
	defer executionLock.Unlock()

	executions := make([]*TaskExecution, 0, len(executionList))
	for _, execution := range executionList {
		executions = append(executions, execution)
	}
	return executions
}

// 启动调度器
func StartScheduler() {
	fmt.Println("Starting scheduler...")
	c.Start()
	fmt.Println("Scheduler started.")
}

// 停止调度器
func StopScheduler() {
	fmt.Println("Stopping scheduler...")
	c.Stop()
	fmt.Println("Scheduler stopped.")
}

func StartExecutionAutomation(isAuto bool, taskRepo *domaintaskcenter.Repository, execRepo *domainExecutionCenter.Repository) {
	if isAuto {
		taskRepos, err := taskRepo.GetAllTasks()
		if err != nil {
			fmt.Println("Failed to get tasks:", err)
			return
		}
		for _, task := range taskRepos {
			// 创建一个新的 TaskExecution 实例
			te := &TaskExecution{
				Name:     task.TaskName,
				Schedule: task.TaskCron,
				Job: func() {
					// 这里可以调用实际的 HTTP 请求逻辑
					// 例如使用 http_task 包中的方法
					headers := make(map[string]string)
					json.Unmarshal([]byte(task.TaskHeader), &headers)
					executeHttpTask(
						task.TaskUrl, 
						task.TaskMethod, 
						headers, 
						task.TaskBody,
						execRepo,
					)
				},
			}
			// 添加任务到调度器
			_, err := AddTask(te)
			if err != nil {
				fmt.Println("Failed to add task:", err)
			}
		}
		StartScheduler()	
	} else {
		fmt.Println("Automation is disabled.")
	}
}

func executeHttpTask(
	url string, 
	method string, 
	header map[string]string, 
	body string,
	execRepo *domainExecutionCenter.Repository,
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
		defer res.Body.Close()
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
		TaskID:   0, // 这里需要设置为实际的任务ID
		Status:   strconv.Itoa(statusCode),
		StartTime: startTime,
		EndTime:   endTime,
		ErrorLog: ErrLog,
	}
	
	if execRepo != nil {
		if saveErr := execRepo.CreateExecution(execRecord); saveErr != nil {
			fmt.Println("Failed to save execution record:", saveErr)
		}
	}

	fmt.Println(string(responseBody))
}
