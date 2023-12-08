package types

import "fmt"

type Category string

const (
	CategoryOther     Category = "other"
	CategoryVegetable Category = "vegetable"
	CategoryFruit     Category = "fruit"
	CategoryMeat      Category = "meat"
	CategoryDairy     Category = "dairy"
	CategoryGrain     Category = "grain"
	CategorySpice     Category = "spice"
	CategoryCondiment Category = "condiment"
	CategorySweetener Category = "sweetener"
	CategoryOil       Category = "oil"
	CategoryFat       Category = "fat"
	CategoryLiquid    Category = "liquid"
	CategoryAlcohol   Category = "alcohol"
)

// CategoryValues is a slice of all valid categories
var CategoryValues = []Category{CategoryOther, CategoryVegetable, CategoryFruit, CategoryMeat, CategoryDairy, CategoryGrain, CategorySpice, CategoryCondiment, CategorySweetener, CategoryOil, CategoryFat, CategoryLiquid, CategoryAlcohol}

type InvalidCategoryError struct {
	Category Category
}

func (e InvalidCategoryError) Error() string {
	return fmt.Sprintf("invalid category %s. Valid categories are: %v", e.Category, CategoryValues)
}

func (c Category) Valid() bool {
	for _, v := range CategoryValues {
		if v == c {
			return true
		}
	}
	return false
}
