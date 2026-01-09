# Mock Data 使用说明

## 文件说明

`mock_data_100.sql` 包含了 Liaison 系统的完整测试数据，所有实体之间都有正确的关联关系。

**文件位置**: `test/mock/mock_data_100.sql`

## 数据概览

### 1. Users (用户) - 1个
- `admin@liaison.com` - 管理员账户（已登录）

**注意**: 用户密码是 `password123`（已使用 bcrypt 加密）

### 2. Devices (设备) - 100个
设备分布在不同的环境（生产、测试、开发、边缘），每个设备都有：
- 随机配置的 CPU（2-16核）、内存（4GB-32GB）、磁盘（128GB-1TB）
- 不同的操作系统（Linux、Windows、macOS）
- 1-3个网卡接口
- 随机的在线/离线状态（每10个设备有1个离线）

示例设备：
- **设备-001**: 生产服务器，8核CPU，16GB内存，500GB磁盘，Linux Ubuntu 22.04
- **设备-002**: 测试服务器，4核CPU，8GB内存，250GB磁盘，Linux CentOS 7.9
- **设备-003**: 开发服务器，16核CPU，32GB内存，1TB磁盘，Linux Debian 11
- ... 等等

### 3. Ethernet Interfaces (网卡接口) - 约200个
每个设备有1-3个网卡接口，包括：
- 以太网接口（eth0, eth1）
- 无线网卡（wlan0）
- 随机生成的 MAC 地址
- 不同网段的 IP 地址（192.168.x.x）

### 4. Edges (连接器) - 100个
- 每个设备可以关联1-2个连接器
- 连接器名称格式：`连接器-001`, `连接器-002`, ...
- 状态：大部分运行中，每20个有1个停止
- 在线状态：与关联设备的状态相关

### 5. Access Keys (访问密钥) - 100个
每个 Edge 都有一个对应的 AccessKey：
- 格式：`ak-001-xxxxxxxx` / `sk-001-secret-xxxxxxxx`
- 每个连接器都有唯一的访问密钥对

### 6. Applications (应用) - 100个
- 应用名称格式：`应用-001`, `应用-002`, ...
- 每个应用关联到一个设备和一个或多个连接器
- 端口范围：8000-9999
- 应用类型：TCP
- 每10个应用有1个关联多个Edge（高可用场景）

### 7. Proxies (代理) - 100个
每个 Application 都有一个对应的 Proxy：
- 代理名称格式：`代理-001`, `代理-002`, ...
- 外部端口范围：10000-59999
- 状态：大部分运行中，每15个有1个停止
- 每个代理都有描述信息

## 关联关系图

```
User (1个，独立)
  ↓
Device (100个)
  ├─→ EthernetInterface (约200个，每个设备1-3个)
  └─→ Edge (100个，每个设备1-2个)
      └─→ AccessKey (100个，每个Edge一个)
          └─→ Application (100个)
              └─→ Proxy (100个，每个App一个)
```

## 数据特点

### 关联关系
- **Device → Edge**: 通过 `device_id` 关联，每个设备可以有1-2个连接器
- **Edge → AccessKey**: 通过 `edge_id` 关联，每个连接器一个密钥
- **Device → Application**: 通过 `device_id` 关联
- **Edge → Application**: 通过 `edge_ids` JSON 数组关联（支持多Edge）
- **Application → Proxy**: 通过 `application_id` 关联，每个应用一个代理

### 数据多样性
- **操作系统**: Linux (Ubuntu, CentOS, Debian)、Windows Server、macOS
- **状态分布**: 
  - 设备：90%在线，10%离线
  - 连接器：95%运行中，5%停止
  - 代理：93%运行中，7%停止
- **IP地址**: 分布在不同的网段（192.168.10.x - 192.168.250.x）
- **端口**: 应用端口8000-9999，代理端口10000-59999
- **资源使用**: 随机的CPU、内存、磁盘使用率

### 真实性
- 合理的资源使用率（CPU: 10-80%, Memory: 20-90%, Disk: 15-70%）
- 时间戳按创建顺序递增
- 符合规范的MAC地址格式
- 密码已加密（所有用户密码都是 `password123`）

## 使用方法

### 方式1: 直接执行 SQL 文件

```bash
# MySQL/MariaDB
mysql -u your_user -p your_database < test/mock/mock_data_100.sql

# SQLite
sqlite3 your_database.db < test/mock/mock_data_100.sql

# PostgreSQL
psql -U your_user -d your_database -f test/mock/mock_data_100.sql
```

### 方式2: 使用 Python 脚本重新生成

如果需要修改数据或重新生成：

```bash
# 修改 test/mock/generate_mock_data.py 中的参数
python3 test/mock/generate_mock_data.py
```

脚本会生成新的 `test/mock_data_100.sql` 文件。

### 方式3: 在应用中使用

如果使用 GORM，可以创建一个 Go 脚本来插入数据：

