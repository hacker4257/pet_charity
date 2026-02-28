import { useEffect, useState } from 'react';
import { Row, Col, Card, Statistic, Button, Typography } from 'antd';
import { HeartOutlined, SmileOutlined, AlertOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import request from '../../api/request';

const { Title, Paragraph } = Typography;

const Home = () => {
    const navigate = useNavigate();
    const { t } = useTranslation();
    const [stats, setStats] = useState({
        adoptable_pets: 0,
        adopted_pets: 0,
        rescue_count: 0,
    });

    useEffect(() => {
        request.get('/pets/stats').then((res: any) => {
            setStats(res.data);
        }).catch(() => {});
    }, []);

    return (
        <div>
            {/* Banner */}
            <div style={{
                textAlign: 'center',
                padding: '60px 20px',
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                borderRadius: 8,
                marginBottom: 40,
            }}>
                <Title level={1} style={{ color: '#fff', marginBottom: 16 }}>
                    感谢有您 - Pet Charity
                </Title>
                <Paragraph style={{ color: 'rgba(255,255,255,0.85)', fontSize: 18,marginBottom: 32 }}>
                    让每一个生命都被温柔以待
                </Paragraph>
                <Button type="primary" size="large" onClick={() =>navigate('/pets')}
                    style={{ marginRight: 16 }}
                >
                    领养宠物
                </Button>
                <Button size="large" ghost onClick={() => navigate('/rescue')}>
                    参与救助
                </Button>
            </div>

            {/* 统计数据 */}
            <Row gutter={24} style={{ marginBottom: 40 }}>
                <Col span={8}>
                    <Card>
                        <Statistic
                            title="待领养宠物"
                            value={stats.adoptable_pets}
                            prefix={<HeartOutlined />}
                            suffix="只"
                        />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card>
                        <Statistic
                            title="已成功领养"
                            value={stats.adopted_pets}
                            prefix={<SmileOutlined />}
                            suffix="只"
                        />
                    </Card>
                </Col>
                <Col span={8}>
                    <Card>
                        <Statistic
                            title="救助记录"
                            value={stats.rescue_count}
                            prefix={<AlertOutlined />}
                            suffix="条"
                        />
                    </Card>
                </Col>
            </Row>

            {/* 快捷入口 */}
            <Row gutter={24}>
                <Col span={8}>
                    <Card hoverable onClick={() => navigate('/pets')}>
                        <Title level={4}>🐾 宠物领养</Title>

                        <Paragraph>浏览所有待领养的宠物，给它们一个温暖的家</Paragraph>
                    </Card>
                </Col>
                <Col span={8}>
                    <Card hoverable onClick={() => navigate('/rescue')}>
                        <Title level={4}>🆘 流浪救助</Title>
                        <Paragraph>发现流浪动物？立即上报，一起参与救助</Paragraph>
                    </Card>
                </Col>
                <Col span={8}>
                    <Card hoverable onClick={() => navigate('/donation')}>
                        <Title level={4}>💝 爱心捐赠</Title>

                        <Paragraph>为流浪动物贡献一份力量，每一分都有意义</Paragraph>
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default Home;
