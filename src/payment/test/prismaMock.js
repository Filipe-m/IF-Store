// test/prismaMock.js
import sinon from 'sinon'

const prismaMock = {
  paymentMethod: {
    findMany: sinon.stub(),
    create: sinon.stub(),
    delete: sinon.stub(),
    findFirst: sinon.stub()
  },
  payment: {
    create: sinon.stub(),
    delete: sinon.stub()
  }
}

export default prismaMock
