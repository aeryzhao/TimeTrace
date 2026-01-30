<script setup>
import { ref, onMounted } from 'vue'
import { getCategories, createCategory, deleteCategory, getActivities, createActivity, deleteActivity, pinActivity } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete, PriceTag, Collection } from '@element-plus/icons-vue'

const categories = ref([])
const activities = ref([])
const activeTab = ref('categories')

// Predefined colors (Designer approved palette)
const predefinedColors = [
  '#409EFF', // Element Blue
  '#67C23A', // Element Green
  '#E6A23C', // Element Orange
  '#F56C6C', // Element Red
  '#909399', // Element Gray
  '#3B82F6', // Tailwind Blue 500
  '#10B981', // Tailwind Emerald 500
  '#F59E0B', // Tailwind Amber 500
  '#EF4444', // Tailwind Red 500
  '#8B5CF6', // Tailwind Violet 500
  '#EC4899', // Tailwind Pink 500
  '#06B6D4', // Tailwind Cyan 500
  '#6366F1', // Tailwind Indigo 500
  '#14B8A6', // Tailwind Teal 500
  '#F97316'  // Tailwind Orange 500
]

// Forms
const catForm = ref({ name: '', color: '#409EFF' })
const actForm = ref({ name: '', category_id: null, pinned: 0 })

const fetchAll = async () => {
  const [c, a] = await Promise.all([getCategories(), getActivities()])
  categories.value = c.data
  activities.value = a.data
}

onMounted(fetchAll)

const handleCreateCategory = async () => {
  if (!catForm.value.name) return
  try {
    await createCategory(catForm.value)
    catForm.value.name = ''
    ElMessage.success('分类已创建')
    fetchAll()
  } catch (e) {
    ElMessage.error('创建失败')
  }
}

const handleDeleteCategory = async (id) => {
  try {
    await ElMessageBox.confirm('确定删除该分类吗？关联的活动可能会受到影响。', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteCategory(id)
    ElMessage.success('已删除')
    fetchAll()
  } catch (e) {
    // Cancelled
  }
}

const handleCreateActivity = async () => {
  if (!actForm.value.name || !actForm.value.category_id) return
  try {
    await createActivity(actForm.value)
    actForm.value.name = ''
    ElMessage.success('活动已创建')
    fetchAll()
  } catch (e) {
    ElMessage.error('创建失败')
  }
}

const handleDeleteActivity = async (id) => {
  try {
    await ElMessageBox.confirm('确定删除该活动吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteActivity(id)
    ElMessage.success('已删除')
    fetchAll()
  } catch (e) {}
}

const handlePin = async (act) => {
  await pinActivity(act.id)
  act.pinned = act.pinned ? 0 : 1 // Optimistic update
}
</script>

<template>
  <div class="categories-page">
    <el-tabs v-model="activeTab" class="custom-tabs">
      <!-- Categories Tab -->
      <el-tab-pane name="categories">
        <template #label>
          <span class="custom-tab-label">
            <el-icon><Collection /></el-icon>
            <span>分类管理</span>
          </span>
        </template>
        
        <el-card class="box-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>新建分类</span>
            </div>
          </template>
          <div class="form-inline">
            <el-input 
              v-model="catForm.name" 
              placeholder="新分类名称" 
              class="form-input"
              clearable
            >
              <template #prefix>
                <el-icon><PriceTag /></el-icon>
              </template>
            </el-input>
            <el-color-picker v-model="catForm.color" show-alpha :predefine="predefinedColors" />
            <el-button type="primary" @click="handleCreateCategory" :icon="Plus" class="action-btn">添加</el-button>
          </div>
        </el-card>
        
        <el-card class="box-card list-card" shadow="never">
          <el-table :data="categories" stripe style="width: 100%" :header-cell-style="{ background: '#f5f7fa' }">
            <el-table-column prop="id" label="ID" width="80" align="center" />
            <el-table-column label="颜色" width="100" align="center">
              <template #default="{ row }">
                <div class="color-preview" :style="{ backgroundColor: row.color }"></div>
              </template>
            </el-table-column>
            <el-table-column prop="name" label="名称" />
            <el-table-column label="操作" align="right">
              <template #default="{ row }">
                <el-button type="danger" text bg size="small" :icon="Delete" @click="handleDeleteCategory(row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      <!-- Activities Tab -->
      <el-tab-pane name="activities">
        <template #label>
          <span class="custom-tab-label">
            <el-icon><PriceTag /></el-icon>
            <span>活动管理</span>
          </span>
        </template>
        
        <el-card class="box-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>新建活动</span>
            </div>
          </template>
          <div class="form-inline">
            <el-select v-model="actForm.category_id" placeholder="选择分类" class="form-select">
              <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id">
                <span class="option-color" :style="{ backgroundColor: c.color }"></span>
                {{ c.name }}
              </el-option>
            </el-select>
            <el-input 
              v-model="actForm.name" 
              placeholder="活动名称" 
              class="form-input"
              clearable
            />
            <el-button type="primary" @click="handleCreateActivity" :icon="Plus" class="action-btn">添加</el-button>
          </div>
        </el-card>

        <el-card class="box-card list-card" shadow="never">
          <el-table :data="activities" stripe style="width: 100%" :header-cell-style="{ background: '#f5f7fa' }">
            <el-table-column prop="id" label="ID" width="80" align="center" />
            <el-table-column label="分类" width="150">
              <template #default="{ row }">
                <el-tag :color="row.category?.color" effect="dark" class="category-tag">{{ row.category?.name }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="name" label="名称" />
            <el-table-column label="置顶" width="100" align="center">
              <template #default="{ row }">
                <el-switch 
                  :model-value="!!row.pinned" 
                  @change="handlePin(row)" 
                  active-color="#13ce66"
                />
              </template>
            </el-table-column>
            <el-table-column label="操作" align="right">
              <template #default="{ row }">
                <el-button type="danger" text bg size="small" :icon="Delete" @click="handleDeleteActivity(row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped>
.categories-page {
  max-width: 1000px;
  margin: 0 auto;
}

.custom-tabs :deep(.el-tabs__item) {
  font-size: 16px;
  padding: 0 20px;
}

.custom-tab-label {
  display: flex;
  align-items: center;
  gap: 5px;
}

.box-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.list-card {
  border: none;
  box-shadow: none;
}

.card-header {
  font-weight: bold;
  color: var(--el-text-color-primary);
}

.form-inline {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.form-input {
  width: 250px;
}

.form-select {
  width: 200px;
}

.action-btn {
  padding-left: 20px;
  padding-right: 20px;
}

.color-preview {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  margin: 0 auto;
  border: 1px solid #dcdfe6;
}

.category-tag {
  border: none;
  font-weight: 500;
}

.option-color {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-right: 8px;
  vertical-align: middle;
}
</style>
