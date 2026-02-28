import request from "./request";



export function getPets(params:{
    page?:number;
    pageSize?:number;
    species?:string;
    gender?:string;
    status?:string;
    org_id?:number;
}) {
    return request.get('/pets', {params});
}

export function getPetDetail(id: number) {
    return request.get(`/pets/${id}`);
}

export function createPet(data:{
    name: string;
    species?:string;
    breed?: string;
    age: number;
    gender: string;
    description:string;
    cover_image?:string;
}) {
    return request.post('/pets', data);
}

// 删除宠物
export function deletePet(id: number) {
    return request.delete(`/pets/${id}`);
}

// 更新宠物
export function updatePet(id: number, data: {
    name?: string;
    species?: string;
    breed?: string;
    age?: number;
    gender?: string;
    description?: string;
    cover_image?: string;
    status?: string;
}) {
    return request.put(`/pets/${id}`, data);
}

// 上传宠物图片（修正路径）
export function uploadPetImage(petId: number, file: File, sortOrder?: number) {
    const formData = new FormData();
    formData.append('image', file);
    if (sortOrder !== undefined) {
        formData.append('sort_order', String(sortOrder));
    }
    return request.post(`/pets/${petId}/images`, formData);
}