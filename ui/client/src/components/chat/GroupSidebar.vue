<template>
  <div class="group-sidebar">
    <div class="section">
      <div class="section-title">群公告</div>
      <p class="announcement">{{ groupDetail?.announcement || '暂无公告' }}</p>
    </div>
    <div class="section">
      <div class="section-title">群聊成员 {{ members.length }}</div>
      <n-spin :show="loadingMembers">
        <ul class="member-list">
          <li v-for="m in members" :key="m.userId" class="member-item">
            <n-avatar round size="small">{{ (m.userId || '?').slice(-1) }}</n-avatar>
            <span class="member-name">{{ m.userId }}</span>
            <n-tag v-if="m.role === 'owner'" size="small" type="info">群主</n-tag>
            <n-tag v-else-if="m.role === 'admin'" size="small">管理员</n-tag>
          </li>
        </ul>
      </n-spin>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { NAvatar, NTag, NSpin } from 'naive-ui'
import { useWebSocket } from '@/composables/useWebSocket'

const props = defineProps<{ conversationId: string }>()
const { send, connected } = useWebSocket()

const groupDetail = ref<{ name: string; memberCount: number; announcement: string } | null>(null)
const members = ref<{ userId: string; role: string }[]>([])
const loadingMembers = ref(false)

function load() {
  if (!connected.value || !props.conversationId) return
  loadingMembers.value = true
  send('conversation.get', { id: props.conversationId }, (env) => {
    if (env.type === 'conversation.get.ok' && env.payload) {
      const p = env.payload as { name: string; memberCount: number; announcement: string }
      groupDetail.value = p
    }
    loadingMembers.value = false
  })
  send('conversation.listMembers', { conversationId: props.conversationId }, (env) => {
    if (env.type === 'conversation.listMembers.ok' && env.payload) {
      const p = env.payload as { members: { userId: string; role: string }[] }
      members.value = p.members || []
    }
  })
}

watch(
  () => props.conversationId,
  (id) => {
    if (id) load()
    else {
      groupDetail.value = null
      members.value = []
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.group-sidebar {
  padding: 16px;
  overflow: auto;
  border-left: 1px solid var(--n-border-color);
}
.section {
  margin-bottom: 20px;
}
.section-title {
  font-weight: 600;
  margin-bottom: 8px;
  font-size: 14px;
}
.announcement {
  font-size: 13px;
  color: var(--n-textColor2);
  white-space: pre-wrap;
  margin: 0;
}
.member-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.member-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0;
}
.member-name {
  flex: 1;
  font-size: 13px;
}
</style>
