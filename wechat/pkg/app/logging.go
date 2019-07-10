package app

import (
	"git.sdkeji.top/share/sdlib/log"
)

func NewNSQLogger(logger log.Logger) NSQLogger {
	return NSQLogger{logger}
}

type NSQLogger struct {
	logger log.Logger
}

func (logger NSQLogger) Output(calldepth int, s string) error {
	logger.logger.Info(s, "service", "nsq_logger")
	return nil
}
