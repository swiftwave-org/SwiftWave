package cronjob

import (
	"bytes"
	"strings"
	"sync"
	"time"

	"github.com/swiftwave-org/swiftwave/pkg/ssh_toolkit"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/logger"
)

func (m Manager) MonitorServerStatus() {
	logger.CronJobLogger.Println("Starting server status monitor [cronjob]")
	for {
		m.monitorServerStatus()
		time.Sleep(2 * time.Second)
	}
}

func (m Manager) monitorServerStatus() {
	logger.CronJobLogger.Println("Triggering Server Status Monitor Job")
	// Fetch all servers
	servers, err := core.FetchAllServers(&m.ServiceManager.DbClient)
	if err != nil {
		logger.CronJobLoggerError.Println("Failed to fetch server list")
		logger.CronJobLoggerError.Println(err)
		return
	}
	if len(servers) == 0 {
		logger.CronJobLogger.Println("Skipping ! No server found")
		return
	}

	var wg sync.WaitGroup
	for _, server := range servers {
		if server.Status == core.ServerNeedsSetup || server.Status == core.ServerPreparing {
			continue
		}
		wg.Add(1)
		go func(server core.Server) {
			defer wg.Done()
			m.checkAndUpdateServerStatus(server)
		}(server)
	}
	wg.Wait()
}

func (m Manager) checkAndUpdateServerStatus(server core.Server) {
	if m.isServerOnline(server) {
		if server.Status != core.ServerOnline {
			err := core.MarkServerAsOnline(&m.ServiceManager.DbClient, &server)
			if err != nil {
				logger.CronJobLoggerError.Println("DB Error : Failed to mark server as online >", server.HostName, err)
			} else {
				logger.CronJobLogger.Println("Server marked as online >", server.HostName)
			}
		}
	} else {
		if server.Status != core.ServerOffline {
			err := core.MarkServerAsOffline(&m.ServiceManager.DbClient, &server)
			if err != nil {
				logger.CronJobLoggerError.Println("DB Error : Failed to mark server as offline >", server.HostName, err)
			} else {
				logger.CronJobLogger.Println("Server marked as offline >", server.HostName)
			}
		} else {
			logger.CronJobLogger.Println("Server already offline >", server.HostName)
		}
	}
}

func (m Manager) isServerOnline(server core.Server) bool {
	retries := 3 // try for 3 times before giving up
	if server.Status == core.ServerOffline {
		/**
		* If server is offline, try only once
		* Else, it will take total 30 seconds (3 retries * 10 seconds of default SSH timeout)
		 */
		retries = 1
	}
	// try for 3 times
	for i := 0; i < retries; i++ {
		cmd := "echo ok"
		stdoutBuf := new(bytes.Buffer)
		stderrBuf := new(bytes.Buffer)
		err := ssh_toolkit.ExecCommandOverSSHWithOptions(cmd, stdoutBuf, stderrBuf, 3, server.IP, server.SSHPort, server.User, m.Config.SystemConfig.SshPrivateKey, false)
		if err != nil {
			logger.CronJobLoggerError.Println("Error while checking if server is online", server.HostName, err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		if strings.Compare(strings.TrimSpace(stdoutBuf.String()), "ok") == 0 {
			return true
		}
	}
	return false
}
