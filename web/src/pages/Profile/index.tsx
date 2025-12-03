import { getProfile } from '@/services/iam';
import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { useRequest } from 'ahooks';
import {
  Avatar,
  Button,
  Card,
  Col,
  Descriptions,
  Row,
  Space,
  Spin,
} from 'antd';
import React, { useState } from 'react';
import ChangePasswordModal from './components/ChangePasswordModal';
import EditProfileModal from './components/EditProfileModal';

const Profile: React.FC = () => {
  const [editModalVisible, setEditModalVisible] = useState(false);
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);

  // 获取用户信息
  const { data, loading, refresh } = useRequest(() => getProfile());

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <>
      <Row gutter={16}>
        <Col span={24}>
          <Card>
            <div style={{ textAlign: 'center', marginBottom: 24 }}>
              <Avatar
                size={120}
                src={data?.avatar}
                icon={<UserOutlined />}
                style={{ marginBottom: 16 }}
              />
              <h2>{data?.username}</h2>
              <p style={{ color: '#999' }}>{data?.email}</p>
              <Space style={{ marginTop: 16 }}>
                <Button
                  type="primary"
                  onClick={() => setEditModalVisible(true)}
                >
                  编辑资料
                </Button>
                <Button
                  icon={<LockOutlined />}
                  onClick={() => setPasswordModalVisible(true)}
                >
                  修改密码
                </Button>
              </Space>
            </div>

            <Descriptions column={2} bordered>
              <Descriptions.Item label="用户ID">{data?.id}</Descriptions.Item>
              <Descriptions.Item label="角色">
                {data?.role || '普通用户'}
              </Descriptions.Item>
              <Descriptions.Item label="用户名">
                {data?.username}
              </Descriptions.Item>
              <Descriptions.Item label="邮箱">{data?.email}</Descriptions.Item>
              <Descriptions.Item label="创建时间">
                {data?.created_at}
              </Descriptions.Item>
              <Descriptions.Item label="更新时间">
                {data?.updated_at}
              </Descriptions.Item>
            </Descriptions>
          </Card>
        </Col>
      </Row>

      {/* 编辑资料弹框 */}
      <EditProfileModal
        visible={editModalVisible}
        user={data}
        onCancel={() => setEditModalVisible(false)}
        onSuccess={() => {
          setEditModalVisible(false);
          refresh();
        }}
      />

      {/* 修改密码弹框 */}
      <ChangePasswordModal
        visible={passwordModalVisible}
        onCancel={() => setPasswordModalVisible(false)}
        onSuccess={() => setPasswordModalVisible(false)}
      />
    </>
  );
};

export default Profile;
