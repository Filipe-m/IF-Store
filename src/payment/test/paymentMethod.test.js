import { expect } from 'chai'
import { randomUUID } from 'crypto'
import sinon from 'sinon'
import supertest from 'supertest'
import app from '../app.js' // Adjust the import path as needed
import prisma from '../libs/prisma.js'
import prismaMock from './prismaMock.js' // Adjust path if necessary

const request = supertest(app)

// Replace actual prisma methods with mocks
sinon.replace(
  prisma.paymentMethod,
  'findMany',
  prismaMock.paymentMethod.findMany
)
sinon.replace(prisma.paymentMethod, 'create', prismaMock.paymentMethod.create)
sinon.replace(prisma.paymentMethod, 'delete', prismaMock.paymentMethod.delete)

// Set default mock behavior
prismaMock.paymentMethod.findMany.resolves([])
prismaMock.paymentMethod.create.resolves({})
prismaMock.paymentMethod.delete.resolves({})

afterEach(() => {
  sinon.reset()
})

describe('GET /paymentMethod/:userId', () => {
  it('should create BOLETO and PIX if no methods are found', async () => {
    prismaMock.paymentMethod.findMany.resolves([])
    const uuid = randomUUID()
    const res = await request.get(`/paymentMethod/${uuid}`)
    expect(res.status).to.equal(200)
    expect(prismaMock.paymentMethod.create.callCount).to.equal(2)
  })

  it('should not create if an invalid id is provided', async () => {
    const res = await request.get('/paymentMethod/invalid-id')
    expect(res.status).to.equal(400)
    expect(prismaMock.paymentMethod.create.callCount).to.equal(0)
  })

  it('should not create new payment method if one already exists', async () => {
    const uuid = randomUUID()
    prismaMock.paymentMethod.findMany.resolves([
      {
        id: '27d99428-0a23-4c73-ab37-01cb7ce9a1e5',
        paymentType: 'BOLETO'
      },
      {
        id: '61e853ac-cb3c-4fe0-ab44-a05d7e97ecc8',
        paymentType: 'PIX'
      }
    ])

    const res = await request.get(`/paymentMethod/${uuid}`)
    expect(res.status).to.equal(200)
    expect(prismaMock.paymentMethod.create.callCount).to.equal(0)
  })
})

describe('POST /paymentMethod', () => {
  const validObject = [
    {
      number: '5203 0081 9897 5523',
      expiration: '02/12',
      cvv: 231,
      card_holder: 'John Doe'
    }
  ]

  it('should create a new payment method', async function () {
    const uuid = randomUUID()
    const res = await request.post(`/paymentMethod/${uuid}`).send(validObject)

    expect(res.status).to.equal(200)
    expect(prismaMock.paymentMethod.create.callCount).to.equal(1)
  })

  it('should not create if an invalid id is provided', async () => {
    const res = await request
      .post('/paymentMethod/invalid-id')
      .send(validObject)
    expect(res.status).to.equal(400)
    expect(prismaMock.paymentMethod.create.callCount).to.equal(0)
  })

  it('should not create if missing parameters', async () => {
    const res = await request.post(`/paymentMethod/${randomUUID()}`).send([
      {
        number: '5203 0081 9897 5523',
        expiration: '02/12'
      }
    ])
    expect(res.status).to.equal(400)
    expect(prismaMock.paymentMethod.create.callCount).to.equal(0)
  })
  it('should not create if invalid card number', async () => {
    const res = await request.post(`/paymentMethod/${randomUUID()}`).send([
      {
        number: '5203 0081 9897 552',
        expiration: '02/12',
        cvv: 231,
        card_holder: 'John Doe'
      }
    ])
    expect(res.status).to.equal(400)
    expect(prismaMock.paymentMethod.create.callCount).to.equal(0)
  })
  it('should not create if invalid expiration date', async () => {
    const res = await request.post(`/paymentMethod/${randomUUID()}`).send([
      {
        number: '5203 0081 9897 5523',
        expiration: '02-12',
        cvv: 231,
        card_holder: 'John Doe'
      }
    ])
    expect(res.status).to.equal(400)
    expect(prismaMock.paymentMethod.create.callCount).to.equal(0)
  })
  it('should not create if invalid CVV', async () => {
    const res = await request.post(`/paymentMethod/${randomUUID()}`).send([
      {
        number: '5203 0081 9897 5523',
        expiration: '02/12',
        cvv: 23,
        card_holder: 'John Doe'
      }
    ])
    expect(res.status).to.equal(400)
    expect(prismaMock.paymentMethod.create.callCount).to.equal(0)
  })
  it('should not create if card_holder is not provided', async () => {
    const res = await request.post(`/paymentMethod/${randomUUID()}`).send([
      {
        number: '5203 0081 9897 5523',
        expiration: '02/12',
        cvv: 235
      }
    ])
    expect(res.status).to.equal(400)
    expect(prismaMock.paymentMethod.create.callCount).to.equal(0)
  })

  it('should return 409 if prisma throws an error', async () => {
    prismaMock.paymentMethod.create.rejects({ meta: { cause: 'error' } })
    const res = await request
      .post(`/paymentMethod/${randomUUID()}`)
      .send(validObject)
    expect(res.status).to.equal(409)
  })
})

describe('DELETE /paymentMethod/:paymentId', () => {
  it('should delete a payment method', async () => {
    const uuid = randomUUID()
    const res = await request.delete(`/paymentMethod/${uuid}`)
    expect(res.status).to.equal(200)
    expect(prismaMock.paymentMethod.delete.callCount).to.equal(1)
  })

  it('should not delete if an invalid id is provided', async () => {
    const res = await request.delete('/paymentMethod/invalid-id')
    expect(res.status).to.equal(400)
    expect(prismaMock.paymentMethod.delete.callCount).to.equal(0)
  })
})
