package fast

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"go-edgetts/logger"
)

// HttpPostJson 发送POST请求，并将args对象序列化为JSON格式作为请求的Body
func HttpPostJson(url string, args interface{}) ([]byte, error) {
	// 获取一个请求对象
	req := fasthttp.AcquireRequest()
	// 请求完成后释放请求对象
	defer fasthttp.ReleaseRequest(req)

	// 设置请求的Content-Type为application/json
	// 默认是application/x-www-form-urlencoded
	req.Header.SetContentType("application/json")
	// 设置请求方法为POST
	req.Header.SetMethod("POST")
	// 设置请求的URI
	req.SetRequestURI(url)

	// 将args对象序列化为JSON格式
	marshal, _ := jsoniter.Marshal(args)
	// 设置请求的Body为序列化后的JSON数据
	req.SetBody(marshal)

	// 获取一个响应对象
	resp := &fasthttp.Response{}
	// 响应完成后释放响应对象
	//defer fasthttp.ReleaseResponse(resp)

	// 发送请求并处理可能的错误
	if err := fasthttp.Do(req, resp); err != nil {
		// 请求失败时记录日志并返回空字符串
		logger.SugarLogger.Debugln("请求失败:", err.Error())
		return nil, errors.New("请求失败")
	}
	// 记录响应的状态码
	logger.SugarLogger.Debugf("后台返回结果:%d", resp.StatusCode())
	// 返回响应的Body内容
	return resp.Body(), nil
}
