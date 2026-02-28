import request from "./request";

export function toggleFavorite(petId: number) {
    return request.post(`/pets/${petId}/favorite`);
}

export function getFavoriteStatus(petId: number) {
    return request.get(`/pets/${petId}/favorite`);
}

export function getMyFavorites() {
    return request.get('/favorites');
}
