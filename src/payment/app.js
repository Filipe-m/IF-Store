import { PrismaClient } from '@prisma/client'
import express from 'express'
import valid_credit_card from './helpers/creditCardValidator.js'
const prisma = new PrismaClient()
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
        res.status(400).send('Bad Request')
      }
    }
  })
)

app.get('/paymentMethod/:userId', async (req, res) => {
  try {
    const userId = req.params.userId

    if (isNaN(parseInt(userId))) {
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
        paymentType: true,
        cardNumber: true,
        cardExpiration: true,
        cardCvv: true
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

  if (isNaN(parseInt(userId))) {
    return res.status(400).json({ message: 'Invalid user ID' })
  }

  if (req.body[0] === undefined) {
    return res
      .status(400)
      .json({ message: 'Invalid number of payment methods' })
  }
  const { number, expiration, cvv } = req.body[0]

  if (!number || !expiration || !cvv) {
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
        cardCvv: cvv
      }
    })
  } catch (e) {
    return res.status(409).json(e.message)
  }

  return res.status(200).json({ message: 'Card created' })
})

export default app
