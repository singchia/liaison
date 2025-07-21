<template>
  <el-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="device ? '编辑设备' : '添加设备'"
    width="500px"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="设备名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入设备名称" />
      </el-form-item>
      
      <el-form-item label="设备类型" prop="type">
        <el-select v-model="form.type" placeholder="请选择设备类型" style="width: 100%">
          <el-option label="IoT设备" value="IoT" />
          <el-option label="传感器" value="Sensor" />
          <el-option label="摄像头" value="Camera" />
          <el-option label="控制器" value="Controller" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="设备状态" prop="status">
        <el-select v-model="form.status" placeholder="请选择设备状态" style="width: 100%">
          <el-option label="在线" value="online" />
          <el-option label="离线" value="offline" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="设备位置" prop="location">
        <el-input v-model="form.location" placeholder="请输入设备位置" />
      </el-form-item>
      
      <el-form-item label="设备描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入设备描述"
        />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="$emit('update:modelValue', false)">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          {{ device ? '更新' : '创建' }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'

const props = defineProps({
  modelValue: Boolean,
  device: Object
})

const emit = defineEmits(['update:modelValue', 'submit'])

const formRef = ref()
const loading = ref(false)

const form = ref({
  name: '',
  type: '',
  status: 'online',
  location: '',
  description: ''
})

const rules = {
  name: [
    { required: true, message: '请输入设备名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择设备类型', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择设备状态', trigger: 'change' }
  ],
  location: [
    { required: true, message: '请输入设备位置', trigger: 'blur' }
  ]
}

watch(() => props.device, (newDevice) => {
  if (newDevice) {
    form.value = { ...newDevice }
  } else {
    form.value = {
      name: '',
      type: '',
      status: 'online',
      location: '',
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