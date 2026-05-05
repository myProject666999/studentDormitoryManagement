import React from 'react';
import { Row, Col, Card, Statistic } from 'antd';
import { 
  UserOutlined, 
  HomeOutlined, 
  ToolOutlined, 
  NotificationOutlined 
} from '@ant-design/icons';

const Dashboard = () => {
  return (
    <div>
      <h2 style={{ marginBottom: 24 }}>仪表盘</h2>
      
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="学生总数"
              value={0}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="寝室总数"
              value={0}
              prefix={<HomeOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="维修工总数"
              value={0}
              prefix={<ToolOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="公告数量"
              value={0}
              prefix={<NotificationOutlined />}
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} lg={12}>
          <Card title="系统说明">
            <p>欢迎使用学生寝室管理系统！</p>
            <p>本系统提供以下功能：</p>
            <ul>
              <li>公告信息管理：发布、编辑、删除公告</li>
              <li>寝室管理：寝室信息的增删改查</li>
              <li>学生管理：学生信息管理</li>
              <li>寝室安排：学生寝室分配管理</li>
              <li>报修管理：学生报修、维修工维修</li>
              <li>用户管理：管理员、维修工信息管理</li>
            </ul>
          </Card>
        </Col>
        
        <Col xs={24} lg={12}>
          <Card title="快速链接">
            <p>您可以通过左侧菜单快速访问各个功能模块：</p>
            <ul>
              <li><b>仪表盘</b>：查看系统概览</li>
              <li><b>公告信息管理</b>：发布和查看公告</li>
              <li><b>寝室管理</b>：管理寝室信息</li>
              <li><b>寝室安排</b>：分配学生寝室</li>
              <li><b>学生信息管理</b>：管理学生信息</li>
              <li><b>维修工管理</b>：管理维修工信息</li>
              <li><b>寝室报修管理</b>：处理报修申请</li>
              <li><b>维修情况管理</b>：查看维修记录</li>
              <li><b>管理员管理</b>：管理系统管理员</li>
              <li><b>个人中心</b>：修改个人信息和密码</li>
            </ul>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
