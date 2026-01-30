<script setup>
import { ref, onMounted, watch } from 'vue'
import { getDailyReport } from '../api'
import dayjs from 'dayjs'
import * as echarts from 'echarts'
import { Calendar, Timer } from '@element-plus/icons-vue'

const date = ref(dayjs().format('YYYY-MM-DD'))
const reportData = ref(null)
const loading = ref(false)

const pieChartRef = ref(null)
const barChartRef = ref(null)
let pieChart = null
let barChart = null

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getDailyReport(date.value)
    reportData.value = res.data
    renderCharts()
  } finally {
    loading.value = false
  }
}

const renderCharts = () => {
  if (!reportData.value) return

  // Pie Chart (Categories)
  if (pieChartRef.value) {
    if (!pieChart) pieChart = echarts.init(pieChartRef.value)
    const data = reportData.value.by_category.map(c => ({
      value: (c.duration / 3600).toFixed(2), // Hours
      name: c.name,
      itemStyle: { color: c.color || null }
    }))
    
    pieChart.setOption({
      title: { text: '分类时间分布 (小时)', left: 'center' },
      tooltip: { trigger: 'item', formatter: '{b}: {c}小时 ({d}%)' },
      legend: { bottom: '0', left: 'center' },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 20,
            fontWeight: 'bold'
          }
        },
        data
      }]
    })
  }

  // Bar Chart (Activities)
  if (barChartRef.value) {
    if (!barChart) barChart = echarts.init(barChartRef.value)
    // Top 10
    const sorted = [...reportData.value.by_activity].sort((a, b) => b.duration - a.duration).slice(0, 10)
    
    barChart.setOption({
      title: { text: '热门活动 (小时)', left: 'center' },
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { 
        type: 'category', 
        data: sorted.map(i => i.name), 
        axisLabel: { interval: 0, rotate: 30 } 
      },
      yAxis: { type: 'value' },
      series: [{
        data: sorted.map(i => (i.duration / 3600).toFixed(2)),
        type: 'bar',
        showBackground: true,
        backgroundStyle: {
          color: 'rgba(180, 180, 180, 0.2)'
        },
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#83bff6' },
            { offset: 0.5, color: '#188df0' },
            { offset: 1, color: '#188df0' }
          ])
        }
      }]
    })
  }
}

watch(date, fetchData)

onMounted(() => {
  fetchData()
  window.addEventListener('resize', () => {
    pieChart && pieChart.resize()
    barChart && barChart.resize()
  })
})
</script>

<template>
  <div class="reports-page">
    <el-card class="toolbar-card" shadow="never">
      <div class="toolbar-content">
        <span class="toolbar-label">
          <el-icon><Calendar /></el-icon>
          统计日期
        </span>
        <el-date-picker 
          v-model="date" 
          type="date" 
          value-format="YYYY-MM-DD" 
          placeholder="选择日期"
          :clearable="false"
        />
      </div>
    </el-card>

    <div v-if="reportData" class="dashboard-grid">
      <el-card class="summary-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>总览</span>
          </div>
        </template>
        <div class="summary-content">
          <el-statistic title="总时长" :value="(reportData.total_duration / 3600).toFixed(2)">
            <template #suffix>小时</template>
          </el-statistic>
          <div class="summary-icon">
            <el-icon><Timer /></el-icon>
          </div>
        </div>
      </el-card>

      <div class="charts-container" v-loading="loading">
        <el-card shadow="hover" class="chart-card">
          <div ref="pieChartRef" class="chart"></div>
        </el-card>
        <el-card shadow="hover" class="chart-card">
          <div ref="barChartRef" class="chart"></div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<style scoped>
.reports-page {
  max-width: 1200px;
  margin: 0 auto;
}

.toolbar-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.toolbar-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.toolbar-label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-weight: bold;
  color: var(--el-text-color-regular);
}

.summary-card {
  margin-bottom: 20px;
  border-radius: 8px;
  width: 300px;
}

.summary-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-icon {
  font-size: 48px;
  color: var(--el-color-primary-light-5);
}

.charts-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
  gap: 20px;
}

.chart-card {
  border-radius: 8px;
}

.chart {
  height: 400px;
  width: 100%;
}
</style>
