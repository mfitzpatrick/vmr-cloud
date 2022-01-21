package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMergeMaps(t *testing.T) {
	getTime := func(t *testing.T, str string) time.Time {
		tm, err := time.Parse(time.RFC3339, str)
		assert.Equal(t, nil, err)
		return tm
	}
	testMerge := func(m1, m2 map[string]interface{}) map[string]interface{} {
		mergeMaps(m1, m2)
		return m1
	}
	assert.Equal(t, map[string]interface{}{
		"ts1":  getTime(t, "2022-01-01T05:05:05Z"),
		"ts2":  getTime(t, "2022-01-01T05:05:05Z"),
		"num1": 5,
		"num2": 9,
		"num3": 0,
		"str1": "hi 1",
		"str2": "hi 2",
	}, testMerge(map[string]interface{}{
		"ts1":  getTime(t, "2022-01-01T05:05:05Z"),
		"ts2":  getTime(t, "0001-01-01T00:00:00Z"),
		"num2": 9,
		"num3": 0,
		"str1": "hi 1",
	}, map[string]interface{}{
		"ts1":  getTime(t, "0001-01-01T00:00:00Z"),
		"ts2":  getTime(t, "2022-01-01T05:05:05Z"),
		"num1": 5,
		"str2": "hi 2",
	}))
}

func TestToInterfaceMap(t *testing.T) {
	getTime := func(t *testing.T, str string) time.Time {
		tm, err := time.Parse(time.RFC3339, str)
		assert.Equal(t, nil, err)
		return tm
	}
	expect := toInterfaceMap(voyage{}) // empty interface map
	override := map[string]interface{}{
		"start-time": getTime(t, "2022-01-01T11:12:13Z"),
		"title":      "test title",
		"weather": map[string]interface{}{
			"seaway-tide": map[string]interface{}{
				"height-metres": 1.1,
				"time":          getTime(t, "2022-01-01T17:13:14Z"),
			},
		},
	}
	mergeMaps(expect, override)
	assert.Equal(t, expect, toInterfaceMap(voyage{
		StartTime: getTime(t, "2022-01-01T11:12:13Z"),
		Title:     "test title",
		Weather: weather{Tide: tide{
			Height: 1.1,
			Time:   getTime(t, "2022-01-01T17:13:14Z"),
		}},
	}))
}

func equalVoyage(t *testing.T, vExpect, v1, v2 voyage) bool {
	vOut, err := mergeVoyageStructs(v1, v2)
	if !assert.Equal(t, nil, err) {
		return false
	}
	if !assert.Equal(t, true, reflect.DeepEqual(vExpect, vOut)) {
		return false
	}
	return true
}

func TestMergeVoyageStructs(t *testing.T) {
	getTime := func(t *testing.T, str string) time.Time {
		tm, err := time.Parse(time.RFC3339, str)
		assert.Equal(t, nil, err)
		return tm
	}
	v1 := voyage{
		VoyageID: 1,
	}
	equalVoyage(t, v1, v1, voyage{})

	equalVoyage(t,
		voyage{
			VoyageID: 2,
		}, voyage{
			VoyageID: 0,
		}, voyage{
			VoyageID: 2,
		},
	)

	equalVoyage(t,
		voyage{
			VoyageID: 2,
		}, voyage{
			VoyageID: 2,
		}, voyage{
			VoyageID: 0,
		},
	)

	equalVoyage(t,
		voyage{
			VoyageID: 2,
			Weather: weather{
				Wind: wind{
					Speed: 10,
				},
			},
		}, voyage{
			VoyageID: 0,
			Weather: weather{
				Wind: wind{
					Speed: 10,
				},
			},
		}, voyage{
			VoyageID: 2,
		},
	)

	equalVoyage(t,
		voyage{
			VoyageID: 2,
			Weather: weather{
				Wind: wind{
					Speed: 10,
				},
			},
		}, voyage{
			VoyageID: 0,
		}, voyage{
			VoyageID: 2,
			Weather: weather{
				Wind: wind{
					Speed: 10,
				},
			},
		},
	)

	equalVoyage(t,
		voyage{
			VoyageID: 2,
			Weather: weather{
				Wind: wind{
					Speed: 10,
				},
				Tide: tide{
					Height: 3.72,
				},
			},
			Desc: "sample description",
		}, voyage{
			VoyageID: 0,
			Weather: weather{
				Wind: wind{
					Speed: 10,
				},
			},
			Desc: "sample description",
		}, voyage{
			VoyageID: 2,
			Weather: weather{
				Tide: tide{
					Height: 3.72,
				},
			},
		},
	)

	equalVoyage(t,
		voyage{
			VoyageID: 2,
			Weather: weather{
				Wind: wind{
					Speed: 10,
					Dir:   319,
				},
				Tide: tide{
					Height: 3.72,
				},
			},
			Skipper: crew{
				Name: "Alice",
				Rank: "INSHORE_SKIPPER",
			},
		}, voyage{
			VoyageID: 0,
			Weather: weather{
				Wind: wind{
					Speed: 10,
				},
			},
			Skipper: crew{
				Name: "Alice",
				Rank: "INSHORE_SKIPPER",
			},
		}, voyage{
			VoyageID: 2,
			Weather: weather{
				Wind: wind{
					Dir: 319,
				},
				Tide: tide{
					Height: 3.72,
				},
			},
		},
	)

	equalVoyage(t,
		voyage{
			VoyageID:  3,
			StartTime: getTime(t, "2022-01-01T11:12:13Z"),
		}, voyage{
			VoyageID:  0,
			StartTime: getTime(t, "2022-01-01T11:12:13Z"),
		}, voyage{
			VoyageID: 3,
		},
	)
}

