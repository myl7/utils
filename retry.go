// Copyright (C) 2022 myl7
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"math"
	"net/http"
	"time"
)

type RetryOpt struct {
	Block  func(i int)
	GenArg func(i int) any
}

func WithRetry(f func(any) (any, error), arg any, opt *RetryOpt) any {
	for i := 0; ; i++ {
		a := arg
		if opt != nil && opt.GenArg != nil {
			a = opt.GenArg(i)
		}

		r, err := f(a)
		if err == nil {
			return r
		}

		if opt != nil && opt.Block != nil {
			opt.Block(i)
		}
	}
}

func (o *RetryOpt) BlockInterval(t time.Duration) *RetryOpt {
	o.Block = func(_ int) {
		time.Sleep(t)
	}
	return o
}

func (o *RetryOpt) BlockExpInterval(t time.Duration, max time.Duration) *RetryOpt {
	o.Block = func(i int) {
		nt := int64(math.Pow(2, float64(i))) * int64(t)
		if int64(max) > 0 && nt > int64(max) {
			nt = int64(max)
		}

		time.Sleep(time.Duration(nt))
	}
	return o
}

func HTTPGetWithRetry(url string, opt *RetryOpt) any {
	f := func(arg any) (any, error) {
		return http.Get(arg.(string))
	}
	return WithRetry(f, url, opt)
}

func HTTPPostWithRetry(url, contentType, body string, opt *RetryOpt) any {
	f := func(arg any) (any, error) {
		args := arg.([]string)
		return http.Post(args[0], args[1], bytes.NewReader([]byte(args[2])))
	}
	return WithRetry(f, []string{url, contentType, body}, opt)
}
