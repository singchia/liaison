import { PageContainer } from '@ant-design/pro-components';
import { Card, Row, Col, Spin } from 'antd';
import { Pie } from '@ant-design/plots';
import { useEffect, useState } from 'react';
import { getDeviceList, getApplicationList, getEdgeList } from '@/services/api';

interface PieData {
  type: string;
  value: number;
}

const DashboardPage: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [deviceData, setDeviceData] = useState<PieData[]>([]);
  const [applicationData, setApplicationData] = useState<PieData[]>([]);
  const [edgeData, setEdgeData] = useState<PieData[]>([]);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      // 获取所有数据
      const [devicesRes, applicationsRes, edgesRes] = await Promise.all([
        getDeviceList({ page_size: 1000 }),
        getApplicationList({ page_size: 1000 }),
        getEdgeList({ page_size: 1000 }),
      ]);

      // 处理设备数据 - 按操作系统分类
      const devices = devicesRes.data?.devices || [];
      const deviceStats: Record<string, number> = {};
      devices.forEach((device: API.Device) => {
        const os = device.os?.toLowerCase() || 'unknown';
        let osType = '其他';
        if (os.includes('linux')) {
          osType = 'Linux';
        } else if (os.includes('darwin') || os.includes('mac')) {
          osType = 'macOS';
        } else if (os.includes('windows')) {
          osType = 'Windows';
        }
        deviceStats[osType] = (deviceStats[osType] || 0) + 1;
      });
      const deviceDataList = Object.entries(deviceStats).map(([type, value]) => ({ type, value }));
      setDeviceData(deviceDataList);

      // 处理应用数据 - 按应用类型分类
      const applications = applicationsRes.data?.applications || [];
      const appStats: Record<string, number> = {};
      const typeMap: Record<string, string> = {
        web: 'Web',
        tcp: 'TCP',
        udp: 'UDP',
        ssh: 'SSH',
        rdp: 'RDP',
        database: '数据库',
        mysql: 'MySQL',
        postgresql: 'PostgreSQL',
        redis: 'Redis',
        mongodb: 'MongoDB',
      };
      applications.forEach((app: API.Application) => {
        const type = app.application_type || 'unknown';
        const displayType = typeMap[type.toLowerCase()] || type.toUpperCase();
        appStats[displayType] = (appStats[displayType] || 0) + 1;
      });
      const applicationDataList = Object.entries(appStats).map(([type, value]) => ({ type, value }));
      setApplicationData(applicationDataList);

      // 处理连接器数据 - 按在线状态分类
      const edges = edgesRes.data?.edges || [];
      const edgeStats: Record<string, number> = {
        在线: 0,
        离线: 0,
      };
      edges.forEach((edge: API.Edge) => {
        if (edge.online === 1) {
          edgeStats['在线']++;
        } else {
          edgeStats['离线']++;
        }
      });
      const edgeDataList = Object.entries(edgeStats).map(([type, value]) => ({ type, value }));
      setEdgeData(edgeDataList);

      // 调试信息
      console.log('设备数据:', deviceDataList);
      console.log('应用数据:', applicationDataList);
      console.log('连接器数据:', edgeDataList);
    } catch (error) {
      console.error('加载数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const getPieConfig = (data: PieData[]) => {
    return {
      data,
      angleField: 'value',
      colorField: 'type',
      // 标签显示在外部，显示类型名称
      label: {
        text: 'type',
        position: 'outside',
      },
      legend: false, // 隐藏默认图例，我们手动添加
      // @ant-design/plots 默认有 tooltip，显示类型和数值
      // 使用更丰富的颜色方案
      color: ['#1890ff', '#52c41a', '#faad14', '#f5222d', '#722ed1', '#13c2c2', '#eb2f96', '#fa8c16'],
      height: 250,
    };
  };

  return (
    <PageContainer>
      <Spin spinning={loading}>
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={24} md={8}>
            <Card title="设备统计" variant="outlined">
              <div style={{ display: 'flex', flexDirection: 'column' }}>
                {deviceData.length > 0 ? (
                  <>
                    <div style={{ flex: 1, minHeight: 250 }}>
                      <Pie {...getPieConfig(deviceData)} />
                    </div>
                    <div style={{ marginTop: 16, textAlign: 'center' }}>
                      {deviceData.map((item, index) => {
                        const colors = ['#1890ff', '#52c41a', '#faad14', '#f5222d', '#722ed1', '#13c2c2'];
                        const color = colors[index % colors.length];
                        return (
                          <span key={item.type} style={{ margin: '0 8px', fontSize: '12px' }}>
                            <span
                              style={{
                                display: 'inline-block',
                                width: 12,
                                height: 12,
                                backgroundColor: color,
                                marginRight: 4,
                                verticalAlign: 'middle',
                              }}
                            />
                            {item.type}: {item.value}
                          </span>
                        );
                      })}
                    </div>
                  </>
                ) : (
                  <div style={{ textAlign: 'center', padding: '40px 0' }}>
                    暂无数据
                  </div>
                )}
              </div>
            </Card>
          </Col>
          <Col xs={24} sm={24} md={8}>
            <Card title="应用统计" variant="outlined">
              <div style={{ display: 'flex', flexDirection: 'column' }}>
                {applicationData.length > 0 ? (
                  <>
                    <div style={{ flex: 1, minHeight: 250 }}>
                      <Pie {...getPieConfig(applicationData)} />
                    </div>
                    <div style={{ marginTop: 16, textAlign: 'center' }}>
                      {applicationData.map((item, index) => {
                        const colors = ['#1890ff', '#52c41a', '#faad14', '#f5222d', '#722ed1', '#13c2c2'];
                        const color = colors[index % colors.length];
                        return (
                          <span key={item.type} style={{ margin: '0 8px', fontSize: '12px' }}>
                            <span
                              style={{
                                display: 'inline-block',
                                width: 12,
                                height: 12,
                                backgroundColor: color,
                                marginRight: 4,
                                verticalAlign: 'middle',
                              }}
                            />
                            {item.type}: {item.value}
                          </span>
                        );
                      })}
                    </div>
                  </>
                ) : (
                  <div style={{ textAlign: 'center', padding: '40px 0' }}>
                    暂无数据
                  </div>
                )}
              </div>
            </Card>
          </Col>
          <Col xs={24} sm={24} md={8}>
            <Card title="连接器统计" variant="outlined">
              <div style={{ display: 'flex', flexDirection: 'column' }}>
                {edgeData.length > 0 ? (
                  <>
                    <div style={{ flex: 1, minHeight: 250 }}>
                      <Pie {...getPieConfig(edgeData)} />
                    </div>
                    <div style={{ marginTop: 16, textAlign: 'center' }}>
                      {edgeData.map((item, index) => {
                        const colors = ['#1890ff', '#52c41a', '#faad14', '#f5222d', '#722ed1', '#13c2c2'];
                        const color = colors[index % colors.length];
                        return (
                          <span key={item.type} style={{ margin: '0 8px', fontSize: '12px' }}>
                            <span
                              style={{
                                display: 'inline-block',
                                width: 12,
                                height: 12,
                                backgroundColor: color,
                                marginRight: 4,
                                verticalAlign: 'middle',
                              }}
                            />
                            {item.type}: {item.value}
                          </span>
                        );
                      })}
                    </div>
                  </>
                ) : (
                  <div style={{ textAlign: 'center', padding: '40px 0' }}>
                    暂无数据
                  </div>
                )}
              </div>
            </Card>
          </Col>
        </Row>
      </Spin>
    </PageContainer>
  );
};

export default DashboardPage;
