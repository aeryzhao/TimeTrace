# PRD：时间花费记录工具（Lyubishchev Time Log）Web MVP

## 0. 产品定位
- **定位**：面向个人的“实时计时 + 结构化分类 + 复盘统计”的时间日志工具。
- **核心价值**：用极低操作成本记录“我现在在做什么”，并在日/周维度看见真实时间分配。
- **MVP 原则**：记录链路极短、数据结构稳、统计够用、可持续迭代。

---

## 1. 目标与指标

### 1.1 MVP 目标
1) 3 秒内开始计时（从打开页面到开始一个活动）
2) 1 秒内完成活动切换（自动结束上一段）
3) 支持当天补修正（编辑/拆分/合并）
4) 有“日报”统计和“周报”统计（不含目标）

### 1.2 验收标准（可测试）
- 能创建领域分类体系（类目/活动）
- 任意时刻只有**一个 running 记录**
- 切换活动时：上一条自动填 end_time，新条 start_time=now
- 日视图能看到时间轴，且统计总时长一致
- SQLite 数据持久化，服务重启后计时状态可恢复（running 的处理见技术点）

---

## 2. 用户与使用场景

### 2.1 目标用户
- 需要深度工作/多项目切换的个人（研发、产品、写作者、独立开发者等）

### 2.2 核心场景
- **开始**：上班打开网页 → 点击“编码”开始计时
- **切换**：来会了 → 点击“会议” → 自动结束“编码”并开始“会议”
- **结束**：午休/下班 → 点“停止”
- **补改**：发现某段分类错了 → 在时间轴编辑成正确活动
- **复盘**：晚上看日报、周末看周报

---

## 3. 信息架构与页面清单（Vue + Element UI）

### 3.1 页面结构
1) **首页 / 计时台（Dashboard）**
2) **时间轴（日视图）（Timeline Day）**
3) **统计（日报/周报）（Reports）**
4) **分类管理（Categories & Activities）**
5) **设置（Settings：时区、一天起始、数据导出等）**

---

## 4. 核心功能需求（MVP）

## 4.1 实时计时（核心主链路）
### 功能点
- 当前计时状态展示：
  - 当前活动名、所属类目、开始时间、已持续时长
- 操作按钮：
  - **开始/切换活动**：从“常用活动列表”一键开始
  - **停止**：结束当前记录（写入 end_time）
- 常用活动：
  - 最近使用（最近10个）
  - 置顶（用户在分类管理中设定 pin）

### 交互要求（降低成本）
- 点击活动卡片即“切换并开始”，不再二次确认
- 可选快速备注（不强制）：开始后可在当前记录上点“添加备注”

### 规则
- 系统全局同一时间最多一个 running 段
- running 段没有 end_time
- 如果用户开始一个活动且已有 running：
  - 自动结束旧 running（end_time=now）
  - 创建新 running（start_time=now）

---

## 4.2 时间轴（日视图）与编辑
### 展示
- 以当天为单位的记录列表（按时间倒序/正序切换）
- 显示：起止时间、时长、活动、类目、备注

### 编辑能力（MVP 必要）
- 修改活动/类目（改错）
- 修改开始/结束时间
- 拆分（Split）：把一段拆成两段（常见：中间被打断）
- 合并（Merge）：相邻且同活动的两段合并
- 删除（Delete）：删除一段记录

### 校验规则
- start_time < end_time
- 不允许同一日内出现记录重叠（保存时校验；如重叠给出提示并阻止保存）
- 可允许跨天记录吗？
  - **MVP 建议：允许，但展示按 start_time 所在日归属**（实现简单）
  - 或者强制切分跨天段（更严谨但复杂）。MVP 先不强制切分。

---

## 4.3 分类体系（领域化）
### 数据结构
- **类目 Category**：领域大类（如：研发/产品/运营/学习/生活维护）
- **活动 Activity**：可直接点击开始计时的最小单元（如：编码/评审/写方案/读论文）

### 管理功能（MVP）
- 类目 CRUD
- 活动 CRUD（绑定类目）
- 活动置顶 pin（出现在首页快捷）
- 活动颜色（用于统计图表，可选）

### 预置建议（首次进入引导）
提供“知识工作者领域模板”，用户可一键导入再改：
- 研发：编码/调试/Review/技术方案
- 产品：写PRD/对齐评审/数据分析/用户访谈
- 沟通：会议/即时沟通
- 学习：阅读/课程/总结输出
- 生活：通勤/运动/家务/休息

---

## 4.4 统计报表（日报 + 周报）
### 日报（必做）
- 类目占比（饼图）
- 活动 Top N（柱状图）
- 总记录时长、最大单段时长、碎片段数量（例如 <10min）
- 时间分布（按小时聚合，可选）

### 周报（MVP 做基础版）
- 本周总时长、每日总时长趋势（折线）
- 类目对比（堆叠柱状）
- 本周 Top 活动
- 碎片化指标趋势

> 暂不做目标，但周报要能让用户“看出问题”，比如碎片段过多、会议占比过高。

---

## 4.5 数据导出/备份（与你的 SQLite 诉求一致）
- 导出 CSV（时间记录表）
- 导出 JSON（包含类目/活动/记录，便于迁移）
- （可选）导入 JSON（非 MVP 也可做，投入不大但要校验）

---

## 5. 非功能需求（技术与体验）

### 5.1 性能
- 首页加载 < 1s（本地 SQLite + 简单 API 基本没问题）
- 报表查询：按日/周聚合 SQL，避免全表拉取后前端算

