<script setup>
import { ref, onMounted } from 'vue'
import { getCategories, createCategory, deleteCategory } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete, PriceTag, Collection } from '@element-plus/icons-vue'

const categories = ref([])

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

const fetchAll = async () => {
  const c = await getCategories()
  categories.value = c.data
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

</script>

<template>
  <div class="categories-page">
    <div class="page-title">
      <el-icon><Collection /></el-icon>
      <span>分类管理</span>
    </div>

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
  </div>
</template>

<style scoped>
.categories-page {
  max-width: 1000px;
  margin: 0 auto;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-size: 20px;
  font-weight: 700;
  color: var(--el-text-color-primary);
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
</style>
