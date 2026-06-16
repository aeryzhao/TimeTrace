<script setup>
import { ref, onMounted, watch } from 'vue'
import { getDailyReport } from '../api'
import dayjs from 'dayjs'
import * as echarts from 'echarts'
import { Calendar } from '@element-plus/icons-vue'

const date = ref(dayjs().format('YYYY-MM-DD'))
const reportData = ref(null)
const loading = ref(false)
const pieChartRef = ref(null)
const barChartRef = ref(null)
let pieChart = null
let barChart = null

const brand = '#4F46E5'
const brandPalette = ['#4F46E5','#818CF8','#10B981','#F59E0B','#EF4444','#06B6D4','#EC4899','#8B5CF6']

const renderCharts = () => {
  if (!reportData.value) return

  if (pieChartRef.value) {
    if (!pieChart) pieChart = echarts.init(pieChartRef.value)
    const data = reportData.value.by_category.map((c, i) => ({
      value: (c.duration / 3600).toFixed(2),
      name: c.name,
      itemStyle: { color: c.color || brandPalette[i % brandPalette.length] }
    }))
    pieChart.setOption({
      tooltip: { trigger: 'item', formatter: '{b}: {c}h ({d}%)' },
      legend: { bottom: 0, left: 'center', textStyle: { color: '#64748b', fontSize: 12 } },
      series: [{
        type: 'pie',
        radius: ['42%', '68%'],
        center: ['50%', '45%'],
        avoidLabelOverlap: true,
        itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
        label: { show: false },
        emphasis: { label: { show: true, fontSize: 16, fontWeight: 700 } },
        data
      }]
    })
  }

  if (barChartRef.value) {
    if (!barChart) barChart = echarts.init(barChartRef.value)
    const sorted = [...reportData.value.by_activity]
      .sort((a, b) => b.duration - a.duration)
      .slice(0, 10)
    barChart.setOption({
      tooltip: { trigger: 'axis', axisPointer: { type: 'none' } },
      grid: { left: 0, right: 16, bottom: 40, top: 16, containLabel: true },
      xAxis: {
        type: 'category',
        data: sorted.map(i => i.name),
        axisLabel: { interval: 0, rotate: 30, color: '#64748b', fontSize: 11 },
        axisLine: { lineStyle: { color: '#e2e8f0' } }
      },
      yAxis: {
        type: 'value',
        axisLabel: { color: '#94a3b8', fontSize: 11 },
        splitLine: { lineStyle: { color: '#f1f5f9' } }
      },
      series: [{
        data: sorted.map(i => (i.duration / 3600).toFixed(2)),
        type: 'bar',
        barMaxWidth: 48,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#818CF8' },
            { offset: 1, color: brand }
          ]),
          borderRadius: [6, 6, 0, 0]
        }
      }]
    })
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getDailyReport(date.value)
    reportData.value = res.data
    setTimeout(renderCharts, 50)
  } finally {
    loading.value = false
  }
}

watch(date, fetchData)
onMounted(() => {
  fetchData()
  window.addEventListener('resize', () => {
    pieChart?.resize()
    barChart?.resize()
  })
})

const totalHours = () => (reportData.value?.total_duration / 3600).toFixed(1)
</script>

<template>
  <div class="reports-page">
    <div class="page-header">
      <h1 class="page-title">报表</h1>
      <p class="page-subtitle">每日时间分布与活动统计</p>
    </div>

    <div class="toolbar">
      <el-date-picker
        v-model="date"
        type="date"
        value-format="YYYY-MM-DD"
        placeholder="选择日期"
        :clearable="false"
        size="large"
      >
        <template #prefix><el-icon><Calendar /></el-icon></template>
      </el-date-picker>
    </div>

    <div v-if="reportData" v-loading="loading">
      <!-- Summary stat -->
      <div class="stat-row">
        <div class="stat-card">
          <div class="stat-label">总时长</div>
          <div class="stat-value">{{ totalHours() }}<span class="stat-unit">h</span></div>
        </div>
        <div class="stat-card">
          <div class="stat-label">分类数</div>
          <div class="stat-value">{{ reportData.by_category.length }}<span class="stat-unit">个</span></div>
        </div>
        <div class="stat-card">
          <div class="stat-label">活动数</div>
          <div class="stat-value">{{ reportData.by_activity.length }}<span class="stat-unit">项</span></div>
        </div>
      </div>

      <!-- Charts -->
      <div class="charts-grid">
        <div class="chart-card">
          <div class="chart-title">分类分布</div>
          <div ref="pieChartRef" class="chart"></div>
        </div>
        <div class="chart-card">
          <div class="chart-title">热门活动 Top 10</div>
          <div ref="barChartRef" class="chart"></div>
        </div>
      </div>
    </div>

    <div v-else-if="!loading" class="empty-state">
      <svg width="40" height="40" fill="none" viewBox="0 0 24 24" stroke="#cbd5e1" stroke-width="1.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z" />
      </svg>
      <p>当天暂无数据</p>
    </div>
  </div>
</template>

<style scoped>
.reports-page { max-width: 1080px; margin: 0 auto; }

.page-header { margin-bottom: 28px; }
.page-title { font-size: 1.75rem; font-weight: 800; color: #1e293b; margin: 0 0 4px; letter-spacing: -0.5px; }
.page-subtitle { margin: 0; color: #64748b; font-size: 0.9rem; }

.toolbar { margin-bottom: 28px; }

.stat-row { display: flex; gap: 16px; margin-bottom: 24px; flex-wrap: wrap; }
.stat-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 20px 28px;
  min-width: 140px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}
.stat-label { font-size: 0.78rem; font-weight: 600; color: #94a3b8; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 8px; }
.stat-value { font-size: 2rem; font-weight: 800; color: #1e293b; line-height: 1; }
.stat-unit { font-size: 1rem; font-weight: 500; color: #94a3b8; margin-left: 3px; }

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(420px, 1fr));
  gap: 20px;
}
.chart-card {
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 20px 24px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}
.chart-title { font-size: 0.85rem; font-weight: 700; color: #374151; margin-bottom: 16px; text-transform: uppercase; letter-spacing: 0.4px; }
.chart { height: 380px; width: 100%; }

.empty-state { text-align: center; padding: 80px 20px; color: #94a3b8; }
.empty-state p { margin: 12px 0 0; font-size: 0.9rem; }

@media (max-width: 640px) {
  .charts-grid { grid-template-columns: 1fr; }
}
</style>
