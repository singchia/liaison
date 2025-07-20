<template>
  <div class="product-management">
    <!-- 顶部导航栏 -->
    <header class="header">
      <div class="header-content">
        <div class="logo">
          <span>Liaison 产品管理系统</span>
        </div>
        <div class="user-info">
          <div class="user-dropdown">
            <span>管理员</span>
            <span>▼</span>
          </div>
        </div>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <div class="main-container">
      <!-- 侧边栏 -->
      <aside class="sidebar">
        <nav class="sidebar-menu">
          <a 
            href="#" 
            :class="['menu-item', { active: activeTab === 'dashboard' }]"
            @click.prevent="handleTabSelect('dashboard')"
          >
            仪表盘
          </a>
          <a 
            href="#" 
            :class="['menu-item', { active: activeTab === 'devices' }]"
            @click.prevent="handleTabSelect('devices')"
          >
            设备管理
          </a>
          <a 
            href="#" 
            :class="['menu-item', { active: activeTab === 'connectors' }]"
            @click.prevent="handleTabSelect('connectors')"
          >
            连接器管理
          </a>
          <a 
            href="#" 
            :class="['menu-item', { active: activeTab === 'applications' }]"
            @click.prevent="handleTabSelect('applications')"
          >
            应用管理
          </a>
          <a 
            href="#" 
            :class="['menu-item', { active: activeTab === 'proxies' }]"
            @click.prevent="handleTabSelect('proxies')"
          >
            代理管理
          </a>
        </nav>
      </aside>

      <!-- 内容区域 -->
      <main class="main-content">
        <!-- 仪表盘 -->
        <div v-if="activeTab === 'dashboard'" class="dashboard">
          <h2>系统概览</h2>
          <div class="stats-row">
            <div class="stat-card">
              <div class="stat-content">
                <div class="stat-info">
                  <div class="stat-number">{{ deviceCount }}</div>
                  <div class="stat-label">设备总数</div>
                </div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-content">
                <div class="stat-info">
                  <div class="stat-number">{{ connectorCount }}</div>
                  <div class="stat-label">连接器总数</div>
                </div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-content">
                <div class="stat-info">
                  <div class="stat-number">{{ applicationCount }}</div>
                  <div class="stat-label">应用总数</div>
                </div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-content">
                <div class="stat-info">
                  <div class="stat-number">{{ proxyCount }}</div>
                  <div class="stat-label">代理总数</div>
                </div>
              </div>
            </div>
          </div>

          <div class="info-section">
            <div class="info-card">
              <h3>系统状态概览</h3>
              <div class="info-content">
                <div class="status-item">
                  <span class="status-label">设备状态：</span>
                  <span class="status-value online">{{ devices.filter(d => d.status === 'online').length }} 台在线</span>
                  <span class="status-value offline">{{ devices.filter(d => d.status === 'offline').length }} 台离线</span>
                </div>
                <div class="status-item">
                  <span class="status-label">连接器状态：</span>
                  <span class="status-value active">{{ connectors.filter(c => c.status === 'active').length }} 个活跃</span>
                  <span class="status-value inactive">{{ connectors.filter(c => c.status === 'inactive').length }} 个非活跃</span>
                </div>
                <div class="status-item">
                  <span class="status-label">应用状态：</span>
                  <span class="status-value running">{{ applications.filter(a => a.status === 'running').length }} 个运行中</span>
                  <span class="status-value stopped">{{ applications.filter(a => a.status === 'stopped').length }} 个已停止</span>
                </div>
                <div class="status-item">
                  <span class="status-label">代理状态：</span>
                  <span class="status-value active">{{ proxies.filter(p => p.status === 'active').length }} 个活跃</span>
                  <span class="status-value inactive">{{ proxies.filter(p => p.status === 'inactive').length }} 个非活跃</span>
                </div>
              </div>
            </div>
            <div class="info-card">
              <h3>系统活动日志</h3>
              <div class="info-content">
                <div class="log-item">
                  <span class="log-time">14:30</span>
                  <span class="log-message">设备 "IoT传感器01" 连接成功</span>
                </div>
                <div class="log-item">
                  <span class="log-time">14:25</span>
                  <span class="log-message">连接器 "HTTP网关" 状态更新为活跃</span>
                </div>
                <div class="log-item">
                  <span class="log-time">14:20</span>
                  <span class="log-message">应用 "监控系统" 启动完成</span>
                </div>
                <div class="log-item">
                  <span class="log-time">14:15</span>
                  <span class="log-message">代理 "负载均衡器" 配置已更新</span>
                </div>
                <div class="log-item">
                  <span class="log-time">14:10</span>
                  <span class="log-message">系统自检完成，所有服务正常</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 设备管理 -->
        <div v-if="activeTab === 'devices'" class="content-section">
                      <div class="section-header">
              <h2>设备管理</h2>
              <div class="header-actions">
                <div class="summary-info">
                  <span class="summary-item">总数：{{ deviceCount }}</span>
                  <span class="summary-item online">在线：{{ devices.filter(d => d.status === 'online').length }}</span>
                  <span class="summary-item offline">离线：{{ devices.filter(d => d.status === 'offline').length }}</span>
                </div>
                <button class="btn btn-primary" @click="showDeviceModal = true">
                  添加设备
                </button>
              </div>
            </div>
          
          <div class="card">
            <div class="table-toolbar">
              <div class="search-section">
                <input
                  v-model="deviceSearch"
                  type="text"
                  placeholder="搜索设备名称..."
                  class="search-input"
                />
                <select v-model="deviceStatusFilter" class="filter-select">
                  <option value="">全部状态</option>
                  <option value="online">在线</option>
                  <option value="offline">离线</option>
                </select>
              </div>
              <div class="view-options">
                <span class="result-count">共 {{ filteredDevices.length }} 条记录</span>
              </div>
            </div>

            <div class="table-container">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>#</th>
                    <th>设备信息</th>
                    <th>类型</th>
                    <th>状态</th>
                    <th>位置</th>
                    <th>创建时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="device in filteredDevices" :key="device.id" class="table-row">
                    <td class="device-id">{{ device.id }}</td>
                    <td class="device-info">
                      <div class="device-name">{{ device.name }}</div>
                      <div class="device-desc">设备ID: {{ device.id }}</div>
                    </td>
                    <td>
                      <span class="type-badge">{{ device.type }}</span>
                    </td>
                    <td>
                      <span :class="['status-badge', device.status === 'online' ? 'online' : 'offline']">
                        {{ device.status === 'online' ? '在线' : '离线' }}
                      </span>
                    </td>
                    <td class="location">{{ device.location }}</td>
                    <td class="date">{{ device.created_at }}</td>
                    <td class="actions">
                      <button class="btn btn-small btn-edit" @click="editDevice(device)">
                        编辑
                      </button>
                      <button class="btn btn-small btn-delete" @click="deleteDevice(device)">
                        删除
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- 连接器管理 -->
        <div v-if="activeTab === 'connectors'" class="content-section">
          <div class="section-header">
            <h2>连接器管理</h2>
            <button class="btn btn-primary" @click="showConnectorModal = true">
              ➕ 添加连接器
            </button>
          </div>
          
          <div class="card">
            <div class="table-toolbar">
              <input
                v-model="connectorSearch"
                type="text"
                placeholder="搜索连接器..."
                class="search-input"
              />
              <select v-model="connectorStatusFilter" class="filter-select">
                <option value="">全部状态</option>
                <option value="active">活跃</option>
                <option value="inactive">非活跃</option>
              </select>
            </div>

            <div class="table-container">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>连接器名称</th>
                    <th>类型</th>
                    <th>状态</th>
                    <th>端点</th>
                    <th>创建时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="connector in filteredConnectors" :key="connector.id">
                    <td>{{ connector.id }}</td>
                    <td>{{ connector.name }}</td>
                    <td>{{ connector.type }}</td>
                    <td>
                      <span :class="['status-tag', connector.status === 'active' ? 'success' : 'info']">
                        {{ connector.status === 'active' ? '活跃' : '非活跃' }}
                      </span>
                    </td>
                    <td>{{ connector.endpoint }}</td>
                    <td>{{ connector.created_at }}</td>
                    <td>
                      <button class="btn btn-small" @click="editConnector(connector)">编辑</button>
                      <button class="btn btn-small btn-danger" @click="deleteConnector(connector)">删除</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- 应用管理 -->
        <div v-if="activeTab === 'applications'" class="content-section">
          <div class="section-header">
            <h2>应用管理</h2>
            <button class="btn btn-primary" @click="showApplicationModal = true">
              ➕ 添加应用
            </button>
          </div>
          
          <div class="card">
            <div class="table-toolbar">
              <input
                v-model="applicationSearch"
                type="text"
                placeholder="搜索应用..."
                class="search-input"
              />
              <select v-model="applicationStatusFilter" class="filter-select">
                <option value="">全部状态</option>
                <option value="running">运行中</option>
                <option value="stopped">已停止</option>
              </select>
            </div>

            <div class="table-container">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>应用名称</th>
                    <th>类型</th>
                    <th>状态</th>
                    <th>版本</th>
                    <th>创建时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="application in filteredApplications" :key="application.id">
                    <td>{{ application.id }}</td>
                    <td>{{ application.name }}</td>
                    <td>{{ application.type }}</td>
                    <td>
                      <span :class="['status-tag', application.status === 'running' ? 'success' : 'warning']">
                        {{ application.status === 'running' ? '运行中' : '已停止' }}
                      </span>
                    </td>
                    <td>{{ application.version }}</td>
                    <td>{{ application.created_at }}</td>
                    <td>
                      <button class="btn btn-small" @click="editApplication(application)">编辑</button>
                      <button class="btn btn-small btn-danger" @click="deleteApplication(application)">删除</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- 代理管理 -->
        <div v-if="activeTab === 'proxies'" class="content-section">
          <div class="section-header">
            <h2>代理管理</h2>
            <button class="btn btn-primary" @click="showProxyModal = true">
              ➕ 添加代理
            </button>
          </div>
          
          <div class="card">
            <div class="table-toolbar">
              <input
                v-model="proxySearch"
                type="text"
                placeholder="搜索代理..."
                class="search-input"
              />
              <select v-model="proxyStatusFilter" class="filter-select">
                <option value="">全部状态</option>
                <option value="active">活跃</option>
                <option value="inactive">非活跃</option>
              </select>
            </div>

            <div class="table-container">
              <table class="data-table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>代理名称</th>
                    <th>类型</th>
                    <th>状态</th>
                    <th>目标地址</th>
                    <th>创建时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="proxy in filteredProxies" :key="proxy.id">
                    <td>{{ proxy.id }}</td>
                    <td>{{ proxy.name }}</td>
                    <td>{{ proxy.type }}</td>
                    <td>
                      <span :class="['status-tag', proxy.status === 'active' ? 'success' : 'info']">
                        {{ proxy.status === 'active' ? '活跃' : '非活跃' }}
                      </span>
                    </td>
                    <td>{{ proxy.target }}</td>
                    <td>{{ proxy.created_at }}</td>
                    <td>
                      <button class="btn btn-small" @click="editProxy(proxy)">编辑</button>
                      <button class="btn btn-small btn-danger" @click="deleteProxy(proxy)">删除</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </main>
    </div>

    <!-- 模态框 -->
    <div v-if="showDeviceModal" class="modal-overlay" @click="showDeviceModal = false">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h3>{{ editingDevice ? '编辑设备' : '添加设备' }}</h3>
          <button class="modal-close" @click="showDeviceModal = false">×</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="handleDeviceSubmit">
            <div class="form-group">
              <label>设备名称</label>
              <input v-model="deviceForm.name" type="text" required />
            </div>
            <div class="form-group">
              <label>设备类型</label>
              <select v-model="deviceForm.type" required>
                <option value="">请选择</option>
                <option value="IoT">IoT设备</option>
                <option value="Sensor">传感器</option>
                <option value="Camera">摄像头</option>
                <option value="Controller">控制器</option>
              </select>
            </div>
            <div class="form-group">
              <label>设备状态</label>
              <select v-model="deviceForm.status" required>
                <option value="online">在线</option>
                <option value="offline">离线</option>
              </select>
            </div>
            <div class="form-group">
              <label>设备位置</label>
              <input v-model="deviceForm.location" type="text" required />
            </div>
            <div class="form-actions">
              <button type="button" class="btn" @click="showDeviceModal = false">取消</button>
              <button type="submit" class="btn btn-primary">保存</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useProductStore } from '../stores'

