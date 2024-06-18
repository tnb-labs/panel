package plugins

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TheTNB/panel/pkg/tools"
	"github.com/goravel/framework/contracts/http"

	"github.com/TheTNB/panel/app/http/controllers"
)

type SupervisorController struct {
	service string
}

func NewSupervisorController() *SupervisorController {
	var service string
	if tools.IsRHEL() {
		service = "supervisord"
	} else {
		service = "supervisor"
	}

	return &SupervisorController{
		service: service,
	}
}

// Service 获取服务名称
func (r *SupervisorController) Service(ctx http.Context) http.Response {
	return controllers.Success(ctx, r.service)
}

// Log 日志
func (r *SupervisorController) Log(ctx http.Context) http.Response {
	log, err := tools.Exec(`tail -n 200 /var/log/supervisor/supervisord.log`)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, log)
	}

	return controllers.Success(ctx, log)
}

// ClearLog 清空日志
func (r *SupervisorController) ClearLog(ctx http.Context) http.Response {
	if out, err := tools.Exec(`echo "" > /var/log/supervisor/supervisord.log`); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// Config 获取配置
func (r *SupervisorController) Config(ctx http.Context) http.Response {
	var config string
	var err error
	if tools.IsRHEL() {
		config, err = tools.Read(`/etc/supervisord.conf`)
	} else {
		config, err = tools.Read(`/etc/supervisor/supervisord.conf`)
	}

	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, config)
}

// SaveConfig 保存配置
func (r *SupervisorController) SaveConfig(ctx http.Context) http.Response {
	config := ctx.Request().Input("config")
	var err error
	if tools.IsRHEL() {
		err = tools.Write(`/etc/supervisord.conf`, config, 0644)
	} else {
		err = tools.Write(`/etc/supervisor/supervisord.conf`, config, 0644)
	}

	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	if err = tools.ServiceRestart(r.service); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("重启 %s 服务失败", r.service))
	}

	return controllers.Success(ctx, nil)
}

// Processes 进程列表
func (r *SupervisorController) Processes(ctx http.Context) http.Response {
	type process struct {
		Name   string `json:"name"`
		Status string `json:"status"`
		Pid    string `json:"pid"`
		Uptime string `json:"uptime"`
	}

	out, err := tools.Exec(`supervisorctl status | awk '{print $1}'`)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	var processes []process
	for _, line := range strings.Split(out, "\n") {
		if len(line) == 0 {
			continue
		}

		var p process
		p.Name = line
		if status, err := tools.Exec(`supervisorctl status ` + line + ` | awk '{print $2}'`); err == nil {
			p.Status = status
		}
		if p.Status == "RUNNING" {
			pid, _ := tools.Exec(`supervisorctl status ` + line + ` | awk '{print $4}'`)
			p.Pid = strings.ReplaceAll(pid, ",", "")
			uptime, _ := tools.Exec(`supervisorctl status ` + line + ` | awk '{print $6}'`)
			p.Uptime = uptime
		} else {
			p.Pid = "-"
			p.Uptime = "-"
		}
		processes = append(processes, p)
	}

	paged, total := controllers.Paginate(ctx, processes)

	return controllers.Success(ctx, http.Json{
		"total": total,
		"items": paged,
	})
}

// StartProcess 启动进程
func (r *SupervisorController) StartProcess(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	if out, err := tools.Exec(`supervisorctl start ` + process); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// StopProcess 停止进程
func (r *SupervisorController) StopProcess(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	if out, err := tools.Exec(`supervisorctl stop ` + process); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// RestartProcess 重启进程
func (r *SupervisorController) RestartProcess(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	if out, err := tools.Exec(`supervisorctl restart ` + process); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// ProcessLog 进程日志
func (r *SupervisorController) ProcessLog(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	var logPath string
	var err error
	if tools.IsRHEL() {
		logPath, err = tools.Exec(`cat '/etc/supervisord.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
	} else {
		logPath, err = tools.Exec(`cat '/etc/supervisor/conf.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
	}

	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "无法从进程"+process+"的配置文件中获取日志路径")
	}

	log, err := tools.Exec(`tail -n 200 ` + logPath)
	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, log)
	}

	return controllers.Success(ctx, log)
}

