import request from 'supertest'
import app from '../app'

describe("GET /paymentMethod/-1", () => {
  it("O resultado da consulta deve ser 200", async () => {
    const res = await request(app).get("/paymentMethod/-1");
    expect(res.statusCode).toEqual(200);
  });
});