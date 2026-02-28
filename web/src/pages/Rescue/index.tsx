import { useEffect, useState } from 'react';
import { Card, Row, Col, Button, List, Tag, Select, Spin, Empty, message } from
'antd';
import { PlusOutlined, EnvironmentOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getRescues, getRescueMapData, getNearbyOrgs } from '../../api/rescue';
import useAuthStore from '../../store/useAuthStore';
import RescueMap from '../../components/RescueMap';
import RescueReport from '../../components/RescueReport';

interface Rescue {
    id: number;
    title: string;
    description: string;
    urgency: string;
    status: string;
    address: string;
    created_at: string;
}

interface MapPoint {
    id: number;
    title: string;
    longitude: number;
    latitude: number;
    urgency: string;
    status: string;
}

interface NearbyOrg {
    id: number;
    name: string;
    address: string;
    distance: number;
    longitude: number;
    latitude: number;
}

const urgencyTagColor: Record<string, string> = {
    critical: 'red',
    high: 'orange',
    medium: 'blue',
    low: 'green',
};

const Rescue = () => {
    const { t } = useTranslation();
    const { isLoggedIn } = useAuthStore();

    const [rescues, setRescues] = useState<Rescue[]>([]);
    const [mapPoints, setMapPoints] = useState<MapPoint[]>([]);
    const [nearbyOrgs, setNearbyOrgs] = useState<NearbyOrg[]>([]);
    const [status, setStatus] = useState('');
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(false);
    const [reportOpen, setReportOpen] = useState(false);

    // 拉取列表
    const fetchList = async () => {
        setLoading(true);
        try {
            const params: Record<string, any> = { page, page_size: 10 };
            if (status) params.status = status;
            const res: any = await getRescues(params);
            setRescues(res.data.list || []);
            setTotal(res.data.total || 0);
        } catch {
            //
        } finally {
            setLoading(false);
        }
    };

    // 拉取地图标记点
    const fetchMapData = async () => {
        try {
            const res: any = await getRescueMapData();
            setMapPoints(res.data || []);
        } catch {
            //
        }
    };

    // 查找附近救助站
    const findNearbyOrgs = () => {
        if (!navigator.geolocation) {
            message.warning('浏览器不支持定位');
            return;
        }
        navigator.geolocation.getCurrentPosition(
            async (pos) => {
                try {
                    const res: any = await getNearbyOrgs({
                        lng: pos.coords.longitude,
                        lat: pos.coords.latitude,
                        radius: 10,
                    });
                    setNearbyOrgs(res.data || []);
                    if ((res.data || []).length === 0) {
                        message.info('附近暂无救助站');
                    }
                } catch {
                    //
                }
            },
            () => {
                message.error('获取位置失败');
            }
        );
    };

    useEffect(() => {
        fetchList();
    }, [status, page]);

    useEffect(() => {
        fetchMapData();
    }, []);

    return (
        <div>
            {/* 顶部操作栏 */}
            <Card style={{ marginBottom: 16 }}>
                <Row justify="space-between" align="middle">
                    <Col>
                        <Button
                            type="primary"
                            icon={<PlusOutlined />}
                            onClick={() => {
                                if (!isLoggedIn) {
                                    message.warning('请先登录');
                                    return;
                                }
                                setReportOpen(true);
                            }}
                        >
                            {t('rescue.report')}
                        </Button>
                        <Button
                            icon={<EnvironmentOutlined />}
                            style={{ marginLeft: 12 }}
                            onClick={findNearbyOrgs}
                        >
                            {t('rescue.nearbyOrgs')}
                        </Button>
                    </Col>
                    <Col>
                        <span>状态：</span>
                        <Select
                            value={status}
                            onChange={(v) => { setStatus(v); setPage(1); }}
                            style={{ width: 120 }}
                        >
                            <Select.Option value="">全部</Select.Option>
                            <Select.Option value="pending">待救助</Select.Option>
                            <Select.Option value="claimed">已认领</Select.Option>
                            <Select.Option value="completed">已完成</Select.Option>
                        </Select>
                    </Col>
                </Row>
            </Card>

            {/* 主体：左列表 + 右地图 */}
            <Row gutter={16}>
                {/* 左侧：救助列表 */}
                <Col xs={24} md={10}>
                    <Spin spinning={loading}>
                        <List
                            dataSource={rescues}
                            locale={{ emptyText: <Empty description={t('common.noData')} /> }}
                            pagination={{
                                current: page,
                                pageSize: 10,
                                total,
                                onChange: (p) => setPage(p),
                                size: 'small',
                            }}
                            renderItem={(item) => (
                                <Card
                                    size="small"
                                    style={{ marginBottom: 12 }}
                                    hoverable
                                >
                                    <div style={{ display: 'flex', justifyContent:'space-between' }}>
                                        <strong>{item.title}</strong>
                                        <Tag color={urgencyTagColor[item.urgency]}>
                                            {t(`rescue.${item.urgency}`)}
                                        </Tag>
                                    </div>
                                    <p style={{
                                        color: '#666', margin: '8px 0 4px',
                                        overflow: 'hidden',
                                        textOverflow: 'ellipsis',
                                        whiteSpace: 'nowrap',
                                    }}>
                                        {item.description}
                                    </p>
                                    <div style={{ color: '#999', fontSize: 12 }}>
                                        <EnvironmentOutlined /> {item.address ||'未知地点'}
                                        <span style={{ marginLeft: 12 }}>
                                            {new Date(item.created_at).toLocaleDateString()}
                                        </span>
                                    </div>
                                </Card>
                            )}
                        />
                    </Spin>

                    {/* 附近救助站列表 */}
                    {nearbyOrgs.length > 0 && (
                        <Card
                            title={t('rescue.nearbyOrgs')}
                            size="small"
                            style={{ marginTop: 16 }}
                        >
                            {nearbyOrgs.map((org) => (
                                <div
                                    key={org.id}
                                    style={{
                                        display: 'flex',
                                        justifyContent: 'space-between',
                                        alignItems: 'center',
                                        padding: '8px 0',
                                        borderBottom: '1px solid #f0f0f0',
                                    }}
                                >
                                    <div>
                                        <div><strong>{org.name}</strong></div>
                                        <div style={{ color: '#999', fontSize: 12}}>
                                            {org.address}
                                        </div>
                                    </div>
                                    <div style={{ textAlign: 'right' }}>
                                        <div style={{ color: '#1890ff' }}>
                                            {org.distance.toFixed(1)} km
                                        </div>
                                        <a
                                            href={`https://uri.amap.com/navigation?to=${org.longitude},${org.latitude},${org.name}`}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            style={{ fontSize: 12 }}
                                        >
                                            {t('rescue.navigate')}
                                        </a>
                                    </div>
                                </div>
                            ))}
                        </Card>
                    )}
                </Col>

                {/* 右侧：地图 */}
                <Col xs={24} md={14}>
                    <Card bodyStyle={{ padding: 0, overflow: 'hidden', borderRadius:8 }}>
                        <RescueMap points={mapPoints} />
                    </Card>
                </Col>
            </Row>

            {/* 上报弹窗 */}
            <RescueReport
                open={reportOpen}
                onClose={() => setReportOpen(false)}
                onSuccess={() => {
                    setReportOpen(false);
                    fetchList();
                    fetchMapData();
                }}
            />
        </div>
    );
};

export default Rescue;
