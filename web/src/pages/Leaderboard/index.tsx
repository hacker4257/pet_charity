import { useEffect, useState } from 'react';
import { Card, Table, Avatar, Tag, Timeline, Row, Col, Spin } from 'antd';
import { TrophyOutlined, CrownOutlined, FireOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getLeaderboard, getMyRank, getFeed } from '../../api/activity';
import useAuthStore from '../../store/useAuthStore';

const actionLabelMap: Record<string, string> = {
    adopt: '领养宠物',
    donate: '爱心捐赠',
    rescue_report: '发起救助',
    rescue_claim: '认领救助',
    follow_up: '跟进救助',
};

const Leaderboard = () => {
    const { t } = useTranslation();
    const { isLoggedIn } = useAuthStore();
    const [list, setList] = useState<any[]>([]);
    const [myRank, setMyRank] = useState<any>(null);
    const [feed, setFeed] = useState<any[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        setLoading(true);
        Promise.all([
            getLeaderboard({ page: 1, page_size: 50 }),
            getFeed({ page: 1, page_size: 20 }),
            isLoggedIn ? getMyRank().catch(() => null) : Promise.resolve(null),
        ])
            .then(([lbRes, feedRes, rankRes]: any[]) => {
                setList(lbRes?.data || []);
                setFeed(feedRes?.data || []);
                if (rankRes?.data) setMyRank(rankRes.data);
            })
            .catch(() => {})
            .finally(() => setLoading(false));
    }, [isLoggedIn]);

    const columns = [
        {
            title: '排名',
            dataIndex: 'rank',
            key: 'rank',
            width: 80,
            render: (rank: number) => {
                if (rank === 1) return <CrownOutlined style={{ color: '#faad14', fontSize: 20 }} />;
                if (rank === 2) return <CrownOutlined style={{ color: '#bfbfbf', fontSize: 18 }} />;
                if (rank === 3) return <CrownOutlined style={{ color: '#d48806', fontSize: 16 }} />;
                return rank;
            },
        },
        {
            title: '用户',
            key: 'user',
            render: (_: any, record: any) => (
                <span>
                    <Avatar src={record.avatar} size="small" style={{ marginRight: 8 }}>
                        {record.nickname?.[0]}
                    </Avatar>
                    {record.nickname || `用户${record.user_id}`}
                </span>
            ),
        },
        {
            title: '积分',
            dataIndex: 'score',
            key: 'score',
            render: (score: number) => (
                <Tag color="orange" icon={<FireOutlined />}>
                    {score}
                </Tag>
            ),
        },
    ];

    if (loading) {
        return <div style={{ textAlign: 'center', padding: 60 }}><Spin size="large" /></div>;
    }

    return (
        <div style={{ maxWidth: 1200, margin: '0 auto' }}>
            <h2><TrophyOutlined /> {t('nav.leaderboard') || '爱心排行榜'}</h2>

            {myRank && (
                <Card style={{ marginBottom: 16 }}>
                    <span>我的排名：<strong>第 {myRank.rank} 名</strong></span>
                    <Tag color="orange" style={{ marginLeft: 12 }} icon={<FireOutlined />}>
                        {myRank.score} 积分
                    </Tag>
                </Card>
            )}

            <Row gutter={24}>
                <Col xs={24} lg={14}>
                    <Card title="排行榜">
                        <Table
                            columns={columns}
                            dataSource={list}
                            rowKey="user_id"
                            pagination={false}
                            size="middle"
                        />
                    </Card>
                </Col>
                <Col xs={24} lg={10}>
                    <Card title="最新动态" style={{ marginTop: 0 }}>
                        {feed.length === 0 ? (
                            <div style={{ color: '#999', textAlign: 'center' }}>暂无动态</div>
                        ) : (
                            <Timeline
                                items={feed.map((item: any, idx: number) => ({
                                    key: idx,
                                    children: (
                                        <div>
                                            <span style={{ fontWeight: 500 }}>
                                                用户{item.user_id}
                                            </span>{' '}
                                            {actionLabelMap[item.action] || item.action}
                                            <div style={{ color: '#999', fontSize: 12 }}>
                                                {new Date(item.timestamp * 1000).toLocaleString()}
                                            </div>
                                        </div>
                                    ),
                                }))}
                            />
                        )}
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default Leaderboard;