// ClearProcessLog 清空进程日志
func (r *SupervisorController) ClearProcessLog(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	var logPath string
	var err error
	if tools.IsRHEL() {
		logPath, err = tools.Exec(`cat '/etc/supervisord.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
	} else {
		logPath, err = tools.Exec(`cat '/etc/supervisor/conf.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
	}

	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "无法从进程"+process+"的配置文件中获取日志路径")
	}

	if out, err := tools.Exec(`echo "" > ` + logPath); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	return controllers.Success(ctx, nil)
}

// ProcessConfig 获取进程配置
func (r *SupervisorController) ProcessConfig(ctx http.Context) http.Response {
	process := ctx.Request().Query("process")
	var config string
	var err error
	if tools.IsRHEL() {
		config, err = tools.Read(`/etc/supervisord.d/` + process + `.conf`)
	} else {
		config, err = tools.Read(`/etc/supervisor/conf.d/` + process + `.conf`)
	}

	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}

	return controllers.Success(ctx, config)
}

// SaveProcessConfig 保存进程配置
func (r *SupervisorController) SaveProcessConfig(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	config := ctx.Request().Input("config")
	var err error
	if tools.IsRHEL() {
		err = tools.Write(`/etc/supervisord.d/`+process+`.conf`, config, 0644)
	} else {
		err = tools.Write(`/etc/supervisor/conf.d/`+process+`.conf`, config, 0644)
	}

	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}

	_, _ = tools.Exec(`supervisorctl reread`)
	_, _ = tools.Exec(`supervisorctl update`)
	_, _ = tools.Exec(`supervisorctl restart ` + process)

	return controllers.Success(ctx, nil)
}

// AddProcess 添加进程
func (r *SupervisorController) AddProcess(ctx http.Context) http.Response {
	if sanitize := controllers.Sanitize(ctx, map[string]string{
		"name":    "required|alpha_dash",
		"user":    "required|alpha_dash",
		"path":    "required",
		"command": "required",
		"num":     "required",
	}); sanitize != nil {
		return sanitize
	}

	name := ctx.Request().Input("name")
	user := ctx.Request().Input("user")
	path := ctx.Request().Input("path")
	command := ctx.Request().Input("command")
	num := ctx.Request().InputInt("num", 1)
	config := `[program:` + name + `]
command=` + command + `
process_name=%(program_name)s
directory=` + path + `
autostart=true
autorestart=true
user=` + user + `
numprocs=` + strconv.Itoa(num) + `
redirect_stderr=true
stdout_logfile=/var/log/supervisor/` + name + `.log
stdout_logfile_maxbytes=2MB
`

	var err error
	if tools.IsRHEL() {
		err = tools.Write(`/etc/supervisord.d/`+name+`.conf`, config, 0644)
	} else {
		err = tools.Write(`/etc/supervisor/conf.d/`+name+`.conf`, config, 0644)
	}

	if err != nil {
		return controllers.Error(ctx, http.StatusUnprocessableEntity, err.Error())
	}

	_, _ = tools.Exec(`supervisorctl reread`)
	_, _ = tools.Exec(`supervisorctl update`)
	_, _ = tools.Exec(`supervisorctl start ` + name)

	return controllers.Success(ctx, nil)
}

// DeleteProcess 删除进程
func (r *SupervisorController) DeleteProcess(ctx http.Context) http.Response {
	process := ctx.Request().Input("process")
	if out, err := tools.Exec(`supervisorctl stop ` + process); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, out)
	}

	var logPath string
	var err error
	if tools.IsRHEL() {
		logPath, err = tools.Exec(`cat '/etc/supervisord.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
		if err := tools.Remove(`/etc/supervisord.d/` + process + `.conf`); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
		}
	} else {
		logPath, err = tools.Exec(`cat '/etc/supervisor/conf.d/` + process + `.conf' | grep stdout_logfile= | awk -F "=" '{print $2}'`)
		if err := tools.Remove(`/etc/supervisor/conf.d/` + process + `.conf`); err != nil {
			return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
		}
	}

	if err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, "无法从进程"+process+"的配置文件中获取日志路径")
	}

	if err := tools.Remove(logPath); err != nil {
		return controllers.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	_, _ = tools.Exec(`supervisorctl reread`)
	_, _ = tools.Exec(`supervisorctl update`)

	return controllers.Success(ctx, nil)
}
