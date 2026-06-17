<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useTimerStore } from '../stores/timer'
import { useSettingsStore } from '../stores/settings'
import { VideoPlay, VideoPause, Plus, EditPen, Search } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

const store = useTimerStore()
const settings = useSettingsStore()
const elapsedTime = ref('00:00:00')
const quickForm = ref({ name: '', category_id: null })
const searchQuery = ref('')
let timerInterval = null
let lastReminderMinute = 0

const filteredActivities = computed(() => {
  if (!searchQuery.value) return store.activities
  const q = searchQuery.value.toLowerCase()
  return store.activities.filter(a =>
    a.name.toLowerCase().includes(q) || a.category?.name?.toLowerCase().includes(q)
  )
})

const quickSuggestions = computed(() => store.activities.slice(0, 8))

const playReminderSound = () => {
  const audioContext = new (window.AudioContext || window.webkitAudioContext)()
  const oscillator = audioContext.createOscillator()
  const gainNode = audioContext.createGain()
  oscillator.connect(gainNode)
  gainNode.connect(audioContext.destination)
  oscillator.frequency.setValueAtTime(800, audioContext.currentTime)
  oscillator.frequency.setValueAtTime(600, audioContext.currentTime + 0.1)
  oscillator.frequency.setValueAtTime(800, audioContext.currentTime + 0.2)
  gainNode.gain.setValueAtTime(0.3, audioContext.currentTime)
  gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.3)
  oscillator.start(audioContext.currentTime)
  oscillator.stop(audioContext.currentTime + 0.3)
}

const updateTicker = () => {
  if (store.currentEntry?.start_time) {
    const diff = dayjs().diff(dayjs(store.currentEntry.start_time), 'second')
    const h = Math.floor(diff / 3600)
    const m = Math.floor((diff % 3600) / 60)
    const s = diff % 60
    elapsedTime.value = `${String(h).padStart(2,'0')}:${String(m).padStart(2,'0')}:${String(s).padStart(2,'0')}`
    
    const activityName = store.currentEntry.activity?.name || '计时中'
    const totalMinutes = h * 60 + m
    document.title = `${String(totalMinutes).padStart(2,'0')}:${String(s).padStart(2,'0')} - ${activityName} | TimeTrace`
    
    const interval = settings.reminderInterval
    if (m > 0 && interval > 0 && m % interval === 0 && s === 0 && m !== lastReminderMinute) {
      lastReminderMinute = m
      if (settings.soundEnabled) playReminderSound()
      const totalMin = h * 60 + m
      if (settings.notificationEnabled && Notification.permission === 'granted') {
        new Notification('TimeTrace 提醒', {
          body: `「${activityName}」已进行 ${totalMin} 分钟`,
          icon: '/favicon.svg'
        })
      }
    }
  } else {
    elapsedTime.value = '00:00:00'
    document.title = 'TimeTrace - 时间追踪'
    lastReminderMinute = 0
  }
}

onMounted(async () => {
  await Promise.all([store.fetchMeta(), store.fetchCurrentTimer(), settings.fetchSettings()])
  timerInterval = setInterval(updateTicker, 1000)
  updateTicker()
  document.title = 'TimeTrace - 时间追踪'
  if (settings.notificationEnabled && 'Notification' in window && Notification.permission === 'default') {
    Notification.requestPermission()
  }
})

onUnmounted(() => clearInterval(timerInterval))

const handleStart = (activity) => store.start(activity.name, activity.category_id || null)

