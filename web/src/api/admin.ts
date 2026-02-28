import request from './request';

// 后台统计数据
export function getAdminStats() {
    return request.get('/admin/stats');
}

// 待审核机构列表
export function getPendingOrgs(params: { page?: number; page_size?: number }) {
    return request.get('/admin/organizations/pending', { params });
}

// 审核机构
export function reviewOrg(id: number, data: { status: string; reject_reason?: string
}) {
    return request.put(`/admin/organizations/${id}/review`, data);
}

// 用户列表
export function getUserList(params: { page?: number; page_size?: number; role?:
string }) {
    return request.get('/admin/users', { params });
}

// 修改用户角色
export function updateUserRole(id: number, data: { role: string }) {
    return request.put(`/admin/users/${id}/role`, data);
}

// 禁用/启用用户
export function updateUserStatus(id: number, data: { status: string }) {
    return request.put(`/admin/users/${id}/status`, data);
}