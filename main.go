package main

import (
	"encoding/json"
	"github.com/cy18cn/micro-svc-common/handlers"
	"github.com/cy18cn/zlog"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Start() {
	err := zlog.InitZapLogger()
	if err != nil {
		return
	}
	log := zlog.GetLogger()
	defer log.Sync()
	r := httprouter.New()

	r.GET("/test", handlers.NewHandlers(test))
	r.POST("/testpost", handlers.NewHandlers(test))

	http.ListenAndServe(":8080", handlers.HttpRouterMiddleware(log, r))
}

type response struct {
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

func test(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	res := make(map[string]interface{})
	res["name"] = "test"
	res["password"] = "123465"
	resp := response{
		Msg:    "success",
		Result: res,
	}
	b, err := json.Marshal(&resp)
	if err != nil {
		panic(err)
	}
	w.Write(b)
}

func main() {
	Start()
}
