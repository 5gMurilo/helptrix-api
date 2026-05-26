package seeder

import (
	"github.com/5gMurilo/helptrix-api/core/domain"
	categoryinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/category"
)

func SeedCategories(svc categoryinterfaces.ICategoryService) error {
	rows := []domain.Category{
		{Name: "Residential Cleaning", Description: "Professional housekeeping, organization, and sanitizing of domestic spaces"},
		{Name: "Plumbing", Description: "Installation and repair of pipes, faucets, and drainage systems"},
		{Name: "Electrical Services", Description: "Installation, maintenance, and repair of electrical panels, outlets, and lighting"},
		{Name: "Gardening & Landscaping", Description: "Lawn care, pruning, planting, and landscape maintenance services"},
		{Name: "IT Support", Description: "Technical assistance for computers, networks, software setup, and troubleshooting"},
		{Name: "Childcare & Babysitting", Description: "Professional childcare and supervised babysitting for families"},
		{Name: "Furniture Assembly", Description: "Assembly, installation, and arrangement of modular and custom furniture"},
		{Name: "Pet Care & Dog Walking", Description: "Dog walking, pet sitting, and basic care sessions for animals"},
		{Name: "Private Tutoring", Description: "Personalized academic support and homework assistance for students of all levels"},
		{Name: "Local Delivery & Errands", Description: "Pickup and delivery of parcels and documents within the city"},
		{Name: "House Painting & Renovation", Description: "Interior and exterior painting, wall repair, and home renovation services"},
		{Name: "Personal Chef & Meal Prep", Description: "Custom meal preparation and cooking services for individuals and events"},
	}

	return svc.Seed(rows)
}
