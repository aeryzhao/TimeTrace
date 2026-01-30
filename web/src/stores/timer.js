import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getCurrentTimer, startTimer, stopTimer, getCategories, getActivities } from '../api'
import { ElMessage } from 'element-plus'

export const useTimerStore = defineStore('timer', () => {
  const currentEntry = ref(null)
  const categories = ref([])
  const activities = ref([])
  const loading = ref(false)

  // Computed
  const isRunning = computed(() => !!currentEntry.value)
  const pinnedActivities = computed(() => {
    return activities.value.filter(a => a.pinned === 1)
  })

  // Actions
  async function fetchMeta() {
    const [cRes, aRes] = await Promise.all([getCategories(), getActivities()])
    categories.value = cRes.data
    activities.value = aRes.data
  }

  async function fetchCurrentTimer() {
    try {
      const res = await getCurrentTimer()
      currentEntry.value = res.data // null or object
    } catch (e) {
      console.error(e)
    }
  }

  async function start(activityId, note = '') {
    try {
      loading.value = true
      const res = await startTimer(activityId, note)
      currentEntry.value = res.data
      ElMessage.success('已开始计时')
    } catch (e) {
      ElMessage.error('开始计时失败')
    } finally {
      loading.value = false
    }
  }

  async function stop() {
    try {
      loading.value = true
      await stopTimer()
      currentEntry.value = null
      ElMessage.success('已停止计时')
    } catch (e) {
      ElMessage.error('停止计时失败')
    } finally {
      loading.value = false
    }
  }

  return {
    currentEntry,
    categories,
    activities,
    loading,
    isRunning,
    pinnedActivities,
    fetchMeta,
    fetchCurrentTimer,
    start,
    stop
  }
})
