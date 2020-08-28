import axios from 'axios';

export function shortLink(link) {
    return axios.post('/api/v1/link', {
        link: link
    })
}

export function getLinkHash(link) {
    return axios.get('/api/v1/link', {
        params: {
            link: link
        }
    })
}
