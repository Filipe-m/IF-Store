
# IF Store

Um servidor de um marktplace

## Documentação da API

### Retorna todos os metodos de pagamento que um usuário possui

```http
  GET /paymentMethod/${id}
```

| Parâmetro   | Tipo       | Descrição                           |
| :---------- | :--------- | :---------------------------------- |
| `id` | `string` | **Obrigatório**. ID do usuário |

### Cadastra um novo cartão para o usuário

```http
  POST /paymentMethod/${id}
```

| Parâmetro   | Tipo       | Descrição                                   |
| :---------- | :--------- | :------------------------------------------ |
| `id` | `string` | **Obrigatório**. ID do usuário |

#### Estrutura Json a ser enviada

```json
[{
  "number": "5203 0081 9897 5523",
  "expiration": "02/12",
  "cvv": 231
}]
```

## Stack utilizada

**Back-end:** Node, Go

## Autores

- [@Filipe-m](https://www.github.com/octokatherine)
- [@victorvcruz](https://github.com/victorvcruz)
