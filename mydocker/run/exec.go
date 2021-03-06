package run

import (
	// 这个很关键，引入而不使用，但其在启动的时候后自动调用
	_ "dockerDemo/mydocker/nsenter"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

// EnvExecPid
/**
前面的C代码中已经出现了mydocker_pid 和 mydocker_cmd 这两个key
主要是为了控制是否执行c代码里面的setns
*/
const EnvExecPid = "mydocker_pid"
const EnvExecCmd = "mydocker_cmd"

func ExecContainer(containerName string, commandArray []string) error {
	// 根据传过来的容器名获取宿主机对应的pid
	pid, err := getContainerPidByName(containerName)
	if err != nil {
		return err
	}

	// 把命令以空格为分隔符拼接成一个字符串，便于传递
	cmdStr := strings.Join(commandArray, " ")
	log.Infof("container pid %s", pid)
	log.Infof("command %s", cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := os.Setenv(EnvExecPid, pid); err != nil {
		return fmt.Errorf("setenv %s err: %v", EnvExecPid, err)
	}
	if err := os.Setenv(EnvExecCmd, cmdStr); err != nil {
		return fmt.Errorf("setenv %s err: %v", EnvExecCmd, err)
	}

	envs, err := getEnvsByPid(pid)
	if err != nil {
		return err
	}
	cmd.Env = append(os.Environ(), envs...)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("exec container %s err: %v", containerName, err)
	}
	return nil
}
