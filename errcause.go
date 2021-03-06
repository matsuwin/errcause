package errcause

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"reflect"
	"strings"
	"time"
)

/**
 * (recover) 错误恢复逻辑
 *
 * Not recover:
 *     触发 panic 后开始向上 (函数调用链) 传递错误，到达当前 goroutine 顶层时会退出整个进程！！！
 *
 * recover:
 *     触发 panic 后开始向上传递错误，遇见第一个 recover 后结束传递，达到恢复的效果。
 */

// TurnOff 关闭此功能
func TurnOff() {
	turnOff = true
}

// Cause 从 error 中获取包含堆栈记录的错误根本原因
func Cause(err error) string {
	message := fmt.Sprintf("panic: %+v", errors.Cause(err))
	if strings.Count(message, "runtime.goexit") == 0 {
		message = fmt.Sprintf("(Not github.com/pkg/errors.New) panic: %+v", errors.New(err.Error()))
	}
	return message
}

// Recover panic! Error recovery
func Recover() {
	if turnOff {
		return
	}
	if err := recover(); err != nil {
		RFC3339Nano := time.Now().Local().Format(time.RFC3339Nano)
		if witch {
			if reflect.TypeOf(err).Kind() == reflect.String {
				message := fmt.Sprintf("[  ERROR  ] %s -> %s\n", RFC3339Nano, err)
				save(RFC3339Nano, message)
			} else {
				message := fmt.Sprintf("[  ERROR  ] %s\n%s\n\n", RFC3339Nano, Cause(err.(error)))
				save(RFC3339Nano, message)
			}
			go func() {
				time.Sleep(time.Second)
				witch = true
			}()
		}
		witch = false
	}
}

func save(RFC3339Nano, message string) {
	fmt.Println(message)
	f, _ := os.OpenFile("panic."+RFC3339Nano[:10]+".log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	_, _ = f.WriteString(message)
	_ = f.Close()
}

var witch, turnOff = true, false
