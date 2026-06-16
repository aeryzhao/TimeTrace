<script setup>
import { ref, onMounted } from 'vue'
import { getCategories, createCategory, deleteCategory } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'

const categories = ref([])
const catForm = ref({ name: '', color: '#4F46E5' })

const predefinedColors = [
  '#4F46E5','#818CF8','#10B981','#F59E0B','#EF4444',
  '#06B6D4','#EC4899','#8B5CF6','#3B82F6','#F97316',
  '#14B8A6','#6366F1','#67C23A','#E6A23C','#F56C6C'
]

const fetchAll = async () => {
  const c = await getCategories()
  categories.value = c.data
}

onMounted(fetchAll)

const handleCreate = async () => {
  if (!catForm.value.name.trim()) return
  try {
    await createCategory(catForm.value)
    catForm.value.name = ''
    ElMessage.success('分类已创建')
    fetchAll()
  } catch {
    ElMessage.error('创建失败')
  }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定删除该分类吗？关联的活动可能受影响。', '确认删除', {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteCategory(id)
    ElMessage.success('已删除')
    fetchAll()
  } catch {}
}
</script>

<template>
  <div class="categories-page">
    <div class="page-header">
      <h1 class="page-title">分类管理</h1>
      <p class="page-subtitle">创建和管理活动分类，用颜色区分不同类型的工作</p>
    </div>

    <!-- Create form -->
    <div class="create-card">
      <div class="create-card-title">新建分类</div>
      <div class="create-form">
        <el-input
          v-model="catForm.name"
          placeholder="输入分类名称，例如：工作、学习、休息"
          size="large"
          clearable
          @keyup.enter="handleCreate"
          style="flex:1"
        />
        <el-color-picker
          v-model="catForm.color"
          show-alpha
          :predefine="predefinedColors"
          size="large"
        />
        <el-button type="primary" size="large" :icon="Plus" @click="handleCreate">
          添加分类
        </el-button>
      </div>
    </div>

    <!-- Category list -->
    <div class="list-header">
      <span class="list-label">已有分类</span>
      <span class="list-count">{{ categories.length }} 个</span>
    </div>

    <div v-if="categories.length === 0" class="empty-state">
      <svg width="36" height="36" fill="none" viewBox="0 0 24 24" stroke="#cbd5e1" stroke-width="1.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9.568 3H5.25A2.25 2.25 0 0 0 3 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 0 0 5.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 0 0 9.568 3Z" />
        <path stroke-linecap="round" stroke-linejoin="round" d="M6 6h.008v.008H6V6Z" />
      </svg>
      <p>还没有任何分类，新建一个吧</p>
    </div>

    <div class="cat-grid">
      <div v-for="cat in categories" :key="cat.id" class="cat-card">
        <div class="cat-color-strip" :style="{ background: cat.color }"></div>
        <div class="cat-info">
          <span class="cat-name">{{ cat.name }}</span>
          <span class="cat-id">#{{ cat.id }}</span>
        </div>
        <div
          class="cat-swatch"
          :style="{ background: cat.color }"
          :title="cat.color"
        ></div>
        <button class="delete-btn" @click="handleDelete(cat.id)" title="删除">
          <el-icon><Delete /></el-icon>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.categories-page { max-width: 760px; margin: 0 auto; }

.page-header { margin-bottom: 28px; }
.page-title { font-size: 1.75rem; font-weight: 800; color: #1e293b; margin: 0 0 4px; letter-spacing: -0.5px; }
.page-subtitle { margin: 0; color: #64748b; font-size: 0.9rem; }

/* Create card */
.create-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 20px 24px;
  margin-bottom: 28px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}
.create-card-title {
  font-size: 0.78rem;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 14px;
}
.create-form {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

/* List header */
.list-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;
}
.list-label {
  font-size: 0.78rem;
  font-weight: 700;
  color: #94a3b8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.list-count {
  font-size: 0.75rem;
  background: #f1f5f9;
  color: #64748b;
  padding: 2px 8px;
  border-radius: 100px;
  font-weight: 600;
}

/* Category grid */
.cat-grid { display: flex; flex-direction: column; gap: 8px; }
.cat-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  display: flex;
  align-items: center;
  overflow: hidden;
  transition: box-shadow 0.15s;
}
.cat-card:hover { box-shadow: 0 4px 10px rgba(0,0,0,0.07); }

.cat-color-strip { width: 4px; align-self: stretch; flex-shrink: 0; }
.cat-info {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 16px;
}
.cat-name { font-size: 0.95rem; font-weight: 600; color: #1e293b; }
.cat-id { font-size: 0.75rem; color: #cbd5e1; font-family: 'SF Mono','Consolas',monospace; }

.cat-swatch {
  width: 22px; height: 22px;
  border-radius: 50%;
  border: 2px solid #fff;
  box-shadow: 0 0 0 1px rgba(0,0,0,0.1);
  flex-shrink: 0;
  margin-right: 12px;
}

.delete-btn {
  display: flex; align-items: center; justify-content: center;
  width: 36px; height: 36px; margin-right: 10px;
  border-radius: 8px; border: 1px solid #e2e8f0;
  background: #fff; color: #94a3b8;
  cursor: pointer; transition: background 0.12s, color 0.12s, border-color 0.12s;
  font-size: 1rem;
}
.delete-btn:hover { background: #FEF2F2; color: #DC2626; border-color: #FECACA; }

.empty-state { text-align: center; padding: 60px 20px; color: #94a3b8; }
.empty-state p { margin: 12px 0 0; font-size: 0.9rem; }

@media (max-width: 560px) {
  .create-form { flex-direction: column; align-items: stretch; }
}
</style>
