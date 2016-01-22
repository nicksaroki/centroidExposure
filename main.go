package main

import (
	"encoding/csv"
	"github.com/kellydunn/golang-geo"
	"os"
	"sort"
	"strconv"
)

type centroid struct {
	exposure, lat, lon float64
}

type Centroids []centroid

type dealer struct {
	id          int
	coordinates *geo.Point
	exposure    float64
}

func (slice Centroids) Len() int {
	return len(slice)
}

func (slice Centroids) Less(i, j int) bool {
	return slice[i].exposure < slice[j].exposure
}

func (slice Centroids) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func main() {
	latMin, _ := strconv.ParseFloat(os.Args[1], 64) //20.0
	latMax, _ := strconv.ParseFloat(os.Args[2], 64) //61.0
	lonMin, _ := strconv.ParseFloat(os.Args[3], 64) //-159.0
	lonMax, _ := strconv.ParseFloat(os.Args[4], 64) //-67.0
	step, _ := strconv.ParseFloat(os.Args[5], 64)   //.1
	radius, _ := strconv.ParseFloat(os.Args[6], 64) //7.0
	dealers := []dealer{}
	centroids := Centroids{}
	centroidsString := [][]string{}

	csvInputFile, _ := os.Open("dealers.csv")
	r := csv.NewReader(csvInputFile)
	dealerCSVRecords, _ := r.ReadAll()
	csvInputFile.Close()

	for _, importedDealer := range dealerCSVRecords {
		id, _ := strconv.Atoi(importedDealer[0])
		lat, _ := strconv.ParseFloat(importedDealer[1], 64)
		lon, _ := strconv.ParseFloat(importedDealer[2], 64)
		coordinates := geo.NewPoint(lat, lon)
		exposure, _ := strconv.ParseFloat(importedDealer[3], 64)
		currImportedDealer := dealer{id, coordinates, exposure}
		dealers = append(dealers, currImportedDealer)
	}

	//Step through all longitudes
	for lonCurr := lonMin; lonCurr < lonMax; lonCurr += step {
		//Step through all latitudes
		for latCurr := latMin; latCurr < latMax; latCurr += step {
			//Reset total exposure for new centroid
			totalExposure := 0.0
			//Convert centroid into a "point" usable by the distance function
			currCentroidPoint := geo.NewPoint(latCurr, lonCurr)
			//Step through all dealers at each centroid
			for _, currDealer := range dealers {
				//Evaluate if the distance between the current centroid and current dealer is less than the radius in miles
				if currCentroidPoint.GreatCircleDistance(currDealer.coordinates)*.62137 <= radius {
					totalExposure += currDealer.exposure
				}
			}
			//If centroid has any exposure, add it to our output data set
			if totalExposure > 0 {
				currCentroid := centroid{totalExposure, latCurr, lonCurr}
				centroids = append(centroids, currCentroid)
			}
		}
	}
	sort.Sort(sort.Reverse(centroids))

	for _, currCentroid := range centroids {
		centroidsSlice := []string{strconv.FormatFloat(currCentroid.exposure, 'f', -1, 64), strconv.FormatFloat(currCentroid.lat, 'f', 10, 64), strconv.FormatFloat(currCentroid.lon, 'f', 10, 64)}
		centroidsString = append(centroidsString, centroidsSlice)
	}
	csvOutputFile, _ := os.Create("output.csv")
	defer csvOutputFile.Close()
	w := csv.NewWriter(csvOutputFile)
	w.WriteAll(centroidsString)
	w.Flush()
}
