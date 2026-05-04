<template>
  <div class="page">
    <NuxtRouteAnnouncer />

    <main class="card">
      <p class="eyebrow">Event Driven SNS</p>
      <h1 class="title">Login</h1>
      <p class="description">
        バックエンドの <code>/auth/login</code> へ接続する簡易ログイン画面です。
      </p>

      <form class="form" @submit.prevent="submitLogin">
        <label class="field">
          <span class="label">Login ID</span>
          <input
            v-model="loginId"
            class="input"
            type="text"
            name="login_id"
            autocomplete="username"
            placeholder="Login ID or email"
            required
          />
        </label>

        <label class="field">
          <span class="label">Password</span>
          <input
            v-model="password"
            class="input"
            type="password"
            name="password"
            autocomplete="current-password"
            placeholder="Password"
            required
          />
        </label>

        <button class="button" type="submit" :disabled="pending">
          {{ pending ? 'ログイン中...' : 'ログイン' }}
        </button>

        <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
      </form>
    </main>
  </div>
</template>

<script setup lang="ts">
const config = useRuntimeConfig()

const loginId = ref('')
const password = ref('')
const pending = ref(false)
const errorMessage = ref('')

async function submitLogin() {
  pending.value = true
  errorMessage.value = ''

  try {
    await $fetch('/auth/login', {
      baseURL: config.public.apiBaseUrl,
      method: 'POST',
      credentials: 'include',
      body: {
        login_id: loginId.value,
        password: password.value,
      },
    })

    window.alert('ログインが完了しました。')
  } catch (error) {
    console.error('Login failed:', error)
    errorMessage.value = 'ログインに失敗しました。'
  } finally {
    pending.value = false
  }
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
  background:
    radial-gradient(circle at top, #f6d8b8 0%, transparent 32%),
    linear-gradient(180deg, #f6f1e8 0%, #ece4d7 100%);
  color: #1f1a17;
}

.card {
  width: min(100%, 420px);
  padding: 32px;
  border: 1px solid #d8c9b4;
  border-radius: 24px;
  background: rgba(255, 251, 246, 0.92);
  box-shadow: 0 24px 60px rgba(74, 48, 24, 0.12);
}

.eyebrow {
  margin: 0 0 8px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #9b6a36;
}

.title {
  margin: 0;
  font-size: 36px;
  line-height: 1.1;
}

.description {
  margin: 12px 0 24px;
  color: #5b4d41;
  line-height: 1.7;
}

.form {
  display: grid;
  gap: 16px;
}

.field {
  display: grid;
  gap: 8px;
}

.label {
  font-size: 14px;
  font-weight: 600;
}

.input {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid #cdbba4;
  border-radius: 12px;
  background: #fffdfa;
  font: inherit;
}

.input:focus {
  outline: 2px solid #d69957;
  outline-offset: 2px;
}

.button {
  border: 0;
  border-radius: 999px;
  padding: 14px 18px;
  background: #1f1a17;
  color: #fff7ef;
  font: inherit;
  font-weight: 700;
  cursor: pointer;
}

.button:disabled {
  cursor: wait;
  opacity: 0.7;
}

.error {
  margin: 0;
  color: #b42318;
  font-size: 14px;
}

</style>
