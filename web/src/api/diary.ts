import request from "./request";

// 创建日记
export function createDiary(data: { pet_id: number; content: string }) {
    return request.post('/diaries', data);
}

// 日记详情
export function getDiaryDetail(id: number) {
    return request.get(`/diaries/${id}`);
}

// 更新日记
export function updateDiary(id: number, data: { content: string }) {
    return request.put(`/diaries/${id}`, data);
}

// 删除日记
export function deleteDiary(id: number) {
    return request.delete(`/diaries/${id}`);
}

// 上传日记图片
export function uploadDiaryImage(diaryId: number, file: File, sortOrder?: number) {
    const formData = new FormData();
    formData.append('image', file);
    if (sortOrder !== undefined) {
        formData.append('sort_order', String(sortOrder));
    }
    return request.post(`/diaries/${diaryId}/images`, formData);
}

// 删除日记图片
export function deleteDiaryImage(diaryId: number, imageId: number) {
    return request.delete(`/diaries/${diaryId}/images/${imageId}`);
}

// 点赞/取消点赞
export function toggleDiaryLike(diaryId: number) {
    return request.post(`/diaries/${diaryId}/like`);
}

// 公开日记列表
export function getDiariesPublic(params: { page?: number; page_size?: number }) {
    return request.get('/diaries', { params });
}

// 按宠物查日记
export function getDiariesByPet(petId: number, params: { page?: number; page_size?: number }) {
    return request.get(`/diaries/pet/${petId}`, { params });
}
