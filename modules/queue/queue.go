package queue

type Type string

type Data interface{}

type HandleFunc func(data ...Data)
