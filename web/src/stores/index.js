import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useProductStore = defineStore('product', () => {
  // 状态
  const devices = ref([
    { id: 1, name: 'IoT传感器01', type: 'IoT', status: 'online', location: '北京', created_at: '2024-01-01' },
    { id: 2, name: '温度传感器02', type: 'Sensor', status: 'offline', location: '上海', created_at: '2024-01-02' },
    { id: 3, name: '摄像头03', type: 'Camera', status: 'online', location: '深圳', created_at: '2024-01-03' }
  ])
  const connectors = ref([
    { id: 1, name: 'HTTP网关', type: 'HTTP', status: 'active', endpoint: 'http://edge1.example.com', created_at: '2024-01-01' },
    { id: 2, name: 'WebSocket连接器', type: 'WebSocket', status: 'inactive', endpoint: 'ws://edge2.example.com', created_at: '2024-01-02' },
    { id: 3, name: 'MQTT代理', type: 'MQTT', status: 'active', endpoint: 'mqtt://edge3.example.com', created_at: '2024-01-03' }
  ])
  const applications = ref([
    { id: 1, name: '监控系统', type: 'Web', status: 'running', version: '1.0.0', created_at: '2024-01-01' },
    { id: 2, name: '移动应用', type: 'Mobile', status: 'stopped', version: '2.0.0', created_at: '2024-01-02' },
    { id: 3, name: '数据分析服务', type: 'Microservice', status: 'running', version: '1.5.0', created_at: '2024-01-03' }
  ])
  const proxies = ref([
    { id: 1, name: '负载均衡器', type: 'HTTP', status: 'active', target: 'http://target1.com', created_at: '2024-01-01' },
    { id: 2, name: 'TCP代理', type: 'TCP', status: 'inactive', target: 'tcp://target2.com:8080', created_at: '2024-01-02' },
    { id: 3, name: 'HTTPS代理', type: 'HTTPS', status: 'active', target: 'https://target3.com', created_at: '2024-01-03' }
  ])
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
      // 模拟API调用延迟
      await new Promise(resolve => setTimeout(resolve, 500))
      // 保持现有数据，不重新获取
    } catch (error) {
      console.error('获取设备列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  const fetchConnectors = async () => {
    loading.value = true
    try {
      await new Promise(resolve => setTimeout(resolve, 500))
    } catch (error) {
      console.error('获取连接器列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  const fetchApplications = async () => {
    loading.value = true
    try {
      await new Promise(resolve => setTimeout(resolve, 500))
    } catch (error) {
      console.error('获取应用列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  const fetchProxies = async () => {
    loading.value = true
    try {
      await new Promise(resolve => setTimeout(resolve, 500))
    } catch (error) {
      console.error('获取代理列表失败:', error)
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