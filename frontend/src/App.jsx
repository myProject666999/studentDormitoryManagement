import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import Login from './pages/Login';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import Profile from './pages/Profile';
import NoticeList from './pages/NoticeList';
import { isAuthenticated } from './utils/auth';

const PrivateRoute = ({ children }) => {
  return isAuthenticated() ? children : <Navigate to="/login" />;
};

function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route 
            path="/" 
            element={
              <PrivateRoute>
                <Layout />
              </PrivateRoute>
            }
          >
            <Route index element={<Navigate to="/dashboard" replace />} />
            <Route path="dashboard" element={<Dashboard />} />
            <Route path="profile" element={<Profile />} />
            <Route path="notices" element={<NoticeList />} />
            <Route path="dormitories" element={<NoticeList />} />
            <Route path="assignments" element={<NoticeList />} />
            <Route path="my-assignment" element={<NoticeList />} />
            <Route path="students" element={<NoticeList />} />
            <Route path="workers" element={<NoticeList />} />
            <Route path="repair-requests" element={<NoticeList />} />
            <Route path="my-repair-requests" element={<NoticeList />} />
            <Route path="repair-records" element={<NoticeList />} />
            <Route path="my-repair-records" element={<NoticeList />} />
            <Route path="admins" element={<NoticeList />} />
          </Route>
        </Routes>
      </Router>
    </ConfigProvider>
  );
}

export default App;
