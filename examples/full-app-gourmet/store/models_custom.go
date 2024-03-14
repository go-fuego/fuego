package store

type RecipeWithDosings struct {
	Recipe  `json:"recipe"`
	Dosings []GetIngredientsOfRecipeRow `json:"dosings"`
}
