#!/usr/bin/env python3
"""
生成 Liaison 系统的 Mock 数据
每个实体生成 100 条记录，并保持合理的关联关系
"""

import random
import hashlib
from datetime import datetime, timedelta

def generate_users(count=1):
    """生成用户数据（只保留admin用户）"""
    users = []
    password_hash = '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'  # password123
    
    users.append({
        'email': 'admin@liaison.com',
        'password': password_hash,
        'status': 'active',
        'last_login': '2026-01-09 08:30:00',
        'login_ip': '192.168.1.100',
        'created_at': '2025-01-01 10:00:00'
    })
    
    return users

def generate_devices(count=100):
    """生成设备数据"""
    devices = []
    base_date = datetime(2025, 1, 5)
    os_types = ['Linux', 'Windows', 'macOS']
    os_versions = {
        'Linux': ['Ubuntu 22.04', 'CentOS 7.9', 'Debian 11', 'Ubuntu 20.04'],
        'Windows': ['Windows Server 2022', 'Windows Server 2019', 'Windows 11'],
        'macOS': ['macOS 13.0', 'macOS 12.0', 'macOS 14.0']
    }
    
    for i in range(1, count + 1):
        created_date = base_date + timedelta(days=i)
        fingerprint = f'fp-device-{i:03d}-{hashlib.md5(str(i).encode()).hexdigest()[:6]}'
        name = f'设备-{i:03d}'
        host_name = f'host-{i:03d}'
        online = 2 if i % 10 == 0 else 1
        cpu = random.randint(2, 16)
        memory = random.randint(4096, 32768)
        disk = random.randint(128000, 1024000)
        os_type = os_types[i % 3]
        os_version = random.choice(os_versions[os_type])
        
        devices.append({
            'fingerprint': fingerprint,
            'name': name,
            'description': f'[MOCK] {name}的描述信息',
            'host_name': host_name,
            'online': online,
            'heartbeat_at': datetime.now().strftime('%Y-%m-%d %H:%M:%S') if online == 1 else (datetime.now() - timedelta(days=2)).strftime('%Y-%m-%d %H:%M:%S'),
            'cpu': cpu,
            'memory': memory,
            'disk': disk,
            'os': os_type,
            'os_version': os_version,
            'cpu_usage': round(random.uniform(10, 80), 1),
            'memory_usage': round(random.uniform(20, 90), 1),
            'disk_usage': round(random.uniform(15, 70), 1),
            'created_at': created_date.strftime('%Y-%m-%d %H:%M:%S')
        })
    
    return devices

def generate_ethernet_interfaces(devices):
    """为每个设备生成网卡接口"""
    interfaces = []
    base_date = datetime(2025, 1, 5)
    
    for idx, device in enumerate(devices, 1):
        interface_count = random.randint(1, 3)
        for i in range(interface_count):
            interface_name = ['eth0', 'eth1', 'wlan0'][i] if i < 3 else f'eth{i}'
            mac_parts = [f'{random.randint(0, 255):02X}' for _ in range(6)]
            mac = ':'.join(mac_parts)
            
            ip_third = 10 + (idx % 240)
            ip_fourth = 10 + (i * 10) + (idx % 10)
            ip = f'192.168.{ip_third}.{ip_fourth}'
            
            created_date = base_date + timedelta(days=idx)
            interfaces.append({
                'device_id': idx,
                'name': interface_name,
                'mac': mac,
                'ip': ip,
                'netmask': '255.255.255.0',
                'created_at': created_date.strftime('%Y-%m-%d %H:%M:%S')
            })
    
    return interfaces

def generate_edges(devices, count=100):
    """生成连接器数据"""
    edges = []
    base_date = datetime(2025, 1, 10)
    
    device_idx = 0
    for i in range(1, count + 1):
        device_idx = (device_idx % len(devices)) + 1
        name = f'连接器-{i:03d}'
        status = 2 if i % 20 == 0 else 1
        online = 2 if device_idx % 10 == 0 else 1
        
        created_date = base_date + timedelta(days=i)
        edges.append({
            'name': name,
            'status': status,
            'online': online,
            'heartbeat_at': datetime.now().strftime('%Y-%m-%d %H:%M:%S') if online == 1 else (datetime.now() - timedelta(days=1)).strftime('%Y-%m-%d %H:%M:%S'),
            'description': f'[MOCK] {name}的描述信息',
            'device_id': device_idx,
            'created_at': created_date.strftime('%Y-%m-%d %H:%M:%S')
        })
    
    return edges

