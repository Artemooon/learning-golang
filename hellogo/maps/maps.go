package main

import "fmt"

type Vertex struct {
	Lat float64
	Long float64
}

var m map[string]Vertex

func main() {
	m := make(map[string]Vertex)
	m2 := make(map[string]int)
	m["BellLabs"] = Vertex{
		40.68433, -74.39967,
	}
    m["KFC"] = Vertex{
        33.4342, -4034.42342,
    }
	fmt.Println(m["Bell Labs"])
	fmt.Println(m["KFC"])

	var m1 = map[string]Vertex{
    	"Bell Labs": Vertex{
    		40.68433, -74.39967,
    	},
    	"Google": Vertex{
    		37.42202, -122.08408,
    	},
        "Macdondalds": Vertex{
            503.4324, -234.3432,
        },
        "Apple": {
            7945.43,
            -953.3432,
        },
    }
    delete(m1, "Apple")
	fmt.Println(m1)
	fmt.Println(m1["Apple"])
    fmt.Println(m1["Ford"])

    for _, v := range m1 {
        fmt.Println("mapItem", v)
    }
    fmt.Println("m2", m2["test"])
}
