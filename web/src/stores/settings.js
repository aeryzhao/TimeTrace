import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getSettings, updateSettings } from '../api'

export const useSettingsStore = defineStore('settings', () => {
  const soundEnabled = ref(true)
  const notificationEnabled = ref(true)
  const reminderInterval = ref(25)
  const loading = ref(false)

  async function fetchSettings() {
    try {
      loading.value = true
      const res = await getSettings()
      soundEnabled.value = res.data.sound_enabled
      notificationEnabled.value = res.data.notification_enabled
      reminderInterval.value = res.data.reminder_interval
    } catch (e) {
      console.error('Failed to fetch settings:', e)
    } finally {
      loading.value = false
    }
  }

  async function saveSettings() {
    try {
      await updateSettings({
        sound_enabled: soundEnabled.value,
        notification_enabled: notificationEnabled.value,
        reminder_interval: reminderInterval.value,
      })
    } catch (e) {
      console.error('Failed to save settings:', e)
    }
  }

  return {
    soundEnabled,
    notificationEnabled,
    reminderInterval,
    loading,
    fetchSettings,
    saveSettings,
  }
})
