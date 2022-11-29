package lrucache

import "errors"

var ErrNotFound = errors.New("value not found")

var ErrQueueEmpty = errors.New("queue is empty")
