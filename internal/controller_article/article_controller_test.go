package controller_article

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	testing.Init()
	// err := utils.LoadCfg("../../conf/message_compagin.yaml", &serverConfig)
	// if err != nil {
	// 	fmt.Printf("cfg error:%v\n", err)
	// 	return
	// }
	// if InitRes() != err {
	// 	fmt.Printf("InitRes error:%v\n", err)
	// 	return
	// }
}

func postJson(uri string, param map[string]interface{}, router *gin.Engine) *httptest.ResponseRecorder {
	jsonByte, _ := json.Marshal(param)
	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))
	req.Header.Set("Content-Type", "application/json")
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}
func getJson(uri string, router *gin.Engine) *httptest.ResponseRecorder {
	// 构造post请求，json数据以请求body的形式传递
	req := httptest.NewRequest("GET", uri, nil)
	req.Header.Set("Content-Type", "application/json")
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}

func TestInsert(t *testing.T) {
	gin := InitEngine()

	//w := performRequest(gin, "POST", "/campaign")
	w := getJson("/campaign", gin)
	assert.Equal(t, http.StatusOK, w.Code)

	body, err := io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatalf("read body err: %v", err)
	}
	defer w.Result().Body.Close()

	var resp map[string]interface{}
	err = json.Unmarshal([]byte(body), &resp)
	fmt.Println(resp)

	// assert.Nil(t, err)
	// value := resp.Retcode

	// assert.Equal(t, value, 0)
}
