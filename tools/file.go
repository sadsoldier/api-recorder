/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package tools

import (
    "os"
    "fmt"
    "path/filepath"
    "strings"
    "io/ioutil"
    "errors"
)

/* Return true if file exists and not directory */
func FileExists(name string) bool {
    fi, err := os.Stat(name)
    if err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return !fi.IsDir()
}


func PathLength(path string) int {
    path = filepath.Clean(path)
    if len(path) == 0 {
        return 0
    }
    return len(strings.Split(path, "/"))
}

/* Return list of directories in tree */
func PathWalkDir(basePath string, depth int) ([]string, error) {
    depth = depth + PathLength(basePath)
    var list []string
    err := filepath.Walk(basePath,
        func(filePath string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if PathLength(filePath) > depth  {
                return filepath.SkipDir
            }
            if info.IsDir() {
                list = append(list, filePath)
            }
            return nil
        })
    if err != nil {
        return list, err
    }
    return list, err
}

/* Return total size of data in file tree */
func DirSize(path string) (int64, error) {
    var size int64
    err := filepath.Walk(path,
        func(_ string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if !info.IsDir() {
                size += info.Size()
            }
            return err
        })
    return size, err
}

/* Calculate total file size in the directory */
func BucketSize(directoryPath string) (int64, error) {
    directoryPath = filepath.Clean(directoryPath)
    fi, err := os.Stat(directoryPath)
    if err != nil {
        return 0, err
    }
    if !fi.IsDir() {
        return 0, errors.New(fmt.Sprintf("%s is not dir", directoryPath))
    }
    files, err := ioutil.ReadDir(directoryPath)
    if err != nil {
        return 0, err
    }
    var size int64
    for _, file := range files {
        if !file.IsDir() {
            size = size + file.Size()
        }
    }
    return size, nil
}

/* Calculate total file size in the directory
 * with filename patttern */
func BucketSizeGlob(directoryPath string, pattern string) (int64, error) {
    directoryPath = filepath.Clean(directoryPath)

    fileInfo, err := os.Stat(directoryPath)
    if err != nil {
        return 0, err
    }
    if !fileInfo.IsDir() {
        return 0, errors.New(fmt.Sprintf("%s is not dir", directoryPath))
    }

    files, err := filepath.Glob(filepath.Join(directoryPath, pattern))
    if err != nil {
        return 0, err
    }
    var size int64
    for _, file := range files {
        fileInfo, err := os.Stat(file)
        if !fileInfo.IsDir() && err == nil {
            size = size + fileInfo.Size()
        }
    }
    return size, nil
}
