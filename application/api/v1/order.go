package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	bc "application/blockchain"
	"application/pkg/app"

	"github.com/gin-gonic/gin"
)

type OrderQueryRequestBody struct {
	OrderId string `json:"orderId"` //orderId
}

func QueryOrderList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(OrderQueryRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.OrderId != "" {
		bodyBytes = append(bodyBytes, []byte(body.OrderId))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("queryOrderList", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

func QueryOrderHistory(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(OrderQueryRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}

	order001Info := [3]string{"d4735e3a265e", "d4735e3a265e", "001"}
	var bodyBytes [][]byte
	// if body.OrderId != "" {
	// 	bodyBytes = append(bodyBytes, []byte(body.OrderId))
	// }

	for _, val := range order001Info {
		bodyBytes = append(bodyBytes, []byte(val))
	}

	//调用智能合约
	resp, err := bc.ChannelQuery("queryOrderHistory", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}
