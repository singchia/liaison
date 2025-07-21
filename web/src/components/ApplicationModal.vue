<template>
  <el-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="application ? '编辑应用' : '添加应用'"
    width="500px"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="应用名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入应用名称" />
      </el-form-item>
      
      <el-form-item label="应用类型" prop="type">
        <el-select v-model="form.type" placeholder="请选择应用类型" style="width: 100%">
          <el-option label="Web应用" value="Web" />
          <el-option label="移动应用" value="Mobile" />
          <el-option label="桌面应用" value="Desktop" />
          <el-option label="微服务" value="Microservice" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="应用状态" prop="status">
        <el-select v-model="form.status" placeholder="请选择应用状态" style="width: 100%">
          <el-option label="运行中" value="running" />
          <el-option label="已停止" value="stopped" />
          <el-option label="维护中" value="maintenance" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="版本号" prop="version">
        <el-input v-model="form.version" placeholder="请输入版本号" />
      </el-form-item>
      
      <el-form-item label="部署环境" prop="environment">
        <el-select v-model="form.environment" placeholder="请选择部署环境" style="width: 100%">
          <el-option label="开发环境" value="dev" />
          <el-option label="测试环境" value="test" />
          <el-option label="生产环境" value="prod" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="应用描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入应用描述"
        />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="$emit('update:modelValue', false)">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          {{ application ? '更新' : '创建' }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: Boolean,
  application: Object
})

const emit = defineEmits(['update:modelValue', 'submit'])

const formRef = ref()
const loading = ref(false)

const form = ref({
  name: '',
  type: '',
  status: 'running',
  version: '',
  environment: 'dev',
  description: ''
})

const rules = {
  name: [
    { required: true, message: '请输入应用名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择应用类型', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择应用状态', trigger: 'change' }
  ],
  version: [
    { required: true, message: '请输入版本号', trigger: 'blur' }
  ],
  environment: [
    { required: true, message: '请选择部署环境', trigger: 'change' }
  ]
}

watch(() => props.application, (newApplication) => {
  if (newApplication) {
    form.value = { ...newApplication }
  } else {
    form.value = {
      name: '',
      type: '',
      status: 'running',
      version: '',
      environment: 'dev',
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