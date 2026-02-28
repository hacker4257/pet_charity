import request from './request';

// 注册机构
export function createOrg(data: {
    name: string;
    description: string;
    address: string;
    longitude?: number;
    latitude?: number;
    contact_phone: string;
    capacity?: number;
}) {
    return request.post('/organizations', data);
}

// 更新机构信息
export function updateOrg(id: number, data: {
    description?: string;
    address?: string;
    contact_phone?: string;
    capacity?: number;
}) {
    return request.put(`/organizations/${id}`, data);
}

// 获取机构详情
export function getOrgDetail(id: number) {
    return request.get(`/organizations/${id}`);
}

// 机构列表（公开）
export function getOrgList(params: { page?: number; page_size?: number }) {
    return request.get('/organizations', { params });
}


// 获取当前用户的机构信息
export function getMyOrg() {
    return request.get('/organizations/mine');
}


