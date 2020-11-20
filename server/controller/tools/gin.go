/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package controllerTools

import (
    "fmt"
    "net/http"
    "errors"

    "github.com/gin-gonic/gin"
)

type Response struct {
    Error       bool        `json:"error"`
    Message     string      `json:"message,omitempty"`
    Result      interface{} `json:"result,omitempty"`
}

func SendError(context *gin.Context, err error) {
    if err == nil {
        err = errors.New("undefined")
    }
    response := Response{
        Error: true,
        Message: fmt.Sprintf("%s", err),
        Result: nil,
    }
    context.IndentedJSON(http.StatusOK, response)
}

func SendOk(context *gin.Context) {
    response := Response{
        Error: false,
        Message: "",
        Result: nil,
    }
    context.IndentedJSON(http.StatusOK, response)
}

func SendMessage(context *gin.Context, message string) {
    response := Response{
        Error: false,
        Message: fmt.Sprintf("%s", message),
        Result: nil,
    }
    context.IndentedJSON(http.StatusOK, response)
}

func SendResult(context *gin.Context, result interface{}) {
    response := Response{
        Error: false,
        Result: result,
    }
    context.IndentedJSON(http.StatusOK, &response)
}

func AbortContext(context *gin.Context, code int, err error) {
    if err == nil {
        err = errors.New("undefined")
    }
    response := Response{
        Error: true,
        Message: fmt.Sprintf("%s", err),
        Result: nil,
    }
    context.IndentedJSON(code, response)
    context.Abort()
}
