import { check } from 'k6';
import http from 'k6/http';

export const options = {
    discardResponseBodies: true,
    scenarios: {
        local: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '5s', target: 1 },
                { duration: '5s', target: 3 },
                { duration: '5s', target: 1 },
                { duration: '10s', target: 5 },
                { duration: '5s', target: 1 },
            ],
            gracefulRampDown: '0s',
        },
    },
}

const shortCodes = [
    'abc123',
    'def456',
    'mno345',
]

export default function () {
    const randomShortCode = shortCodes[Math.floor(Math.random() * shortCodes.length)];
    const resp = http.get(`http://localhost:4242/${randomShortCode}`, { redirects: 0 });
    console.log(`Requested short code: ${randomShortCode}, Status: ${resp.status}, Body: ${resp.body}`);
    check(resp, { 'is status 302': (r) => r.status === 302 } );
}
