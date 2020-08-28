import axios from "axios";

export function getLinkStats(link) {
    return axios.get('/api/v1/link/stats', {
        params: {
            link: link
        }
    })
}

export function getTopLinkStats() {
    return axios.get('/api/v1/link/stats/top')
}