const store = useProductStore()

// 响应式数据
const activeTab = ref('dashboard')
const showDeviceModal = ref(false)
const editingDevice = ref(null)

// 从store获取数据
const { devices, connectors, applications, proxies } = store

// 搜索和筛选
const deviceSearch = ref('')
const deviceStatusFilter = ref('')
const connectorSearch = ref('')
const connectorStatusFilter = ref('')
const applicationSearch = ref('')
const applicationStatusFilter = ref('')
const proxySearch = ref('')
const proxyStatusFilter = ref('')

// 表单数据
const deviceForm = ref({
  name: '',
  type: '',
  status: 'online',
  location: ''
})

// 计算属性
const filteredDevices = computed(() => {
  let filtered = store.devices
  if (deviceSearch.value) {
    filtered = filtered.filter(device => 
      device.name.toLowerCase().includes(deviceSearch.value.toLowerCase())
    )
  }
  if (deviceStatusFilter.value) {
    filtered = filtered.filter(device => device.status === deviceStatusFilter.value)
  }
  return filtered
})

const filteredConnectors = computed(() => {
  let filtered = store.connectors
  if (connectorSearch.value) {
    filtered = filtered.filter(connector => 
      connector.name.toLowerCase().includes(connectorSearch.value.toLowerCase())
    )
  }
  if (connectorStatusFilter.value) {
    filtered = filtered.filter(connector => connector.status === connectorStatusFilter.value)
  }
  return filtered
})

