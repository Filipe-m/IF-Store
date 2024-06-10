describe('end2end', () => {
  it('add products', () => {
    cy.request('POST', 'http://localhost:9094/product/register', {
      name: "Sanduicheira",
      description: "Sanduicheira Elétrica Grill Click 220V",
      price: 50.5
    }).then(
        (response) => {
          const productId = response.body.id;
          cy.request('PUT', `http://localhost:9094/stock/${productId}/add`, { quantity: 20 })
        }
    )

    cy.request('POST', 'http://localhost:9094/product/register', {
      name: "Churrasqueira",
      description: "Churrasqueira Elétrica Grill Click 220V",
      price: 69.5
    }).then(
        (response) => {
          const productId = response.body.id;
          cy.request('PUT', `http://localhost:9094/stock/${productId}/add`, { quantity: 20 })
        }
    )

    cy.request('POST', 'http://localhost:9094/product/register', {
      name: "Faca de Cozinha",
      description: "Faca de Cozinha Tramontina",
      price: 140.99
    }).then(
        (response) => {
          const productId = response.body.id;
          cy.request('PUT', `http://localhost:9094/stock/${productId}/add`, { quantity: 20 })
        }
    )
  })

  it('complete flow', () => {
    cy.visit('http://localhost:3000')
    // create user
    cy.get('#username').type('testuser');
    cy.get('#username').should('have.value', 'testuser')

    cy.get('#email').type('testuser@example.com');
    cy.get('#email').should('have.value', 'testuser@example.com')

    cy.get('form').submit();

    // add product to cart
    cy.url().should('include', '/products');

    cy.get('.card-body button').first().click();
    cy.get('.alert-success').should('be.visible');

    cy.get('.card-body button').eq(1).click();
    cy.get('.alert-success').should('be.visible');

    cy.get('.card-body button').eq(2).click();
    cy.get('.alert-success').should('be.visible');

    // finish order
    cy.get('.btn-outline-primary').click();
    cy.url().should('include', '/cart');

    cy.get('.card-title').should('have.length', 3);

    cy.get('.d-flex > .btn').click();
    cy.get('.modal').should('be.visible');

    cy.get('#address').type('Rua Teste, 123');
    cy.get('#address').should('have.value', 'Rua Teste, 123')

    cy.get('#neighborhood').type('Bairro Teste');
    cy.get('#neighborhood').should('have.value', 'Bairro Teste')

    cy.get('#city').type('Cidade Teste');
    cy.get('#city').should('have.value', 'Cidade Teste')

    cy.get('#zipCode').type('12345-678');
    cy.get('#zipCode').should('have.value', '12345-678')

    cy.get(':nth-child(2) > .nav-link').click();

    cy.get('.show > .btn').click();

    cy.get('#creditCardNumber').type('5203 0081 9897 5523');
    cy.get('#creditCardNumber').should('have.value', '5203 0081 9897 5523')

    cy.get('#cardOwner').type('Teste Testador');
    cy.get('#cardOwner').should('have.value', 'Teste Testador')

    cy.get('#securityCode').type('123');
    cy.get('#securityCode').should('have.value', '123')

    cy.get('#expirationDate').type('12/26');
    cy.get('#expirationDate').should('have.value', '12/26')

    cy.get(':nth-child(3) > .btn').click();

    cy.get('#paymentMethods').select('CREDIT_CARD');

    cy.get('.modal-footer > .btn-primary').last().click();
    cy.get('.alert').should('contain', 'Compra finalizada com sucesso!');
  })
})
