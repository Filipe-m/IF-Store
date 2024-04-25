import http from 'k6/http';
import {check} from 'k6';

export const options = {
    stages: [
        {target: 1, duration: '5s'},
        {target: 5, duration: '5s'},
        {target: 10, duration: '10s'},
    ],
    thresholds: {
        http_req_failed: ['rate<0.01'],
        http_req_duration: ['p(95)<500'],
        checks: ['rate>0.99'],
    },
};

export default function () {
    const randomString = Math.random().toString(36).substring(2);

    const username = 'example_user_' + randomString;
    const email = 'user_' + randomString + '@example.com';

    const body = {username: username, email: email}

    const res = http.post('http://localhost:9091/users', JSON.stringify(body), {
        headers: {'Content-Type': 'application/json'},
    });

    check(res, {
        'status is 201': (r) => r.status === 201,
        'response body': (r) => {
            const body = JSON.parse(r.body);
            return body.id !== "";
        },
    });
}