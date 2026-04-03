<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'

const router = useRouter()
const route = useRoute()
const { t, locale } = useI18n()
const authStore = useAuthStore()
const themeStore = useThemeStore()

const mobileMenuOpen = ref(false)
const userMenuOpen = ref(false)

const navLinks = computed(() => [
  { name: 'Home', label: t('nav.home'), path: '/' },
  { name: 'ProblemList', label: t('nav.problems'), path: '/problems' },
  { name: 'ContestList', label: t('nav.contests'), path: '/contests' },
  { name: 'CTFPractice', label: t('nav.ctf'), path: '/ctf' },
])

function isActive(path: string): boolean {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}

function navigate(path: string) {
  router.push(path)
  mobileMenuOpen.value = false
}

function toggleLanguage() {
  const newLocale = locale.value === 'zh' ? 'en' : 'zh'
  locale.value = newLocale
  themeStore.setLocale(newLocale)
}

function toggleTheme() {
  themeStore.toggleTheme()
}

async function handleLogout() {
  await authStore.logout()
  userMenuOpen.value = false
  router.push({ name: 'Home' })
}

function toggleMobileMenu() {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

function toggleUserMenu() {
  userMenuOpen.value = !userMenuOpen.value
}
</script>

<template>
  <header class="app-header">
    <div class="header-inner container">
      <!-- Logo -->
      <div class="header-left">
        <router-link to="/" class="logo" @click="mobileMenuOpen = false">
          <span class="logo-icon">Y</span>
          <span class="logo-text">YogduOJ</span>
        </router-link>
      </div>

      <!-- Mobile menu button -->
      <button class="mobile-menu-btn" @click="toggleMobileMenu">
        <span :class="['hamburger', { active: mobileMenuOpen }]">
          <span></span>
          <span></span>
          <span></span>
        </span>
      </button>

      <!-- Navigation -->
      <nav :class="['header-nav', { open: mobileMenuOpen }]">
        <router-link
          v-for="link in navLinks"
          :key="link.name"
          :to="link.path"
          :class="['nav-link', { active: isActive(link.path) }]"
          @click="mobileMenuOpen = false"
        >
          {{ link.label }}
        </router-link>
      </nav>

      <!-- Right side -->
      <div class="header-right">
        <!-- Language switcher -->
        <button class="icon-btn" @click="toggleLanguage" :title="locale === 'zh' ? 'English' : '中文'">
          {{ locale === 'zh' ? 'EN' : '中' }}
        </button>

        <!-- Theme toggle -->
        <button class="icon-btn" @click="toggleTheme" :title="themeStore.isDark ? 'Light Mode' : 'Dark Mode'">
          {{ themeStore.isDark ? '☀' : '☾' }}
        </button>

        <!-- User menu -->
        <div v-if="authStore.isLoggedIn" class="user-menu">
          <button class="user-btn" @click="toggleUserMenu">
            <span class="user-avatar">{{ authStore.user?.username?.charAt(0)?.toUpperCase() || 'U' }}</span>
            <span class="user-name">{{ authStore.user?.username }}</span>
          </button>
          <div v-if="userMenuOpen" class="user-dropdown">
            <router-link to="/profile" class="dropdown-item" @click="userMenuOpen = false">
              {{ t('nav.profile') }}
            </router-link>
            <router-link v-if="authStore.isAdmin" to="/admin" class="dropdown-item" @click="userMenuOpen = false">
              {{ t('nav.admin') }}
            </router-link>
            <button class="dropdown-item logout" @click="handleLogout">
              {{ t('nav.logout') }}
            </button>
          </div>
        </div>

        <!-- Login button -->
        <router-link v-else to="/login" class="btn btn-primary login-btn">
          {{ t('nav.login') }}
        </router-link>
      </div>
    </div>
  </header>
</template>

<style scoped lang="scss">
.app-header {
  position: sticky;
  top: 0;
  z-index: 1000;
  height: var(--header-height);
  background-color: var(--color-bg-secondary);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(10px);
}

.header-inner {
  display: flex;
  align-items: center;
  height: 100%;
  gap: 24px;
}

.header-left {
  flex-shrink: 0;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-text);
  font-size: 20px;
  font-weight: 700;
  text-decoration: none;

  &:hover {
    color: var(--color-primary);
  }
}

.logo-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
  border-radius: 8px;
  color: white;
  font-size: 18px;
  font-weight: 800;
}

.logo-text {
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-nav {
  display: flex;
  align-items: center;
  gap: 4px;
  flex: 1;
}

.nav-link {
  padding: 8px 16px;
  border-radius: var(--border-radius);
  color: var(--color-text-secondary);
  font-size: 14px;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.2s ease;

  &:hover {
    color: var(--color-text);
    background-color: rgba(255, 255, 255, 0.05);
  }

  &.active {
    color: var(--color-primary);
    background-color: rgba(0, 212, 255, 0.1);
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: var(--border-radius);
  background-color: rgba(255, 255, 255, 0.05);
  color: var(--color-text-secondary);
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s ease;

  &:hover {
    background-color: rgba(255, 255, 255, 0.1);
    color: var(--color-text);
  }
}

.user-menu {
  position: relative;
}

.user-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 12px 4px 4px;
  border: none;
  border-radius: var(--border-radius);
  background-color: rgba(255, 255, 255, 0.05);
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }
}

.user-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
  color: white;
  font-size: 14px;
  font-weight: 600;
}

.user-name {
  color: var(--color-text);
  font-size: 14px;
  font-weight: 500;
}

.user-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 160px;
  background-color: var(--color-bg-secondary);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow);
  overflow: hidden;
  z-index: 1001;
}

.dropdown-item {
  display: block;
  width: 100%;
  padding: 10px 16px;
  color: var(--color-text);
  font-size: 14px;
  text-decoration: none;
  text-align: left;
  background: none;
  border: none;
  cursor: pointer;
  transition: background-color 0.2s ease;

  &:hover {
    background-color: rgba(255, 255, 255, 0.05);
  }

  &.logout {
    color: var(--color-error);
  }
}

.login-btn {
  padding: 6px 16px;
  font-size: 14px;
  text-decoration: none;
}

.mobile-menu-btn {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
}

.hamburger {
  display: flex;
  flex-direction: column;
  gap: 4px;
  width: 20px;

  span {
    display: block;
    height: 2px;
    background-color: var(--color-text);
    border-radius: 2px;
    transition: all 0.3s ease;
  }

  &.active {
    span:nth-child(1) {
      transform: rotate(45deg) translate(4px, 4px);
    }
    span:nth-child(2) {
      opacity: 0;
    }
    span:nth-child(3) {
      transform: rotate(-45deg) translate(4px, -4px);
    }
  }
}

@media (max-width: 768px) {
  .mobile-menu-btn {
    display: block;
  }

  .header-nav {
    position: absolute;
    top: var(--header-height);
    left: 0;
    right: 0;
    flex-direction: column;
    background-color: var(--color-bg-secondary);
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
    padding: 8px;
    display: none;

    &.open {
      display: flex;
    }
  }

  .nav-link {
    width: 100%;
    padding: 12px 16px;
  }

  .user-name {
    display: none;
  }
}
</style>
