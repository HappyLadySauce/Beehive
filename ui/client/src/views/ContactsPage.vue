<template>
  <div class="contacts-page">
    <n-card title="好友与群通知" style="margin-bottom: 16px">
      <n-space>
        <n-button quaternary @click="showFriendRequests = true">好友通知</n-button>
        <n-button quaternary @click="showGroupRequests = true">群通知</n-button>
      </n-space>
    </n-card>
    <n-modal v-model:show="showFriendRequests" preset="card" title="好友通知" style="width: 480px">
      <n-spin :show="loadingFriendRequests">
        <n-list v-if="friendRequestItems.length > 0">
          <n-list-item v-for="r in friendRequestItems" :key="r.requestId">
            <n-thing>
              <template #header>来自 {{ r.fromUserId }}</template>
              <template #default>{{ r.message || '请求添加你为好友' }}</template>
              <template #action>
                <n-space>
                  <n-button size="small" type="primary" @click="acceptFriendRequest(r.requestId)">通过</n-button>
                  <n-button size="small" @click="declineFriendRequest(r.requestId)">拒绝</n-button>
                </n-space>
              </template>
            </n-thing>
          </n-list-item>
        </n-list>
        <n-empty v-else description="暂无待处理好友申请" />
      </n-spin>
    </n-modal>
    <n-modal v-model:show="showGroupRequests" preset="card" title="群通知" style="width: 480px">
      <n-spin :show="loadingGroupRequests">
        <n-list v-if="groupRequestItems.length > 0">
          <n-list-item v-for="r in groupRequestItems" :key="r.requestId">
            <n-thing>
              <template #header>{{ r.userId }} 申请加入群 {{ r.conversationId }}</template>
              <template #default>{{ r.message || '申请加入群聊' }}</template>
              <template #action>
                <n-space>
                  <n-button size="small" type="primary" @click="approveGroupRequest(r)">通过</n-button>
                  <n-button size="small" @click="declineGroupRequest(r)">拒绝</n-button>
                </n-space>
              </template>
            </n-thing>
          </n-list-item>
        </n-list>
        <n-empty v-else description="暂无待审批入群申请" />
      </n-spin>
    </n-modal>
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

    <n-card title="添加好友" style="margin-top: 16px">
      <n-form :model="addForm" inline>
        <n-form-item label="用户名或 10 位账号">
          <n-input
            v-model:value="addForm.toUsernameOrAccount"
            placeholder="输入对方用户名或 10 位账号"
            style="width: 220px"
          />
        </n-form-item>
        <n-form-item label="验证消息（选填）">
          <n-input
            v-model:value="addForm.message"
            placeholder="如：你好，我是..."
            style="width: 220px"
          />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" :loading="adding" @click="addContact">发送好友申请</n-button>
        </n-form-item>
      </n-form>
      <n-alert v-if="addError" type="error" closable @close="addError = ''">{{ addError }}</n-alert>
    </n-card>

    <n-card title="申请加群" style="margin-top: 16px">
      <n-form :model="groupApplyForm" inline>
        <n-form-item label="群号（11 位）">
          <n-input
            v-model:value="groupApplyForm.conversationId"
            placeholder="输入群号"
            style="width: 160px"
          />
        </n-form-item>
        <n-form-item label="验证消息（选填）">
          <n-input
            v-model:value="groupApplyForm.message"
            placeholder="如：申请加入"
            style="width: 220px"
          />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" :loading="applyingGroup" @click="applyJoinGroup">申请加入</n-button>
        </n-form-item>
      </n-form>
      <n-alert v-if="groupApplyError" type="error" closable @close="groupApplyError = ''">{{ groupApplyError }}</n-alert>
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
            <n-button size="small" type="error" quaternary @click="confirmRemoveContact(c.userId)">删除好友</n-button>
          </li>
        </ul>
      </n-spin>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
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
  NSpace,
  NModal,
  NList,
  NListItem,
  NThing,
  useMessage,
  useDialog,
} from 'naive-ui'
import { useAuthStore } from '@/stores/auth'
import { useConversationStore } from '@/stores/conversation'
import { useWebSocket } from '@/composables/useWebSocket'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()
const authStore = useAuthStore()
const conversationStore = useConversationStore()
const { send, connected } = useWebSocket()

const formRef = ref(null)
const creating = ref(false)
const createError = ref('')
const form = reactive({ toUsernameOrAccount: '' })

const addForm = reactive({ toUsernameOrAccount: '', message: '' })
const adding = ref(false)
const addError = ref('')
const groupApplyForm = reactive({ conversationId: '', message: '' })
const applyingGroup = ref(false)
const groupApplyError = ref('')

const contactList = ref<{ userId: string }[]>([])
const loadingContacts = ref(false)
const showFriendRequests = ref(false)
const showGroupRequests = ref(false)
const friendRequestItems = ref<{ requestId: string; fromUserId: string; message: string }[]>([])
const groupRequestItems = ref<{ requestId: string; conversationId: string; userId: string; message: string }[]>([])
const loadingFriendRequests = ref(false)
const loadingGroupRequests = ref(false)

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
  const payload: { toUserId?: string; toUsername?: string; toAccount?: string; message?: string } = {}
  if (isTenDigitAccount(raw)) {
    payload.toAccount = raw
  } else {
    payload.toUsername = raw
  }
  const msg = addForm.message?.trim()
  if (msg) payload.message = msg
  send('contact.request', payload, (env) => {
    adding.value = false
    if (env.type === 'contact.request.ok') {
      message.success('已发送好友申请，等待对方通过')
      addForm.toUsernameOrAccount = ''
      addForm.message = ''
    } else if (env.error) {
      const errMsg = env.error.message || ''
      addError.value = /already|已是好友/i.test(errMsg) ? '已是好友，无需重复申请' : errMsg || '发送申请失败'
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
      message.success('已删除好友')
      loadContacts()
    } else if (env.error) {
      message.error(env.error.message || '删除失败')
    }
  })
}

