import React, { useState, useEffect } from 'react';
import { Card, Form, Input, Button, Tabs, message, Descriptions, Divider } from 'antd';
import { getCurrentUser, changePassword, updateProfile, getRole } from '../utils/auth';

const Profile = () => {
  const [user, setUser] = useState(null);
  const [role, setRole] = useState(null);
  const [loading, setLoading] = useState(false);
  const [form] = Form.useForm();
  const [passwordForm] = Form.useForm();

  useEffect(() => {
    const currentUser = getCurrentUser();
    setUser(currentUser);
    setRole(getRole());
    if (currentUser) {
      form.setFieldsValue(currentUser);
    }
  }, [form]);

  const handleUpdateProfile = async (values) => {
    setLoading(true);
    try {
      const response = await updateProfile(values);
      if (response.code === 200) {
        message.success('个人信息更新成功');
        localStorage.setItem('user', JSON.stringify(response.data));
        setUser(response.data);
      }
    } catch (error) {
      message.error(error.message || '更新失败');
    } finally {
      setLoading(false);
    }
  };

  const handleChangePassword = async (values) => {
    if (values.newPassword !== values.confirmPassword) {
      message.error('两次输入的密码不一致');
      return;
    }
    
    setLoading(true);
    try {
      const response = await changePassword(values.oldPassword, values.newPassword);
      if (response.code === 200) {
        message.success('密码修改成功');
        passwordForm.resetFields();
      }
    } catch (error) {
      message.error(error.message || '密码修改失败');
    } finally {
      setLoading(false);
    }
  };

  const tabItems = [
    {
      key: 'info',
      label: '个人信息',
      children: (
        <div>
          <Descriptions title="当前信息" bordered column={2}>
            <Descriptions.Item label="用户名">{user?.username}</Descriptions.Item>
            <Descriptions.Item label="角色">
              {role === 'admin' ? '管理员' : role === 'student' ? '学生' : '维修工'}
            </Descriptions.Item>
            <Descriptions.Item label="姓名">{user?.name}</Descriptions.Item>
            {role === 'student' && (
              <>
                <Descriptions.Item label="学号">{user?.studentNumber}</Descriptions.Item>
                <Descriptions.Item label="班级">{user?.class}</Descriptions.Item>
                <Descriptions.Item label="专业">{user?.major}</Descriptions.Item>
              </>
            )}
            {role === 'worker' && (
              <Descriptions.Item label="工号">{user?.workerNumber}</Descriptions.Item>
            )}
            <Descriptions.Item label="电话">{user?.phone || '-'}</Descriptions.Item>
            <Descriptions.Item label="邮箱">{user?.email || '-'}</Descriptions.Item>
          </Descriptions>

          <Divider />

          <Card title="修改信息" size="small">
            <Form
              form={form}
              layout="vertical"
              onFinish={handleUpdateProfile}
              initialValues={user || {}}
            >
              <Form.Item label="姓名" name="name">
                <Input placeholder="请输入姓名" />
              </Form.Item>
              <Form.Item label="电话" name="phone">
                <Input placeholder="请输入电话" />
              </Form.Item>
              <Form.Item label="邮箱" name="email">
                <Input placeholder="请输入邮箱" />
              </Form.Item>
              {role === 'student' && (
                <Form.Item label="性别" name="gender">
                  <Input placeholder="请输入性别" />
                </Form.Item>
              )}
              {role === 'worker' && (
                <>
                  <Form.Item label="性别" name="gender">
                    <Input placeholder="请输入性别" />
                  </Form.Item>
                  <Form.Item label="专长" name="specialty">
                    <Input placeholder="请输入专长" />
                  </Form.Item>
                </>
              )}
              <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading}>
                  保存修改
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </div>
      ),
    },
    {
      key: 'password',
      label: '修改密码',
      children: (
        <Card title="密码修改" size="small" style={{ maxWidth: 500 }}>
          <Form
            form={passwordForm}
            layout="vertical"
            onFinish={handleChangePassword}
          >
            <Form.Item
              label="原密码"
              name="oldPassword"
              rules={[{ required: true, message: '请输入原密码' }]}
            >
              <Input.Password placeholder="请输入原密码" />
            </Form.Item>
            <Form.Item
              label="新密码"
              name="newPassword"
              rules={[
                { required: true, message: '请输入新密码' },
                { min: 6, message: '密码至少6位' },
              ]}
            >
              <Input.Password placeholder="请输入新密码" />
            </Form.Item>
            <Form.Item
              label="确认新密码"
              name="confirmPassword"
              rules={[
                { required: true, message: '请确认新密码' },
              ]}
            >
              <Input.Password placeholder="请再次输入新密码" />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loading}>
                修改密码
              </Button>
            </Form.Item>
          </Form>
        </Card>
      ),
    },
  ];

  return (
    <div>
      <h2 style={{ marginBottom: 24 }}>个人中心</h2>
      <Tabs items={tabItems} />
    </div>
  );
};

export default Profile;
