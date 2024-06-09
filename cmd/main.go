package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/csh0101/netagent.git/internal/agent"
	"github.com/csh0101/netagent.git/internal/controller"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {

	root := &cobra.Command{
		Use: "netgent",
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(1)
		},
	}

	root.AddCommand(ControllerCmd())
	agentCmd, err := AgentCmd()
	if err != nil {
		panic(err)
	}
	root.AddCommand(agentCmd)

	return root
}

var (
	dport     string
	cport     string
	agentName string
)

func AgentCmd() (*cobra.Command, error) {
	agentCmd := &cobra.Command{
		Use: "agent",
		Run: func(cmd *cobra.Command, args []string) {
			if err := agentCmd(); err != nil {
				panic(err)
			}
		},
	}
	agentCmd.Flags().StringVar(&dport, "dport", "9999", "data tunnel address")
	agentCmd.Flags().StringVar(&cport, "cport", "8888", "control tunnel address")
	agentCmd.Flags().StringVar(&agentName, "agentName", "agent-dev0", "netagent")
	if err := agentCmd.Flags().Parse(nil); err != nil {
		return nil, err
	}
	return agentCmd, nil
}

func agentCmd() error {

	a := &agent.Agent{}

	dataPort, err := strconv.Atoi(dport)
	if err != nil {
		return err
	}

	controlPort, err := strconv.Atoi(cport)
	if err != nil {
		return err
	}

	if err := a.Run(&agent.Config{
		Name:                 "agent-dev1",
		DataTunnelAddress:    "127.0.0.1",
		DataTunnelPort:       dataPort,
		ControlTunnelAddress: "127.0.0.1",
		ControlTunnelPort:    controlPort,
		MaxRetries:           5,
	}); err != nil {
		return err
	}
	// 创建一个通道用于接收信号
	sigs := make(chan os.Signal, 1)
	// 监听指定的信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// 创建一个通道用于等待程序退出
	done := make(chan bool, 1)
	// 启动一个 goroutine 以异步接收信号
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		// 在这里添加任何清理操作，例如关闭文件、释放资源等
		done <- true
	}()
	fmt.Println("Waiting for signal")
	<-done
	fmt.Println("Exiting")

	return nil
}

func ControllerCmd() *cobra.Command {
	control := &cobra.Command{
		Use: "control",
		Run: func(cmd *cobra.Command, args []string) {
			runController()
		},
	}
	return control
}

func runController() error {
	if err := new(controller.Controller).Run(&controller.Config{
		ControlPort: 8888,
		Cidr:        "172.2.1.0/24",
		DataPort:    9999,
	}); err != nil {
		return err
	}
	// 创建一个通道用于接收信号
	sigs := make(chan os.Signal, 1)
	// 监听指定的信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// 创建一个通道用于等待程序退出
	done := make(chan bool, 1)
	// 启动一个 goroutine 以异步接收信号
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		// 在这里添加任何清理操作，例如关闭文件、释放资源等
		done <- true
	}()
	fmt.Println("Waiting for signal")
	<-done
	fmt.Println("Exiting")
	return nil
}
