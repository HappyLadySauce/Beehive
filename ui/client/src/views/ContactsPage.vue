<template>
  <div class="contacts-page">
    <n-card title="发起单聊">
      <n-form ref="formRef" :model="form" inline>
        <n-form-item label="对方用户 ID">
          <n-input v-model:value="form.otherUserId" placeholder="输入用户 ID，如 u_xxx" style="width: 200px" />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" :loading="creating" @click="createSingleChat">发起会话</n-button>
        </n-form-item>
      </n-form>
      <n-alert v-if="createError" type="error" closable @close="createError = ''">{{ createError }}</n-alert>
    </n-card>
    <n-card title="联系人" style="margin-top: 16px">
      <n-empty description="联系人列表（可输入用户 ID 上方发起单聊）" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NForm, NFormItem, NInput, NButton, NAlert, NEmpty, useMessage } from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { useConversationStore } from '@/stores/conversation'
import { useWebSocket } from '@/composables/useWebSocket'

const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()
const conversationStore = useConversationStore()
const { send, connected } = useWebSocket()

const formRef = ref(null)
const creating = ref(false)
const createError = ref('')
const form = reactive({ otherUserId: '' })

function createSingleChat() {
  const otherId = form.otherUserId.trim()
  if (!otherId) {
    message.warning('请输入对方用户 ID')
    return
  }
  if (otherId === authStore.user?.id) {
    message.warning('不能与自己发起会话')
    return
  }
  if (!connected.value) {
    message.error('未连接，请先登录')
    return
  }
  creating.value = true
  createError.value = ''
  send(
    'conversation.create',
    { type: 'single', memberIds: [authStore.user!.id, otherId] },
    (env) => {
      creating.value = false
      if (env.type === 'conversation.create.ok' && env.payload) {
        const p = env.payload as { conversationId: string }
        conversationStore.setCurrentConversation(p.conversationId)
        message.success('会话已创建')
        router.push(`/app/chats/${p.conversationId}`)
      } else if (env.error) {
        createError.value = env.error.message || '创建会话失败'
      }
    }
  )
}
</script>

<style scoped>
.contacts-page {
  padding: 16px;
}
</style>
