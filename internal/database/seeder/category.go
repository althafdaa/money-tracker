package seeder

import (
	"encoding/json"
	"money-tracker/internal/utils"
	"os"
	"strconv"

	"gorm.io/gorm"
)

type Category struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Seeder struct {
	db *gorm.DB
}

func (s *Seeder) seedCategory() error {
	is_seeding, err := strconv.ParseBool(os.Getenv("SEED"))

	if !is_seeding || err != nil {
		return nil
	}

	file, err := os.Open("public/default_category.json")

	if err != nil {
		return err
	}

	categories := []Category{}
	err = json.NewDecoder(file).Decode(&categories)
	if err != nil {
		return err
	}

	for _, category := range categories {
		slug, err := utils.Slugify(category.Name)
		if err != nil {
			println(err)
			return err
		}

		println("Seeding category: ", category.Name)
		err = s.db.Exec("insert into category (name, slug, type) values (?, ?, ?) returning *", category.Name, slug, category.Type).Error

		if err != nil {
			println(err)
			return err
		}
		println("Seeding success: ", category.Name)
	}

	return nil
}

func (s *Seeder) SeedEverything() error {
	is_seeding, err := strconv.ParseBool(os.Getenv("SEED"))

	if !is_seeding || err != nil {
		return nil
	}

	err = s.seedCategory()
	if err != nil {
		return err
	}

	return nil
}

func NewSeeder(
	db *gorm.DB,
) *Seeder {
	is_seeding, err := strconv.ParseBool(os.Getenv("SEED"))

	if err != nil {
		is_seeding = false
	}
	println("ENV SEED: ", is_seeding)

	return &Seeder{
		db: db,
	}
}
