import api from './axios';

const getRolePrefix = () => {
  const role = localStorage.getItem('role');
  if (role === 'admin') return '/admin';
  if (role === 'student') return '/student';
  if (role === 'worker') return '/worker';
  return '';
};

export const noticeApi = {
  getAll: (params) => api.get('/notices', { params }),
  getById: (id) => api.get(`/notices/${id}`),
  create: (data) => api.post(`${getRolePrefix()}/notices`, data),
  update: (id, data) => api.put(`${getRolePrefix()}/notices/${id}`, data),
  delete: (id) => api.delete(`${getRolePrefix()}/notices/${id}`),
};

export const adminApi = {
  getAll: (params) => api.get('/admin/admins', { params }),
  create: (data) => api.post('/admin/admins', data),
  update: (id, data) => api.put(`/admin/admins/${id}`, data),
  delete: (id) => api.delete(`/admin/admins/${id}`),
};

export const studentApi = {
  getAll: (params) => api.get('/admin/students', { params }),
  create: (data) => api.post('/admin/students', data),
  update: (id, data) => api.put(`/admin/students/${id}`, data),
  delete: (id) => api.delete(`/admin/students/${id}`),
};

export const workerApi = {
  getAll: (params) => api.get('/admin/workers', { params }),
  getAllActive: () => api.get(`${getRolePrefix()}/workers/all`),
  create: (data) => api.post('/admin/workers', data),
  update: (id, data) => api.put(`/admin/workers/${id}`, data),
  delete: (id) => api.delete(`/admin/workers/${id}`),
};

export const dormitoryApi = {
  getAll: (params) => api.get('/admin/dormitories', { params }),
  getAvailable: (params) => api.get('/admin/dormitories/available', { params }),
  getById: (id) => api.get(`/admin/dormitories/${id}`),
  create: (data) => api.post('/admin/dormitories', data),
  update: (id, data) => api.put(`/admin/dormitories/${id}`, data),
  delete: (id) => api.delete(`/admin/dormitories/${id}`),
};

export const assignmentApi = {
  getAll: (params) => api.get('/admin/assignments', { params }),
  getById: (id) => api.get(`/admin/assignments/${id}`),
  getMyAssignment: () => api.get('/student/my-assignment'),
  create: (data) => api.post('/admin/assignments', data),
  update: (id, data) => api.put(`/admin/assignments/${id}`, data),
  delete: (id) => api.delete(`/admin/assignments/${id}`),
};

export const repairRequestApi = {
  getAll: (params) => api.get('/admin/repair-requests', { params }),
  getMyRequests: (params) => api.get(`${getRolePrefix()}/my-repair-requests`, { params }),
  getById: (id) => api.get(`/admin/repair-requests/${id}`),
  create: (data) => api.post(`${getRolePrefix()}/repair-requests`, data),
  update: (id, data) => api.put(`${getRolePrefix()}/repair-requests/${id}`, data),
  delete: (id) => api.delete(`/admin/repair-requests/${id}`),
};

export const repairRecordApi = {
  getAll: (params) => api.get('/admin/repair-records', { params }),
  getMyRecords: (params) => api.get(`${getRolePrefix()}/my-repair-records`, { params }),
  getById: (id) => api.get(`/admin/repair-records/${id}`),
  create: (data) => api.post(`${getRolePrefix()}/repair-records`, data),
  update: (id, data) => api.put(`${getRolePrefix()}/repair-records/${id}`, data),
  delete: (id) => api.delete(`/admin/repair-records/${id}`),
};

export const uploadApi = {
  uploadImage: (file) => {
    const formData = new FormData();
    formData.append('file', file);
    return api.post('/upload/image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },
  uploadAttachment: (file) => {
    const formData = new FormData();
    formData.append('file', file);
    return api.post('/upload/attachment', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },
};
