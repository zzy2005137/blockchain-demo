package routers

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/lib"
	"github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/utils"
)

//QueryOrderHistory 获取指定订单的历史记录
func QueryOrderHistoryPrivate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 获取用户证书

	userCert := GetUserCert(stub, args)
	if userCert != "Admin@org1.blockchainrealestate.com" {
		return shim.Error("Error: private data access failed")
	}

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

func GetUserCert(stub shim.ChaincodeStubInterface, args []string) string {
	creatorByte, _ := stub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		fmt.Errorf("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Errorf("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Errorf("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName
	fmt.Println("Name:" + uname)
	//logger.Infof("===========Name======Name=========%s", uname)
	// return shim.Success([]byte(uname))
	return uname
}
