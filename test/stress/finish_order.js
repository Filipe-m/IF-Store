import http from 'k6/http';
import {check} from 'k6';

export const options = {
    stages: [
        {target: 20, duration: '10s'},
        {target: 30, duration: '10s'},
        {target: 50, duration: '10s'},
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
    const productName = 'example_product_' + randomString;
    const productDescription = 'example product description ' + randomString;
    const price = Math.floor(Math.random() * 1000) + 1;
    const quantity = Math.floor(Math.random() * 100) + 1;

    const bodyUser = {username: username, email: email};
    const resUser = http.post('http://localhost:9091/users', JSON.stringify(bodyUser), {
        headers: {'Content-Type': 'application/json'},
    });

    check(resUser, {
        'status is 201': (r) => r.status === 201,
    });

    const userId = JSON.parse(resUser.body).id;

    const bodyProduct = {name: productName, description: productDescription, price: price};
    const resProduct = http.post('http://localhost:9094/product/register', JSON.stringify(bodyProduct), {
        headers: {'Content-Type': 'application/json'},
    });

    check(resProduct, {
        'status is 201': (r) => r.status === 201,
    });

    const productId = JSON.parse(resProduct.body).id;

    const bodyStock = {quantity: quantity};
    const resStock = http.put(`http://localhost:9094/stock/${productId}/add`, JSON.stringify(bodyStock), {
        headers: {'Content-Type': 'application/json'},
    });

    check(resStock, {
        'status is 201': (r) => r.status === 201,
    });

    const bodyOrderItem = {product_id: productId, quantity: quantity};
    const resOrderItem = http.post('http://localhost:9095/order-item', JSON.stringify(bodyOrderItem), {
        headers: {'Content-Type': 'application/json', 'USER-ID': userId},
    });

    check(resOrderItem, {
        'status is 201': (r) => r.status === 201,
    });

    const orderId = JSON.parse(resOrderItem.body).id;

    const bodyFinishOrder = {order_id: orderId, payment_data: "lorem"};
    const resFinishOrder = http.post('http://localhost:9095/order/finish', JSON.stringify(bodyFinishOrder), {
        headers: {'Content-Type': 'application/json', 'USER-ID': userId},
    });

    check(resFinishOrder, {
        'status is 201': (r) => r.status === 201,
    });
}