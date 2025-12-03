import RightContent from '@/layouts/RightContent';
import React from 'react';
import styles from './index.less';

interface HeaderRenderProps {
  logo?: React.ReactNode;
  title?: React.ReactNode;
  menuHeaderRender?: React.ReactNode;
}

const HeaderRender: React.FC<HeaderRenderProps> = ({ logo, title }) => {
  return (
    <div className={styles.headerRender}>
      {/* 左侧：logo 和标题 */}
      <div className={styles.left}>
        <img height={40} src={logo as string} alt="logo" />
        {title}
      </div>

      {/* 右侧：用户信息等 */}
      <div className={styles.right}>
        <RightContent />
      </div>
    </div>
  );
};

export default HeaderRender;
