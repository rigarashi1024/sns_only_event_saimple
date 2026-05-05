<template>
  <header class="header">
    <div class="inner">
      <NuxtLink class="brand" to="/timeline">
        <span class="brand-mark">ED</span>
        <span class="brand-copy">
          <strong>Event Driven SNS</strong>
          <span>After Login</span>
        </span>
      </NuxtLink>

      <nav class="nav" aria-label="Primary">
        <NuxtLink class="nav-link" to="/timeline">Timeline</NuxtLink>
        <NuxtLink class="profile-link" to="/profile" aria-label="プロフィールへ移動">
          <span class="profile-icon" aria-hidden="true">👤</span>
        </NuxtLink>
      </nav>
    </div>
  </header>
</template>

<style scoped>
.header {
  position: sticky;
  top: 0;
  z-index: 20;
  padding: 16px 20px 0;
}

.inner {
  width: min(100%, 920px);
  margin: 0 auto;
  padding: 14px 16px;
  border: 1px solid rgba(115, 92, 62, 0.16);
  border-radius: 24px;
  background: rgba(255, 251, 245, 0.86);
  box-shadow: 0 20px 50px rgba(86, 62, 32, 0.08);
  backdrop-filter: blur(12px);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  color: inherit;
  text-decoration: none;
}

.brand-mark {
  width: 42px;
  height: 42px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, #8e5f2b 0%, #d8a15a 100%);
  color: #fff8ef;
  font-size: 13px;
  font-weight: 800;
  letter-spacing: 0.08em;
}

.brand-copy {
  display: grid;
  gap: 2px;
}

.brand-copy strong {
  font-size: 15px;
}

.brand-copy span {
  color: #76695c;
  font-size: 12px;
}

.nav {
  display: flex;
  align-items: center;
  gap: 10px;
}

.nav-link,
.profile-link {
  color: inherit;
  text-decoration: none;
}

.nav-link {
  padding: 10px 12px;
  border-radius: 999px;
  color: #5e5347;
  font-weight: 600;
}

.nav-link.router-link-active {
  background: rgba(142, 95, 43, 0.1);
  color: #6d471d;
}

.profile-link {
  width: 44px;
  height: 44px;
  border: 1px solid #d7c2a8;
  border-radius: 999px;
  display: grid;
  place-items: center;
  background: #fffdf9;
  transition: transform 0.18s ease, box-shadow 0.18s ease;
}

.profile-link:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 24px rgba(86, 62, 32, 0.1);
}

.profile-link.router-link-active {
  border-color: #8e5f2b;
  background: #f8ede0;
}

.profile-icon {
  font-size: 20px;
  line-height: 1;
}

@media (max-width: 640px) {
  .header {
    padding-inline: 14px;
  }

  .inner {
    padding: 12px 14px;
  }

  .brand-copy span {
    display: none;
  }
}
</style>
