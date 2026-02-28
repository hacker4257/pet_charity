import request from "./request";


export function getRescues(params:{
    page?:number;
    page_size?:number;
    species?:string;
    urgency?:string;
    status?:string;
}){
    return request.get('/rescues', { params });
}

export function getRescueDatil(id: number) {
    return request.get(`/rescues/${id}`);
}

export function getRescueMapData() {
    return request.get('/rescues/map');
}

export function createRescue(data:{
    title: string;
    description: string;
    species: string;
    urgency: string;
    longitude: number;
    latitude: number;
    address: string;
    contact_phone?:string;
}){
    return request.post('/rescues', data);
}


export function addFollow(rescueId: number, data:{content:string}) {
    return request.post(`/rescues/${rescueId}/follow`, data);
}

export function claimRescue(rescueId: number) {
    return request.post(`/rescues/${rescueId}/claim`);
}

export function updateClaimStatus(rescueId: number, data:{status: string}) {
    return request.put(`/rescues/${rescueId}/claim/status`, data);
}

export function getNearbyOrgs(params:{
    lng:number;
    lat:number;
    radius?:number;
}){
    return request.get('/organizations/nearby', { params });
}