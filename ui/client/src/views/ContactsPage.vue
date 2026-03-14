<template>
  <div class="contacts-page">
    <n-card title="发起单聊">
      <n-form ref="formRef" :model="form" inline>
        <n-form-item label="用户名或 10 位账号">
          <n-input
            v-model:value="form.toUsernameOrAccount"
            placeholder="输入用户名或 10 位数字账号"
            style="width: 220px"
          />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" :loading="creating" @click="createSingleChat">发起会话</n-button>
        </n-form-item>
      </n-form>
      <n-alert v-if="createError" type="error" closable @close="createError = ''">{{ createError }}</n-alert>
    </n-card>

    <n-card title="添加联系人" style="margin-top: 16px">
      <n-form :model="addForm" inline>
        <n-form-item label="用户名或 10 位账号">
          <n-input
            v-model:value="addForm.toUsernameOrAccount"
            placeholder="输入对方用户名或 10 位账号"
            style="width: 220px"
          />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" :loading="adding" @click="addContact">添加</n-button>
        </n-form-item>
      </n-form>
      <n-alert v-if="addError" type="error" closable @close="addError = ''">{{ addError }}</n-alert>
    </n-card>

    <n-card title="联系人" style="margin-top: 16px">
      <n-spin :show="loadingContacts">
        <div v-if="contactList.length === 0 && !loadingContacts" class="empty-hint">
          <n-empty description="暂无联系人，可上方添加或发起单聊" />
        </div>
        <ul v-else class="contact-list">
          <li v-for="c in contactList" :key="c.userId" class="contact-item">
            <span class="contact-id">{{ c.userId }}</span>
            <n-button size="small" @click="startChatWith(c.userId)">发消息</n-button>
            <n-button size="small" type="error" quaternary @click="removeContact(c.userId)">移除</n-button>
          </li>
        </ul>
      </n-spin>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  NCard,
  NForm,
  NFormItem,
  NInput,
  NButton,
  NAlert,
  NEmpty,
  NSpin,
  useMessage,
} from 'naive-ui'
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
const form = reactive({ toUsernameOrAccount: '' })

const addForm = reactive({ toUsernameOrAccount: '' })
const adding = ref(false)
const addError = ref('')

const contactList = ref<{ userId: string }[]>([])
const loadingContacts = ref(false)

function isTenDigitAccount(s: string): boolean {
  return /^\d{10}$/.test(s)
}

function createSingleChat() {
  const raw = form.toUsernameOrAccount.trim()
  if (!raw) {
    message.warning('请输入用户名或 10 位账号')
    return
  }
  if (!connected.value) {
    message.error('未连接，请先登录')
    return
  }
  creating.value = true
  createError.value = ''
  const payload: { type: string; toUsername?: string; toAccount?: string } = { type: 'single' }
  if (isTenDigitAccount(raw)) {
    payload.toAccount = raw
  } else {
    payload.toUsername = raw
  }
  send('conversation.create', payload, (env) => {
    creating.value = false
    if (env.type === 'conversation.create.ok' && env.payload) {
      const p = env.payload as { conversationId: string }
      conversationStore.setCurrentConversation(p.conversationId)
      message.success('会话已创建')
      router.push(`/app/chats/${p.conversationId}`)
    } else if (env.error) {
      createError.value = env.error.message || '创建会话失败'
    }
  })
}

function addContact() {
  const raw = addForm.toUsernameOrAccount.trim()
  if (!raw) {
    message.warning('请输入对方用户名或 10 位账号')
    return
  }
  if (!connected.value) {
    message.error('未连接，请先登录')
    return
  }
  adding.value = true
  addError.value = ''
  const payload: { toUserId?: string; toUsername?: string; toAccount?: string } = {}
  if (isTenDigitAccount(raw)) {
    payload.toAccount = raw
  } else {
    payload.toUsername = raw
  }
  send('contact.add', payload, (env) => {
    adding.value = false
    if (env.type === 'contact.add.ok') {
      message.success('已添加联系人')
      addForm.toUsernameOrAccount = ''
      loadContacts()
    } else if (env.error) {
      addError.value = env.error.message || '添加失败'
    }
  })
}

function loadContacts() {
  if (!connected.value) return
  loadingContacts.value = true
  send('contact.list', {}, (env) => {
    loadingContacts.value = false
    if (env.type === 'contact.list.ok' && env.payload) {
      const ids = (env.payload as { contactUserIds?: string[] }).contactUserIds || []
      contactList.value = ids.map((userId) => ({ userId }))
    }
  })
}

function startChatWith(userId: string) {
  if (!connected.value) return
  send(
    'conversation.create',
    { type: 'single', memberIds: [authStore.user!.id, userId] },
    (env) => {
      if (env.type === 'conversation.create.ok' && env.payload) {
        const p = env.payload as { conversationId: string }
        conversationStore.setCurrentConversation(p.conversationId)
        router.push(`/app/chats/${p.conversationId}`)
      } else if (env.error) {
        message.error(env.error.message || '发起会话失败')
      }
    }
  )
}

function removeContact(userId: string) {
  if (!connected.value) return
  send('contact.remove', { contactUserId: userId }, (env) => {
    if (env.type === 'contact.remove.ok') {
      message.success('已移除')
      loadContacts()
    } else if (env.error) {
      message.error(env.error.message || '移除失败')
    }
  })
}

onMounted(() => {
  loadContacts()
})
</script>

<style scoped>
.contacts-page {
  padding: 16px;
}
.empty-hint {
  padding: 24px 0;
}
.contact-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.contact-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--n-border-color);
}
.contact-item:last-child {
  border-bottom: none;
}
.contact-id {
  flex: 1;
  font-family: monospace;
}
</style>
