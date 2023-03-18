# go-clean-architecture

Desafio #3 do curso [Go Expert](https://goexpert.fullcycle.com.br/curso/).

## Instruções

Após clonar o projeto, instale as dependências rodando `go mod tidy` na pasta raiz.

Inicie o banco de dados com o comando `docker compose up -d`

Rode o comando `go run main.go wire_gen.go` na pasta `/cmd/ordersystem`

O app sobe 3 servidores:
* REST na porta 8000
* GraphQL na porta 8080
* gRPC na porta 50051

## Uso
### Para interagir com os endpoints REST:

Requisição `POST` para `http://localhost:8000/order` para criar uma Order. Exemplo de body da requisição:
```
{
	"id": "rest",
	"price": 100.5,
	"tax": 0.5
}
```

Requisição `GET` para `http://localhost:8000/order` para listar todas as Orders.


### Para interagir com o servidor GraphQL:

Abrindo o navegador em `http://localhost:8080` o GraphQL Playground estará aberto. Exemplo de mutation para criar uma Order:
```
mutation createOrder {
  	createOrder(input: { id: "graphql", Price: 50, Tax: 0.7 }) {
		id
		Price
		Tax
		FinalPrice
	}
}
```

Query para buscar listar todas as Orders:
```
query listOrders {
  listOrders {
    id
    Price
    Tax
    FinalPrice
  }
}
```


### Para interagir com o servidor gRPC:

Utilizando o [evans](https://github.com/ktr0731/evans), basta rodar `evans -r repl` na raiz do projeto e usar `call CreateOrder` para criar uma Order e `call ListOrders` para listar todas as Orders.

