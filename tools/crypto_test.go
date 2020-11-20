/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package tools

import (
    "testing"
)

func TestCreateHash(t *testing.T) {
    hash, err := CreateHash("123456")
    if err != nil {
        t.Error(err)
    }
    if len(hash) != 59 {
        t.Error(err)
    }
}


func TestCheckHash(t *testing.T) {
    hash := "$5$rfBd67ti3SMt$V4Gi.A19xLBWMSZDFw4JvfotiQdF.nDmMsFB3DyPvI0"
    password := "123456"

    err := CheckHash(hash, password)
    if err != nil {
        t.Error(err)
    }
}
