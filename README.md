# Full Cycle - Desafio GoLang Pós Go Expert - Clean Architecture

Este projeto é parte de um desafio da Pós Go Expert, nele são cobertos os conhecimentos em Go Rotines, channels, contextos, tratamentos de erros, packages, Clean Architecture, gRPC, GraphQL, APIRest, Eventos, DI.


### 📋  Usaremos o Docker para subir nossa imagem com o comando abaixo

* Build
```sh
docker compose up -d
```  


### Informações dos Serviços

**APIRestful - [Porta 8000]**

```plaintext
GET /order  - Listagem de todas as `orders`
POST /order - Criação de uma `order`
```

**GraphQL - [Porta 8080]**

```plaintext
Query
    - orders: [Order!]!
Mutation
    - createOrder(input: OrderInput): Order
```

**gRPC**

```bash
evans -r repl -p 50051
package pb
service OrderService
```

- Create order:

```bash
call CreateOrder
```

- List all orders:

```bash
call ListOrders
```