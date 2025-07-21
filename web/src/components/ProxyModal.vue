<template>
  <el-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="proxy ? '编辑代理' : '添加代理'"
    width="500px"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="代理名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入代理名称" />
      </el-form-item>
      
      <el-form-item label="代理类型" prop="type">
        <el-select v-model="form.type" placeholder="请选择代理类型" style="width: 100%">
          <el-option label="HTTP代理" value="HTTP" />
          <el-option label="HTTPS代理" value="HTTPS" />
          <el-option label="TCP代理" value="TCP" />
          <el-option label="UDP代理" value="UDP" />
          <el-option label="SOCKS代理" value="SOCKS" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="代理状态" prop="status">
        <el-select v-model="form.status" placeholder="请选择代理状态" style="width: 100%">
          <el-option label="活跃" value="active" />
          <el-option label="非活跃" value="inactive" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="目标地址" prop="target">
        <el-input v-model="form.target" placeholder="请输入目标地址" />
      </el-form-item>
      
      <el-form-item label="监听端口" prop="port">
        <el-input-number v-model="form.port" :min="1" :max="65535" style="width: 100%" />
      </el-form-item>
      
      <el-form-item label="认证方式" prop="auth_type">
        <el-select v-model="form.auth_type" placeholder="请选择认证方式" style="width: 100%">
          <el-option label="无认证" value="none" />
          <el-option label="用户名密码" value="basic" />
          <el-option label="Token" value="token" />
        </el-select>
      </el-form-item>
      
      <el-form-item label="代理描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入代理描述"
        />
      </el-form-item>
    </el-form>
    
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="$emit('update:modelValue', false)">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          {{ proxy ? '更新' : '创建' }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: Boolean,
  proxy: Object
})

const emit = defineEmits(['update:modelValue', 'submit'])

const formRef = ref()
const loading = ref(false)

const form = ref({
  name: '',
  type: '',
  status: 'active',
  target: '',
  port: 8080,
  auth_type: 'none',
  description: ''
})

const rules = {
  name: [
    { required: true, message: '请输入代理名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择代理类型', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择代理状态', trigger: 'change' }
  ],
  target: [
    { required: true, message: '请输入目标地址', trigger: 'blur' }
  ],
  port: [
    { required: true, message: '请输入监听端口', trigger: 'blur' }
  ],
  auth_type: [
    { required: true, message: '请选择认证方式', trigger: 'change' }
  ]
}

watch(() => props.proxy, (newProxy) => {
  if (newProxy) {
    form.value = { ...newProxy }
  } else {
    form.value = {
      name: '',
      type: '',
      status: 'active',
      target: '',
      port: 8080,
      auth_type: 'none',
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