<script setup>
import { ref, onMounted, watch } from 'vue'
import { getTimeEntries, updateTimeEntry, deleteTimeEntry } from '../api'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'

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
    if (!end) return '进行中...'
    const diff = dayjs(end).diff(dayjs(start), 'second')
    const h = Math.floor(diff / 3600)
    const m = Math.floor((diff % 3600) / 60)
    return `${h}小时 ${m}分`
  }

  // Editing
  const dialogVisible = ref(false)
  const currentEdit = ref({})

  const handleEdit = (row) => {
    currentEdit.value = { ...row } // Copy
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
    } catch (e) {
      ElMessage.error('保存失败')
    }
  }

  const handleDelete = async (id) => {
    try {
      await ElMessageBox.confirm('确定删除这条记录吗？', '警告', { 
        type: 'warning',
        confirmButtonText: '删除',
        cancelButtonText: '取消'
      })
      await deleteTimeEntry(id)
      fetchData()
    } catch (e) {}
  }
</script>

<template>
  <div class="timeline-page">
    <div class="toolbar">
      <el-date-picker v-model="date" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" />
      <el-button @click="fetchData" icon="Refresh" circle></el-button>
    </div>

    <el-card shadow="never">
      <el-table :data="entries" v-loading="loading" style="width: 100%">
        <el-table-column label="时间" width="140">
          <template #default="{ row }">
            <div>{{ dayjs(row.start_time).format('HH:mm') }}</div>
            <div class="end-time" v-if="row.end_time">{{ dayjs(row.end_time).format('HH:mm') }}</div>
            <div v-else class="running">进行中</div>
          </template>
        </el-table-column>
        
        <el-table-column label="时长" width="120">
          <template #default="{ row }">
            {{ formatDuration(row.start_time, row.end_time) }}
          </template>
        </el-table-column>

        <el-table-column label="活动">
          <template #default="{ row }">
            <div class="activity-cell">
              <el-tag size="small" :color="row.category?.color" effect="dark" style="border:none; margin-right: 5px">
                {{ row.category?.name || '未分类' }}
              </el-tag>
              <span class="act-name">{{ row.activity?.name }}</span>
            </div>
            <div class="note" v-if="row.note">{{ row.note }}</div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" align="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)" icon="Edit" circle></el-button>
            <el-button size="small" type="danger" @click="handleDelete(row.id)" icon="Delete" circle></el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Edit Dialog -->
    <el-dialog v-model="dialogVisible" title="编辑记录" width="500px">
      <el-form :model="currentEdit" label-width="80px">
        <el-form-item label="开始时间">
          <el-time-picker v-model="currentEdit.start_time" format="HH:mm:ss" />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-time-picker v-model="currentEdit.end_time" format="HH:mm:ss" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="currentEdit.note" type="textarea" />
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
.toolbar {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}
.end-time {
  color: #909399;
  font-size: 0.9em;
}
.running {
  color: var(--el-color-success);
  font-weight: bold;
}
.act-name {
  font-weight: bold;
}
.note {
  font-size: 0.85em;
  color: #606266;
  margin-top: 4px;
}
</style>
