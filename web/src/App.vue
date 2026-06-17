<script setup>
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { Timer, Calendar, PieChart, Collection, Setting } from '@element-plus/icons-vue'
import { computed } from 'vue'

const route = useRoute()
const navItems = [
  { path: '/', icon: Timer, label: '计时台' },
  { path: '/timeline', icon: Calendar, label: '时间轴' },
  { path: '/reports', icon: PieChart, label: '报表' },
  { path: '/categories', icon: Collection, label: '分类管理' },
  { path: '/settings', icon: Setting, label: '设置' },
]
</script>

<template>
  <div class="app-shell">
    <aside class="sidebar">
      <div class="sidebar-logo">
        <el-icon class="logo-icon"><Timer /></el-icon>
        <span class="logo-text">TimeTrace</span>
      </div>
      <nav class="sidebar-nav">
        <RouterLink
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ active: route.path === item.path }"
        >
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>
    </aside>
    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>

<style>
body { background: #f0f2f7; }
.el-card {
  --el-card-border-radius: 14px !important;
  --el-card-border-color: #e2e8f0 !important;
  box-shadow: 0 1px 3px 0 rgba(0,0,0,0.06), 0 1px 2px -1px rgba(0,0,0,0.04) !important;
}
.el-table { font-size: 14px; }
.el-button { border-radius: 8px !important; }
.el-tag { border-radius: 6px !important; }
.el-input__wrapper { border-radius: 8px !important; }
.el-select__wrapper { border-radius: 8px !important; }
</style>

<style scoped>
.app-shell {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  width: 220px;
  min-width: 220px;
  background: #fff;
  border-right: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  padding: 0;
  box-shadow: 1px 0 8px rgba(0,0,0,0.04);
  z-index: 10;
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 24px 20px 20px;
  font-size: 1.25rem;
  font-weight: 800;
  color: #4F46E5;
  letter-spacing: -0.3px;
  border-bottom: 1px solid #f1f5f9;
}

.logo-icon {
  font-size: 1.5rem;
  background: #EEF2FF;
  padding: 6px;
  border-radius: 10px;
  color: #4F46E5;
}

.sidebar-nav {
  padding: 14px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border-radius: 10px;
  text-decoration: none;
  font-size: 0.93rem;
  font-weight: 500;
  color: #64748b;
  transition: background 0.15s, color 0.15s;
  cursor: pointer;
}

.nav-item:hover {
  background: #f1f5f9;
  color: #1e293b;
}

.nav-item.active {
  background: #EEF2FF;
  color: #4F46E5;
  font-weight: 600;
}

.nav-item .el-icon {
  font-size: 1.1rem;
}

.app-main {
  flex: 1;
  overflow-y: auto;
  padding: 32px 36px;
  background: #f0f2f7;
}

@media (max-width: 768px) {
  .sidebar { display: none; }
  .app-main { padding: 16px; }
}
</style>
