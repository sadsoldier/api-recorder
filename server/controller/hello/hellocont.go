/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package helloController

import (
    "github.com/gin-gonic/gin"

    "dummyapi/server/config"
    "dummyapi/server/controller/tools"
    "dummyapi/model/hello"
)

type Controller struct{
    config  *appConfig.Config
}

func (this *Controller) Hello(context *gin.Context) {
    model := helloModel.New()

    result, err := model.Hello()
    if err != nil {
        controllerTools.SendError(context, err)
    }
    controllerTools.SendResult(context, &result)
}

func New(config *appConfig.Config) *Controller {
    return &Controller{
        config: config,
    }
}