def generate_access_keys(edges):
    """为每个连接器生成访问密钥"""
    access_keys = []
    base_date = datetime(2025, 1, 10)
    
    for idx, edge in enumerate(edges, 1):
        access_key = f'ak-{idx:03d}-{hashlib.md5(str(idx).encode()).hexdigest()[:8]}'
        secret_key = f'sk-{idx:03d}-secret-{hashlib.md5(str(idx * 2).encode()).hexdigest()[:16]}'
        
        created_date = base_date + timedelta(days=idx)
        access_keys.append({
            'edge_id': idx,
            'access_key': access_key,
            'secret_key': secret_key,
            'created_at': created_date.strftime('%Y-%m-%d %H:%M:%S')
        })
    
    return access_keys

def generate_applications(devices, edges, count=100):
    """生成应用数据"""
    applications = []
    base_date = datetime(2025, 1, 15)
    
    device_idx = 0
    edge_idx = 0
    for i in range(1, count + 1):
        device_idx = (device_idx % len(devices)) + 1
        edge_idx = (edge_idx % len(edges)) + 1
        
        # 每10个应用有一个关联多个Edge
        if i % 10 == 0:
            edge_ids = f'[{edge_idx},{edge_idx % len(edges) + 1}]'
        else:
            edge_ids = f'[{edge_idx}]'
        
        name = f'应用-{i:03d}'
        ip_third = 10 + (device_idx % 240)
        ip_fourth = 100 + (i % 100)
        ip = f'192.168.{ip_third}.{ip_fourth}'
        port = 8000 + (i % 2000)
        
        created_date = base_date + timedelta(days=i)
        applications.append({
            'edge_ids': edge_ids,
            'device_id': device_idx,
            'name': name,
            'description': f'[MOCK] {name}的描述信息',
            'ip': ip,
            'port': port,
            'heartbeat_at': datetime.now().strftime('%Y-%m-%d %H:%M:%S') if device_idx % 10 != 0 else (datetime.now() - timedelta(days=1)).strftime('%Y-%m-%d %H:%M:%S'),
            'application_type': 'tcp',
            'created_at': created_date.strftime('%Y-%m-%d %H:%M:%S')
        })
    
    return applications

def generate_proxies(applications):
    """为每个应用生成代理"""
    proxies = []
    base_date = datetime(2025, 1, 20)
    
    for idx, app in enumerate(applications, 1):
        name = f'代理-{idx:03d}'
        status = 2 if idx % 15 == 0 else 1
        port = 10000 + (idx % 50000)
        
        created_date = base_date + timedelta(days=idx)
        proxies.append({
            'application_id': idx,
            'name': name,
            'port': port,
            'status': status,
            'description': f'[MOCK] {name}的描述信息，外部端口{port}',
            'created_at': created_date.strftime('%Y-%m-%d %H:%M:%S')
        })
    
    return proxies