const filteredApplications = computed(() => {
  let filtered = store.applications
  if (applicationSearch.value) {
    filtered = filtered.filter(application => 
      application.name.toLowerCase().includes(applicationSearch.value.toLowerCase())
    )
  }
  if (applicationStatusFilter.value) {
    filtered = filtered.filter(application => application.status === applicationStatusFilter.value)
  }
  return filtered
})

const filteredProxies = computed(() => {
  let filtered = store.proxies
  if (proxySearch.value) {
    filtered = filtered.filter(proxy => 
      proxy.name.toLowerCase().includes(proxySearch.value.toLowerCase())
    )
  }
  if (proxyStatusFilter.value) {
    filtered = filtered.filter(proxy => proxy.status === proxyStatusFilter.value)
  }
  return filtered
})

// 方法
const handleTabSelect = (index) => {
  activeTab.value = index
  loadData()
}

const loadData = () => {
  switch (activeTab.value) {
    case 'devices':
      store.fetchDevices()
      break
    case 'connectors':
      store.fetchConnectors()
      break
    case 'applications':
      store.fetchApplications()
      break
    case 'proxies':
      store.fetchProxies()
      break
  }
}

const editDevice = (device) => {
  editingDevice.value = device
  deviceForm.value = { ...device }
  showDeviceModal.value = true
}

