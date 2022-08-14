package benchmarks

func sampleJsonObject() map[string]interface{} {
	return map[string]interface{}{
		"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
		"user":      "Jon Doe",
		"score":     -100,
		"is active": false,
		"nested": map[string]interface{}{
			"a": 1,
			"b": 2,
			"c": map[string]interface{}{
				"d": map[string]interface{}{
					"e": 3,
				},
			},
			"f": []int{4, 5, 6},
			"g": map[string]interface{}{
				"h": "A",
				"i": true,
				"j": 1,
				"k": 1.5,
				"l": nil,
			},
		},
	}
}
