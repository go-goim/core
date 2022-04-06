package log

import "github.com/go-kratos/kratos/v2/log"

// define log api and implement

func init() {
	// TODO: add any other log implement.
	log.SetLogger(log.DefaultLogger)
}
