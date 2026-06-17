<script setup>
import { onMounted, watch } from 'vue'
import { useSettingsStore } from '../stores/settings'
import { Setting, Bell, Headset, Timer } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const settings = useSettingsStore()

onMounted(() => settings.fetchSettings())

watch(
  () => [settings.soundEnabled, settings.notificationEnabled, settings.reminderInterval],
  () => settings.saveSettings()
)

const testSound = () => {
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

const testNotification = () => {
  if (!('Notification' in window)) {
    ElMessage.warning('当前浏览器不支持通知')
    return
  }
  if (Notification.permission === 'denied') {
    ElMessage.warning('通知权限已被拒绝，请在浏览器设置中允许')
    return
  }
  if (Notification.permission === 'default') {
    Notification.requestPermission().then(perm => {
      if (perm === 'granted') {
        new Notification('TimeTrace', { body: '通知已开启，这是测试通知' })
      }
    })
    return
  }
  new Notification('TimeTrace', { body: '通知已开启，这是测试通知' })
}

const intervalOptions = [10, 15, 20, 25, 30, 45, 60]
</script>

<template>
  <div class="settings-page">
    <div class="page-header">
      <h1 class="page-title">设置</h1>
      <p class="page-subtitle">自定义计时提醒方式和间隔</p>
    </div>

    <div class="settings-card">
      <div class="card-header">
        <div>
          <h3 class="card-title">提醒设置</h3>
          <p class="card-desc">配置计时过程中的提醒方式</p>
        </div>
      </div>

      <div class="setting-row">
        <div class="setting-info">
          <el-icon><Headset /></el-icon>
          <div>
            <div class="setting-label">声音提醒</div>
            <div class="setting-desc">计时达到间隔时间时播放提示音</div>
          </div>
        </div>
        <div class="setting-actions">
          <el-button size="small" @click="testSound" :disabled="!settings.soundEnabled">试听</el-button>
          <el-switch v-model="settings.soundEnabled" />
        </div>
      </div>

      <div class="setting-row">
        <div class="setting-info">
          <el-icon><Bell /></el-icon>
          <div>
            <div class="setting-label">浏览器通知</div>
            <div class="setting-desc">计时达到间隔时间时发送系统通知</div>
          </div>
        </div>
        <div class="setting-actions">
          <el-button size="small" @click="testNotification" :disabled="!settings.notificationEnabled">测试</el-button>
          <el-switch v-model="settings.notificationEnabled" />
        </div>
      </div>

      <div class="setting-row">
        <div class="setting-info">
          <el-icon><Timer /></el-icon>
          <div>
            <div class="setting-label">提醒间隔</div>
            <div class="setting-desc">每隔多长时间提醒一次（分钟）</div>
          </div>
        </div>
        <div class="setting-actions">
          <el-select v-model="settings.reminderInterval" style="width: 120px">
            <el-option
              v-for="m in intervalOptions"
              :key="m"
              :label="`${m} 分钟`"
              :value="m"
            />
          </el-select>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page { max-width: 640px; margin: 0 auto; }

.page-header { margin-bottom: 28px; }
.page-title { font-size: 1.75rem; font-weight: 800; color: #1e293b; margin: 0 0 4px; letter-spacing: -0.5px; }
.page-subtitle { margin: 0; color: #64748b; font-size: 0.9rem; }

.settings-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 28px 32px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.06);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 28px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f1f5f9;
}

.card-icon {
  font-size: 1.5rem;
  background: #EEF2FF;
  padding: 10px;
  border-radius: 12px;
  color: #4F46E5;
}

.card-title { font-size: 1.1rem; font-weight: 700; color: #1e293b; margin: 0; }
.card-desc { font-size: 0.85rem; color: #94a3b8; margin: 4px 0 0; }

.setting-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px 0;
  border-bottom: 1px solid #f8fafc;
}
.setting-row:last-child { border-bottom: none; }

.setting-info {
  display: flex;
  align-items: center;
  gap: 14px;
}
.setting-info .el-icon { font-size: 1.2rem; color: #64748b; }

.setting-label { font-size: 0.95rem; font-weight: 600; color: #1e293b; }
.setting-desc { font-size: 0.82rem; color: #94a3b8; margin-top: 2px; }

.setting-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}
</style>
