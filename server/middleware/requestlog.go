/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package middleware

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "log"
    "strings"

    "github.com/gin-gonic/gin"
)

func RequestLogMiddleware() gin.HandlerFunc {
    return func(context *gin.Context) {

        var requestBody []byte
        if context.Request.Body != nil {
            requestBody, _ = ioutil.ReadAll(context.Request.Body)
        }

        contentType := context.GetHeader("Content-Type")
        contentType = strings.ToLower(contentType)

        buffer := bytes.NewBuffer(nil)
        json.Indent(buffer, requestBody, "", "    ")

        if strings.Contains(contentType, "application/json") {
            log.Print("request:\n", buffer.String())
        }

        context.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))
        context.Next()
    }
}
