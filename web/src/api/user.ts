import request from "./request";


export function getMe() {
    return request.get('/users/me');
}

export function updateProfile(data:{
    nickname?: string;
    phone?:string;
    email?:string;
    language?: string;
}){
    return request.put('/users/me', data);
}

export function changePassword(data:{
    old_password: string;
    new_password:string;
}){
    return request.put('/users/me/password',data);
}

export function uploadeAvatar(file:File) {
    const formData = new FormData();
    formData.append('avatar', file);
    return request.put('/users/me/avatar', formData);
}