package routers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/lib"
	"github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/utils"
)

// Createorder 新建订单(管理员)
// CompositeKeys = OrderKey-Owner-OrderId
func CreateOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 4 {
		return shim.Error("参数个数不满足")
	}
	accountId := args[0] //accountId用于验证是否为管理员
	owner := args[1]     //order owner
	orderId := args[2]   //order id
	statusMap := lib.OrderStatusConstant()
	status := statusMap[args[3]] //order status

	if accountId == "" || orderId == "" || owner == "" || status == "" {
		return shim.Error("参数存在空值")
	}

	//判断是否管理员操作
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, lib.AccountKey, []string{accountId})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("操作人权限验证失败%s", err))
	}
	var account lib.Account
	if err = json.Unmarshal(resultsAccount[0], &account); err != nil {
		return shim.Error(fmt.Sprintf("查询操作人信息-反序列化出错: %s", err))
	}
	if account.UserName != "管理员" {
		return shim.Error(fmt.Sprintf("ERROR : Operation failed, not Admin %s", err))
	}

	//判断订单是否存在
	// resultsOrderId, err := utils.GetStateByPartialCompositeKeys(stub, lib.OrderKey, []string{orderId})
	// if err != nil || len(resultsOrderId) >= 1 {
	// 	return shim.Error(fmt.Sprintf("OrderId已存在%s", err))
	// }

	//判断订单是否存在
	orderCheck, err := utils.GetStateByPartialCompositeKeys2(stub, lib.OrderKey, []string{owner, orderId})
	if err != nil || len(orderCheck) != 0 {
		return shim.Error(fmt.Sprintf("订单已存在%s", err))
	}

	//创建新订单
	order := &lib.Order{
		// OrderId: fmt.Sprintf("%d", time.Now().Local().UnixNano()),
		OrderId: orderId,
		Owner:   owner,
		Status:  status,
	}
	// 写入账本
	if err := utils.WriteLedger(order, stub, lib.OrderKey, []string{order.Owner, order.OrderId}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	orderByte, err := json.Marshal(order)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(orderByte)
}

// QueryorderList 查询订单(可查询所有，也可根据所有人查询名下订单)
func QueryOrderList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var orderList []lib.Order
	results, err := utils.GetStateByPartialCompositeKeys2(stub, lib.OrderKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var order lib.Order
			err := json.Unmarshal(v, &order)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryorderList-反序列化出错: %s", err))
			}
			orderList = append(orderList, order)
		}
	}
	orderListByte, err := json.Marshal(orderList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryorderList-序列化出错: %s", err))
	}
	return shim.Success(orderListByte)
}

// UpdateOrder 更新订单状态
func UpdateOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 4 {
		return shim.Error("参数个数不满足")
	}
	accountId := args[0] //accountId
	owner := args[1]     //key1
	orderId := args[2]   //key2
	statusMap := lib.OrderStatusConstant()
	status := statusMap[args[3]]

	if accountId == "" || orderId == "" || owner == "" || status == "" {
		return shim.Error("参数存在空值")
	}

	//判断订单是否存在
	orderCheck, err := utils.GetStateByPartialCompositeKeys2(stub, lib.OrderKey, []string{owner, orderId})
	if err != nil || len(orderCheck) != 1 {
		return shim.Error(fmt.Sprintf("OrderId不存在%s", err))
	}
	//获取目标订单
	var targetOrder lib.Order
	err = json.Unmarshal(orderCheck[0], &targetOrder)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryorderList-反序列化出错: %s", err))
	}

	//判断是否为订单所有方操作
	if targetOrder.Owner != accountId {
		return shim.Error(fmt.Sprintf("操作人权限不足%s", err))
	}

	//更新订单信息
	targetOrder.Status = status

	// 写入账本
	if err := utils.WriteLedger(targetOrder, stub, lib.OrderKey, []string{targetOrder.Owner, targetOrder.OrderId}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	orderByte, err := json.Marshal(targetOrder)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(orderByte)
}

//QueryOrderHistory 获取指定订单的历史记录
func QueryOrderHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 3 {
		return shim.Error("参数个数不满足")
	}
	accountId := args[0] //accountId
	owner := args[1]     //owner
	orderId := args[2]   //id

	if accountId == "" || orderId == "" || owner == "" {
		return shim.Error("参数存在空值")
	}

	//判断订单是否存在
	orderCheck, err := utils.GetStateByPartialCompositeKeys2(stub, lib.OrderKey, []string{owner, orderId})
	if err != nil || len(orderCheck) != 1 {
		return shim.Error(fmt.Sprintf("OrderId不存在%s", err))
	}

	keyArgs := args[1:]

	//创建复合键
	key, err := stub.CreateCompositeKey(lib.OrderKey, keyArgs)
	if err != nil {
		shim.Error(fmt.Sprintf("%s-创建组合键出错: %s", lib.OrderKey, err))
	}

	//获取历史数据
	resultsIterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		shim.Error(fmt.Sprintf("查询历史信息错误: %s", err))
	}
	defer resultsIterator.Close()

	// var results [][]byte
	// for resultIterator.HasNext() {
	// 	val, err := resultIterator.Next()
	// 	if err != nil {
	// 		return shim.Error(fmt.Sprintf("返回的数据出错"))
	// 	}
	// 	results = append(results, val.GetValue())
	// }

	// for _, v := range results {
	// 	if v != nil {
	// 		var order lib.Order
	// 		err := json.Unmarshal(v, &order)
	// 		if err != nil {
	// 			return shim.Error(fmt.Sprintf("QueryorderList-反序列化出错: %s", err))
	// 		}
	// 		orderList = append(orderList, order)
	// 	}
	// }

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForMarble returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
	// return shim.Success(results[0])
}