def generate_sql(users, devices, interfaces, edges, access_keys, applications, proxies):
    """生成SQL语句"""
    sql_lines = [
        '-- Mock Data for Liaison System (Large Dataset)',
        '-- 生成时间: ' + datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
        '-- 说明: 包含1个用户（admin）、100个设备、100个连接器、100个应用、100个代理等相互关联的测试数据',
        '',
        '-- ============================================',
        '-- 1. Users (用户) - 1个（仅admin用户）',
        '-- ============================================',
    ]
    
    # Users
    for user in users:
        last_login = f"'{user['last_login']}'" if user['last_login'] else 'NULL'
        login_ip = f"'{user['login_ip']}'" if user['login_ip'] else 'NULL'
        sql_lines.append(
            f"INSERT INTO users (created_at, updated_at, deleted_at, email, password, status, last_login, login_ip) VALUES "
            f"('{user['created_at']}', '{user['created_at']}', NULL, '{user['email']}', '{user['password']}', "
            f"'{user['status']}', {last_login}, {login_ip});"
        )
    
    sql_lines.extend(['', '-- ============================================', '-- 2. Devices (设备) - 100个', '-- ============================================'])
    
    # Devices
    for device in devices:
        sql_lines.append(
            f"INSERT INTO devices (created_at, updated_at, deleted_at, fingerprint, name, description, host_name, online, heartbeat_at, cpu, memory, disk, os, os_version, cpu_usage, memory_usage, disk_usage) VALUES "
            f"('{device['created_at']}', NOW(), NULL, '{device['fingerprint']}', '{device['name']}', "
            f"'{device['description']}', '{device['host_name']}', {device['online']}, '{device['heartbeat_at']}', "
            f"{device['cpu']}, {device['memory']}, {device['disk']}, '{device['os']}', '{device['os_version']}', "
            f"{device['cpu_usage']}, {device['memory_usage']}, {device['disk_usage']});"
        )
    
    sql_lines.extend(['', '-- ============================================', '-- 3. Ethernet Interfaces (网卡接口)', '-- ============================================'])
    
    # Ethernet Interfaces
    for interface in interfaces:
        sql_lines.append(
            f"INSERT INTO ethernet_interfaces (created_at, updated_at, deleted_at, device_id, name, mac, ip, netmask) VALUES "
            f"('{interface['created_at']}', '{interface['created_at']}', NULL, {interface['device_id']}, "
            f"'{interface['name']}', '{interface['mac']}', '{interface['ip']}', '{interface['netmask']}');"
        )
    
    sql_lines.extend(['', '-- ============================================', '-- 4. Edges (连接器) - 100个', '-- ============================================'])
    
    # Edges
    for edge in edges:
        sql_lines.append(
            f"INSERT INTO edges (created_at, updated_at, deleted_at, name, status, online, heartbeat_at, description, device_id) VALUES "
            f"('{edge['created_at']}', NOW(), NULL, '{edge['name']}', {edge['status']}, {edge['online']}, "
            f"'{edge['heartbeat_at']}', '{edge['description']}', {edge['device_id']});"
        )
    
    sql_lines.extend(['', '-- ============================================', '-- 5. Access Keys (访问密钥) - 100个', '-- ============================================'])
    
    # Access Keys
    for ak in access_keys:
        sql_lines.append(
            f"INSERT INTO access_keys (created_at, updated_at, deleted_at, edge_id, access_key, secret_key) VALUES "
            f"('{ak['created_at']}', '{ak['created_at']}', NULL, {ak['edge_id']}, '{ak['access_key']}', '{ak['secret_key']}');"
        )
    
    sql_lines.extend(['', '-- ============================================', '-- 6. Applications (应用) - 100个', '-- ============================================'])
    
    # Applications
    for app in applications:
        sql_lines.append(
            f"INSERT INTO applications (created_at, updated_at, deleted_at, edge_ids, device_id, name, description, ip, port, heartbeat_at, application_type) VALUES "
            f"('{app['created_at']}', NOW(), NULL, '{app['edge_ids']}', {app['device_id']}, '{app['name']}', "
            f"'{app['description']}', '{app['ip']}', {app['port']}, '{app['heartbeat_at']}', '{app['application_type']}');"
        )
    
    sql_lines.extend(['', '-- ============================================', '-- 7. Proxies (代理) - 100个', '-- ============================================'])
    
    # Proxies
    for proxy in proxies:
        sql_lines.append(
            f"INSERT INTO proxies (created_at, updated_at, deleted_at, application_id, name, port, status, description) VALUES "
            f"('{proxy['created_at']}', NOW(), NULL, {proxy['application_id']}, '{proxy['name']}', "
            f"{proxy['port']}, {proxy['status']}, '{proxy['description']}');"
        )
    
    return '\n'.join(sql_lines)

def main():
    print("正在生成 Mock 数据...")
    
    users = generate_users(1)  # 只生成1个用户（admin）
    devices = generate_devices(100)
    interfaces = generate_ethernet_interfaces(devices)
    edges = generate_edges(devices, 100)
    access_keys = generate_access_keys(edges)
    applications = generate_applications(devices, edges, 100)
    proxies = generate_proxies(applications)
    
    sql = generate_sql(users, devices, interfaces, edges, access_keys, applications, proxies)
    
    output_file = 'test/mock/mock_data_100.sql'
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(sql)
    
    print(f"✅ 数据生成完成！")
    print(f"   - 用户: {len(users)} 个")
    print(f"   - 设备: {len(devices)} 个")
    print(f"   - 网卡接口: {len(interfaces)} 个")
    print(f"   - 连接器: {len(edges)} 个")
    print(f"   - 访问密钥: {len(access_keys)} 个")
    print(f"   - 应用: {len(applications)} 个")
    print(f"   - 代理: {len(proxies)} 个")
    print(f"\n📄 SQL 文件已保存到: {output_file}")

if __name__ == '__main__':
    main()
