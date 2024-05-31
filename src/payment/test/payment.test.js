import { expect } from 'chai'
import { randomUUID } from 'crypto'
import sinon from 'sinon'
import supertest from 'supertest'
import app from '../app.js' // Adjust the import path as needed
import prisma from '../libs/prisma.js'
import prismaMock from './prismaMock.js' // Adjust path if necessary

const request = supertest(app)

// Replace actual prisma methods with mocks
sinon.replace(prisma.payment, 'create', prismaMock.payment.create)
sinon.replace(prisma.payment, 'delete', prismaMock.payment.delete)
sinon.replace(
  prisma.paymentMethod,
  'findFirst',
  prismaMock.paymentMethod.findFirst
)

// Set default mock behavior
prismaMock.payment.create.resolves({})
prismaMock.payment.delete.resolves({})
prismaMock.paymentMethod.findFirst.resolves([])

afterEach(() => {
  sinon.reset()
})

describe('POST /payment/:userId', () => {
  it('should create a payment', async () => {
    const userId = randomUUID()
    const paymentMethodId = randomUUID()
    const orderId = randomUUID()
    prismaMock.paymentMethod.findFirst.resolves([
      {
        id: paymentMethodId,
        userId: userId
      }
    ])

    const response = await request.post(`/payment/${userId}`).send([
      {
        orderId: orderId,
        paymentMethodId: paymentMethodId,
        amount: Math.random()
      }
    ])

    expect(response.status).to.equal(200)
    expect(prismaMock.payment.create.callCount).to.equal(1)
  })
  it('should not create a payment if payment method does not belong to user', async () => {
    const userId = randomUUID()
    const paymentMethodId = randomUUID()
    const orderId = randomUUID()
    prismaMock.paymentMethod.findFirst.resolves(null)

    const response = await request.post(`/payment/${userId}`).send([
      {
        orderId: orderId,
        paymentMethodId: paymentMethodId,
        amount: Math.random()
      }
    ])

    expect(response.status).to.equal(403)
    expect(prismaMock.payment.create.callCount).to.equal(0)
  })
  it('should not create a payment if amount is not a number', async () => {
    const userId = randomUUID()
    const paymentMethodId = randomUUID()
    const orderId = randomUUID()
    prismaMock.paymentMethod.findFirst.resolves([
      {
        id: paymentMethodId,
        userId: userId
      }
    ])

    const response = await request.post(`/payment/${userId}`).send([
      {
        orderId: orderId,
        paymentMethodId: paymentMethodId,
        amount: 'not a number'
      }
    ])

    expect(response.status).to.equal(400)
    expect(prismaMock.payment.create.callCount).to.equal(0)
  })
  it('should not create a payment if missing parameters', async () => {
    const userId = randomUUID()
    const paymentMethodId = randomUUID()
    const orderId = randomUUID()
    prismaMock.paymentMethod.findFirst.resolves([
      {
        id: paymentMethodId,
        userId: userId
      }
    ])

    const response = await request.post(`/payment/${userId}`).send([
      {
        orderId: orderId
      }
    ])

    expect(response.status).to.equal(400)
    expect(prismaMock.payment.create.callCount).to.equal(0)
  })
  it('should not create a payment if invalid user ID', async () => {
    const response = await request.post('/payment/invalid-id').send([
      {
        orderId: randomUUID(),
        paymentMethodId: randomUUID(),
        amount: Math.random()
      }
    ])

    expect(response.status).to.equal(400)
    expect(prismaMock.payment.create.callCount).to.equal(0)
  })
  it('should not create a payment if invalid payment method ID', async () => {
    const userId = randomUUID()
    const orderId = randomUUID()
    const paymentMethodId = 'invalid-id'
    const response = await request.post(`/payment/${userId}`).send([
      {
        orderId: orderId,
        paymentMethodId: paymentMethodId,
        amount: Math.random()
      }
    ])

    expect(response.status).to.equal(400)
    expect(prismaMock.payment.create.callCount).to.equal(0)
  })
  it('should not create a payment if invalid order ID', async () => {
    const userId = randomUUID()
    const orderId = 'invalid-id'
    const paymentMethodId = randomUUID()
    const response = await request.post(`/payment/${userId}`).send([
      {
        orderId: orderId,
        paymentMethodId: paymentMethodId,
        amount: Math.random()
      }
    ])

    expect(response.status).to.equal(400)
    expect(prismaMock.payment.create.callCount).to.equal(0)
  })
  it('should not create a payment if no info is provided', async () => {
    const userId = randomUUID()
    const response = await request.post(`/payment/${userId}`).send([])

    expect(response.status).to.equal(400)
    expect(prismaMock.payment.create.callCount).to.equal(0)
  })
  it('should not create a payment if no body is provided', async () => {
    const userId = randomUUID()
    const response = await request.post(`/payment/${userId}`)

    expect(response.status).to.equal(400)
    expect(prismaMock.payment.create.callCount).to.equal(0)
  })
})

describe('DELETE /payment/:paymentId', () => {
  it('should delete a payment', async () => {
    const paymentId = randomUUID()
    prismaMock.payment.delete.resolves({})

    const response = await request.delete(`/payment/${paymentId}`)

    expect(response.status).to.equal(200)
    expect(prismaMock.payment.delete.callCount).to.equal(1)
  })
  it('should not delete a payment if invalid payment ID', async () => {
    const response = await request.delete('/payment/invalid-id')

    expect(response.status).to.equal(400)
    expect(prismaMock.payment.delete.callCount).to.equal(0)
  })
})
