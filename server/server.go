package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"github.com/yixy/gateway/server/handler"

	"github.com/yixy/gateway/cfg"
	"github.com/yixy/gateway/log"
	"go.uber.org/zap"
)

var errCh chan error

func Start() error {
	logFile := filepath.Join(cfg.Dir, viper.GetString(cfg.LOG_FILE))
	log.InitLogger(logFile)
	log.Logger.Info("gateway init...")
	var errServerCh chan error = make(chan error)
	var errShutCh chan error = make(chan error)
	defer close(errServerCh)
	defer close(errShutCh)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.NotFoundHandler)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      mux,
		ReadTimeout:  time.Second * time.Duration(cfg.Rtimeout),
		WriteTimeout: time.Second * time.Duration(cfg.Wtimeout),
	}

	// Listen stop signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGKILL)
	go func() {
		log.Logger.Info("ready to Listen stop signal")
		sig := <-ch
		log.Logger.Info("receive signal", zap.Any("signal", sig))
		// stop
		signal.Stop(ch)
		// timeout context for shutdown
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(cfg.ShutTimeout)*time.Second)
		errShutCh <- srv.Shutdown(ctx)
	}()

	// open pid file
	lock := path.Join(cfg.Dir, "pid")
	log.Logger.Info("start init pid file", zap.String("pidFile", lock))
	lockFile, err := os.OpenFile(lock, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Logger.Error("create lockFile err.", zap.Error(err))
		return err
	}
	defer lockFile.Close()

	// try to lock pid file
	log.Logger.Info("try to lock pid file")
	err = syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		log.Logger.Error("syscall.Flock err.", zap.Error(err))
		return err
	}
	defer os.Remove(lock)
	defer syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)

	//record pid
	pid := strconv.Itoa(os.Getpid())
	log.Logger.Info("record pid", zap.String("pid", pid))
	lockFile.WriteString(pid)

	//start httpserver
	go func() {
		log.Logger.Info("start http server ...")
		errServerCh <- srv.ListenAndServe()
	}()

	//handle http error
	log.Logger.Info("handle http error")
	select {
	case e := <-errShutCh:
		if e != nil {
			log.Logger.Error("http server stop error.", zap.Error(e))
		} else {
			log.Logger.Info("http server stop")
		}
		return e
	case e := <-errServerCh:
		if e == http.ErrServerClosed {
			log.Logger.Info("close srv.ListenAndServe the service", zap.Error(e))
			//waiting for srv.shutdown
			time.Sleep(time.Duration(cfg.ShutTimeout) * time.Second)
		} else if e == nil {
			log.Logger.Info("srv.ListenAndServe return ok")
		} else {
			log.Logger.Error("srv.ListenAndServe error.", zap.Error(e))
		}
		return e
	}
}

func Stop() error {
	pidFile := path.Join(cfg.Dir, "pid")
	pidStr, err := ioutil.ReadFile(pidFile)
	if err != nil {
		fmt.Println("open pid file error", err)
		return err
	}
	pid, err := strconv.Atoi(string(pidStr))
	if err != nil {
		fmt.Println("parse pid error.", err)
		return err
	}
	err = signalOperation(pid, "SIGINT")
	if err != nil {
		fmt.Println("kill pid error", err)
		return err
	}
	return nil
}

//kill信号操作
func signalOperation(processId int, sig string) error {
	args := []string{"-s", sig, strconv.Itoa(processId)}
	cmd := exec.Command("kill", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
