#!/bin/bash

# echo "验证查询账户信息"
# docker exec cli peer chaincode invoke -C assetschannel -n blockchain-manufacture -c '{"Args":["queryAccountList"]}'


#其他智能合约函数测试


# create order
# echo "创建订单001"
# docker exec cli peer chaincode invoke -C assetschannel -n blockchain-manufacture -c '{"Args":["createOrder","5feceb66ffc8","d4735e3a265e","001","toBeStarted"]}'  #管理员id, owner, orderId, status

# sleep 2
# # query orderList 
# echo "查询所有订单信息"
# docker exec cli peer chaincode invoke -C assetschannel -n blockchain-manufacture -c '{"Args":["queryOrderList"]}'  

# sleep 2
#update order 
# echo "更新订单信息"
# docker exec cli peer chaincode invoke -C assetschannel -n blockchain-manufacture -c '{"Args":["updateOrder","d4735e3a265e","d4735e3a265e","003","inProgress"]}'

# sleep 2
#get history 
echo "组织1查询001号订单历史交易信息"
docker exec cli peer chaincode invoke -C assetschannel -n blockchain-manufacture -c '{"Args":["queryOrderHistory","d4735e3a265e","d4735e3a265e","003"]}'

