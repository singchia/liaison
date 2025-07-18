import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useProductStore = defineStore('product', () => {
  // 状态
  const devices = ref([])
  const connectors = ref([])
  const applications = ref([])
  const proxies = ref([])
  const loading = ref(false)

  // 计算属性
  const deviceCount = computed(() => devices.value.length)
  const connectorCount = computed(() => connectors.value.length)
  const applicationCount = computed(() => applications.value.length)
  const proxyCount = computed(() => proxies.value.length)

  // 方法
  const fetchDevices = async () => {
    loading.value = true
    try {
      // 模拟API调用
      const response = await fetch('/api/v1/devices')
      devices.value = await response.json()
    } catch (error) {
      console.error('获取设备列表失败:', error)
      // 使用模拟数据
      devices.value = [
        { id: 1, name: '设备1', type: 'IoT', status: 'online', location: '北京', created_at: '2024-01-01' },
        { id: 2, name: '设备2', type: 'Sensor', status: 'offline', location: '上海', created_at: '2024-01-02' }
      ]
    } finally {
      loading.value = false
    }
  }

  const fetchConnectors = async () => {
    loading.value = true
    try {
      const response = await fetch('/api/v1/edges')
      connectors.value = await response.json()
    } catch (error) {
      console.error('获取连接器列表失败:', error)
      connectors.value = [
        { id: 1, name: '连接器1', type: 'HTTP', status: 'active', endpoint: 'http://edge1.example.com', created_at: '2024-01-01' },
        { id: 2, name: '连接器2', type: 'WebSocket', status: 'inactive', endpoint: 'ws://edge2.example.com', created_at: '2024-01-02' }
      ]
    } finally {
      loading.value = false
    }
  }

  const fetchApplications = async () => {
    loading.value = true
    try {
      const response = await fetch('/api/v1/applications')
      applications.value = await response.json()
    } catch (error) {
      console.error('获取应用列表失败:', error)
      applications.value = [
        { id: 1, name: '应用1', type: 'Web', status: 'running', version: '1.0.0', created_at: '2024-01-01' },
        { id: 2, name: '应用2', type: 'Mobile', status: 'stopped', version: '2.0.0', created_at: '2024-01-02' }
      ]
    } finally {
      loading.value = false
    }
  }

  const fetchProxies = async () => {
    loading.value = true
    try {
      const response = await fetch('/api/v1/proxies')
      proxies.value = await response.json()
    } catch (error) {
      console.error('获取代理列表失败:', error)
      proxies.value = [
        { id: 1, name: '代理1', type: 'HTTP', status: 'active', target: 'http://target1.com', created_at: '2024-01-01' },
        { id: 2, name: '代理2', type: 'TCP', status: 'inactive', target: 'tcp://target2.com:8080', created_at: '2024-01-02' }
      ]
    } finally {
      loading.value = false
    }
  }

  return {
    devices,
    connectors,
    applications,
    proxies,
    loading,
    deviceCount,
    connectorCount,
    applicationCount,
    proxyCount,
    fetchDevices,
    fetchConnectors,
    fetchApplications,
    fetchProxies
  }
}) 