```go
// 示例代码
func InsertMockData(db *gorm.DB) error {
    // 读取并执行 SQL 文件
    sqlBytes, err := os.ReadFile("test/mock/mock_data_100.sql")
    if err != nil {
        return err
    }
    
    sql := string(sqlBytes)
    // 分割 SQL 语句并执行
    statements := strings.Split(sql, ";")
    for _, stmt := range statements {
        stmt = strings.TrimSpace(stmt)
        if stmt != "" && !strings.HasPrefix(stmt, "--") {
            if err := db.Exec(stmt).Error; err != nil {
                return err
            }
        }
    }
    return nil
}
```

## 注意事项

1. **ID 自增**: 如果数据库表已有数据，可能需要调整 ID 或使用 AUTO_INCREMENT
2. **外键约束**: 确保按照顺序插入（Device → Edge → AccessKey → Application → Proxy）
3. **密码**: 所有用户密码都是 `password123`（bcrypt 哈希值）
4. **时间戳**: 所有时间戳都是示例时间，可以根据需要调整
5. **JSON 字段**: Application 的 `edge_ids` 字段使用 JSON 格式，如 `[1]` 或 `[1,5]`
6. **NOW() 函数**: SQL 文件中使用了 `NOW()` 函数，某些数据库可能需要替换为具体时间

## 测试场景

这些数据可以用于测试以下场景：

### 1. 搜索功能
- ✅ 按设备名搜索 Edge（支持模糊匹配）
- ✅ 按连接器名搜索 Edge（支持模糊匹配）
- ✅ 按代理名搜索 Proxy（支持模糊匹配）
- ✅ 按设备名搜索 Application（支持模糊匹配）
- ✅ 按IP搜索 Device（支持模糊匹配）

### 2. 关联查询
- ✅ Edge 列表显示关联的 Device
- ✅ Application 列表显示关联的 Device 和 Proxy
- ✅ Proxy 列表显示关联的 Application

### 3. 状态管理
- ✅ 测试 Proxy 的启动/停止功能
- ✅ 测试 Edge 的在线/离线状态
- ✅ 测试 Device 的在线/离线状态

### 4. 分页和排序
- ✅ 测试各种实体的列表分页
- ✅ 测试默认倒序排列（按ID倒序）

### 5. 性能测试
- ✅ 测试大量数据下的查询性能
- ✅ 测试关联查询的性能
- ✅ 测试搜索功能的性能

## 数据清理

如果需要清理测试数据：

### 方式1: 硬删除
```sql
-- 注意：按照外键依赖的相反顺序删除
DELETE FROM proxies;
DELETE FROM applications;
DELETE FROM access_keys;
DELETE FROM edges;
DELETE FROM ethernet_interfaces;
DELETE FROM devices;
DELETE FROM users;
```

### 方式2: 软删除（如果启用了 GORM 的软删除）
```sql
UPDATE proxies SET deleted_at = NOW();
UPDATE applications SET deleted_at = NOW();
UPDATE access_keys SET deleted_at = NOW();
UPDATE edges SET deleted_at = NOW();
UPDATE ethernet_interfaces SET deleted_at = NOW();
UPDATE devices SET deleted_at = NOW();
-- 注意：用户数据通常不删除
```

## 数据统计

执行以下 SQL 可以查看数据统计：

```sql
-- 统计各实体数量
SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'devices', COUNT(*) FROM devices
UNION ALL
SELECT 'ethernet_interfaces', COUNT(*) FROM ethernet_interfaces
UNION ALL
SELECT 'edges', COUNT(*) FROM edges
UNION ALL
SELECT 'access_keys', COUNT(*) FROM access_keys
UNION ALL
SELECT 'applications', COUNT(*) FROM applications
UNION ALL
SELECT 'proxies', COUNT(*) FROM proxies;

-- 统计在线/离线设备
SELECT online, COUNT(*) as count FROM devices GROUP BY online;

-- 统计运行中/停止的代理
SELECT status, COUNT(*) as count FROM proxies GROUP BY status;

-- 统计每个设备关联的连接器数量
SELECT d.name, COUNT(e.id) as edge_count 
FROM devices d 
LEFT JOIN edges e ON d.id = e.device_id 
GROUP BY d.id, d.name 
ORDER BY edge_count DESC;
```

## 常见问题

### Q: 如何修改数据量？
A: 编辑 `test/generate_mock_data.py`，修改 `main()` 函数中的数量参数。

### Q: 如何修改用户密码？
A: 使用 Go 的密码工具生成新的 bcrypt 哈希，或直接修改 SQL 文件中的密码哈希值。

### Q: 数据导入失败怎么办？
A: 
1. 检查数据库连接和权限
2. 确认表结构是否正确
3. 检查是否有外键约束冲突
4. 查看数据库错误日志

### Q: 如何生成更多数据？
A: 修改 `generate_mock_data.py` 中的 `count` 参数，例如：
```python
devices = generate_devices(500)  # 生成500个设备
```

## 更新日志

- **2026-01-09**: 初始版本，生成100个实体（用户、设备、连接器、应用、代理）
- **2026-01-09**: 更新为仅保留1个用户（admin），其他实体保持100个
