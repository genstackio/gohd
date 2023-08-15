package gohd

type CreateErrorUrlFn = func(code int, err error, country string) (string, error)
