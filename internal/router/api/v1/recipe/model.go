package recipe

type Group struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	RecipeIDs   []int  `json:"recipe_ids"` // IDs of recipes in the group
}

type Recipe struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	PreparationTime int      `json:"preparation_time"` // in minutes
	CookingTime     int      `json:"cooking_time"`     // in minutes
	Steps           []string `json:"steps"`
	Servings        int      `json:"servings"`
	Difficulty      string   `json:"difficulty"`
}
