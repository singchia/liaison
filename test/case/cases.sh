# 边缘
# 列举边缘
curl 127.0.0.1:8080/v1/edges 

# 创建边缘
curl -X POST 127.0.0.1:8080/v1/edges -H "Content-Type: application/json" -d '{"name": "test", "description": "test"}'

# 获取边缘
curl 127.0.0.1:8080/v1/edges/1

# 更新边缘
curl -X PUT 127.0.0.1:8080/v1/edges/1 -H "Content-Type: application/json" -d '{"name": "test2", "description": "test2"}'

# 删除边缘

# 应用
# 创建应用
curl -X POST 127.0.0.1:8080/v1/applications -H 'Content-type: application/json' -d '{"ip": "127.0.0.1", "name": "ssh", "application_type": "tcp", "port": 22, "edge_id": 1}'

# 代理
# 创建代理
curl -X POST 127.0.0.1:8080/v1/proxies -H 'Content-type: application/json' -d '{"name": "ssh", "port": 65522, "application_id": 1}'

# 更新代理
curl -X PUT 127.0.0.1:8080/v1/proxies/1 -H 'Content-type: application/json' -d '{"name": "ssh2", "description": "ssh2"}'

# 删除代理
curl -X DELETE 127.0.0.1:8080/v1/proxies/1