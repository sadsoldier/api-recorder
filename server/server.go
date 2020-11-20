/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package server

import (
    "errors"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"

    "github.com/gin-gonic/gin"

    "dummyapi/server/config"
    "dummyapi/server/daemon"

    "dummyapi/server/controller/tools"
    //"dummyapi/server/controller/hello"
    "dummyapi/server/controller/dummy"
    "dummyapi/server/middleware"
)

type Server struct {
    Config      *appConfig.Config
}

func (this *Server) Run() error {
    var err error

    /* Daemonize process */
    daemon := daemon.New(this.Config)
    err = daemon.Daemonize()
    if err != nil {
        return err
    }
    /* Set signal handlers */
    daemon.SetSignalHandler()

    /* Setup gin */
    this.setupGin()

    /* Create and setup router */
    router := gin.New()

    if this.Config.Debug {
        router.Use(middleware.RequestLogMiddleware())
        router.Use(middleware.ResponseLogMiddleware())
    }

    router.Use(gin.LoggerWithFormatter(logFormatter()))
    router.Use(gin.Recovery())

    /* Set route handlers */
    //helloControllerIm := helloController.New(this.Config)
    dummyControllerIm := dummyController.New(this.Config)

    //routerGroup := router.Group("/api/v1")
    //routerGroup.GET("/hello", helloControllerIm.Hello)
    //routerGroup.POST("/hello", dummyControllerIm.Writer)

    /* Set noroute handler */
    //router.NoRoute(this.NoRoute)
    router.NoRoute(dummyControllerIm.Writer)

    /* Start run loop */
    log.Printf("start listen on :%d", this.Config.Port)

    return router.Run(":" + fmt.Sprintf("%d", this.Config.Port))

}

func (this *Server) NoRoute(context *gin.Context) {
    requestPath := context.Request.URL.Path
    controllerTools.SendError(context, errors.New(fmt.Sprintf("wrong path %s", requestPath)))
    return
}

func (this *Server) Configure() {

    /* Read configuration file */
    this.Config = appConfig.New()
    this.Config.Read()
    //this.Config.Write()

    /* Parse command line options */
    optForeground := flag.Bool("foreground", false, "run in foreground")
    flag.BoolVar(optForeground, "f", false, "run in foreground")

    optPort := flag.Int("port", this.Config.Port, "listen port")
    flag.IntVar(optPort, "p", this.Config.Port, "listen port")

    optDebug := flag.Bool("debug", this.Config.Debug, "debug mode")
    flag.BoolVar(optDebug, "d", false, "debug mode")

    optDevel := flag.Bool("devel", this.Config.Devel, "devel mode")
    flag.BoolVar(optDebug, "e", false, "devel mode")

    optWrite := flag.Bool("write", false, "write config")
    flag.BoolVar(optWrite, "w", false, "write config")

    exeName := filepath.Base(os.Args[0])

    flag.Usage = func() {
        fmt.Println("")
        fmt.Printf("usage: %s command [option]\n", exeName)
        fmt.Println("")
        flag.PrintDefaults()
        fmt.Println("")
    }
    flag.Parse()

    this.Config.Port = *optPort
    this.Config.Debug = *optDebug
    this.Config.Devel = *optDevel
    this.Config.Foreground = *optForeground
}

func (this *Server) setupGin() error {
    gin.DisableConsoleColor()
    if this.Config.Debug{
        gin.SetMode(gin.DebugMode)
    } else {
        gin.SetMode(gin.ReleaseMode)
    }
    accessLogFile, err := os.OpenFile(this.Config.AccessLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
    if err != nil {
      return err
    }
    gin.DefaultWriter = io.MultiWriter(accessLogFile, os.Stdout)
    //gin.DefaultWriter = ioutil.Discard
    return nil
}

func New() *Server {
    return &Server{
    }
}

func logFormatter() func(param gin.LogFormatterParams) string {
    return func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s %s %s %s %s %d %d %s\n",
            param.TimeStamp.Format(time.RFC3339),
            param.ClientIP,
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.BodySize,
            param.Latency,
        )
    }
}
