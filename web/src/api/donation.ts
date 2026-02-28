import request from "./request";

export function createDonation(data:{
    target_type: string;
    target_id?: number;
    amount: number;
    message?: string;
    payment_method:string;
}) {
    return request.post('/donations', data)
}

export function getDonationStatus(id: number) {
    return request.get(`/donations/${id}/status`);
}

export function getMyDonations(params: { page?: number; page_size?: number }) {
return request.get('/donations/mine', { params });
}

export function getPublicDonations(params: {
    page?: number;
    page_size?: number;
    target_type?: string;
    target_id?: number;
}) {
    return request.get('/donations/public', { params });
}