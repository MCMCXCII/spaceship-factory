package inventory

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/MCMCXCII/spaceship-factory/inventory/internal/repository/record"
)

type Repository struct {
	mu   sync.RWMutex
	data map[uuid.UUID]record.Part
}

// New создаёт хранилище деталей и заполняет его начальными данными.
func New() *Repository {
	parts := seedParts()

	data := make(map[uuid.UUID]record.Part, len(parts))
	for _, p := range parts {
		data[p.UUID] = p
	}

	return &Repository{data: data}
}

// seedParts возвращает начальный набор деталей склада.
func seedParts() []record.Part {
	now := time.Now()

	return []record.Part{
		{
			UUID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			Name:          "Алюминиевый корпус",
			Description:   "Лёгкий корпус для небольших кораблей",
			Price:         500000,
			StockQuantity: 10,
			PartType:      record.PartTypeHull,
			CreatedAt:     now,
		},
		{
			UUID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
			Name:          "Титановый корпус",
			Description:   "Прочный корпус для дальних перелётов",
			Price:         1500000,
			StockQuantity: 5,
			PartType:      record.PartTypeHull,
			CreatedAt:     now,
		},
		{
			UUID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
			Name:          "Ионный двигатель C",
			Description:   "Экономичный двигатель базового класса",
			Price:         300000,
			StockQuantity: 8,
			PartType:      record.PartTypeEngine,
			CreatedAt:     now,
		},
		{
			UUID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
			Name:          "Ионный двигатель B",
			Description:   "Мощный двигатель повышенной тяги",
			Price:         800000,
			StockQuantity: 3,
			PartType:      record.PartTypeEngine,
			CreatedAt:     now,
		},
		{
			UUID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
			Name:          "Энергетический щит",
			Description:   "Защита от метеоритов и лазерного огня",
			Price:         400000,
			StockQuantity: 6,
			PartType:      record.PartTypeShield,
			CreatedAt:     now,
		},
		{
			UUID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"),
			Name:          "Лазерная пушка",
			Description:   "Базовое вооружение для самообороны",
			Price:         250000,
			StockQuantity: 7,
			PartType:      record.PartTypeWeapon,
			CreatedAt:     now,
		},
		{
			UUID:          uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"),
			Name:          "Плазменный корпус",
			Description:   "Экспериментальный корпус с плазменным покрытием",
			Price:         2000000,
			StockQuantity: 0,
			PartType:      record.PartTypeHull,
			CreatedAt:     now,
		},
	}
}
