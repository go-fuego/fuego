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

var CategoriesTranslations = map[Category]map[Locale]string{
	CategoryOther:     {LocaleEn: "Other", LocaleFr: "Autre", LocaleEmoji: "🍽"},
	CategoryVegetable: {LocaleEn: "Vegetable", LocaleFr: "Légume", LocaleEmoji: "🥕"},
	CategoryFruit:     {LocaleEn: "Fruit", LocaleFr: "Fruit", LocaleEmoji: "🍎"},
	CategoryMeat:      {LocaleEn: "Meat", LocaleFr: "Viande", LocaleEmoji: "🥩"},
	CategoryDairy:     {LocaleEn: "Dairy", LocaleFr: "Produit laitier", LocaleEmoji: "🥛"},
	CategoryGrain:     {LocaleEn: "Grain", LocaleFr: "Céréale", LocaleEmoji: "🌾"},
	CategorySpice:     {LocaleEn: "Spice", LocaleFr: "Épice", LocaleEmoji: "🌶"},
	CategoryCondiment: {LocaleEn: "Condiment", LocaleFr: "Condiment", LocaleEmoji: "🧂"},
	CategorySweetener: {LocaleEn: "Sweetener", LocaleFr: "Édulcorant", LocaleEmoji: "🍬"},
	CategoryOil:       {LocaleEn: "Oil", LocaleFr: "Huile", LocaleEmoji: "🥥"},
	CategoryFat:       {LocaleEn: "Fat", LocaleFr: "Graisse", LocaleEmoji: "🥓"},
	CategoryLiquid:    {LocaleEn: "Liquid", LocaleFr: "Liquide", LocaleEmoji: "💧"},
	CategoryAlcohol:   {LocaleEn: "Alcohol", LocaleFr: "Alcool", LocaleEmoji: "🍺"},
}

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
