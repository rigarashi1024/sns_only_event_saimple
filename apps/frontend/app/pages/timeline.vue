<template>
  <div class="page">
    <main class="layout">
      <section class="hero">
        <p class="eyebrow">After Login</p>
        <h1 class="title">Timeline</h1>
        <p class="description">
          保存済みのタイムラインを API から取得して表示しています。
        </p>
      </section>

      <section class="panel">
        <div class="panel-header">
          <h2 class="panel-title">Recent Posts</h2>
          <button class="ghost-button" type="button" :disabled="pending" @click="refreshTimeline">
            {{ pending ? '更新中...' : '再取得' }}
          </button>
        </div>

        <p v-if="pending && items.length === 0" class="status">タイムラインを読み込み中です...</p>
        <p v-else-if="errorMessage" class="status error">{{ errorMessage }}</p>
        <p v-else-if="items.length === 0" class="status">表示できる投稿はまだありません。</p>

        <ul v-else class="timeline-list">
          <li v-for="item in items" :key="item.id" class="timeline-card">
            <div class="meta">
              <span class="author">@{{ item.post_user_id }}</span>
              <time class="time" :datetime="item.post_created_at">
                {{ formatDate(item.post_created_at) }}
              </time>
            </div>
            <p class="content">{{ item.content }}</p>
            <p class="caption">
              timeline owner: <code>{{ item.timeline_owner_user_id }}</code>
            </p>
          </li>
        </ul>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
type TimelineEntry = {
  id: string
  post_id: string
  post_user_id: string
  content: string
  post_created_at: string
  timeline_owner_user_id: string
  created_at: string
  updated_at: string
}

type TimelineListResponse = {
  items: TimelineEntry[]
}

const config = useRuntimeConfig()
const pending = ref(true)
const errorMessage = ref('')
const items = ref<TimelineEntry[]>([])

async function loadTimeline() {
  pending.value = true
  errorMessage.value = ''

  try {
    const response = await $fetch<TimelineListResponse>('timeline', {
      baseURL: config.public.apiBaseUrl,
      method: 'GET',
      credentials: 'include',
      headers: {
        'X-Requested-With': 'XMLHttpRequest',
      },
    })

    items.value = response.items
  } catch (error) {
    console.error('Failed to load timeline:', error)
    errorMessage.value = 'タイムラインの取得に失敗しました。ログイン状態を確認してください。'
  } finally {
    pending.value = false
  }
}

async function refreshTimeline() {
  await loadTimeline()
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('ja-JP', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
}

onMounted(async () => {
  await loadTimeline()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 32px 20px 56px;
  background:
    radial-gradient(circle at top left, rgba(255, 204, 138, 0.65) 0%, transparent 28%),
    radial-gradient(circle at top right, rgba(111, 181, 255, 0.2) 0%, transparent 24%),
    linear-gradient(180deg, #f7f2e8 0%, #efe5d4 48%, #e7dcc8 100%);
  color: #1d1a16;
}

.layout {
  width: min(100%, 920px);
  margin: 0 auto;
  display: grid;
  gap: 24px;
}

.hero {
  padding: 8px 4px;
}

.eyebrow {
  margin: 0 0 8px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #8e5f2b;
}

.title {
  margin: 0;
  font-size: clamp(40px, 8vw, 64px);
  line-height: 0.96;
}

.description {
  max-width: 560px;
  margin: 14px 0 0;
  color: #5e5347;
  line-height: 1.7;
}

.panel {
  padding: 24px;
  border: 1px solid rgba(115, 92, 62, 0.18);
  border-radius: 28px;
  background: rgba(255, 250, 243, 0.82);
  box-shadow: 0 28px 70px rgba(86, 62, 32, 0.12);
  backdrop-filter: blur(10px);
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
}

.panel-title {
  margin: 0;
  font-size: 22px;
}

.ghost-button {
  border: 1px solid #cdbba4;
  border-radius: 999px;
  padding: 10px 14px;
  background: #fffdf9;
  color: #2b241d;
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

.ghost-button:disabled {
  cursor: wait;
  opacity: 0.7;
}

.status {
  margin: 0;
  color: #5e5347;
}

.error {
  color: #b42318;
}

.timeline-list {
  display: grid;
  gap: 16px;
  padding: 0;
  margin: 0;
  list-style: none;
}

.timeline-card {
  padding: 18px 18px 16px;
  border: 1px solid #eadbc5;
  border-radius: 20px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.92) 0%, rgba(250, 245, 237, 0.94) 100%);
}

.meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 10px;
  font-size: 14px;
}

.author {
  font-weight: 700;
  color: #8e5f2b;
}

.time {
  color: #74675a;
}

.content {
  margin: 0;
  font-size: 17px;
  line-height: 1.7;
  white-space: pre-wrap;
}

.caption {
  margin: 14px 0 0;
  color: #74675a;
  font-size: 13px;
}

@media (max-width: 640px) {
  .page {
    padding-inline: 14px;
  }

  .panel {
    padding: 18px;
  }

  .panel-header,
  .meta {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
