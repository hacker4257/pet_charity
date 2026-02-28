import { useEffect, useState } from 'react';
import { Card, List, Avatar, Row, Col, Empty } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { useSearchParams } from 'react-router-dom';
import { getConversations } from '../../api/chat';
import useAuthStore from '../../store/useAuthStore';
import ChatWindow from '../../components/ChatWindow';

const Chat = () => {
    const { user } = useAuthStore();
    const [searchParams, setSearchParams] = useSearchParams();
    const [conversations, setConversations] = useState<any[]>([]);
    const [activeUserId, setActiveUserId] = useState<number | null>(
        searchParams.get('userId') ? Number(searchParams.get('userId')) : null
    );
    const [activeName, setActiveName] = useState(searchParams.get('name') || '');

    useEffect(() => {
        getConversations().then((res: any) => {
            setConversations(res.data || []);
        }).catch(() => {});
    }, []);

    const selectConversation = (msg: any) => {
        const otherId = msg.from_user_id === user?.id ? msg.to_user_id :msg.from_user_id;
        const otherName = msg.from_user?.nickname || msg.from_user?.username ||'用户';
        setActiveUserId(otherId);
        setActiveName(otherName);
        setSearchParams({ userId: String(otherId), name: otherName }, { replace:true });
    };

    return (
        <Row gutter={16} style={{ height: 'calc(100vh - 200px)' }}>
            {/* 左侧会话列表 */}
            <Col xs={24} md={8}>
                <Card title="消息" bodyStyle={{ padding: 0 }} style={{ height:'100%' }}>
                    <List
                        dataSource={conversations}
                        locale={{ emptyText: <Empty description="暂无会话"image={Empty.PRESENTED_IMAGE_SIMPLE} /> }}
                        renderItem={(msg: any) => {
                            const otherId = msg.from_user_id === user?.id ?msg.to_user_id : msg.from_user_id;
                            const isActive = activeUserId === otherId;
                            return (
                                <List.Item
                                    onClick={() => selectConversation(msg)}
                                    style={{
                                        cursor: 'pointer', padding: '12px 16px',
                                        background: isActive ? '#e6f7ff' :'transparent',
                                    }}
                                >
                                    <List.Item.Meta
                                        avatar={<Avatar icon={<UserOutlined />} />}
                                        title={msg.from_user?.nickname ||msg.from_user?.username}
                                        description={
                                            <div style={{
                                                overflow: 'hidden', textOverflow:'ellipsis',
                                                whiteSpace: 'nowrap', maxWidth: 200,
                                            }}>
                                                {msg.content}
                                            </div>
                                        }
                                    />
                                </List.Item>
                            );
                        }}
                    />
                </Card>
            </Col>

            {/* 右侧聊天窗口 */}
            <Col xs={24} md={16}>
                {activeUserId ? (
                    <ChatWindow
                        key={activeUserId}
                        targetUserId={activeUserId}
                        targetName={activeName}
                    />
                ) : (
                    <Card style={{ height: '100%', display: 'flex', alignItems:'center', justifyContent: 'center' }}>
                        <Empty description="选择一个会话开始聊天" />
                    </Card>
                )}
            </Col>
        </Row>
    );
};

export default Chat;
