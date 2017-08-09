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

func (c csvTS) mDist(l int,r int) []float64 {
	bucket := make([]float64,c.width)

	for i:=0; i<c.width; i++ {
		left := c.get(l,i)
		right := c.get(r,i)
		bucket[i] = math.Pow((left-right),2)
	}
	return bucket
}

type record struct {
	row int
	dist float64
	norm string
}

func main() {
	var iterations float64
	var complete float64
	fileName := "test.csv"
	iterations = 0.0
	complete = 0.0
	percent:=0
	x := record{}
	hit := record{row: -1}
	miss := record{row: -1}

	file := csvTS{name: fileName}
	file = file.getDimensions()

	values := make([]float64,file.width)

	for i:=0; i<file.length; i++ {
		x.row = i
		x.norm = file.getNorm(i)
		for j:=0; j<file.length; j++ {
			if j != i {
				iterations += 1.0
				complete = iterations/float64(file.length * file.length-1)
				complete = math.Floor(complete * 100)
				if int(complete) > percent {
					percent = int(complete)
					fmt.Println("Percentage Complete: ", percent)
				}
				dist := file.dist(i, j)
				if file.getNorm(j) == x.norm {
					if hit.row == -1 || dist < hit.dist {
						hit.row = j
						hit.dist = dist
					}
				} else {
					if miss.row == -1 || dist < miss.dist {
						miss.row = j
						miss.dist = dist
					}
				}
			}
		}
		temp := file.mDist(i,hit.row)
		//fmt.Println("\tValues: ", values)
		//fmt.Println("\tTemp: ", temp)
		for k:=0; k<len(values); k++ {
			values[k] = values[k] - temp[k]
		}
		temp1 := file.mDist(i,miss.row)
		//fmt.Println("\tValues: ", values)
		//fmt.Println("\tTemp1: ", temp1)
		for k:=0; k<len(values); k++ {
			values[k] = values[k] + temp1[k]
		}

	}
	m:= float64(file.length)
	for i:=0; i<len(values); i++ {
		values[i] = values[i]/m
	}

	fmt.Println("Results for file: ", fileName)
	fmt.Println("File length: ", file.length)
	fmt.Println("Number of Attributes: ", file.width)
	fmt.Println("Number of cells processed: ", file.length * file.width)
	fmt.Println("This Program has run for ", iterations, " iterations.")

	fmt.Println("The values for each column are:")
	for i:=0;i<len(values); i++ {
		fmt.Println("\tAttribute ", i, " ", values[i])
	}
}