### 5.2 可靠性
- 服务重启：running 段依然存在（end_time NULL），前端重新拉取即可继续展示时长（时长=now-start_time）
- 时间以 ISO8601 存储，统一 UTC 或本地时区策略（建议：**数据库存 UTC，展示按用户时区**）

### 5.3 权限/安全（MVP）
- 单用户本地部署：可先不做登录
- 若考虑多用户：加简单账号（后续迭代），MVP 先保持架构可扩展（user_id 字段预留）

---

# 6. 后端设计（Go + SQLite）

## 6.1 模块拆分
- REST API（Gin/Fiber/chi 均可）
- 数据层（repository）
- 业务层（service：计时切换、重叠校验）
- SQLite 迁移（golang-migrate 或自研简单 migrations）

## 6.2 数据表设计（MVP）
> 先按可扩展设计，预留 user_id，但单用户可固定为 1。

### categories
- id INTEGER PK
- user_id INTEGER NOT NULL DEFAULT 1
- name TEXT NOT NULL
- color TEXT NULL
- sort_order INTEGER NOT NULL DEFAULT 0
- created_at DATETIME
- updated_at DATETIME

### activities
- id INTEGER PK
- user_id INTEGER NOT NULL DEFAULT 1
- category_id INTEGER NOT NULL FK(categories.id)
- name TEXT NOT NULL
- pinned INTEGER NOT NULL DEFAULT 0  (0/1)
- color TEXT NULL
- sort_order INTEGER NOT NULL DEFAULT 0
- created_at DATETIME
- updated_at DATETIME

### time_entries
- id INTEGER PK
- user_id INTEGER NOT NULL DEFAULT 1
- category_id INTEGER NOT NULL
- activity_id INTEGER NOT NULL
- start_time DATETIME NOT NULL
- end_time DATETIME NULL  (NULL = running)
- note TEXT NULL
- created_at DATETIME
- updated_at DATETIME

#### 索引建议
- idx_time_entries_user_start (user_id, start_time)
- idx_time_entries_user_end (user_id, end_time)
- idx_time_entries_activity (user_id, activity_id)

---

## 6.3 核心业务规则（后端强校验）
1) **Start/Switch**：保证只有一条 running
2) **Stop**：结束 running（若无 running 返回 204 或提示）
3) **Update entry**：校验不重叠（同 user 下 time_entries 时间区间不能交叉）
4) **Split/Merge**：服务端完成，避免前端拼装错误

---

## 6.4 API 设计（REST）
### 计时（核心）
- `GET /api/v1/timer/current`
  - 返回当前 running entry（或 null）

- `POST /api/v1/timer/start`
  - body: `{ activity_id, note? }`
  - 行为：若存在 running，先 stop 再 start（原子事务）
  - 返回：新 entry

- `POST /api/v1/timer/stop`
  - body: `{ end_time? }`（默认 now）
  - 返回：被结束的 entry

### Time Entries
- `GET /api/v1/time-entries?from=...&to=...`
- `POST /api/v1/time-entries`（手动创建一段）
- `PATCH /api/v1/time-entries/:id`
- `DELETE /api/v1/time-entries/:id`
- `POST /api/v1/time-entries/:id/split`
  - body: `{ split_time }`
- `POST /api/v1/time-entries/merge`
  - body: `{ ids: [id1,id2,...] }`（要求相邻且同 activity/category 且无间隔或间隔在阈值内）

### Categories & Activities
- `GET/POST/PATCH/DELETE /api/v1/categories`
- `GET/POST/PATCH/DELETE /api/v1/activities`
- `POST /api/v1/activities/:id/pin`（或 PATCH pinned 字段）

### Reports
- `GET /api/v1/reports/daily?date=YYYY-MM-DD`
- `GET /api/v1/reports/weekly?from=YYYY-MM-DD`（周一）

### Export
- `GET /api/v1/export/csv?from=&to=`
- `GET /api/v1/export/json`

---

# 7. 前端设计（Vue + Element UI）

## 7.1 页面关键组件
### Dashboard（计时台）
- 当前计时卡片：ElCard + ElTag（类目）+ ElButton（Stop）
- 常用活动区：ElButtonGroup / 可搜索 ElSelect
- 最近记录列表：用于快速切换

### Timeline（日视图）
- 列表：ElTable（可行内编辑）或自定义 time-block
- 编辑弹窗：ElDialog + 时间选择器
- 拆分：在详情弹窗里输入 split_time
- 合并：勾选多行 → “合并”按钮

### Reports
- 图表：ECharts（Element UI 本身不含图表）
- 维度：类目饼图、活动柱状、每日趋势折线

## 7.2 体验细节（非常影响留存）
- 全局快捷：`Ctrl+K` 打开活动搜索，回车立即切换开始（Web 端很实用）
- 切换成功 toast + 但不要打断操作
- 对 running 高亮显示（列表/时间轴中）

---

# 8. 关键技术点与实现建议（Go + SQLite）

1) **running 唯一性**
- 用事务：查询 running（end_time IS NULL），如有则更新 end_time，再插入新记录
- 可以额外加约束：SQLite 不方便做“部分唯一索引”跨行保证，但业务层事务足够

2) **时间重叠校验（Update/Create）**
- 检查条件：`NOT (new_end <= start_time OR new_start >= end_time)`
- 对 end_time NULL 的 running：视为 now 或禁止在编辑时产生冲突

3) **报表聚合 SQL**
- 日报：按 category_id/activity_id 聚合 sum(duration)
- duration 在 SQLite 可用 `(julianday(end_time)-julianday(start_time))*86400` 计算秒（注意 end_time NULL 的排除/按 now 处理）

---
