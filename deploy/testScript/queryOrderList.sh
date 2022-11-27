echo "验证不同节点数据查询一致性"
echo "通过peer1查询"
docker exec cli peer chaincode invoke -C assetschannel -n blockchain-real-estate -c '{"Args":["queryOrderList"]}'
echo "通过peer0查询"
docker exec cli0 peer chaincode invoke -C assetschannel -n blockchain-real-estate -c '{"Args":["queryOrderList"]}'
echo "通过peer2查询"
docker exec cli2 peer chaincode invoke -C assetschannel -n blockchain-real-estate -c '{"Args":["queryOrderList"]}'