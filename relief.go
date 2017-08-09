package main 

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv")

func check(e error){
	 if e != nil {
	 	panic(e)
	 }
}

type csv struct {
	name 			string
	width, length	int
}

func (c csv) getDimensions() csv {
	file, err := os.Open(c.name)
	check(err)

	input := make([]byte, 1)
	c.width = 0
	c.length = 0
	running := true

	maxWidth := 0
	maxLength := 0

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
			maxWidth += 1
		} else if input[0] == '\n' {
			maxWidth = 0
			maxLength += 1
		}

		if maxWidth > c.width {
			c.width = maxWidth
		}

		if maxLength > c.length {
			c.length = maxLength
		}

		
	}
	file.Close()
	c.width++
	c.length++
	return c
}

func (c csv) get(x int, y int) string {
	var buffer string

	if x < c.length && y < c.width {
		file, err := os.Open(c.name)
		check(err)
		//fmt.Println("X and Y are: " ,x ," ", y)
		for i:=0; i < x; i++ {
			//fmt.Println("X block ran:")
			_,file = getUntil('\n', file)
		}
		for j:=0; j <= y; j++ {
			//fmt.Println("Y block ran:")
			buffer,file = getUntil(' ', file)
		}
		
		file.Close()
	}
	return buffer
}

func getUntil(b byte, f *os.File) (string, *os.File){
	
	input := make([]byte, 1)
	buffer := make([]byte, 0)
	running := true


	for i:=0; running == true; i++ {
		_, err := f.Read(input)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if (b == ' '){
			if input[0] == ',' || input[0] == '\n' {
				running = false
			} else {
				buffer = append(buffer, input[0])
			}
		} else {
			if input[0] == b {
				running = false
				//return string(buffer)
			} else {
				buffer = append(buffer, input[0])
			}
		}
	}
	return string(buffer), f
}

func (c csv) dist(l int, r int) float64 {
	var bucket float64
	bucket = 0.0

	for i:=0; i<c.width-1; i++ {
		lTemp := c.get(l,i)
		rTemp := c.get(r,i)
		fmt.Println("Row indexes are: ", l," ", r)
		fmt.Println("\tItem 1: ", lTemp)
		fmt.Println("\tItem 2: ", rTemp)
		left, err := strconv.ParseFloat(lTemp,64)
		check(err)
		right, err := strconv.ParseFloat(rTemp,64)
		check(err)

		bucket += math.Pow((left-right), 2)
	}
	return math.Sqrt(bucket)
}

func (c csv) mDist(l int,r int) []float64 {
	bucket := make([]float64,c.width-1)

	for i:=0; i<c.width-1; i++ {
		lTemp := c.get(l,i)
		rTemp := c.get(r,i)

		left, err := strconv.ParseFloat(lTemp,64)
		check(err)
		right, err := strconv.ParseFloat(rTemp,64)
		check(err)

		bucket[i] = math.Pow((left-right),2)
	}
	return bucket
}

type record struct {
	row int
	dist float64
	norm string
}

/*func (lhs record) dist(p record) float64 {
	var bucket float64
	bucket = 0.0
	for i:=0; i<len(lhs.coords); i++ {
		bucket += math.Pow((lhs.coords[i]-p.coords[0]), 2)
	}
	return math.Sqrt(bucket)
}*/

/*func (p record) coords() int int {
	return p.x, p,y
}*/

func main() {
	x := record{}
	hit := record{row: -1}
	miss := record{row: -1}

	file := csv{name: "diabetes.csv"}
	file = file.getDimensions()

	values := make([]float64,file.width-1)

	/*for i:=0; i<file.length; i++ {
		for j:=0; j<file.width; j++ {
			fmt.Print(file.get(i,j), " ")
		}
		fmt.Println()
	}*/
	//fmt.Println("CheckPoint 1:")
	for i:=0; i<file.length; i++ {
		x.row = i
		x.norm = file.get(i,file.width-1)
		//fmt.Println("CheckPoint 2:")
		for j:=0; j<file.length; j++ {
			if j != i {
				dist := file.dist(i, j)
				if file.get(j, file.width-1) == x.norm {
					//Take the hit
					if hit.row == -1 || dist < hit.dist {
						hit.row = j
						hit.dist = dist
					}
				} else {
					//Take the miss
					if miss.row == -1 || dist < miss.dist {
						miss.row = j
						miss.dist = dist
					}
				}
			}
		}
		//fmt.Println("CheckPoint 3:")
		temp := file.mDist(i,hit.row)
		//fmt.Println("\tValues: ", values)
		//fmt.Println("\tTemp: ", temp)
		for k:=0; k<len(values); k++ {
			values[k] = values[k] - temp[k]
		}
		//fmt.Println("CheckPoint 4:")
		temp1 := file.mDist(i,miss.row)
		//fmt.Println("\tValues: ", values)
		//fmt.Println("\tTemp1: ", temp1)
		for k:=0; k<len(values); k++ {
			values[k] = values[k] + temp1[k]
		}

	}
	//fmt.Println("CheckPoint 5:")
	m:= float64(file.length * file.width-1)
	for i:=0; i<len(values); i++ {
		values[i] = values[i]/m
	}

	fmt.Println("The values for each column are:")
	fmt.Println(values)
}