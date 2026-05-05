<template>
  <div class="page">
    <AppHeader />

    <main class="layout">
      <section class="hero">
        <p class="eyebrow">Profile</p>
        <h1 class="title">My Profile</h1>
        <p class="description">
          自分のプロフィール情報と過去の投稿を表示しています。
        </p>
      </section>

      <div class="grid">
        <section class="profile-card panel">
          <p v-if="pending && !profile" class="status">プロフィールを読み込み中です...</p>
          <p v-else-if="errorMessage" class="status error">{{ errorMessage }}</p>

          <template v-else-if="profile">
            <div class="profile-top">
              <div class="avatar" aria-hidden="true">
                {{ profile.user.nickname.slice(0, 1).toUpperCase() }}
              </div>
              <div>
                <h2 class="name">{{ profile.user.name }}</h2>
                <p class="nickname">@{{ profile.user.nickname }}</p>
              </div>
            </div>

            <dl class="profile-meta">
              <div>
                <dt>User ID</dt>
                <dd>{{ profile.user.id }}</dd>
              </div>
              <div>
                <dt>Email</dt>
                <dd>{{ profile.user.email }}</dd>
              </div>
              <div>
                <dt>Posts</dt>
                <dd>{{ profile.posts.length }}</dd>
              </div>
            </dl>
          </template>
        </section>

        <section class="panel">
          <div class="panel-header">
            <h2 class="panel-title">My Posts</h2>
            <button class="ghost-button" type="button" :disabled="pending" @click="refreshProfile">
              {{ pending ? '更新中...' : '再取得' }}
            </button>
          </div>

          <p v-if="pending && profilePosts.length === 0" class="status">投稿を読み込み中です...</p>
          <p v-else-if="errorMessage" class="status error">{{ errorMessage }}</p>
          <p v-else-if="profilePosts.length === 0" class="status">まだ投稿はありません。</p>

          <ul v-else class="post-list">
            <li v-for="post in profilePosts" :key="post.id" class="post-card">
              <time class="time" :datetime="post.created_at">
                {{ formatDate(post.created_at) }}
              </time>
              <p class="content">{{ post.content }}</p>
              <p class="stats">likes {{ post.like_count }} / comments {{ post.comment_count }}</p>
            </li>
          </ul>
        </section>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
type UserProfile = {
  id: string
  name: string
  email: string
  nickname: string
}

type UserPostSummary = {
  id: string
  user_id: string
  content: string
  like_count: number
  comment_count: number
  created_at: string
  updated_at: string
}

type ProfileResponse = {
  user: UserProfile
  posts: UserPostSummary[]
}

const pending = ref(true)
const errorMessage = ref('')
const profile = ref<ProfileResponse | null>(null)
const { fetchWithSession } = useAuthenticatedApi()

const profilePosts = computed(() => profile.value?.posts ?? [])

async function loadProfile() {
  pending.value = true
  errorMessage.value = ''

  try {
    const response = await fetchWithSession<ProfileResponse>('profile')
    if (!response) {
      return
    }
    profile.value = response
  } catch (error) {
    console.error('Failed to load profile:', error)
    errorMessage.value = 'プロフィールの取得に失敗しました。'
  } finally {
    pending.value = false
  }
}

async function refreshProfile() {
  await loadProfile()
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('ja-JP', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
}

onMounted(async () => {
  await loadProfile()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding-bottom: 56px;
  background:
    radial-gradient(circle at top left, rgba(255, 204, 138, 0.65) 0%, transparent 28%),
    radial-gradient(circle at top right, rgba(111, 181, 255, 0.2) 0%, transparent 24%),
    linear-gradient(180deg, #f7f2e8 0%, #efe5d4 48%, #e7dcc8 100%);
  color: #1d1a16;
}

.layout {
  width: min(100%, 920px);
  margin: 0 auto;
  padding: 24px 20px 0;
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

.grid {
  display: grid;
  gap: 24px;
  grid-template-columns: minmax(260px, 320px) minmax(0, 1fr);
}

.panel {
  padding: 24px;
  border: 1px solid rgba(115, 92, 62, 0.18);
  border-radius: 28px;
  background: rgba(255, 250, 243, 0.82);
  box-shadow: 0 28px 70px rgba(86, 62, 32, 0.12);
  backdrop-filter: blur(10px);
}

.profile-top {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.avatar {
  width: 68px;
  height: 68px;
  border-radius: 22px;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, #8e5f2b 0%, #d8a15a 100%);
  color: #fff8ef;
  font-size: 30px;
  font-weight: 800;
}

.name {
  margin: 0;
  font-size: 24px;
}

.nickname {
  margin: 6px 0 0;
  color: #765f49;
}

.profile-meta {
  display: grid;
  gap: 14px;
  margin: 0;
}

.profile-meta div {
  padding-top: 14px;
  border-top: 1px solid #eadbc5;
}

.profile-meta dt {
  color: #765f49;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.profile-meta dd {
  margin: 6px 0 0;
  font-size: 16px;
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

.post-list {
  display: grid;
  gap: 16px;
  padding: 0;
  margin: 0;
  list-style: none;
}

.post-card {
  padding: 18px;
  border: 1px solid #eadbc5;
  border-radius: 20px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.92) 0%, rgba(250, 245, 237, 0.94) 100%);
}

.time {
  color: #74675a;
  font-size: 13px;
}

.content {
  margin: 10px 0 0;
  font-size: 17px;
  line-height: 1.7;
  white-space: pre-wrap;
}

.stats {
  margin: 14px 0 0;
  color: #74675a;
  font-size: 13px;
}

@media (max-width: 820px) {
  .grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .layout {
    padding-inline: 14px;
  }

  .panel {
    padding: 18px;
  }

  .panel-header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
