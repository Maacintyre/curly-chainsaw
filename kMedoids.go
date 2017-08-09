package main 

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"math"
)

func check(e error){
	 if e != nil {
	 	panic(e)
	 }
}

type csvTS struct {
	name			string
	width, length 	int
	data 			[][]float64
	norms 			[]string
}

func (c csvTS) getDimensions() csvTS {
	file,err := os.Open(c.name)
	check(err)

	input := make([]byte, 1)
	c.width = 0
	c.length = 0
	running := true

	var buffer 	[]byte
	var line	[]float64

	for i:=0; running==true; i++ {
		_, err = file.Read(input)
		//check(err)
		if err != nil {
			if err == io.EOF {
				break
			} 
			panic(err)
		}


		if input[0] == ',' {
			b, err := strconv.ParseFloat(string(buffer),64)
			check(err)
			line = append(line, b)
			buffer = make([]byte,0)
		} else if input[0] == '\n' {
			c.norms = append(c.norms,string(buffer))
			buffer = make([]byte,0)
			c.data = append(c.data, line)
			line = make([]float64,0)
		} else {
			buffer = append(buffer,input[0])
		}

		
	}
	c.norms = append(c.norms,string(buffer))
	c.data = append(c.data, line)
	file.Close()
	c.width = len(c.data[0])
	c.length = len(c.data)
	return c
}

func (c csvTS) get(x int, y int) float64 {
	if x < c.length && y < c.width {
		return c.data[x][y]
	}
	return float64(0)
}

func (c csvTS) getNorm(x int) string {
	if x < c.length {
		return c.norms[x]
	}
	return ""
}

func (c csvTS) dist(l int, r int) float64 {
	var bucket float64
	bucket = 0.0

	for i:=0; i<c.width; i++ {
		left := c.get(l,i)
		right := c.get(r,i)
		//fmt.Println("Row indexes are: ", l," ", r)
		//fmt.Println("\tItem 1: ", left)
		//fmt.Println("\tItem 2: ", right)
		bucket += math.Pow((left-right), 2)
	}
	return math.Sqrt(bucket)
}

func (c csvTS) mDist(l int,r int) float64 {
	var bucket float64
	bucket = 0.0

	for i:=0; i<c.width; i++ {
		left := c.get(l,i)
		right := c.get(r,i)
		bucket = bucket + math.Abs(left-right)
	}
	return bucket
}

type record struct {
	row int
	dists []float64
	norm string
}

func sum(x []float64) float64{
	bucket := float64(0)

	for _,num := range x {
		bucket += num
	}

	return bucket
}


func main() {
	fileName := "diabetes1.csv"
	var bucket float64
	xCost := make([]float64, 0)
	yCost := make([]float64, 0)

	tCost := make([]float64, 0)

	lowestCost := float64(0)

	//expectedRun := float64(1)
	//iterations := float64(0)
	//percentage := float64(0)
	//complete := float64(0)

	m1 := record{}
	m2 := record{}

	file := csvTS{name: fileName}
	file = file.getDimensions()

	/*for i:=float64(file.length); i>float64(0); i--{
		fmt.Println(i)
		expectedRun = expectedRun * i
	}
	fmt.Println(expectedRun)*/

	for i:=0; i<file.length-1; i++ {
		//x.row = i

		for j:=i+1; j<file.length; j++ {

			if i != j {
				for k:=0; k<file.length; k++ {
					if k != j && k != i {
						bucket = file.mDist(i,k)
						//fmt.Println(bucket)
						xCost = append(xCost, bucket)
						//bucket = file.mDist(j,k)
						//fmt.Println("\t",bucket)
						bucket = file.mDist(j,k)
						yCost = append(yCost, bucket)
					} else {
						xCost = append(xCost, float64(0))
						yCost = append(yCost, float64(0))
					}
				}
				for l:=0; l<len(xCost); l++ {
					if xCost[l] != 0 && yCost[l] != 0 {
						if xCost[l] <= yCost[l] {
							yCost[l] = float64(0)
						} else {
							xCost[l] = float64(0)
						}
					}
				}
				//fmt.Println("Distance Matrix: ")
				//for k:=0; k<len(xCost); k++ {
					//fmt.Println("Row: ", k,"\t", xCost[k], " ", yCost[k])
				//}
				bucket = sum(xCost) + sum(yCost)
				tCost = append(tCost, bucket)
				if (lowestCost > float64(0) &&  bucket < lowestCost) || lowestCost == float64(0){
					lowestCost = bucket
					//fmt.Println(i == j)
					m1.row = i
					m1.dists = xCost
					m2.row = j
					m2.dists = yCost
				}
				xCost = make([]float64, 0)
				yCost = make([]float64, 0)
			}
		}
	}

	fmt.Println("Cost of instances: ")
	side := math.Floor(math.Sqrt(float64(len(tCost))))

	bucket = float64(0)
	for i:=0; i<len(tCost); i++ {
		fmt.Print(math.Floor(tCost[i]), " ")
		if len(tCost) > 10 {
			if bucket > side{
				fmt.Println()
				bucket = float64(0)
			} else {
				bucket += float64(1)
			}
		}
	}

	fmt.Println("\nData of lowest cost: ")
	fmt.Println("Lowest Cost: ", lowestCost)
	fmt.Println("Medoid 1 row", m1.row)
	fmt.Println("Medoid 1 distances\n\t", m1.dists)
	fmt.Println("Medoid 2 row", m2.row)
	fmt.Println("Medoid 2 distances\n\t", m2.dists)
	fmt.Println()
}