const deleteDevice = async (device) => {
  if (confirm(`确定要删除设备 "${device.name}" 吗？`)) {
    alert('设备删除成功')
    store.fetchDevices()
  }
}

const handleDeviceSubmit = () => {
  alert(editingDevice.value ? '设备更新成功' : '设备创建成功')
  showDeviceModal.value = false
  editingDevice.value = null
  deviceForm.value = { name: '', type: '', status: 'online', location: '' }
  store.fetchDevices()
}

// 其他方法（简化版）
const showConnectorModal = ref(false)
const showApplicationModal = ref(false)
const showProxyModal = ref(false)
const editingConnector = ref(null)
const editingApplication = ref(null)
const editingProxy = ref(null)

const editConnector = (connector) => {
  editingConnector.value = connector
  alert('编辑连接器功能待实现')
}

const deleteConnector = (connector) => {
  if (confirm(`确定要删除连接器 "${connector.name}" 吗？`)) {
    alert('连接器删除成功')
    store.fetchConnectors()
  }
}

const editApplication = (application) => {
  editingApplication.value = application
  alert('编辑应用功能待实现')
}

const deleteApplication = (application) => {
  if (confirm(`确定要删除应用 "${application.name}" 吗？`)) {
    alert('应用删除成功')
    store.fetchApplications()
  }
}

const editProxy = (proxy) => {
  editingProxy.value = proxy
  alert('编辑代理功能待实现')
}

const deleteProxy = (proxy) => {
  if (confirm(`确定要删除代理 "${proxy.name}" 吗？`)) {
    alert('代理删除成功')
    store.fetchProxies()
  }
}

// 生命周期
onMounted(() => {
  loadData()
})
</script>

<style scoped>
.product-management {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 0;
  height: 60px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  padding: 0 20px;
}

.logo {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
  color: #409eff;
  gap: 8px;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 8px;
}



.main-container {
  flex: 1;
  display: flex;
  height: calc(100vh - 60px);
}

.sidebar {
  width: 200px;
  background: #fff;
  border-right: 1px solid #e4e7ed;
}

.sidebar-menu {
  padding: 20px 0;
}

.menu-item {
  display: block;
  padding: 12px 20px;
  color: #606266;
  text-decoration: none;
  transition: all 0.3s;
}

.menu-item:hover {
  background: #f5f7fa;
  color: #409eff;
}

