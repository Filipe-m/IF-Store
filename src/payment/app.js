import express from 'express'
import valid_credit_card from './helpers/creditCardValidator.js'
import isUUID from './helpers/uuidValidator.js'
import prisma from './libs/prisma.js'

const app = express()

app.use(
  express.json({
    verify: (req, res, buf, encoding) => {
      if (buf.length === 0) {
        return
      }
      try {
        JSON.parse(buf)
      } catch (e) {
        res.status(400).send({ message: 'Bad Request' })
      }
    }
  })
)

app.get('/paymentMethod/:userId', async (req, res) => {
  try {
    const userId = req.params.userId

    if (!isUUID(userId)) {
      return res.status(400).json({ message: 'Invalid user ID' })
    }

    const getUser = await prisma.paymentMethod.findMany({
      where: {
        userId: userId,
        paymentType: {
          not: 'CREDIT_CARD'
        }
      }
    })

    if (getUser.length === 0) {
      await prisma.paymentMethod.create({
        data: {
          userId: userId,
          paymentType: 'BOLETO'
        }
      })
      await prisma.paymentMethod.create({
        data: {
          userId: userId,
          paymentType: 'PIX'
        }
      })
    }

    const noCards = await prisma.paymentMethod.findMany({
      where: {
        userId: userId,
        active: true,
        paymentType: {
          not: 'CREDIT_CARD'
        }
      },
      select: {
        id: true,
        paymentType: true
      }
    })

    const cards = await prisma.paymentMethod.findMany({
      where: {
        userId: userId,
        active: true,
        paymentType: 'CREDIT_CARD'
      },
      select: {
        id: true,
        paymentType: true,
        cardNumber: true,
        cardExpiration: true,
        cardCvv: true,
        cardHolder: true
      }
    })

    let response

    if (cards.length !== 0) {
      response = [...noCards, ...cards]
    } else {
      response = noCards
    }

    return res.status(200).json(response)
  } catch (error) {
    return res.status(500).json({ message: 'An error occurred' })
  }
})

app.post('/paymentMethod/:userId', async (req, res) => {
  const userId = req.params.userId

  if (!isUUID(userId)) {
    return res.status(400).json({ message: 'Invalid user ID' })
  }

  if (req.body[0] === undefined) {
    return res
      .status(400)
      .json({ message: 'Invalid number of payment methods' })
  }
  const { number, expiration, cvv, card_holder } = req.body[0]

  if (!number || !expiration || !cvv || !card_holder) {
    return res.status(400).json({ message: 'Missing parameters' })
  }

  if (valid_credit_card(number) === false) {
    return res.status(400).json({ message: 'Invalid card number' })
  }
  const dateRegex = /^(0[1-9]|1[0-2])\/\d{2}$/
  if (dateRegex.test(expiration) === false) {
    return res.status(400).json({ message: 'Invalid expiration date' })
  }

  const cvvRegex = /^[0-9]{3,4}$/
  if (cvvRegex.test(cvv) === false) {
    return res.status(400).json({ message: 'Invalid CVV' })
  }

  try {
    await prisma.paymentMethod.create({
      data: {
        userId: userId,
        paymentType: 'CREDIT_CARD',
        cardNumber: number,
        cardExpiration: expiration,
        cardCvv: cvv,
        cardHolder: card_holder
      }
    })
  } catch (e) {
    return res.status(409).json({ message: e })
  }

  return res.status(200).json({ message: 'Card created' })
})

app.delete('/paymentMethod/:paymentId', async (req, res) => {
  const paymentId = req.params.paymentId

  if (!isUUID(paymentId)) {
    return res.status(400).json({ message: 'Invalid payment ID' })
  }

  try {
    await prisma.paymentMethod.delete({
      where: {
        id: paymentId
      }
    })
  } catch (e) {
    return res.status(409).json({ message: e.meta.cause })
  }

  res.status(200).json({ message: 'Payment method deleted' })
})

app.post('/payment/:userId', async (req, res) => {
  const userId = req.params.userId

  if (!isUUID(userId)) {
    return res.status(400).json({ message: 'Invalid user ID' })
  }

  if (req.body[0] === undefined) {
    return res.status(400).json({ message: 'No info provided' })
  }
  const { orderId, paymentMethodId, amount } = req.body[0]

  if (isNaN(amount)) {
    return res
      .status(400)
      .json({ message: 'The provided amount is not a number' })
  }

  if (!isUUID(orderId) || !isUUID(paymentMethodId)) {
    return res.status(400).json({ message: 'Invalid ID' })
  }

  if (!orderId || !paymentMethodId || !amount) {
    return res.status(400).json({ message: 'Missing parameters' })
  }

  const belongsToUser = await prisma.paymentMethod.findFirst({
    where: {
      id: paymentMethodId,
      userId: userId
    }
  })

  if (!belongsToUser) {
    return res
      .status(403)
      .json({ message: 'Payment method does not belong to user' })
  }

  try {
    await prisma.payment.create({
      data: {
        orderId: orderId,
        userId: userId,
        amount: amount,
        paymentMethod: paymentMethodId
      }
    })
    return res.status(200).json({ message: 'Payment created' })
  } catch (e) {
    return res.status(409).json({ message: e.meta.cause })
  }
})

app.delete('/payment/:paymentId', async (req, res) => {
  const paymentId = req.params.paymentId

  if (!isUUID(paymentId)) {
    return res.status(400).json({ message: 'Invalid payment ID' })
  }

  try {
    await prisma.payment.delete({
      where: {
        id: paymentId
      }
    })
  } catch (e) {
    console.log(e)
    return res.status(409).json({ message: e })
  }

  res.status(200).json({ message: 'Payment deleted' })
})

export default app