const handleQuickStart = async () => {
  if (!quickForm.value.name.trim()) return
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
    <!-- Page header -->
    <div class="page-header">
      <h1 class="page-title">计时台</h1>
      <p class="page-subtitle">记录你的每一段专注时间</p>
    </div>

    <!-- Timer hero card -->
    <div class="timer-hero" :class="{ running: store.isRunning }">
      <div class="timer-left">
        <div class="timer-status-badge" :class="store.isRunning ? 'badge-running' : 'badge-idle'">
          <span class="status-dot"></span>
          {{ store.isRunning ? '进行中' : '待机' }}
        </div>
        <div v-if="store.isRunning" class="timer-activity-name">
          {{ store.currentEntry.activity?.name }}
        </div>
        <div v-else class="timer-idle-hint">还没有正在进行的活动</div>
        <div v-if="store.isRunning" class="timer-meta">
          <span class="cat-pill" :style="{ background: store.currentEntry.category?.color || '#818CF8' }">
            {{ currentCategoryName }}
          </span>
          <span class="start-hint">{{ dayjs(store.currentEntry.start_time).format('HH:mm') }} 开始</span>
        </div>
      </div>
      <div class="timer-right">
        <div class="timer-display">{{ elapsedTime }}</div>
        <el-button
          v-if="store.isRunning"
          class="timer-btn stop-btn"
          circle
          :icon="VideoPause"
          @click="store.stop"
        />
      </div>
    </div>

    <!-- Quick start -->
    <div class="section">
      <h2 class="section-title">开始活动</h2>
      <div class="start-card">
        <el-input
          v-model="quickForm.name"
          placeholder="输入当前要做的事情，例如：写方案、开会、读书"
          size="large"
          clearable
          @keyup.enter="handleQuickStart"
        >
          <template #prefix><el-icon><EditPen /></el-icon></template>
        </el-input>
        <el-select
          v-model="quickForm.category_id"
          placeholder="选择分类（可选）"
          clearable
          filterable
          size="large"
        >
          <el-option
            v-for="cat in store.categories"
            :key="cat.id"
            :label="cat.name"
            :value="cat.id"
          />
        </el-select>
        <el-button type="primary" size="large" :icon="VideoPlay" @click="handleQuickStart" class="go-btn">
          开始计时
        </el-button>
      </div>
    </div>

    <!-- Quick suggestions -->
    <div class="section" v-if="quickSuggestions.length">
      <h2 class="section-title">最近活动</h2>
      <div class="quick-grid">
        <button
          v-for="act in quickSuggestions"
          :key="act.id"
          class="quick-card"
          @click="handleStart(act)"
        >
          <span
            class="quick-dot"
            :style="{ background: act.category?.color || '#818CF8' }"
          ></span>
          <span class="quick-name">{{ act.name }}</span>
          <span class="quick-cat">{{ act.category?.name || '未分类' }}</span>
        </button>
        <button class="quick-card dashed" @click="$router.push('/categories')">
          <el-icon class="plus-icon"><Plus /></el-icon>
          <span class="quick-name">管理分类</span>
        </button>
      </div>
    </div>

    <!-- All activities -->
    <div class="section">
      <div class="section-row">
        <h2 class="section-title">所有活动</h2>
        <el-input
          v-model="searchQuery"
          placeholder="搜索活动..."
          size="default"
          style="width:240px"
        >
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
      </div>
      <div class="activity-grid">
        <div
          v-for="act in filteredActivities"
          :key="act.id"
          class="activity-card"
          @click="handleStart(act)"
        >
          <div
            class="act-color-bar"
            :style="{ background: act.category?.color || '#818CF8' }"
          ></div>
          <div class="act-body">
            <span class="act-name">{{ act.name }}</span>
            <span class="act-cat">{{ act.category?.name || '未分类' }}</span>
          </div>
          <el-icon class="act-play"><VideoPlay /></el-icon>
        </div>
      </div>
      <div v-if="filteredActivities.length === 0" class="empty-state">
        <el-icon style="font-size:2rem;color:#cbd5e1"><Search /></el-icon>
        <p>没有找到匹配的活动</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard { max-width: 960px; margin: 0 auto; }

.page-header { margin-bottom: 28px; }
.page-title { font-size: 1.75rem; font-weight: 800; color: #1e293b; margin: 0 0 4px; letter-spacing: -0.5px; }
.page-subtitle { margin: 0; color: #64748b; font-size: 0.9rem; }

/* Timer hero */
.timer-hero {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 28px 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 36px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.06);
  transition: border-color 0.3s, box-shadow 0.3s;
}
.timer-hero.running {
  border-color: #a5b4fc;
  box-shadow: 0 0 0 4px #eef2ff, 0 4px 12px rgba(79,70,229,0.1);
}

.timer-left { display: flex; flex-direction: column; gap: 10px; }

.timer-status-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.8rem;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: 100px;
  width: fit-content;
  letter-spacing: 0.3px;
}
.badge-running { background: #d1fae5; color: #065f46; }
.badge-idle { background: #f1f5f9; color: #64748b; }
.status-dot {
  width: 7px; height: 7px; border-radius: 50%;
  background: currentColor;
}
.badge-running .status-dot { animation: pulse 1.4s infinite; }
@keyframes pulse {
  0%, 100% { opacity: 1; } 50% { opacity: 0.3; }
}

.timer-activity-name { font-size: 1.6rem; font-weight: 700; color: #1e293b; line-height: 1.2; }
.timer-idle-hint { font-size: 1rem; color: #94a3b8; }

.timer-meta { display: flex; align-items: center; gap: 10px; }
.cat-pill {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 100px;
  font-size: 0.78rem;
  font-weight: 600;
  color: #fff;
}
.start-hint { font-size: 0.85rem; color: #64748b; }

.timer-right { display: flex; flex-direction: column; align-items: center; gap: 16px; }
.timer-display {
  font-family: 'SF Mono', 'Monaco', 'Consolas', monospace;
  font-size: 3.2rem;
  font-weight: 700;
  color: #4F46E5;
  letter-spacing: 2px;
  line-height: 1;
}

.timer-btn { width: 52px; height: 52px; font-size: 1.4rem; }
.stop-btn {
  --el-button-bg-color: #FEF2F2 !important;
  --el-button-border-color: #FECACA !important;
  --el-button-text-color: #DC2626 !important;
  --el-button-hover-bg-color: #FEE2E2 !important;
}

/* Section */
.section { margin-bottom: 40px; }
.section-title {
  font-size: 1rem;
  font-weight: 700;
  color: #374151;
  margin: 0 0 16px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  font-size: 0.8rem;
}
.section-row { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.section-row .section-title { margin-bottom: 0; }

/* Start card */
.start-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 20px;
  display: grid;
  grid-template-columns: 1fr auto auto;
  gap: 12px;
  align-items: center;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}
.go-btn { min-width: 120px; }

/* Quick grid */
.quick-grid { display: flex; flex-wrap: wrap; gap: 10px; }
.quick-card {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 6px;
  padding: 14px 16px;
  width: 150px;
  min-height: 80px;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  cursor: pointer;
  transition: border-color 0.15s, box-shadow 0.15s, transform 0.15s;
  text-align: left;
  font-family: inherit;
}
.quick-card:hover {
  border-color: #a5b4fc;
  box-shadow: 0 4px 12px rgba(79,70,229,0.1);
  transform: translateY(-2px);
}
.quick-card.dashed { border-style: dashed; border-color: #cbd5e1; }
.quick-card.dashed:hover { border-color: #818cf8; }
.quick-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
.quick-name { font-size: 0.88rem; font-weight: 600; color: #1e293b; line-height: 1.3; }
.quick-cat { font-size: 0.75rem; color: #94a3b8; }
.plus-icon { font-size: 1.2rem; color: #94a3b8; margin-bottom: 2px; }

/* Activity grid */
.activity-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}
.activity-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 14px 14px 0;
  cursor: pointer;
  transition: border-color 0.15s, box-shadow 0.15s, transform 0.15s;
}
.activity-card:hover {
  border-color: #a5b4fc;
  box-shadow: 0 4px 12px rgba(79,70,229,0.1);
  transform: translateY(-1px);
}
.activity-card:hover .act-play { opacity: 1; color: #4F46E5; }

.act-color-bar { width: 4px; align-self: stretch; flex-shrink: 0; border-radius: 0 2px 2px 0; }
.act-body { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 4px; }
.act-name { font-size: 0.9rem; font-weight: 600; color: #1e293b; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.act-cat { font-size: 0.75rem; color: #94a3b8; }
.act-play { font-size: 1.2rem; opacity: 0; transition: opacity 0.15s; color: #cbd5e1; flex-shrink: 0; }

.empty-state { text-align: center; padding: 40px; color: #94a3b8; }
.empty-state p { margin: 8px 0 0; font-size: 0.9rem; }

@media (max-width: 640px) {
  .timer-hero { flex-direction: column; gap: 20px; }
  .timer-display { font-size: 2.4rem; }
  .start-card { grid-template-columns: 1fr; }
}
</style>
