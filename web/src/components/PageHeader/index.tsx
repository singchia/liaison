import { Breadcrumb, Space, Typography } from 'antd';
import React from 'react';

const { Title } = Typography;

interface PageHeaderProps {
  title?: string;
  breadcrumbs?: Array<{ title: string; href?: string }>;
}

const PageHeader: React.FC<PageHeaderProps> = ({ title, breadcrumbs }) => {
  return (
    <div
      style={{
        background: '#fff',
        padding: '16px 24px',
        borderBottom: '1px solid #f0f0f0',
        position: 'sticky',
        top: 0,
        zIndex: 10,
      }}
    >
      <Space direction="vertical" size={4} style={{ width: '100%' }}>
        {breadcrumbs && breadcrumbs.length > 0 && (
          <Breadcrumb items={breadcrumbs} />
        )}
        {title && (
          <Title level={4} style={{ margin: 0 }}>
            {title}
          </Title>
        )}
      </Space>
    </div>
  );
};

export default PageHeader;
