import { useEffect, useState } from 'react';
import { List, Button, Tag, Card, Badge, Empty, message } from 'antd';
import { BellOutlined, CheckOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { getNotifications, markRead, markAllRead } from '../../api/notification';
import useAuthStore from '../../store/useAuthStore';

const typeColorMap: Record<string, string> = {
    adoption: 'blue',
    donation: 'green',
    rescue: 'orange',
    system: 'default',
};

const Notifications = () => {
    const { t } = useTranslation();
    const { isLoggedIn } = useAuthStore();
    const [list, setList] = useState<any[]>([]);
    const [total, setTotal] = useState(0);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(false);

    const fetchList = async () => {
        if (!isLoggedIn) return;
        setLoading(true);
        try {
            const res: any = await getNotifications({ page, page_size: 20 });
            setList(res.data.list || []);
            setTotal(res.data.total || 0);
        } catch {
            //
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchList();
    }, [page]);

    const handleMarkRead = async (id: number) => {
        try {
            await markRead(id);
            setList(prev => prev.map(item =>
                item.id === id ? { ...item, is_read: true } : item
            ));
        } catch {
            //
        }
    };

    const handleMarkAllRead = async () => {
        try {
            await markAllRead();
            setList(prev => prev.map(item => ({ ...item, is_read: true })));
            message.success('已全部标记为已读');
        } catch {
            //
        }
    };

    const unreadCount = list.filter(item => !item.is_read).length;

    return (
        <div style={{ maxWidth: 800, margin: '0 auto' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                <h2><BellOutlined /> {t('nav.notifications') || '通知中心'}</h2>
                {unreadCount > 0 && (
                    <Button icon={<CheckOutlined />} onClick={handleMarkAllRead}>
                        全部已读
                    </Button>
                )}
            </div>

            <Card>
                <List
                    loading={loading}
                    dataSource={list}
                    locale={{ emptyText: <Empty description="暂无通知" /> }}
                    pagination={{
                        current: page,
                        pageSize: 20,
                        total,
                        onChange: (p) => setPage(p),
                    }}
                    renderItem={(item: any) => (
                        <List.Item
                            style={{
                                background: item.is_read ? 'transparent' : '#f6ffed',
                                padding: '12px 16px',
                                cursor: item.is_read ? 'default' : 'pointer',
                            }}
                            onClick={() => !item.is_read && handleMarkRead(item.id)}
                            extra={
                                !item.is_read && (
                                    <Badge status="processing" text="未读" />
                                )
                            }
                        >
                            <List.Item.Meta
                                title={
                                    <span>
                                        <Tag color={typeColorMap[item.type] || 'default'}>
                                            {item.type}
                                        </Tag>
                                        {item.title}
                                    </span>
                                }
                                description={
                                    <div>
                                        <div>{item.content}</div>
                                        <div style={{ color: '#999', fontSize: 12, marginTop: 4 }}>
                                            {new Date(item.created_at).toLocaleString()}
                                        </div>
                                    </div>
                                }
                            />
                        </List.Item>
                    )}
                />
            </Card>
        </div>
    );
};

export default Notifications;
