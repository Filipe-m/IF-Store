package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"inventory/cmd/handlers"
	"inventory/internal/config"
	"inventory/internal/platform"
	"inventory/internal/product"
	"inventory/internal/stock"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("[LOAD ENVIRONMENT VARIABLES FAIL]: %s\n", err.Error())
	}
}

func main() {
	app := fiber.New()
	app.Use(logger.New())

	cfg := config.Load()

	connect, err := platform.NewPostgresConnect(cfg.Database)
	if err != nil {
		log.Fatalf("[CONNECT DATABASE FAIL]: %s", err.Error())
	}

	err = platform.Migrate(connect, &product.Product{}, &stock.Stock{})
	if err != nil {
		log.Fatalf("[MIGRATE DATABASE FAIL]: %s", err.Error())
	}

	productRepository := product.NewRepository(connect)

	stockRepository := stock.NewRepository(connect)

	productHandler := handlers.NewProductHandler(productRepository, stockRepository)

	stockHandler := handlers.NewStockHandler(stockRepository)

	Default(productRepository, stockRepository)

	// Routes
	app.Use(cors.New())
	app.Post("/product/register", productHandler.RegisterProduct)
	app.Put("/product/:id", productHandler.UpdateProduct)
	app.Get("/product/:id", productHandler.FindProduct)
	app.Delete("/product/:id", productHandler.DeleteProduct)

	app.Put("/stock/:productId/add", stockHandler.AddStock)
	app.Put("/stock/:productId/remove", stockHandler.RemoveStock)
	app.Get("/stock/:productId", stockHandler.FindStock)

	app.Get("/product", productHandler.ListProducts)

	if err = app.Listen(":9094"); err != nil {
		log.Fatalf("[START SERVER FAIL]: %s", err.Error())
	}
}

func Default(productRepository product.Repository, stockRepository stock.Repository) {
	products := []*product.Product{
		{
			Name:        "Jogo de Facas Tramontina Plenus",
			Description: "As lâminas de aço inox tem maior durabilidade do fio devido ao tratamento térmico. Os cabos de polipropileno possuem maior resistência e durabilidade. Podem ir à máquina de lavar louças facilitando seu dia a dia.",
			Price:       decimal.NewFromFloat(85.51),
		},
		{
			Name:        "Kindle 11ª Geração",
			Description: "O Kindle 11ª Geração é 11% mais fino e 16% mais leve do que a geração anterior. Ele é vendido em duas opções de cores: preto e branco. A tela de 6 polegadas é antirreflexo e sensível ao toque, facilitando a leitura em qualquer lugar.",
			Price:       decimal.NewFromFloat(449.10),
		},
		{
			Name:        "Smart TV LED 50” 4K Samsung",
			Description: "A Smart TV LED 50” 4K Samsung possui resolução 4 vezes maior do que as TVs Full HD. Ela é compatível com assistentes virtuais, como Alexa e Google Assistente, e possui conexão Bluetooth para você conectar seus dispositivos sem fio.",
			Price:       decimal.NewFromFloat(2.999),
		},
		{
			Name:        "Fogão 4 Bocas Consul",
			Description: "O Fogão 4 Bocas Consul possui acendimento automático, que acende ao simples girar do botão. Ele possui um forno de 58 litros, com prateleiras deslizantes e desmontáveis, facilitando a limpeza.",
			Price:       decimal.NewFromFloat(1.199),
		},
		{
			Name:        "Ar Condicionado Split 9000 BTUs",
			Description: "O Ar Condicionado Split 9000 BTUs é ideal para ambientes de até 15m². Ele possui filtro antibactéria, que retém até 99% das bactérias, e filtro antipoeira, que retém as partículas de poeira do ar.",
			Price:       decimal.NewFromFloat(1.299),
		},
		{
			Name:        "Geladeira Frost Free 375L",
			Description: "A Geladeira Frost Free 375L possui prateleiras de vidro temperado, que suportam até 35kg. Ela possui um compartimento extra frio, que resfria mais rápido os alimentos e bebidas.",
			Price:       decimal.NewFromFloat(2.499),
		},
		{
			Name:        "Máquina de Lavar 11kg Electrolux",
			Description: "A Máquina de Lavar 11kg Electrolux possui 12 programas de lavagem, que se adaptam a diferentes tipos de roupas. Ela possui um sistema de reutilização de água, que reaproveita a água da lavagem para outros fins.",
			Price:       decimal.NewFromFloat(1.499),
		},
		{
			Name:        "Micro-ondas 20L Electrolux",
			Description: "O Micro-ondas 20L Electrolux possui 10 opções de potência, que se adaptam a diferentes tipos de alimentos. Ele possui um menu de receitas pré-programadas, que facilitam o preparo de alimentos.",
			Price:       decimal.NewFromFloat(399.90),
		},
		{
			Name:        "Liquidificador 2L Philco",
			Description: "O Liquidificador 2L Philco possui 12 velocidades, que se adaptam a diferentes tipos de alimentos. Ele possui uma função pulsar, que tritura os alimentos de forma mais eficiente.",
			Price:       decimal.NewFromFloat(199.90),
		},
		{
			Name:        "Cafeteira 1,2L Britânia",
			Description: "A Cafeteira 1,2L Britânia possui um sistema corta-pingos, que permite servir o café antes do término do preparo. Ela possui um filtro permanente, que dispensa o uso de filtros de papel.",
			Price:       decimal.NewFromFloat(99.90),
		},
		{
			Name:        "Panela de Pressão 4,5L Rochedo",
			Description: "A Panela de Pressão 4,5L Rochedo possui um sistema de segurança, que evita acidentes durante o uso. Ela possui um fundo triplo, que distribui o calor de forma uniforme.",
			Price:       decimal.NewFromFloat(99.90),
		},
	}

	for _, p := range products {
		ctx := context.Background()
		err := productRepository.Create(ctx, p)
		if err != nil {
			log.Fatal(err)
		}

		err = stockRepository.Create(ctx, &stock.Stock{
			ProductID: p.ID,
			Quantity:  100,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
