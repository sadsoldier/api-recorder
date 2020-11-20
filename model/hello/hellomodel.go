/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */


package helloModel

type Model struct {
}

type Hello struct {
    Message string     `json:"message"`
}

func (this *Model) Hello() (*Hello, error) {

    hello := Hello{
        Message: "hello",
    }
    return &hello, nil
}

func New() *Model {
    model := Model{
    }
    return &model
}
