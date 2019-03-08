package handlers

import (
	"github.com/cy18cn/zlog"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

//type errHandler struct {
//	logger *zap.Logger
//	handler http.Handler
//}

//func (self *errHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//	// 如果 发现 panic，判断错误输出错误，否则 继续往上层 panic
//	defer func() {
//		if r := recover(); r != nil {
//			self.logger.Error("failed to handle URL",
//				zap.String("url", req.RequestURI),
//				zap.String())
//			log.Printf("Panic: %v", r)
//			http.Error(
//				w, //向writer汇报错误
//				http.StatusText(http.StatusInternalServerError), //错误描述信息（字符串）
//				http.StatusInternalServerError) //系统内部错误
//		}
//	}()
//	// 执行业务代码操作，上面定义的 defer 就是防止 业务代码中出现 panic
//	self.handler.ServeHTTP(w, req)
//
//	// 如果业务代码执行出错
//	//if err != nil {
//	//	//日志输出错误信息
//	//	log.Printf("Error occurred handling request: %s",err.Error())
//	//
//	//	//判断错误类型是否为 自定义错误
//	//	if userErr, ok := err.(userError); ok {
//	//		http.Error(writer,
//	//			userErr.Message(),
//	//			http.StatusBadRequest)
//	//		return
//	//	}
//	//
//	//	// 判断系统错误的类型
//	//	code := http.StatusOK
//	//	switch {
//	//	case os.IsNotExist(err):
//	//		code = http.StatusNotFound  //文件无法找到错误
//	//	case os.IsPermission(err):
//	//		code = http.StatusForbidden // 权限不够错误
//	//	default:
//	//		code = http.StatusInternalServerError //其他错误
//	//	}
//	//	//向writer 中写入错误信息
//	//	http.Error(writer,http.StatusText(code), code)
//	//}
//}

func errorHandler(next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		defer func() {
			if r := recover(); r != nil {
				zlog.Errorf("failed to handle URL: %s, method: %s, params: %v, err: %v",
					request.RequestURI,
					request.Method,
					request.Form,
					r)
				http.Error(writer,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()
		next(writer, request, params)
	}
}
