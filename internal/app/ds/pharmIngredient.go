package ds

type PharmIngredient struct {
	ID              uint    `gorm:"primarykey"`
	Name            string  `gorm:"column:name"`
	ChemicalFormula string  `gorm:"column:chemical_formula"`
	MolecularWeight float64 `gorm:"column:molecular_weight"`
	Description     string  `gorm:"column:description"`
	Dosage          string  `gorm:"column:dosage"`
	Manufacturer    string  `gorm:"column:manufacturer"`
	Price           float64 `gorm:"column:price"`
	ImagePath       string  `gorm:"column:image_path"`
	Status          string  `gorm:"column:status"`
}
