package main

func updateCarFromTo(from *Car, to *Car) {
	if from.Make != "" {
		to.Make = from.Make
	}

	if from.Model != "" {
		to.Model = from.Model
	}

	if from.Year != 0 {
		to.Year = from.Year
	}

	if from.Color != "" {
		to.Color = from.Color
	}

	if from.Category != "" {
		to.Category = from.Category
	}

	if from.Package != "" {
		to.Package = from.Package
	}

	if from.Mileage != 0 {
		to.Mileage = from.Mileage
	}

	if from.Price != 0 {
		to.Price = from.Price
	}

}
