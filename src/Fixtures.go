package src

func RunFixture() {
	p1 := &Product{
		ID:          1,
		Title:       "Wildcraft Laptop Bags",
		Description: "Wildcraft Laptop Bags with Description",
		Price:       1419,
		WebUrl:      "https://smartdoko.com/shop/wildcraft-gravity-nylon-bag-8903338009566",
		ImageUrl:    "https://smartdoko.com/assets/upload/images/product/detail/IMG-0600cc412ddc0ff29dd2e437b89f9cc1.gif",
		Payload:     "flsymAvu",
	}

	p2 := &Product{
		ID:          2,
		Title:       "Samsung Gear S2 Platinum",
		Description: "Step up your style with the Samsung Gear S2 classic. Genuine leather, precious metal and exceptional finishes come together to create a sophisticated wearable that goes with anything.",
		Price:       2040,
		WebUrl:      "https://smartdoko.com/shop/samsung-galaxy-s2-sports-gear-band",
		ImageUrl:    "https://smartdoko.com/assets/upload/images/product/detail/IMG-23d641578ee1dab84459058b8a02e1ee.jpg",
		Payload:     "krvixZGh",
	}

	p3 := &Product{
		ID:          3,
		Title:       "Multipurpose Water Spray Gun ",
		Description: "Expandable & Flexible Plastic Water Garden Pipe with Spray Nozzle For Car Wash Pet Bath",
		Price:       1395,
		WebUrl:      "https://smartdoko.com/shop/multipurpose-water-spray-gun",
		ImageUrl:    "https://smartdoko.com/assets/upload/images/product/detail/IMG-332b397997496f4c2bcdf1d24ebb814a.gif",
		Payload:     "cQxAVrlD",
	}

	Products[p1.ID] = p1
	Products[p2.ID] = p2
	Products[p3.ID] = p3

}
