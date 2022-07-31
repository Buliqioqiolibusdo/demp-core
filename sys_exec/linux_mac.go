//go:build !windows
// +build !windows

package sys_exec

import (
	"os/exec"
	"syscall"
	"time"

	"github.com/crawlab-team/go-trace"
	"github.com/shirou/gopsutil/process"
)

type KillProcessOptions struct {
	Timeout time.Duration
	Force   bool
}

func BuildCmd(cmdStr string) *exec.Cmd {
	return exec.Command("sh", "-c", cmdStr)
}

func Setpgid(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	} else {
		cmd.SysProcAttr.Setpgid = true
	}
}

// func KillProcess(cmd *exec.Cmd) error {
// 	if cmd == nil || cmd.Process == nil {
// 		return nil
// 	}
// 	return syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
// }

func KillProcess(cmd *exec.Cmd, opts *KillProcessOptions) error {
	// process
	p, err := process.NewProcess(int32(cmd.Process.Pid))
	if err != nil {
		return err
	}

	// kill function
	killFunc := func(p *process.Process) error {
		return killProcessRecursive(p, opts.Force)
	}

	if opts.Timeout != 0 {
		// with timeout
		return killProcessWithTimeout(p, opts.Timeout, killFunc)
	} else {
		// without timeout
		return killFunc(p)
	}
}

func killProcessRecursive(p *process.Process, force bool) (err error) {
	// children processes
	cps, err := p.Children()
	if err != nil {
		return killProcess(p)
	}

	// iterate children processes
	for _, cp := range cps {
		if err := killProcessRecursive(cp, force); err != nil {
			return err
		}
	}

	return nil
}

func killProcessWithTimeout(p *process.Process, timeout time.Duration, killFunc func(*process.Process) error) error {
	go func() {
		if err := killFunc(p); err != nil {
			trace.PrintError(err)
		}
	}()
	for i := 0; i < int(timeout.Seconds()); i++ {
		ok, err := process.PidExists(p.Pid)
		if err == nil && !ok {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return forceKillProcess(p)
}

// func KillProcessWithTimeout(cmd *exec.Cmd, timeout time.Duration) error {
// 	if cmd == nil || cmd.Process == nil {
// 		return nil
// 	}
// 	go func() {
// 		if err := syscall.Kill(cmd.Process.Pid, syscall.SIGTERM); err != nil {
// 			trace.PrintError(err)
// 		}
// 	}()
// 	for i := 0; i < int(timeout.Seconds()); i++ {
// 		ok, err := process.PidExists(int32(cmd.Process.Pid))
// 		if err == nil && !ok {
// 			return nil
// 		}
// 		time.Sleep(1 * time.Second)
// 	}
// 	return ForceKillProcess(cmd)
// }

// func ForceKillProcess(cmd *exec.Cmd) error {
// 	return syscall.Kill(cmd.Process.Pid, syscall.SIGKILL)
// }

func forceKillProcess(p *process.Process) (err error) {
	if err := p.Kill(); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func killProcess(p *process.Process) (err error) {
	if err := p.Terminate(); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

// func ConfigureCmdLogging(cmd *exec.Cmd, fn func(scanner *bufio.Scanner)) {
// 	stdout, _ := (*cmd).StdoutPipe()
// 	stderr, _ := (*cmd).StderrPipe()
// 	scannerStdout := bufio.NewScanner(stdout)
// 	scannerStderr := bufio.NewScanner(stderr)
// 	go fn(scannerStdout)
// 	go fn(scannerStderr)
// }
