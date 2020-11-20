/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package dummyController

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "strings"
    "fmt"
    "os"
    "path"
    "time"
    "net/http"
    "log"

    "github.com/gin-gonic/gin"

    "dummyapi/server/config"
)

type Controller struct{
    config  *appConfig.Config
}

func (this *Controller) Writer(context *gin.Context) {

    // Prepare palce for writing
    requestPath := context.Request.URL.Path
    fullPath := path.Join(this.config.StoreDir, requestPath)

    os.MkdirAll(fullPath, 0770)

    // Create timestamp
    now := time.Now()
    timestamp := fmt.Sprintf("%04d-%02d-%02d-%02d%02d-%02d-%09d",
                                now.Year(),
                                now.Month(),
                                now.Day(),
                                now.Hour(),
                                now.Minute(),
                                now.Second(),
                                now.Nanosecond())

    // Get body
    var requestBody []byte
    if context.Request.Body != nil {
        requestBody, _ = ioutil.ReadAll(context.Request.Body)
    }

    // Write body as is
    dataPath := path.Join(fullPath, timestamp + ".raw")
    dataFile, err := os.OpenFile(dataPath, os.O_RDWR|os.O_CREATE, 0640)
    if err != nil {
        log.Println("unable to open file: ", dataPath, ", err: ", err)
    } else {
        defer dataFile.Close()
        _, err := dataFile.Write(requestBody)
        if err != nil {
            log.Println("unable to write to file: ", dataPath, ", err: ", err)
        }
        dataFile.Sync()
    }

    // Write formatted JSON body if possible
    contentType := context.GetHeader("Content-Type")

    if strings.Contains(strings.ToLower(contentType), "application/json") {
        buffer := bytes.NewBuffer(nil)
        json.Indent(buffer, requestBody, "", "    ")

        dataPath := path.Join(fullPath, timestamp + ".json")
        dataFile, err := os.OpenFile(dataPath, os.O_RDWR|os.O_CREATE, 0640)
        if err != nil {
            log.Println("unable to open file: ", dataPath, ", err: ", err)
        } else {
            defer dataFile.Close()
                _, err = dataFile.WriteString(buffer.String())
                if err != nil {
                    log.Println("unable to write to file: ", dataPath, ", err: ", err)
                }
                dataFile.Sync()
        }
    }

    // Write headers to file
    headPath := path.Join(fullPath, timestamp + ".hdr")
    headFile, err := os.OpenFile(headPath, os.O_RDWR|os.O_CREATE, 0640)
    if err != nil {
        log.Println("unable to open file: ", headPath, ", err: ", err)
    } else {
        defer headFile.Close()
        for name, values := range context.Request.Header {
            _, err := headFile.WriteString(name + ":")
            if err != nil {
                log.Println("unable to write to file: ", headPath, ", err: ", err)
            }
            for _, value := range values {
                _, err := headFile.WriteString(" " + value)
                if err != nil {
                    log.Println("unable to write to file: ", headPath, ", err: ", err)
                }
            }
            _, err = headFile.WriteString("\n")
            if err != nil {
                log.Println("unable to write to file: ", headPath, ", err: ", err)
            }
        }
        _, err = headFile.WriteString("\n\n")
        if err != nil {
            log.Println("unable to write to file: ", headPath, ", err: ", err)
        }
        headFile.Sync()
    }

    // Wite some additional "meta" data
    metaPath := path.Join(fullPath, timestamp + ".meta")
    metaFile, err := os.OpenFile(metaPath, os.O_RDWR|os.O_CREATE, 0640)
    if err != nil {
        log.Println("unable to open file: ", metaPath, ", err: ", err)
    } else {
        defer metaFile.Close()
        // Method
        _, err = metaFile.WriteString("Method: " + context.Request.Method + "\n")
        if err != nil {
            log.Println("unable to write to file: ", metaPath, ", err: ", err)
        }
        // Request-URI
        _, err = metaFile.WriteString("Request-URI: " + context.Request.RequestURI + "\n")
        if err != nil {
            log.Println("unable to write to file: ", metaPath, ", err: ", err)
        }
        // Content-Length
        _, err = metaFile.WriteString("Content-Length: " + fmt.Sprintf("%d", context.Request.ContentLength) + "\n")
        if err != nil {
            log.Println("unable to write to file: ", metaPath, ", err: ", err)
        }
        // Remote-Addr
        _, err = metaFile.WriteString("Remote-Addr: " + context.Request.RemoteAddr + "\n")
        if err != nil {
            log.Println("unable to write to file: ", metaPath, ", err: ", err)
        }
        // Host
        _, err = metaFile.WriteString("Host: " + context.Request.Host + "\n")
        if err != nil {
            log.Println("unable to write to file: ", metaPath, ", err: ", err)
        }
        metaFile.Sync()
    }

    reader := bytes.NewReader([]byte(requestBody))
    requestURL := fmt.Sprintf("%s%s", this.config.Target, context.Request.RequestURI)

    log.Println("request URL:", requestURL)

    request, err := http.NewRequest(context.Request.Method, requestURL, reader)
    // Remap header from one side to another
    for headerName, values := range context.Request.Header {
        headerValue := ""
        for _, value := range values {
            headerValue += value
        }
        request.Header.Set(headerName, headerValue)
    }

    client := &http.Client{}
    response, err := client.Do(request)

    responseBody, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Println("unable to read response body: ", requestURL, ", err: ", err)
        context.AbortWithStatus(http.StatusInternalServerError)
        return
    }

    responseContentType := response.Header.Get("Content-Type")
    context.Data(response.StatusCode, responseContentType, responseBody)
}

func New(config *appConfig.Config) *Controller {
    return &Controller{
        config: config,
    }
}
