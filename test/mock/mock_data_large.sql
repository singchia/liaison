-- Mock Data for Liaison System (Large Dataset)
-- 生成时间: 2026-01-09
-- 说明: 包含100个用户、100个设备、100个连接器、100个应用、100个代理等相互关联的测试数据
-- ⚠️ 注意: 此文件使用 MySQL/MariaDB 语法（包含存储过程、变量、循环等）
-- 如果使用 SQLite，请使用 mock_data_100.sql 文件

-- ============================================
-- 1. Users (用户) - 100个
-- ============================================
INSERT INTO users (created_at, updated_at, deleted_at, email, password, status, last_login, login_ip) VALUES
('2025-01-01 10:00:00', '2025-01-01 10:00:00', NULL, 'admin@liaison.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', '2026-01-09 08:30:00', '192.168.1.100'),
('2025-01-02 11:00:00', '2025-01-02 11:00:00', NULL, 'user1@liaison.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', '2026-01-08 15:20:00', '192.168.1.101'),
('2025-01-03 12:00:00', '2025-01-03 12:00:00', NULL, 'user2@liaison.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'active', NULL, NULL);

-- 生成剩余97个用户
SET @user_counter = 3;
WHILE @user_counter < 100 DO
  SET @user_counter = @user_counter + 1;
  SET @email = CONCAT('user', @user_counter, '@liaison.com');
  SET @created_date = DATE_ADD('2025-01-01', INTERVAL @user_counter DAY);
  INSERT INTO users (created_at, updated_at, deleted_at, email, password, status, last_login, login_ip) VALUES
  (@created_date, @created_date, NULL, @email, '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 
   IF(@user_counter % 10 = 0, 'inactive', 'active'),
   IF(@user_counter % 3 = 0, DATE_ADD(NOW(), INTERVAL -@user_counter HOUR), NULL),
   IF(@user_counter % 3 = 0, CONCAT('192.168.1.', 100 + @user_counter), NULL));
END WHILE;

-- ============================================
-- 2. Devices (设备) - 100个
-- ============================================
-- 使用循环生成100个设备
SET @device_counter = 0;
WHILE @device_counter < 100 DO
  SET @device_counter = @device_counter + 1;
  SET @device_name = CONCAT('设备-', LPAD(@device_counter, 3, '0'));
  SET @fingerprint = CONCAT('fp-device-', LPAD(@device_counter, 3, '0'), '-', SUBSTRING(MD5(RAND()), 1, 6));
  SET @host_name = CONCAT('host-', LPAD(@device_counter, 3, '0'));
  SET @created_date = DATE_ADD('2025-01-05', INTERVAL @device_counter DAY);
  SET @online_status = IF(@device_counter % 10 = 0, 2, 1); -- 每10个有一个离线
  SET @cpu = 2 + (RAND() * 14); -- 2-16核
  SET @memory = 4096 + (RAND() * 28672); -- 4GB-32GB
  SET @disk = 128000 + (RAND() * 896000); -- 128GB-1TB
  SET @os_type = IF(@device_counter % 3 = 0, 'Linux', IF(@device_counter % 3 = 1, 'Windows', 'macOS'));
  SET @os_version = IF(@device_counter % 3 = 0, 'Ubuntu 22.04', IF(@device_counter % 3 = 1, 'Windows Server 2022', 'macOS 13.0'));
  
  INSERT INTO devices (created_at, updated_at, deleted_at, fingerprint, name, description, host_name, online, heartbeat_at, cpu, memory, disk, os, os_version, cpu_usage, memory_usage, disk_usage) VALUES
  (@created_date, NOW(), NULL, @fingerprint, @device_name, 
   CONCAT('[MOCK] ', @device_name, '的描述信息'), @host_name, @online_status, 
   IF(@online_status = 1, NOW(), DATE_ADD(NOW(), INTERVAL -2 DAY)),
   FLOOR(@cpu), FLOOR(@memory), FLOOR(@disk), @os_type, @os_version,
   RAND() * 80, RAND() * 90, RAND() * 70);
END WHILE;

-- ============================================
-- 3. Ethernet Interfaces (网卡接口)
-- ============================================
-- 为每个设备生成1-3个网卡接口
SET @device_counter = 0;
WHILE @device_counter < 100 DO
  SET @device_counter = @device_counter + 1;
  SET @interface_count = 1 + FLOOR(RAND() * 3); -- 1-3个网卡
  SET @interface_idx = 0;
  
  WHILE @interface_idx < @interface_count DO
    SET @interface_idx = @interface_idx + 1;
    SET @interface_name = IF(@interface_idx = 1, 'eth0', IF(@interface_idx = 2, 'eth1', 'wlan0'));
    SET @mac_base = CONCAT('00:', LPAD(FLOOR(RAND() * 255), 2, '0'), ':', LPAD(FLOOR(RAND() * 255), 2, '0'));
    SET @mac = CONCAT(@mac_base, ':', LPAD(FLOOR(RAND() * 255), 2, '0'), ':', LPAD(FLOOR(RAND() * 255), 2, '0'), ':', LPAD(FLOOR(RAND() * 255), 2, '0'));
    SET @ip_third = 10 + (@device_counter % 240);
    SET @ip_fourth = 10 + (@interface_idx * 10) + (@device_counter % 10);
    SET @ip = CONCAT('192.168.', @ip_third, '.', @ip_fourth);
    
    INSERT INTO ethernet_interfaces (created_at, updated_at, deleted_at, device_id, name, mac, ip, netmask) VALUES
    (DATE_ADD('2025-01-05', INTERVAL @device_counter DAY), 
     DATE_ADD('2025-01-05', INTERVAL @device_counter DAY), 
     NULL, @device_counter, @interface_name, @mac, @ip, '255.255.255.0');
  END WHILE;
END WHILE;

-- ============================================
-- 4. Edges (连接器) - 100个
-- ============================================
-- 每个设备可以有1-2个连接器
SET @edge_counter = 0;
SET @device_counter = 0;
WHILE @edge_counter < 100 DO
  SET @device_counter = @device_counter + 1;
  IF @device_counter > 100 THEN SET @device_counter = 1; END IF;
  
  SET @edge_counter = @edge_counter + 1;
  SET @edge_name = CONCAT('连接器-', LPAD(@edge_counter, 3, '0'));
  SET @edge_status = IF(@edge_counter % 20 = 0, 2, 1); -- 每20个有一个停止
  SET @edge_online = IF(@device_counter % 10 = 0, 2, 1); -- 与设备状态相关
  SET @created_date = DATE_ADD('2025-01-10', INTERVAL @edge_counter DAY);
  
  INSERT INTO edges (created_at, updated_at, deleted_at, name, status, online, heartbeat_at, description, device_id) VALUES
  (@created_date, NOW(), NULL, @edge_name, @edge_status, @edge_online,
   IF(@edge_online = 1, NOW(), DATE_ADD(NOW(), INTERVAL -1 DAY)),
   CONCAT('[MOCK] ', @edge_name, '的描述信息'), @device_counter);
END WHILE;

-- ============================================
-- 5. Access Keys (访问密钥) - 100个
-- ============================================
-- 每个Edge一个AccessKey
SET @ak_counter = 0;
WHILE @ak_counter < 100 DO
  SET @ak_counter = @ak_counter + 1;
  SET @access_key = CONCAT('ak-', LPAD(@ak_counter, 3, '0'), '-', SUBSTRING(MD5(RAND()), 1, 8));
  SET @secret_key = CONCAT('sk-', @ak_counter, '-secret-', SUBSTRING(MD5(RAND()), 1, 16));
  
  INSERT INTO access_keys (created_at, updated_at, deleted_at, edge_id, access_key, secret_key) VALUES
  (DATE_ADD('2025-01-10', INTERVAL @ak_counter DAY),
   DATE_ADD('2025-01-10', INTERVAL @ak_counter DAY),
   NULL, @ak_counter, @access_key, @secret_key);
END WHILE;

-- ============================================
-- 6. Applications (应用) - 100个
-- ============================================
SET @app_counter = 0;
SET @device_counter = 0;
SET @edge_counter = 0;
WHILE @app_counter < 100 DO
  SET @app_counter = @app_counter + 1;
  SET @device_counter = @device_counter + 1;
  IF @device_counter > 100 THEN SET @device_counter = 1; END IF;
  
  SET @edge_counter = @edge_counter + 1;
  IF @edge_counter > 100 THEN SET @edge_counter = 1; END IF;
  
  -- 有些应用关联多个Edge
  SET @edge_ids = IF(@app_counter % 10 = 0, 
    CONCAT('[', @edge_counter, ',', IF(@edge_counter + 1 > 100, 1, @edge_counter + 1), ']'),
    CONCAT('[', @edge_counter, ']'));
  
  SET @app_name = CONCAT('应用-', LPAD(@app_counter, 3, '0'));
  SET @app_description = CONCAT('[MOCK] ', @app_name, '的描述信息');
  SET @ip_third = 10 + (@device_counter % 240);
  SET @ip_fourth = 100 + (@app_counter % 100);
  SET @ip = CONCAT('192.168.', @ip_third, '.', @ip_fourth);
  SET @port = 8000 + (@app_counter % 2000); -- 8000-9999端口
  SET @created_date = DATE_ADD('2025-01-15', INTERVAL @app_counter DAY);
  
  INSERT INTO applications (created_at, updated_at, deleted_at, edge_ids, device_id, name, description, ip, port, heartbeat_at, application_type) VALUES
  (@created_date, NOW(), NULL, @edge_ids, @device_counter, @app_name, @app_description, @ip, @port,
   IF(@device_counter % 10 = 0, DATE_ADD(NOW(), INTERVAL -1 DAY), NOW()),
   'tcp');
END WHILE;

-- ============================================
-- 7. Proxies (代理) - 100个
-- ============================================
SET @proxy_counter = 0;
WHILE @proxy_counter < 100 DO
  SET @proxy_counter = @proxy_counter + 1;
  SET @proxy_name = CONCAT('代理-', LPAD(@proxy_counter, 3, '0'));
  SET @proxy_status = IF(@proxy_counter % 15 = 0, 2, 1); -- 每15个有一个停止
  SET @proxy_port = 10000 + (@proxy_counter % 50000); -- 10000-59999端口
  SET @created_date = DATE_ADD('2025-01-20', INTERVAL @proxy_counter DAY);
  
  INSERT INTO proxies (created_at, updated_at, deleted_at, application_id, name, port, status, description) VALUES
  (@created_date, NOW(), NULL, @proxy_counter, @proxy_name, @proxy_port, @proxy_status,
   CONCAT('[MOCK] ', @proxy_name, '的描述信息，外部端口', @proxy_port));
END WHILE;
