/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package middleware

import (
    "bytes"
    "encoding/json"
    "log"
    "strings"

    "github.com/gin-gonic/gin"
)


type LogWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

func (this LogWriter) Write(data []byte) (int, error) {
    this.body.Write(data)
    return this.ResponseWriter.Write(data)
}

func (this LogWriter) WriteString(data string) (int, error) {
    this.body.WriteString(data)
    return this.ResponseWriter.WriteString(data)
}

func ResponseLogMiddleware() gin.HandlerFunc {

    return func(context *gin.Context) {
        contentType := context.GetHeader("Content-Type")
        contentType = strings.ToLower(contentType)

        writer := &LogWriter{
            body: bytes.NewBuffer(nil),
            ResponseWriter: context.Writer,
        }
        context.Writer = writer

        context.Next()

        buffer := bytes.NewBuffer(nil)
        json.Indent(buffer, writer.body.Bytes(), "", "    ")

        if strings.Contains(contentType, "application/json") {
            log.Print("response:\n", buffer.String())
        }
    }
}