func equalAssist(t *testing.T, vExpect, v1, v2 assist) bool {
	vOut, err := mergeAssistStructs(v1, v2)
	if !assert.Equal(t, nil, err) {
		return false
	}
	if !assert.Equal(t, true, reflect.DeepEqual(vExpect, vOut)) {
		return false
	}
	return true
}

func TestMergeAssistStructs(t *testing.T) {
	getTime := func(t *testing.T, str string) time.Time {
		tm, err := time.Parse(time.RFC3339, str)
		assert.Equal(t, nil, err)
		return tm
	}
	v1 := assist{
		AssistID: 1,
	}
	equalAssist(t, v1, v1, assist{})

	equalAssist(t,
		assist{
			AssistID: 2,
		}, assist{
			AssistID: 0,
		}, assist{
			AssistID: 2,
		},
	)

	equalAssist(t,
		assist{
			AssistID: 2,
		}, assist{
			AssistID: 2,
		}, assist{
			AssistID: 0,
		},
	)

	equalAssist(t,
		assist{
			AssistID: 2,
			Client: client{
				Name: "Alice",
			},
		}, assist{
			AssistID: 0,
			Client: client{
				Name: "Alice",
			},
		}, assist{
			AssistID: 2,
		},
	)

	equalAssist(t,
		assist{
			AssistID: 2,
			Client: client{
				Name: "Alice",
			},
		}, assist{
			AssistID: 0,
		}, assist{
			AssistID: 2,
			Client: client{
				Name: "Alice",
			},
		},
	)

	equalAssist(t,
		assist{
			AssistID: 2,
			Client: client{
				Name:  "Alice",
				Phone: "123456",
			},
			Problem: "sample problem",
		}, assist{
			AssistID: 0,
			Client: client{
				Name: "Alice",
			},
			Problem: "sample problem",
		}, assist{
			AssistID: 2,
			Client: client{
				Phone: "123456",
			},
		},
	)

	equalAssist(t,
		assist{
			AssistID: 2,
			Client: client{
				Name:     "Alice",
				Phone:    "123456",
				MemberNo: 57,
			},
			Pickup: scene{
				Loc: location{
					Name: "Coomera Waters",
					GPS: coordinate{
						Lat:  -27.0192739,
						Long: 153.2937465,
					},
				},
				Time: getTime(t, "2022-01-01T01:02:03Z"),
			},
		}, assist{
			AssistID: 0,
			Client: client{
				Name:     "Alice",
				MemberNo: 57,
			},
			Pickup: scene{
				Loc: location{
					Name: "Coomera Waters",
					GPS: coordinate{
						Long: 153.2937465,
					},
				},
				Time: getTime(t, "2022-01-01T01:02:03Z"),
			},
		}, assist{
			AssistID: 2,
			Pickup: scene{
				Loc: location{
					GPS: coordinate{
						Lat: -27.0192739,
					},
				},
			},
			Client: client{
				Phone: "123456",
			},
		},
	)
}
