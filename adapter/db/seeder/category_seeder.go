package seeder

import (
	"fmt"

	"github.com/5gMurilo/helptrix-api/core/domain"
	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) error {
	rows := []domain.Category{
		{Name: "Residential cleaning", Description: "housekeeping, organization, and sanitizing of domestic spaces"},
		{Name: "Plumber", Description: "installation and repair of pipes, faucets, and toilets"},
		{Name: "Electrician", Description: "panels, outlets, lighting, and minor electrical repairs"},
		{Name: "Gardening", Description: "pruning, planting, lawn and planter maintenance"},
		{Name: "Child care", Description: "babysitting and supervision of children at agreed times"},
		{Name: "Furniture assembly", Description: "assembly and disassembly of custom and modular furniture"},
		{Name: "IT support", Description: "network setup, backup, and software guidance"},
		{Name: "Local deliveries", Description: "pickup and delivery of parcels and documents in the city"},
		{Name: "Private tutoring", Description: "homework help and on-demand study support"},
		{Name: "Pet walks", Description: "walks and basic care for dogs in short sessions"},
	}

	for i := range rows {
		row := rows[i]
		res := db.Where("name = ?", row.Name).FirstOrCreate(&row)
		if res.Error != nil {
			return fmt.Errorf("seed category %q: %w", row.Name, res.Error)
		}
	}

	return nil
}
