<template>
  <div class="login-page">
    <n-card title="登录" style="max-width: 360px">
      <n-form ref="formRef" :model="form" :rules="rules">
        <n-form-item path="username" label="用户名">
          <n-input v-model:value="form.username" placeholder="请输入用户名" />
        </n-form-item>
        <n-form-item path="password" label="密码">
          <n-input v-model:value="form.password" type="password" placeholder="请输入密码" show-password-on="click" />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" block :loading="loading" @click="handleSubmit">登录</n-button>
        </n-form-item>
      </n-form>
      <template #footer>
        <router-link to="/auth/register">注册账号</router-link>
      </template>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import type { FormInst, FormRules } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { useWebSocket } from '@/composables/useWebSocket'

const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()
const { send, ensureConnected } = useWebSocket()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const form = reactive({ username: '', password: '' })
const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

function handleSubmit() {
  formRef.value?.validate((err) => {
    if (err) return
    loading.value = true
    ensureConnected().then(() => {
      send('auth.login', { username: form.username, password: form.password }, (env) => {
        loading.value = false
        if (env.type === 'auth.login.ok' && env.payload) {
          const p = env.payload as { userId: string; accessToken?: string; refreshToken?: string; expiresIn?: number }
          authStore.setFromTokenLogin({
            userId: p.userId,
            accessToken: p.accessToken,
            refreshToken: p.refreshToken,
          })
          message.success('登录成功')
          const redirect = (router.currentRoute.value.query.redirect as string) || '/app/chats'
          router.replace(redirect)
        } else if (env.error) {
          message.error(env.error.message || '登录失败')
        }
      })
    }).catch(() => {
      loading.value = false
      message.error('无法连接服务器，请稍后重试')
    })
  })
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}
</style>
