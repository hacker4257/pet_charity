import request from './request';

export function getConversations() {
    return request.get('/chat/conversations');
}

export function getPrivateHistory(userId: number, params: { page?: number;
page_size?: number }) {
    return request.get(`/chat/private/${userId}`, { params });
}

export function getRoomHistory(roomId: string, params: { page?: number; page_size?:
number }) {
    return request.get(`/chat/room/${roomId}`, { params });
}