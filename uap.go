package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// index in input position
const (
	physics = iota
	chemistry
	maths
	cs
	final
)

const (
	Biotech = iota
	Chemistry
	Engineering
	Mathematics
	Physics
)

type Departments map[string][]Applicant
type Results map[string]float64

type Applicant struct {
	id       int
	fullname string
	res      Results
	prior    [3]string
}

var departmentsNames = []string{"Biotech", "Chemistry", "Engineering", "Mathematics", "Physics"}

/*
Initializes map of departments that holds:
  - department name as key
  - slice of accepted applicants
  - zero department keeps track of unaccepted candidates on every iteration of acceptance procedure
*/
func initDepartments() Departments {
	departments := map[string][]Applicant{
		departmentsNames[0]: {},
		departmentsNames[1]: {},
		departmentsNames[2]: {},
		departmentsNames[3]: {},
		departmentsNames[4]: {},
		"0":                 {},
	}
	return departments
}

func chooseBestResult(applicants *[]Applicant) {
	for _, ap := range *applicants {
		for _, dep := range departmentsNames {
			ap.res[dep] = math.Max(ap.res["Final"], ap.res[dep])
		}
	}
}

/*Creates array of Applicants structures that holds information about every student*/
func makeApplicants(file *os.File) ([]Applicant, error) {

	scanner := bufio.NewScanner(file)

	var input [][]string
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		input = append(input, strings.Fields(line))
		i++
	}
	var scores []float64
	totalNumber := len(input)
	Applicants := make([]Applicant, totalNumber)
	for i := 0; i < totalNumber; i++ {
		Applicants[i].id = i + 1
		Applicants[i].fullname = input[i][0] + " " + input[i][1]
		for a := 2; a <= 6; a++ {
			res, err := strconv.ParseFloat(input[i][a], 64)
			if err != nil {
				return nil, err
			}
			scores = append(scores, res)
		}
		Applicants[i].res = make(map[string]float64)
		Applicants[i].res[departmentsNames[Biotech]] = (scores[chemistry] + scores[physics]) / 2
		Applicants[i].res[departmentsNames[Chemistry]] = scores[chemistry]
		Applicants[i].res[departmentsNames[Engineering]] = (scores[cs] + scores[maths]) / 2
		Applicants[i].res[departmentsNames[Mathematics]] = scores[maths]
		Applicants[i].res[departmentsNames[Physics]] = (scores[physics] + scores[maths]) / 2
		Applicants[i].res["Final"] = scores[final]
		scores = nil
		Applicants[i].prior[0] = (input)[i][7]
		Applicants[i].prior[1] = (input)[i][8]
		Applicants[i].prior[2] = (input)[i][9]
	}

	chooseBestResult(&Applicants)
	return Applicants, nil
}

func compareApplicants(applicants *[]Applicant, domain string) {
	gpaI := 0.0
	gpaJ := 0.0
	sort.Slice(*applicants, func(i, j int) bool {
		switch domain {
		case departmentsNames[Biotech]:
			gpaI = (*applicants)[i].res[departmentsNames[Biotech]]
			gpaJ = (*applicants)[j].res[departmentsNames[Biotech]]
		case departmentsNames[Chemistry]:
			gpaI = (*applicants)[i].res[departmentsNames[Chemistry]]
			gpaJ = (*applicants)[j].res[departmentsNames[Chemistry]]
		case departmentsNames[Engineering]:
			gpaI = (*applicants)[i].res[departmentsNames[Engineering]]
			gpaJ = (*applicants)[j].res[departmentsNames[Engineering]]
		case departmentsNames[Mathematics]:
			gpaI = (*applicants)[i].res[departmentsNames[Mathematics]]
			gpaJ = (*applicants)[j].res[departmentsNames[Mathematics]]
		case departmentsNames[Physics]:
			gpaI = (*applicants)[i].res[departmentsNames[Physics]]
			gpaJ = (*applicants)[j].res[departmentsNames[Physics]]
		}
		if gpaI != gpaJ {
			return gpaI > gpaJ
		}
		return (*applicants)[i].fullname < (*applicants)[j].fullname
	})
}

/*
Distributes applicants depending on their priority department
and maximum number of students for each department
*/
func distributeApplicants(applicants []Applicant, maxNumber int, prior int, departments map[string][]Applicant) Departments {

	for i, _ := range departmentsNames {
		compareApplicants(&applicants, departmentsNames[i])
		for _, applicant := range applicants {
			if applicant.prior[prior-1] == departmentsNames[i] {
				departments[departmentsNames[i]] = append(departments[departmentsNames[i]], applicant)
			}
		}
	}

	for domain, applics := range departments {
		if len(applics) > maxNumber {
			departments[domain] = applics[:maxNumber]
			departments["0"] = append(departments["0"], applics[maxNumber:]...)
		}
	}

	return departments
}

/*Runs three waves of acceptance procedure*/
func processApplicants(applicants *[]Applicant, maxNumberOfStudents int) Departments {
	var departments = initDepartments()
	for i := 1; i <= 3; i++ {
		departments = distributeApplicants(*applicants, maxNumberOfStudents, i, departments)
		if len(departments["0"]) > 0 {
			*applicants = nil
			*applicants = departments["0"]
			departments["0"] = nil
		}
	}
	delete(departments, "0")

	return departments
}

/*
Outputs accepted students by departments:
- on stdout
- in files with appropriate names (biotech.txt, physics.txt etc.)
*/
func output(applicants []Applicant, departments Departments) error {

	res := 0.0
	for _, key := range departmentsNames {
		file, err := os.Create(strings.ToLower(key) + ".txt")
		if err != nil {
			return err
		}
		defer file.Close()
		applicants = departments[key]
		compareApplicants(&applicants, key)
		for _, applicant := range applicants {
			fmt.Fprint(file, applicant.fullname+" ")
			switch key {
			case departmentsNames[Biotech]:
				res = applicant.res[departmentsNames[Biotech]]
			case departmentsNames[Chemistry]:
				res = applicant.res[departmentsNames[Chemistry]]
			case departmentsNames[Engineering]:
				res = applicant.res[departmentsNames[Engineering]]
			case departmentsNames[Mathematics]:
				res = applicant.res[departmentsNames[Mathematics]]
			case departmentsNames[Physics]:
				res = applicant.res[departmentsNames[Physics]]
			}
			fmt.Fprintf(file, "%.1f\n", res)
		}
	}
	return nil
}

func main() {
	var maxNumberOfStudents int

	_, err := fmt.Scan(&maxNumberOfStudents)
	if err != nil || maxNumberOfStudents <= 0 {
		log.Fatal("Error: invalid number of students")
	}

	file, err := os.Open("applicants.txt")
	if err != nil {
		log.Fatal("Error: opening file")
	}
	defer file.Close()

	applicants, err := makeApplicants(file)
	if err != nil {
		log.Fatal("Error: creating applicants")
	}
	var departments = processApplicants(&applicants, maxNumberOfStudents)

	err = output(applicants, departments)
	if err != nil {
		log.Fatal("Error: failed to output to files")
	}
}
