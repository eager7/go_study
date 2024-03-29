package gin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eager7/go/mlog"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type BaseResp struct {
	Errno  int64  `json:"errno"`
	Errmsg string `json:"errmsg"`
}

type GetActionTypesRsp struct {
	BaseResp
	Data TypesData `json:"data"`
}

type TypesData struct {
	Types []string `json:"types"`
}

type Ping struct {
	who string
}

func InitializeGin() {
	router := gin.Default()
	router.GET("/:who/ping", PingHandle)
	router.POST("/file_update", FileHandle)
	if err := router.Run(); err != nil {
		panic(err)
	}
}

func PingHandle(context *gin.Context) {
	who := Ping{who: context.Param("who")}
	if context.BindQuery(&who) != nil {
		resp := GetActionTypesRsp{
			BaseResp: BaseResp{
				Errno:  1,
				Errmsg: "bind failed",
			},
			Data: TypesData{
				Types: nil,
			},
		}
		context.JSON(200, resp)
	}
	context.JSON(200, gin.H{"message": fmt.Sprintf("%s pong", who.who)})
}
/*
curl 测试：
curl -X POST http://localhost:8080/file_update -H "Content-Type: multipart/form-data" -F "account=pct" -F "file=@./send_post.sh"
*/
func FileHandle(context *gin.Context) {
	form, err := context.MultipartForm()//多文件上传
	if err != nil {
		mlog.L.Error(err)
		context.JSON(500, gin.H{"message": "param error"})
		return
	}
	mlog.L.Debug(form.Value, form.File)
	mlog.L.Debug(form.Value["account"], len(form.File["file"]))
	if len(form.File["file"]) > 0 {
		f := form.File["file"][0]
		mlog.L.Debug(f.Filename)
		if fi, err := f.Open(); err != nil {
			mlog.L.Error(err)
		} else {
			if buffer, err := ioutil.ReadAll(fi); err != nil {
				mlog.L.Error(err)
			} else {
				mlog.L.Info("context:", string(buffer))
			}
		}
		if err := context.SaveUploadedFile(f, fmt.Sprintf("/tmp/file/%s", f.Filename)); err != nil {
			mlog.L.Error("save file error:", err)
		}
	}

	context.JSON(200, gin.H{"message": "post resp"})
}

func SendAndRecv(reqUrl string, recvBody interface{}) error {
	body, err := SendRequest("GET", nil, reqUrl)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, recvBody); err != nil {
		fmt.Println("Unmarshal err:", reqUrl, " body:", string(body))
		return err
	}
	return nil
}

func SendRequest(method string, tReq interface{}, reqUrl string) (body []byte, err error) {
	Client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			MaxIdleConns:        3,
			MaxIdleConnsPerHost: 3,
			IdleConnTimeout:     time.Duration(10) * time.Second,
		},
		Timeout: time.Duration(300) * time.Second,
	}
	//通常是采用strings.NewReader函数，将一个string类型转化为io.Reader类型，或者bytes.NewBuffer函数，将[]byte类型转化为io.Reader类型。
	reqBuffers, _ := json.Marshal(tReq)
	reqByte := bytes.NewBuffer(reqBuffers)

	req, err := http.NewRequest("GET", "http://baidu.com", nil)
	if tReq == nil {
		req, err = http.NewRequest(method, reqUrl, nil)
	} else {
		req, err = http.NewRequest(method, reqUrl, reqByte)
	}
	if err != nil {
		fmt.Println("SendRequest  failed.  req_buf:", string(reqBuffers), "req_url:", reqUrl, " err:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := Client.Do(req)
	if err != nil && resp == nil {
		fmt.Println("SendRequest  failed.  req_buf:", string(reqBuffers), " req_url:", reqUrl, " err:", err)
		return nil, err
	}
	defer checkError(resp.Body.Close)

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("SendRequest  failed.  req_buf:", string(reqBuffers), "req_url:", reqUrl, " err:", err)
		return nil, err
	}
	return body, nil
}

func checkError(callBacks ...func() error) {
	for _, callBack := range callBacks {
		if err := callBack(); err != nil {
			fmt.Println(err)
		}
	}
}
