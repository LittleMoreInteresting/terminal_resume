package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	tea "github.com/charmbracelet/bubbletea"
	"terminal_resume/internal/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "23234"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	// 创建 Wish SSH 服务器
	server, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%s", host, port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			logging.Middleware(),
			activeterm.Middleware(),
			bubbletea.Middleware(teaHandler),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// 优雅关闭
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting SSH server on %s:%s", host, port)
	log.Printf("Connect with: ssh %s -p %s", host, port)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

// teaHandler 为每个 SSH 会话创建一个新的 Bubble Tea 程序
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// 创建 Bubble Tea 模型
	m := app.NewModel()

	return m, []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}
}
