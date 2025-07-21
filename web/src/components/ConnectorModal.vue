<template>
  <el-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="connector ? '编辑连接器' : '添加连接器'"
    width="500px"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="连接器名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入连接器名称" />
      </el-form-item>
      
      <el-form-item label="连接器类型" prop="type">
        <el-select v-model="form.type" placeholder="请选择连接器类型" style="width: 100%">
          <el-option label="HTTP" value="HTTP" />
          <el-option label="WebSocket" value="WebSocket" />
          <el-option label="MQTT" value="MQTT" />
          <el-option label="TCP" value="TCP" />
          <el-option label="UDP" value="UDP" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="连接器状态" prop="status">
        <el-select v-model="form.status" placeholder="请选择连接器状态" style="width: 100%">
          <el-option label="活跃" value="active" />
          <el-option label="非活跃" value="inactive" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="端点地址" prop="endpoint">
        <el-input v-model="form.endpoint" placeholder="请输入端点地址" />
      </el-form-item>
      
      <el-form-item label="端口" prop="port">
        <el-input-number v-model="form.port" :min="1" :max="65535" style="width: 100%" />
      </el-form-item>
      
      <el-form-item label="连接器描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入连接器描述"
        />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="$emit('update:modelValue', false)">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          {{ connector ? '更新' : '创建' }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: Boolean,
  connector: Object
})

const emit = defineEmits(['update:modelValue', 'submit'])

const formRef = ref()
const loading = ref(false)

const form = ref({
  name: '',
  type: '',
  status: 'active',
  endpoint: '',
  port: 8080,
  description: ''
})

const rules = {
  name: [
    { required: true, message: '请输入连接器名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择连接器类型', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择连接器状态', trigger: 'change' }
  ],
  endpoint: [
    { required: true, message: '请输入端点地址', trigger: 'blur' }
  ],
  port: [
    { required: true, message: '请输入端口号', trigger: 'blur' }
  ]
}

watch(() => props.connector, (newConnector) => {
  if (newConnector) {
    form.value = { ...newConnector }
  } else {
    form.value = {
      name: '',
      type: '',
      status: 'active',
      endpoint: '',
      port: 8080,
      description: ''
    }
  }
}, { immediate: true })

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    emit('submit', { ...form.value })
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style> 