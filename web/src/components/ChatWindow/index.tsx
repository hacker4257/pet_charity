import { useEffect, useState, useRef, useCallback } from 'react';
import { Input, Button, Avatar, Spin, Empty } from 'antd';
import { SendOutlined, UserOutlined } from '@ant-design/icons';
import { useWebSocket } from '../../hooks/useWebsocket';
import { getPrivateHistory, getRoomHistory } from '../../api/chat';
import useAuthStore from '../../store/useAuthStore';

interface Message {
    from_user_id: number;
    to_user_id?: number;
    room_id?: string;
    content: string;
    nickname?: string;
    created_at?: string;
    from_user?: { nickname: string; username: string; avatar: string };
}

interface ChatWindowProps {
    // 私聊传 targetUserId，群聊传 roomId
    targetUserId?: number;
    targetName?: string;
    roomId?: string;
    roomTitle?: string;
}

const ChatWindow = ({ targetUserId, targetName, roomId, roomTitle }:
ChatWindowProps) => {
    const { user } = useAuthStore();
    const [messages, setMessages] = useState<Message[]>([]);
    const [inputVal, setInputVal] = useState('');
    const [loading, setLoading] = useState(false);
    const listRef = useRef<HTMLDivElement>(null);

    // 收到 WebSocket 消息
    const handleWsMessage = useCallback((msg: any) => {
        if (msg.type !== 'chat') return;

        // 判断是否属于当前窗口
        const isPrivate = targetUserId && !roomId;
        const isRoom = !!roomId;

        if (isPrivate) {
            const related = (msg.from_user_id === targetUserId && msg.to_user_id ===user?.id)
                || (msg.from_user_id === user?.id && msg.to_user_id ===targetUserId);
            if (!related) return;
        }
        if (isRoom && msg.room_id !== roomId) return;

        setMessages((prev) => [...prev, {
            from_user_id: msg.from_user_id,
            content: msg.content,
            nickname: msg.nickname,
            created_at: msg.created_at,
        }]);
    }, [targetUserId, roomId, user?.id]);

    const { send, connected } = useWebSocket(handleWsMessage);

    // 加载历史消息
    useEffect(() => {
        const loadHistory = async () => {
            setLoading(true);
            try {
                let res: any;
                if (roomId) {
                    res = await getRoomHistory(roomId, { page: 1, page_size: 50 });
                } else if (targetUserId) {
                    res = await getPrivateHistory(targetUserId, { page: 1,page_size: 50 });
                }
                const list = (res?.data?.list || []).reverse(); // 倒序存的，翻回来
                setMessages(list.map((m: any) => ({
                    from_user_id: m.from_user_id,
                    content: m.content,
                    nickname: m.from_user?.nickname || m.from_user?.username || '',
                    created_at: m.created_at,
                })));
            } catch {
                //
            } finally {
                setLoading(false);
            }
        };
        loadHistory();

        // 如果是房间，加入
        if (roomId) {
            send({ type: 'join_room', room_id: roomId, content: '' });
            return () => {
                send({ type: 'leave_room', room_id: roomId, content: '' });
            };
        }
    }, [targetUserId, roomId]);

    // 自动滚到底部
    useEffect(() => {
        if (listRef.current) {
            listRef.current.scrollTop = listRef.current.scrollHeight;
        }
    }, [messages]);

    // 发送消息
    const handleSend = () => {
        const text = inputVal.trim();
        if (!text) return;

        if (roomId) {
            send({ type: 'chat', room_id: roomId, content: text, msg_type: 'text'});
        } else if (targetUserId) {
            send({ type: 'chat', to_user_id: targetUserId, content: text, msg_type:'text'});
        }
        setInputVal('');
    };

    const title = roomId ? (roomTitle || `房间 ${roomId}`) : (targetName || '聊天');

    return (
        <div style={{
            display: 'flex', flexDirection: 'column',
            height: 500, border: '1px solid #f0f0f0', borderRadius: 8, overflow:'hidden',
        }}>
            {/* 标题栏 */}
            <div style={{
                padding: '12px 16px', borderBottom: '1px solid #f0f0f0',
                fontWeight: 'bold', display: 'flex', justifyContent:'space-between',
            }}>
                <span>{title}</span>
                <span style={{ fontSize: 12, color: connected ? '#52c41a' : '#999'}}>
                    {connected ? '已连接' : '连接中...'}
                </span>
            </div>

            {/* 消息列表 */}
            <div ref={listRef} style={{ flex: 1, overflow: 'auto', padding: 16 }}>
                {loading ? (
                    <div style={{ textAlign: 'center', padding: 40 }}><Spin /></div>
                ) : messages.length === 0 ? (
                    <Empty description="暂无消息" image={Empty.PRESENTED_IMAGE_SIMPLE} />
                ) : (
                    messages.map((msg, i) => {
                        const isMine = msg.from_user_id === user?.id;
                        return (
                            <div key={i} style={{
                                display: 'flex',
                                justifyContent: isMine ? 'flex-end' : 'flex-start',
                                marginBottom: 12,
                            }}>
                                {!isMine && <Avatar size="small" icon={<UserOutlined/>} style={{ marginRight: 8 }} />}
                                <div>
                                    {!isMine && (
                                        <div style={{ fontSize: 12, color: '#999',marginBottom: 2 }}>
                                            {msg.nickname}
                                        </div>
                                    )}
                                    <div style={{
                                        padding: '8px 12px',
                                        borderRadius: 8,
                                        maxWidth: 300,
                                        wordBreak: 'break-word',
                                        background: isMine ? '#1890ff' : '#f5f5f5',
                                        color: isMine ? '#fff' : '#333',
                                    }}>
                                        {msg.content}
                                    </div>
                                    <div style={{ fontSize: 11, color: '#bbb',marginTop: 2, textAlign: isMine ? 'right' : 'left' }}>
                                        {msg.created_at ? new Date(msg.created_at).toLocaleTimeString() : ''}
                                    </div>
                                </div>
                                {isMine && <Avatar size="small" icon={<UserOutlined/>} style={{ marginLeft: 8 }} />}
                            </div>
                        );
                    })
                )}
            </div>

            {/* 输入栏 */}
            <div style={{ display: 'flex', padding: 12, borderTop: '1px solid#f0f0f0', gap: 8 }}>
                <Input
                    value={inputVal}
                    onChange={(e) => setInputVal(e.target.value)}
                    onPressEnter={handleSend}
                    placeholder="输入消息..."
                    disabled={!connected}
                />
                <Button
                    type="primary"
                    icon={<SendOutlined />}
                    onClick={handleSend}
                    disabled={!connected || !inputVal.trim()}
                />
            </div>
        </div>
    );
};

export default ChatWindow;

