<template>
  <div class="register-page">
    <n-card title="注册" style="max-width: 360px">
      <n-form ref="formRef" :model="form" :rules="rules">
        <n-form-item path="username" label="用户名">
          <n-input v-model:value="form.username" placeholder="请输入用户名" />
        </n-form-item>
        <n-form-item path="password" label="密码">
          <n-input v-model:value="form.password" type="password" placeholder="请输入密码" show-password-on="click" />
        </n-form-item>
        <n-form-item path="confirmPassword" label="确认密码">
          <n-input v-model:value="form.confirmPassword" type="password" placeholder="再次输入密码" show-password-on="click" />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" block disabled>注册（后端接口就绪后开放）</n-button>
        </n-form-item>
      </n-form>
      <template #footer>
        <router-link to="/auth/login">已有账号？去登录</router-link>
      </template>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { NCard, NForm, NFormItem, NInput, NButton } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'

const formRef = ref<FormInst | null>(null)
const form = reactive({ username: '', password: '', confirmPassword: '' })
const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  confirmPassword: [{ required: true, message: '请确认密码', trigger: 'blur' }],
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}
</style>
