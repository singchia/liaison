import { PageContainer } from '@ant-design/pro-components';
import React from 'react';

interface BasicLayoutProps {
  children: React.ReactNode;
}

const BasicLayout: React.FC<BasicLayoutProps> = ({ children }) => {
  return <PageContainer fixedHeader>{children}</PageContainer>;
};

export default BasicLayout;
