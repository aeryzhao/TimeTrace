<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useTimerStore } from '../stores/timer'
import { VideoPlay, VideoPause, Plus, EditPen } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

const store = useTimerStore()
const elapsedTime = ref('00:00:00')
const quickForm = ref({
  name: '',
  category_id: null
})
const searchQuery = ref('')
let timerInterval = null

const filteredActivities = computed(() => {
  if (!searchQuery.value) return store.activities
  const q = searchQuery.value.toLowerCase()
  return store.activities.filter(a => a.name.toLowerCase().includes(q) || a.category?.name?.toLowerCase().includes(q))
})

const quickSuggestions = computed(() => store.activities.slice(0, 8))

const updateTicker = () => {
  if (store.currentEntry && store.currentEntry.start_time) {
    const start = dayjs(store.currentEntry.start_time)
    const now = dayjs()
    const diff = now.diff(start, 'second')
    const h = Math.floor(diff / 3600)
    const m = Math.floor((diff % 3600) / 60)
    const s = diff % 60
    elapsedTime.value = `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
  } else {
    elapsedTime.value = '00:00:00'
  }
}

onMounted(async () => {
  await store.fetchMeta()
  await store.fetchCurrentTimer()
  timerInterval = setInterval(updateTicker, 1000)
  updateTicker()
})

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval)
})

const handleStart = (activity) => {
  store.start(activity.name, activity.category_id || null)
}

const handleQuickStart = async () => {
  await store.start(quickForm.value.name, quickForm.value.category_id)
  if (store.currentEntry) {
    quickForm.value.name = ''
    quickForm.value.category_id = null
  }
}

const currentCategoryName = computed(() => store.currentEntry?.category?.name || '未分类')
</script>

<template>
  <div class="dashboard">
    <el-card class="timer-card" :class="{ 'is-running': store.isRunning }" shadow="hover">
      <template #header>
        <div class="card-header">
          <span>{{ store.isRunning ? '正在进行' : '准备开始' }}</span>
          <span class="timer-display">{{ elapsedTime }}</span>
        </div>
      </template>

      <div v-if="store.isRunning" class="current-activity">
        <div class="info">
          <h3>{{ store.currentEntry.activity?.name }}</h3>
          <el-tag>{{ currentCategoryName }}</el-tag>
          <p class="start-time">开始于 {{ dayjs(store.currentEntry.start_time).format('HH:mm') }}</p>
        </div>
        <el-button type="danger" size="large" circle :icon="VideoPause" @click="store.stop"></el-button>
      </div>
      <div v-else class="placeholder">
        <p>当前没有活动。请在下方输入活动名称后开始。</p>
      </div>
    </el-card>

    <div class="section">
      <h2>开始活动</h2>
      <el-card class="start-card" shadow="hover">
        <div class="start-grid">
          <el-input
            v-model="quickForm.name"
            placeholder="输入当前要做的事情，例如：写方案、开会、吃饭"
            class="start-input"
            size="large"
            clearable
            @keyup.enter="handleQuickStart"
          >
            <template #prefix>
              <el-icon><EditPen /></el-icon>
            </template>
          </el-input>

          <el-select
            v-model="quickForm.category_id"
            placeholder="选择分类（可选）"
            clearable
            filterable
            class="start-select"
            size="large"
          >
            <el-option
              v-for="category in store.categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>

          <el-button
            type="primary"
            size="large"
            :icon="VideoPlay"
            class="start-btn"
            @click="handleQuickStart"
          >
            开始计时
          </el-button>
        </div>
      </el-card>
    </div>

    <div class="section" v-if="quickSuggestions.length">
      <h2>最近活动</h2>
      <div class="quick-grid">
        <el-button
          v-for="act in quickSuggestions"
          :key="act.id"
          class="quick-btn"
          type="primary"
          plain
          @click="handleStart(act)"
        >
          <span class="quick-btn-name">{{ act.name }}</span>
          <span v-if="act.category?.name" class="quick-btn-cat">{{ act.category.name }}</span>
          <span v-else class="quick-btn-cat">未分类</span>
        </el-button>
        <el-button :icon="Plus" class="quick-btn dashed" @click="$router.push('/categories')">管理分类</el-button>
      </div>
    </div>

    <div class="section">
      <h2>所有活动</h2>
      <el-input
        v-model="searchQuery"
        placeholder="搜索活动..."
        prefix-icon="Search"
        class="search-input"
        size="large"
      />
      <div class="activity-list">
        <el-card
          v-for="act in filteredActivities"
          :key="act.id"
          shadow="hover"
          class="activity-item"
          @click="handleStart(act)"
        >
          <div class="activity-inner">
            <div class="act-content">
              <span class="act-name">{{ act.name }}</span>
              <div class="act-meta">
                <el-tag size="small" type="info" effect="plain">{{ act.category?.name || '未分类' }}</el-tag>
              </div>
            </div>
            <div class="act-action">
              <el-icon class="play-icon"><VideoPlay /></el-icon>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<style scoped>
.timer-card {
  margin-bottom: 30px;
  text-align: center;
  transition: all 0.3s;
  border-radius: 12px;
}

.timer-card.is-running {
  border-color: var(--el-color-primary);
  background-color: var(--el-color-primary-light-9);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 1.2rem;
  font-weight: 600;
}

.timer-display {
  font-family: 'Monaco', 'Menlo', 'Consolas', monospace;
  font-size: 2.2rem;
  font-weight: bold;
  color: var(--el-color-primary);
}

.current-activity {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
}

.current-activity h3 {
  margin: 0 0 10px 0;
  font-size: 1.8rem;
  color: var(--el-text-color-primary);
}

.start-time {
  margin-top: 10px;
  color: var(--el-text-color-secondary);
  font-size: 0.9rem;
}

.placeholder {
  padding: 30px;
  color: var(--el-text-color-placeholder);
  font-size: 1rem;
}

.section {
  margin-bottom: 40px;
}

.section h2 {
  font-size: 1.2rem;
  margin-bottom: 20px;
  color: var(--el-text-color-regular);
  border-left: 4px solid var(--el-color-primary);
  padding-left: 10px;
}

.start-card {
  border-radius: 12px;
}

.start-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.8fr) minmax(180px, 0.9fr) auto;
  gap: 14px;
  align-items: center;
}

.start-input,
.start-select {
  width: 100%;
}

.start-btn {
  min-width: 132px;
}

.quick-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

.quick-btn {
  width: 160px;
  min-height: 90px;
  white-space: normal;
  line-height: 1.4;
  font-size: 1rem;
  border-radius: 10px;
  transition: transform 0.2s;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
}

.quick-btn-name {
  font-weight: 600;
}

.quick-btn-cat {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  margin-top: 6px;
}

.quick-btn:hover {
  transform: translateY(-2px);
}

.dashed {
  border-style: dashed;
  color: var(--el-text-color-secondary);
}

.search-input {
  margin-bottom: 20px;
  max-width: 400px;
}

.activity-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 20px;
}

.activity-item {
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s;
  border: 1px solid var(--el-border-color-lighter);
}

.activity-item:hover {
  transform: translateY(-2px);
  border-color: var(--el-color-primary-light-5);
}

.activity-inner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.act-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
  min-width: 0;
}

.act-name {
  font-weight: 600;
  font-size: 1.1rem;
  color: var(--el-text-color-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.act-meta {
  display: flex;
}

.act-action {
  display: flex;
  align-items: center;
}

.play-icon {
  font-size: 1.8rem;
  color: var(--el-color-primary);
  opacity: 0;
  transition: all 0.2s;
  background-color: var(--el-color-primary-light-9);
  padding: 8px;
  border-radius: 50%;
}

.activity-item:hover .play-icon {
  opacity: 1;
  transform: scale(1.1);
}

@media (max-width: 768px) {
  .start-grid {
    grid-template-columns: 1fr;
  }

  .start-btn {
    width: 100%;
  }
}
</style>
