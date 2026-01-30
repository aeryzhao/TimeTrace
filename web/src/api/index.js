import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1'
})

export const getCategories = () => api.get('/categories')
export const createCategory = (data) => api.post('/categories', data)
export const deleteCategory = (id) => api.delete(`/categories/${id}`)

export const getActivities = () => api.get('/activities')
export const createActivity = (data) => api.post('/activities', data)
export const deleteActivity = (id) => api.delete(`/activities/${id}`)
export const pinActivity = (id) => api.post(`/activities/${id}/pin`)

export const getCurrentTimer = () => api.get('/timer/current')
export const startTimer = (activityId, note) => api.post('/timer/start', { activity_id: activityId, note })
export const stopTimer = (endTime) => api.post('/timer/stop', { end_time: endTime })

export const getTimeEntries = (from, to) => api.get('/time-entries', { params: { from, to } })
export const createTimeEntry = (data) => api.post('/time-entries', data)
export const updateTimeEntry = (id, data) => api.patch(`/time-entries/${id}`, data)
export const deleteTimeEntry = (id) => api.delete(`/time-entries/${id}`)

export const getDailyReport = (date) => api.get('/reports/daily', { params: { date } })

export default api
