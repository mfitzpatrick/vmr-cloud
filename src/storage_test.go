package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			Pickup: location{
				Name: "Coomera Waters",
				GPS: coordinate{
					Lat:  -27.0192739,
					Long: 153.2937465,
				},
			},
		}, assist{
			AssistID: 0,
			Client: client{
				Name:     "Alice",
				MemberNo: 57,
			},
			Pickup: location{
				Name: "Coomera Waters",
				GPS: coordinate{
					Long: 153.2937465,
				},
			},
		}, assist{
			AssistID: 2,
			Pickup: location{
				GPS: coordinate{
					Lat: -27.0192739,
				},
			},
			Client: client{
				Phone: "123456",
			},
		},
	)
}
