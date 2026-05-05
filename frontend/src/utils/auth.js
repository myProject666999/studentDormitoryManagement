import api from './axios';

export const login = async (username, password) => {
  const response = await api.post('/login', { username, password });
  if (response.code === 200) {
    const { token, user, role } = response.data;
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    localStorage.setItem('role', role);
    return response.data;
  }
  throw new Error(response.message);
};

export const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  localStorage.removeItem('role');
};

export const getCurrentUser = () => {
  const userStr = localStorage.getItem('user');
  return userStr ? JSON.parse(userStr) : null;
};

export const getRole = () => {
  return localStorage.getItem('role');
};

export const isAuthenticated = () => {
  return !!localStorage.getItem('token');
};

export const changePassword = async (oldPassword, newPassword) => {
  return await api.put('/user/password', { old_password: oldPassword, new_password: newPassword });
};

export const updateProfile = async (data) => {
  return await api.put('/user/profile', data);
};

export const getCurrentUserInfo = async () => {
  return await api.get('/user');
};
