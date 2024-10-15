package utils

import "ff/web/app/services"

func FindPersonByID(id int, people *[]services.Person) services.Person {
	var person services.Person
	for _, p := range *people {
		if p.ID == uint(id) {
			person = p
		}
	}

	return person
}

func FindFeatureFlagByID(id int, featureFlags *[]services.FeatureFlag) services.FeatureFlag {
	var featureFlag services.FeatureFlag
	for _, ff := range *featureFlags {
		if ff.ID == id {
			featureFlag = ff
		}
	}

	return featureFlag
}
