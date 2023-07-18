package core

import "time"

type Options struct {

	// 异步消息任务名
	AsyncQueueName string

	// 协程数
	PoolSize int32

	// 消费任务
	TaskHandler TaskHandler

	// 异常处理方法
	PanicHandler PanicHandler

	// todo 回收空闲协程等待时间
	ExpriyTime time.Time
}

type TaskHandler func(string) error

type PanicHandler func(string, []byte)

type Option func(opts *Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

func WithAsyncQueueName(asyncQueueName string) Option {
	return func(opts *Options) {
		opts.AsyncQueueName = asyncQueueName
	}
}

func WithPoolSize(poolSize int32) Option {
	return func(opts *Options) {
		opts.PoolSize = poolSize
	}
}

func WithTaskHandler(taskHandler TaskHandler) Option {
	return func(opts *Options) {
		opts.TaskHandler = taskHandler
	}
}

func WithPanicHandler(panicHandler func(string, []byte)) Option {
	return func(opts *Options) {
		opts.PanicHandler = panicHandler
	}
}

func WithExpriyTime(expriyTime time.Time) Option {
	return func(opts *Options) {
		opts.ExpriyTime = expriyTime
	}
}
