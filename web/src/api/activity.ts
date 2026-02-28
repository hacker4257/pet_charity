import request from "./request";

export function getLeaderboard(params?: { page?: number; page_size?: number }) {
    return request.get('/leaderboard', { params });
}

export function getMyRank() {
    return request.get('/leaderboard/me');
}

export function getFeed(params?: { page?: number; page_size?: number }) {
    return request.get('/feed', { params });
}
