/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package daemon

import (
    "log"
    "os"
    "io"
    "syscall"
    "strconv"
    "os/signal"
    "time"
    "fmt"
    //"os/user"
    "path/filepath"
    "errors"
    "dummyapi/server/config"
)

type Daemon struct {
    Config *appConfig.Config
}

func (this *Daemon) Daemonize() error {
    var err error

    /* Daemonize process */
    if !this.Config.Foreground {
        this.ForkProcess()
    }

    /* Lookup user system info */
    //user, err := user.Lookup(this.Config.User)
    //if err != nil {
    //    return errors.New(fmt.Sprintf("user lookup error: %s\n", err))
    //}

    /* Make process ID directory */
    err = os.MkdirAll(filepath.Dir(this.Config.PidPath), 0750)
    if err != nil {
        return errors.New(fmt.Sprintf("unable create rundir: %s\n", err))
    }

    /* Save process ID to file */
    if err := this.saveProcessID(this.Config.PidPath); err != nil {
        return errors.New(fmt.Sprintf("unable save process id: %s\n", err))
    }
    defer os.Remove(this.Config.PidPath)

    //uid, err := strconv.Atoi(user.Uid)

    /* Make log directory */
    err = os.MkdirAll(filepath.Dir(this.Config.MessageLogPath), 0750)
    if err != nil {
        return errors.New(fmt.Sprintf("unable create message log dir: %s\n", err))
    }
    //err = os.Chown(filepath.Dir(this.Config.MessageLogPath), uid, os.Getgid())
    //if err != nil {
    //    return errors.New(fmt.Sprintf("unable chown log dir: %s\n", err))
    //}

    /* Make store directory */
    err = os.MkdirAll(this.Config.StoreDir, 0750)
    if err != nil {
        return errors.New(fmt.Sprintf("unable create store dir: %s\n", err))
    }

    //err = os.Chown(this.Config.StoreDir, uid, os.Getgid())
    //if err != nil {
    //    return errors.New(fmt.Sprintf("unable chown store dir: %s\n", err))
    //}

    if _, err := os.Stat(this.Config.StoreDir); err != nil {
        return errors.New(fmt.Sprintf("store dir not exists: %s\n", err))
    }

    /* Change effective user ID */
    //if uid != 0 {
    //    err = syscall.Setuid(uid)
    //    if err != nil {
    //        return errors.New(fmt.Sprintf("set process user id error: %s\n", err))
    //    }
    //    if syscall.Getuid() != uid {
    //        return errors.New(fmt.Sprintf("set process user id error: %s\n", err))
    //    }
    //}

    /* Redirect log to message file */
    _, err = this.redirectLog(this.Config.MessageLogPath, this.Config.Debug)
    if err != nil {
        return errors.New(fmt.Sprintf("unable redirect log to message file: %s\n", err))
    }
    //defer file.Close()

    /* Redirect standard IO */
    if !this.Config.Foreground {
        if _, err := this.redirectIO(); err != nil {
            return errors.New(fmt.Sprintf("unable redirect stdio: %s\n", err))
        }
    }
    return nil
}

func (this *Daemon) saveProcessID(filename string) error {
    pid := os.Getpid()
    file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0640)
    if err != nil {
         return err
    }
    defer file.Close()
    if _, err := file.WriteString(strconv.Itoa(pid)); err != nil {
        return err
    }
    file.Sync()
    return nil
}

func (this *Daemon) redirectLog(filename string, debug bool) (*os.File, error) {
    file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
    if err != nil {
        return nil, err
    }
    writer := io.MultiWriter(os.Stdout, file)
    //writer := io.Writer(file)
    if debug {
        log.SetFlags(log.LstdFlags | log.Lshortfile)
    } else {
        log.SetFlags(log.LstdFlags)
    }
    log.SetOutput(writer)
    return file, nil
}

func (this *Daemon) ForkProcess() error {

    if _, exists := os.LookupEnv("GOGOFORK"); !exists {
        os.Setenv("GOGOFORK", "yes")

        cwd, err := os.Getwd()
        if err != nil {
            return err
        }

        procAttr := syscall.ProcAttr{}
        procAttr.Files = []uintptr{ uintptr(syscall.Stdin), uintptr(syscall.Stdout), uintptr(syscall.Stderr) }
        procAttr.Env = os.Environ()
        procAttr.Dir = cwd
        syscall.ForkExec(os.Args[0], os.Args, &procAttr)
        os.Exit(0)
    }
    _, err := syscall.Setsid()
    if err != nil {
        return err
    }
    os.Chdir("/")
    return nil
}

func (this *Daemon) redirectIO() (*os.File, error) {
    file, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
    if err != nil {
        return nil, err
    }
    syscall.Dup2(int(file.Fd()), int(os.Stdin.Fd()))
    syscall.Dup2(int(file.Fd()), int(os.Stdout.Fd()))
    syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
    return file, nil
}

func (this *Daemon) SetSignalHandler() {

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

    go func() {
        for {
            log.Printf("signal handler start")
            sig := <- sigs
            log.Printf("receive signal %s", sig.String())

            switch sig {
                case syscall.SIGINT, syscall.SIGTERM:
                    log.Printf("exit process by signal %s", sig.String())
                    time.Sleep(time.Millisecond * 100)
                    os.Exit(0)

                case syscall.SIGHUP:
                    log.Printf("restart program")
                    this.ForkProcess()
            }
        }
    }()
}

func New(config *appConfig.Config) *Daemon {
    return &Daemon{
        Config: config,
    }
}
