import React, { useState, useEffect } from 'react';
import { Table, Button, Space, Modal, Form, Input, message, Tag, Popconfirm, Card } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined, UploadOutlined } from '@ant-design/icons';
import { noticeApi, uploadApi } from '../utils/api';
import { getRole } from '../utils/auth';
import dayjs from 'dayjs';

const NoticeList = () => {
  const [list, setList] = useState([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [role, setRole] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const [currentNotice, setCurrentNotice] = useState(null);
  const [form] = Form.useForm();

  const isAdmin = role === 'admin';

  useEffect(() => {
    setRole(getRole());
    fetchNotices();
  }, [page, pageSize]);

  const fetchNotices = async () => {
    setLoading(true);
    try {
      const response = await noticeApi.getAll({ page, page_size: pageSize });
      if (response.code === 200) {
        setList(response.data.list || []);
        setTotal(response.data.total || 0);
      }
    } catch (error) {
      message.error('获取公告列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = () => {
    setCurrentNotice(null);
    form.resetFields();
    setIsModalOpen(true);
  };

  const handleEdit = (record) => {
    setCurrentNotice(record);
    form.setFieldsValue(record);
    setIsModalOpen(true);
  };

  const handleDelete = async (id) => {
    try {
      const response = await noticeApi.delete(id);
      if (response.code === 200) {
        message.success('删除成功');
        fetchNotices();
      }
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleView = (record) => {
    setCurrentNotice(record);
    setIsDetailModalOpen(true);
  };

  const handleSubmit = async (values) => {
    try {
      if (currentNotice) {
        const response = await noticeApi.update(currentNotice.id, values);
        if (response.code === 200) {
          message.success('更新成功');
        }
      } else {
        const response = await noticeApi.create(values);
        if (response.code === 200) {
          message.success('创建成功');
        }
      }
      setIsModalOpen(false);
      fetchNotices();
    } catch (error) {
      message.error('操作失败');
    }
  };

  const handleUploadImage = async (info) => {
    if (info.file.status === 'done') {
      message.success('图片上传成功');
    } else if (info.file.status === 'error') {
      message.error('图片上传失败');
    }
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
    },
    {
      title: '作者',
      dataIndex: 'authorName',
      key: 'authorName',
      width: 100,
    },
    {
      title: '阅读量',
      dataIndex: 'viewCount',
      key: 'viewCount',
      width: 80,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (status) => (
        <Tag color={status === 1 ? 'green' : 'red'}>
          {status === 1 ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 180,
      render: (time) => dayjs(time).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_, record) => (
        <Space size="small">
          <Button type="link" icon={<EyeOutlined />} onClick={() => handleView(record)}>
            查看
          </Button>
          {isAdmin && (
            <>
              <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
                编辑
              </Button>
              <Popconfirm
                title="确定要删除吗？"
                onConfirm={() => handleDelete(record.id)}
                okText="确定"
                cancelText="取消"
              >
                <Button type="link" danger icon={<DeleteOutlined />}>
                  删除
                </Button>
              </Popconfirm>
            </>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <h2>{isAdmin ? '公告信息管理' : '公告查询'}</h2>
        {isAdmin && (
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
            发布公告
          </Button>
        )}
      </div>

      <Table
        columns={columns}
        dataSource={list}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: false,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条记录`,
          onChange: (page) => setPage(page),
        }}
      />

      <Modal
        title={currentNotice ? '编辑公告' : '发布公告'}
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        width={700}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            label="标题"
            name="title"
            rules={[{ required: true, message: '请输入标题' }]}
          >
            <Input placeholder="请输入公告标题" />
          </Form.Item>
          <Form.Item
            label="内容"
            name="content"
            rules={[{ required: true, message: '请输入内容' }]}
          >
            <Input.TextArea
              placeholder="请输入公告内容"
              rows={8}
              showCount
              maxLength={5000}
            />
          </Form.Item>
          <Form.Item label="图片" name="image">
            <Input placeholder="请输入图片URL（或上传图片）" />
          </Form.Item>
          <Form.Item label="状态" name="status">
            <Input type="number" placeholder="1=启用, 0=禁用" defaultValue={1} />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                保存
              </Button>
              <Button onClick={() => setIsModalOpen(false)}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="公告详情"
        open={isDetailModalOpen}
        onCancel={() => setIsDetailModalOpen(false)}
        footer={null}
        width={700}
      >
        {currentNotice && (
          <div>
            <h3 style={{ marginBottom: 16 }}>{currentNotice.title}</h3>
            <div style={{ marginBottom: 16, color: '#666' }}>
              <span>作者: {currentNotice.authorName}</span>
              <span style={{ marginLeft: 24 }}>
                创建时间: {dayjs(currentNotice.createdAt).format('YYYY-MM-DD HH:mm:ss')}
              </span>
              <span style={{ marginLeft: 24 }}>阅读量: {currentNotice.viewCount}</span>
            </div>
            {currentNotice.image && (
              <div style={{ marginBottom: 16 }}>
                <img
                  src={currentNotice.image.startsWith('http') ? currentNotice.image : `http://localhost:8080${currentNotice.image}`}
                  alt="公告图片"
                  style={{ maxWidth: '100%', maxHeight: 300 }}
                />
              </div>
            )}
            <div style={{ whiteSpace: 'pre-wrap', lineHeight: 1.8 }}>
              {currentNotice.content}
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default NoticeList;
