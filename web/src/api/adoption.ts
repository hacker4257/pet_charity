import request from "./request";


export function createAdoption(data: {
    pet_id: number;
    reason: string;
    living_condition: string;
    experience?: string;
}) {
    return request.post('/adoptions', data);
}

export function getMyAdoptions(params: { page?: number; page_size?: number }) {
    return request.get('/adoptions/mine', { params });
}

export function getAdoptionDetail(id: number) {
    return request.get(`/adoptions/${id}`);
}

// 机构：收到的申请
export function getOrgAdoptions(params: {
    page?: number;
    page_size?: number;
    status?: string;
}) {
    return request.get('/adoptions/org', { params });
}

  // 机构：审核
export function reviewAdoption(
    id: number,
    data: { status: string; reject_reason?: string }
) {
    return request.put(`/adoptions/${id}/review`, data);
}

  // 机构：确认完成
export function completeAdoption(id: number) {
    return request.put(`/adoptions/${id}/complete`);
}