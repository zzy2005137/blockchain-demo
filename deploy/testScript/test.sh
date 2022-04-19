#!/bin/bash

# echo "验证查询账户信息"
# docker exec cli peer chaincode invoke -C assetschannel -n blockchain-real-estate -c '{"Args":["queryAccountList"]}'


#其他智能合约函数测试


#create order
# docker exec cli peer chaincode invoke -C assetschannel -n blockchain-real-estate -c '{"Args":["createOrder","5feceb66ffc8","d4735e3a265e","002","toBeStarted"]}'  #管理员id, owner, orderId, status

# query orderList 
echo "查询所有订单信息"
docker exec cli peer chaincode invoke -C assetschannel -n blockchain-real-estate -c '{"Args":["queryOrderList"]}'  

#update order 
# docker exec cli peer chaincode invoke -C assetschannel -n blockchain-real-estate -c '{"Args":["updateOrder","d4735e3a265e","d4735e3a265e","001","inProgress"]}'

#get history 
echo "查询001号订单历史交易信息"
docker exec cli peer chaincode invoke -C assetschannel -n blockchain-real-estate -c '{"Args":["queryOrderHistory","d4735e3a265e","d4735e3a265e","001"]}'