import { useEffect, useRef, useCallback, useState } from 'react';

interface WsMessage {
    type: string;
    from_user_id?: number;
    to_user_id?: number;
    room_id?: string;
    content: string;
    msg_type: string;
    nickname?: string;
    created_at?: string;
}

export function useWebSocket(onMessage: (msg: WsMessage) => void) {
    const wsRef = useRef<WebSocket | null>(null);
    const [connected, setConnected] = useState(false);
    const reconnectTimer = useRef<ReturnType<typeof setTimeout>>();

    const connect = useCallback(() => {
        const token = localStorage.getItem('access_token');
        if (!token) return;

        // WebSocket 不能自定义 Header，通过 query 传 token
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const ws = new
WebSocket(`${protocol}//${window.location.host}/api/v1/chat/ws?token=${token}`);

        ws.onopen = () => {
            setConnected(true);
            // 心跳保活
            const heartbeat = setInterval(() => {
                if (ws.readyState === WebSocket.OPEN) {
                    ws.send(JSON.stringify({ type: 'ping' }));
                }
            }, 25000);
            ws.addEventListener('close', () => clearInterval(heartbeat));
        };

        ws.onmessage = (e) => {
            try {
                const msg: WsMessage = JSON.parse(e.data);
                onMessage(msg);
            } catch {
                //
            }
        };

        ws.onclose = () => {
            setConnected(false);
            // 3 秒后自动重连
            reconnectTimer.current = setTimeout(connect, 3000);
        };

        wsRef.current = ws;
    }, [onMessage]);

    useEffect(() => {
        connect();
        return () => {
            if (reconnectTimer.current) clearTimeout(reconnectTimer.current);
            wsRef.current?.close();
        };
    }, [connect]);

    const send = useCallback((msg: Partial<WsMessage>) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify(msg));
        }
    }, []);

    return { send, connected };
}