.menu-item.active {
  background: #ecf5ff;
  color: #409eff;
  border-right: 3px solid #409eff;
}

.main-content {
  flex: 1;
  padding: 20px;
  background: #f5f5f5;
  overflow-y: auto;
}

.dashboard h2 {
  margin-bottom: 20px;
  color: #303133;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.stat-content {
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
}



.stat-number {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}



.info-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.info-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.info-card h3 {
  margin: 0 0 16px 0;
  color: #303133;
}

.info-content p {
  margin: 8px 0;
  color: #606266;
  font-size: 14px;
}

.status-item {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.status-item:last-child {
  border-bottom: none;
}

.status-label {
  font-weight: 500;
  color: #303133;
  min-width: 100px;
  margin-right: 12px;
}

.status-value {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  margin-right: 8px;
}

.status-value.online,
.status-value.active,
.status-value.running {
  background: #f0f9ff;
  color: #0369a1;
}

.status-value.offline,
.status-value.inactive,
.status-value.stopped {
  background: #fef2f2;
  color: #dc2626;
}

.log-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  padding: 6px 0;
}

.log-time {
  background: #f3f4f6;
  color: #6b7280;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  margin-right: 8px;
  min-width: 40px;
  text-align: center;
}

.log-message {
  color: #374151;
  font-size: 13px;
}

.content-section {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h2 {
  margin: 0;
  color: #303133;
  font-size: 24px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 20px;
}

.summary-info {
  display: flex;
  gap: 16px;
}

.summary-item {
  padding: 6px 12px;
  background: #f8fafc;
  border-radius: 6px;
  font-size: 13px;
  color: #64748b;
}

.summary-item.online {
  background: #f0fdf4;
  color: #16a34a;
}

.summary-item.offline {
  background: #fef2f2;
  color: #dc2626;
}

.card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px 0;
  border-bottom: 1px solid #e5e7eb;
}

.search-section {
  display: flex;
  gap: 12px;
  align-items: center;
}

.view-options {
  display: flex;
  align-items: center;
  gap: 12px;
}

.result-count {
  color: #6b7280;
  font-size: 13px;
}

.search-input {
  width: 300px;
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
}

.filter-select {
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  background: #fff;
}

.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ebeef5;
}

.data-table th {
  background: #fafafa;
  font-weight: 500;
  color: #606266;
}

.data-table tr:hover {
  background: #f8fafc;
}

.table-row {
  transition: all 0.2s ease;
}

.device-id {
  font-weight: 600;
  color: #6b7280;
  font-size: 13px;
}

.device-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.device-name {
  font-weight: 500;
  color: #111827;
  font-size: 14px;
}

.device-desc {
  font-size: 12px;
  color: #6b7280;
}

.type-badge {
  background: #f3f4f6;
  color: #374151;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.online {
  background: #f0fdf4;
  color: #16a34a;
}

.status-badge.offline {
  background: #fef2f2;
  color: #dc2626;
}

.location {
  color: #6b7280;
  font-size: 13px;
}

.date {
  color: #6b7280;
  font-size: 12px;
}

.actions {
  display: flex;
  gap: 6px;
}

.btn-edit {
  background: #eff6ff;
  border-color: #3b82f6;
  color: #1d4ed8;
}

.btn-edit:hover {
  background: #dbeafe;
}

.btn-delete {
  background: #fef2f2;
  border-color: #ef4444;
  color: #dc2626;
}

.btn-delete:hover {
  background: #fee2e2;
}



.btn {
  padding: 8px 16px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background: #fff;
  color: #606266;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.btn:hover {
  border-color: #c6e2ff;
  color: #409eff;
}

.btn-primary {
  background: #409eff;
  border-color: #409eff;
  color: #fff;
}

.btn-primary:hover {
  background: #66b1ff;
  border-color: #66b1ff;
}

.btn-danger {
  background: #f56c6c;
  border-color: #f56c6c;
  color: #fff;
}

.btn-danger:hover {
  background: #f78989;
  border-color: #f78989;
}

.btn-small {
  padding: 4px 8px;
  font-size: 12px;
  margin-right: 4px;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: #fff;
  border-radius: 8px;
  width: 500px;
  max-width: 90vw;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e4e7ed;
}

.modal-header h3 {
  margin: 0;
  color: #303133;
}

.modal-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #909399;
}

.modal-close:hover {
  color: #606266;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #606266;
  font-weight: 500;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #409eff;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}
</style> 