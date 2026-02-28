import request from "./request";

export function getNotifications(params: { page?: number; page_size?: number }) {
    return request.get('/notifications', { params });
}

export function getUnreadCount() {
    return request.get('/notifications/unread');
}

export function markRead(id: number) {
    return request.put(`/notifications/${id}/read`);
}

export function markAllRead() {
    return request.put('/notifications/read-all');
}