function confirmRemoveContact(userId: string) {
  dialog.warning({
    title: '删除好友',
    content: '确定要删除该好友吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => removeContact(userId),
  })
}

function loadFriendRequests() {
  if (!connected.value) return
  loadingFriendRequests.value = true
  send('contact.requestList', {}, (env) => {
    loadingFriendRequests.value = false
    if (env.type === 'contact.requestList.ok' && env.payload) {
      const p = env.payload as { items: { requestId: string; fromUserId: string; message: string }[] }
      friendRequestItems.value = p.items || []
    }
  })
}

function acceptFriendRequest(requestId: string) {
  if (!connected.value) return
  send('contact.accept', { requestId }, (env) => {
    if (env.type === 'contact.accept.ok') {
      message.success('已通过')
      loadFriendRequests()
      loadContacts()
    } else if (env.error) message.error(env.error.message)
  })
}

function declineFriendRequest(requestId: string) {
  if (!connected.value) return
  send('contact.decline', { requestId }, (env) => {
    if (env.type === 'contact.decline.ok') {
      message.success('已拒绝')
      loadFriendRequests()
    } else if (env.error) message.error(env.error.message)
  })
}

function applyJoinGroup() {
  const cid = groupApplyForm.conversationId.trim()
  if (!cid) {
    message.warning('请输入群号')
    return
  }
  if (!connected.value) {
    message.error('未连接，请先登录')
    return
  }
  applyingGroup.value = true
  groupApplyError.value = ''
  const payload: { conversationId: string; message?: string } = { conversationId: cid }
  const msg = groupApplyForm.message?.trim()
  if (msg) payload.message = msg
  send('group.apply', payload, (env) => {
    applyingGroup.value = false
    if (env.type === 'group.apply.ok' && env.payload) {
      const p = env.payload as { joined?: boolean }
      if (p.joined) {
        message.success('已加入群聊')
        groupApplyForm.conversationId = ''
        groupApplyForm.message = ''
        send('conversation.list', { limit: 50 }, (e) => {
          if (e.type === 'conversation.list.ok' && e.payload) {
            const list = (e.payload as { items?: typeof conversationStore.conversations }).items || []
            conversationStore.setConversations(list)
          }
        })
        router.push(`/app/chats/${cid}`)
      } else {
        message.success('已提交申请，等待管理员通过')
        groupApplyForm.conversationId = ''
        groupApplyForm.message = ''
      }
    } else if (env.error) {
      groupApplyError.value = env.error.message || '申请失败'
    }
  })
}

function loadGroupRequests() {
  if (!connected.value) return
  loadingGroupRequests.value = true
  const groups = conversationStore.conversations.filter((c) => c.type === 'group')
  if (groups.length === 0) {
    groupRequestItems.value = []
    loadingGroupRequests.value = false
    return
  }
  let pending = groups.length
  const all: { requestId: string; conversationId: string; userId: string; message: string }[] = []
  groups.forEach((g) => {
    send('group.joinRequestList', { conversationId: g.id }, (env) => {
      if (env.type === 'group.joinRequestList.ok' && env.payload) {
        const p = env.payload as { items: { requestId: string; conversationId: string; userId: string; message: string }[] }
        const items = p.items || []
        items.forEach((it) => all.push(it))
      }
      pending--
      if (pending === 0) {
        groupRequestItems.value = all
        loadingGroupRequests.value = false
      }
    })
  })
}

function approveGroupRequest(r: { requestId: string; conversationId: string }) {
  if (!connected.value) return
  send('group.approve', { conversationId: r.conversationId, requestId: r.requestId }, (env) => {
    if (env.type === 'group.approve.ok') {
      message.success('已通过')
      loadGroupRequests()
    } else if (env.error) message.error(env.error.message)
  })
}

function declineGroupRequest(r: { requestId: string; conversationId: string }) {
  if (!connected.value) return
  send('group.decline', { conversationId: r.conversationId, requestId: r.requestId }, (env) => {
    if (env.type === 'group.decline.ok') {
      message.success('已拒绝')
      loadGroupRequests()
    } else if (env.error) message.error(env.error.message)
  })
}

watch(showFriendRequests, (v) => {
  if (v) loadFriendRequests()
})
watch(showGroupRequests, (v) => {
  if (v) {
    if (conversationStore.conversations.length === 0 && connected.value) {
      send('conversation.list', { limit: 50 }, (env) => {
        if (env.type === 'conversation.list.ok' && env.payload) {
          const p = env.payload as { items: typeof conversationStore.conversations }
          conversationStore.setConversations(p.items || [])
        }
        loadGroupRequests()
      })
    } else {
      loadGroupRequests()
    }
  }
})

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
