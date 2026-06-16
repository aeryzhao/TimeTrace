<script setup>
import { ref, onMounted, watch } from 'vue'
import { getTimeEntries, updateTimeEntry, deleteTimeEntry } from '../api'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, Delete, Refresh } from '@element-plus/icons-vue'

const date = ref(dayjs().format('YYYY-MM-DD'))
const entries = ref([])
const loading = ref(false)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getTimeEntries(date.value, date.value)
    entries.value = res.data
  } finally {
    loading.value = false
  }
}

watch(date, fetchData)
onMounted(fetchData)

const formatDuration = (start, end) => {
  if (!end) return '进行中'
  const diff = dayjs(end).diff(dayjs(start), 'second')
  const h = Math.floor(diff / 3600)
  const m = Math.floor((diff % 3600) / 60)
  return h > 0 ? `${h}h ${m}m` : `${m}m`
}

const dialogVisible = ref(false)
const currentEdit = ref({})

const handleEdit = (row) => {
  currentEdit.value = { ...row }
  dialogVisible.value = true
}

const saveEdit = async () => {
  try {
    await updateTimeEntry(currentEdit.value.id, {
      start_time: currentEdit.value.start_time,
      end_time: currentEdit.value.end_time,
      note: currentEdit.value.note
    })
    ElMessage.success('已保存')
    dialogVisible.value = false
    fetchData()
  } catch {
    ElMessage.error('保存失败')
  }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定删除这条记录吗？', '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消'
    })
    await deleteTimeEntry(id)
    ElMessage.success('已删除')
    fetchData()
  } catch {}
}
</script>

<template>
  <div class="timeline-page">
    <div class="page-header">
      <h1 class="page-title">时间轴</h1>
      <p class="page-subtitle">查看和管理每天的时间记录</p>
    </div>

    <div class="toolbar">
      <el-date-picker
        v-model="date"
        type="date"
        value-format="YYYY-MM-DD"
        placeholder="选择日期"
        :clearable="false"
        size="large"
      />
      <el-button :icon="Refresh" circle @click="fetchData" size="large" />
    </div>

    <div class="entries-list" v-loading="loading">
      <div v-if="entries.length === 0 && !loading" class="empty-state">
        <div class="empty-icon">
          <svg width="40" height="40" fill="none" viewBox="0 0 24 24" stroke="#cbd5e1" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
          </svg>
        </div>
        <p>当天暂无记录</p>
      </div>

      <div v-for="entry in entries" :key="entry.id" class="entry-row">
        <div class="entry-timeline">
          <div class="entry-dot" :style="{ background: entry.category?.color || '#818CF8' }"></div>
          <div class="entry-line"></div>
        </div>
        <div class="entry-card">
          <div class="entry-header">
            <div class="entry-time-block">
              <span class="entry-start">{{ dayjs(entry.start_time).format('HH:mm') }}</span>
              <span class="entry-arrow">→</span>
              <span class="entry-end" v-if="entry.end_time">{{ dayjs(entry.end_time).format('HH:mm') }}</span>
              <span class="entry-running" v-else>进行中</span>
            </div>
            <span class="entry-duration">{{ formatDuration(entry.start_time, entry.end_time) }}</span>
          </div>
          <div class="entry-body">
            <span
              class="cat-badge"
              :style="{ background: entry.category?.color || '#818CF8' }"
            >{{ entry.category?.name || '未分类' }}</span>
            <span class="entry-act-name">{{ entry.activity?.name }}</span>
          </div>
          <div class="entry-note" v-if="entry.note">{{ entry.note }}</div>
          <div class="entry-actions">
            <button class="action-btn" @click="handleEdit(entry)">
              <el-icon><Edit /></el-icon> 编辑
            </button>
            <button class="action-btn danger" @click="handleDelete(entry.id)">
              <el-icon><Delete /></el-icon> 删除
            </button>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="dialogVisible" title="编辑记录" width="440px" :close-on-click-modal="false">
      <el-form :model="currentEdit" label-width="80px" style="padding: 0 8px">
        <el-form-item label="开始时间">
          <el-time-picker v-model="currentEdit.start_time" format="HH:mm:ss" style="width:100%" />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-time-picker v-model="currentEdit.end_time" format="HH:mm:ss" style="width:100%" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="currentEdit.note" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.timeline-page { max-width: 760px; margin: 0 auto; }

.page-header { margin-bottom: 28px; }
.page-title { font-size: 1.75rem; font-weight: 800; color: #1e293b; margin: 0 0 4px; letter-spacing: -0.5px; }
.page-subtitle { margin: 0; color: #64748b; font-size: 0.9rem; }

.toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 28px;
}

.entries-list { display: flex; flex-direction: column; }

.entry-row {
  display: flex;
  gap: 0;
  margin-bottom: 4px;
}

.entry-timeline {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 32px;
  flex-shrink: 0;
  padding-top: 18px;
}
.entry-dot {
  width: 12px; height: 12px;
  border-radius: 50%;
  border: 2px solid #fff;
  box-shadow: 0 0 0 2px currentColor;
  flex-shrink: 0;
  z-index: 1;
}
.entry-line {
  width: 2px;
  flex: 1;
  background: #e2e8f0;
  margin-top: 4px;
  min-height: 20px;
}

.entry-card {
  flex: 1;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 14px 18px;
  margin-left: 8px;
  margin-bottom: 10px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
  transition: box-shadow 0.15s;
}
.entry-card:hover { box-shadow: 0 4px 10px rgba(0,0,0,0.07); }

.entry-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}
.entry-time-block { display: flex; align-items: center; gap: 6px; font-size: 0.85rem; }
.entry-start { font-weight: 700; color: #1e293b; font-family: 'SF Mono','Consolas',monospace; }
.entry-arrow { color: #cbd5e1; font-size: 0.75rem; }
.entry-end { color: #64748b; font-family: 'SF Mono','Consolas',monospace; }
.entry-running { color: #10B981; font-weight: 700; font-size: 0.8rem; }
.entry-duration {
  font-size: 0.8rem; font-weight: 700;
  background: #f1f5f9; color: #475569;
  padding: 2px 8px; border-radius: 100px;
}

.entry-body { display: flex; align-items: center; gap: 10px; margin-bottom: 6px; }
.cat-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 100px;
  font-size: 0.72rem;
  font-weight: 600;
  color: #fff;
  flex-shrink: 0;
}
.entry-act-name { font-size: 0.95rem; font-weight: 600; color: #1e293b; }

.entry-note { font-size: 0.82rem; color: #64748b; margin-bottom: 8px; padding: 6px 8px; background: #f8fafc; border-radius: 6px; }

.entry-actions { display: flex; gap: 8px; }
.action-btn {
  display: inline-flex; align-items: center; gap: 4px;
  font-size: 0.78rem; font-weight: 500;
  padding: 4px 10px; border-radius: 6px;
  border: 1px solid #e2e8f0; background: #fff;
  color: #64748b; cursor: pointer;
  transition: background 0.12s, color 0.12s, border-color 0.12s;
  font-family: inherit;
}
.action-btn:hover { background: #f1f5f9; color: #1e293b; }
.action-btn.danger:hover { background: #FEF2F2; color: #DC2626; border-color: #FECACA; }

.empty-state { text-align: center; padding: 60px 20px; color: #94a3b8; }
.empty-icon { margin-bottom: 12px; }
.empty-state p { margin: 0; font-size: 0.9rem; }
</style>
