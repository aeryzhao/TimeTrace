<script setup>
import { RouterLink, RouterView } from 'vue-router'
import { Timer, Calendar, PieChart, Collection } from '@element-plus/icons-vue'
</script>

<template>
  <div class="common-layout">
    <el-container>
      <el-header class="header">
        <div class="logo">
          <el-icon class="logo-icon"><Timer /></el-icon>
          TimeTrace
        </div>
        <el-menu
          mode="horizontal"
          :router="true"
          :default-active="$route.path"
          class="nav-menu"
          :ellipsis="false"
        >
          <el-menu-item index="/">
            <el-icon><Timer /></el-icon>
            <span>计时台</span>
          </el-menu-item>
          <el-menu-item index="/timeline">
            <el-icon><Calendar /></el-icon>
            <span>时间轴</span>
          </el-menu-item>
          <el-menu-item index="/reports">
            <el-icon><PieChart /></el-icon>
            <span>报表</span>
          </el-menu-item>
          <el-menu-item index="/categories">
            <el-icon><Collection /></el-icon>
            <span>分类管理</span>
          </el-menu-item>
        </el-menu>
      </el-header>
      <el-main class="main-content">
        <RouterView />
      </el-main>
    </el-container>
  </div>
</template>

<style>
body {
  margin: 0;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  background-color: #f5f7fa;
}
</style>

<style scoped>
.header {
  display: flex;
  align-items: center;
  border-bottom: 1px solid var(--el-menu-border-color);
  padding: 0 40px;
  background-color: white;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  z-index: 10;
  height: 60px;
}
.logo {
  font-size: 1.5rem;
  font-weight: bold;
  margin-right: 60px;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
  gap: 10px;
}
.logo-icon {
  font-size: 1.8rem;
}
.nav-menu {
  flex: 1;
  border-bottom: none;
  background-color: transparent;
}
.main-content {
  padding: 30px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  min-height: calc(100vh - 60px);
}
</style>
