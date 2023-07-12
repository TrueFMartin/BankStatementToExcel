package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type day struct {
	day   int
	month int
}
type expense struct {
	date        day
	description string
	isDebit     bool
	amount      float64
}
type balance struct {
	income     float64
	auto       float64
	housing    float64
	food       float64
	medical    float64
	education  float64
	recreation float64
	donations  float64
	other      float64
}

func removeElement(slice []string, index int) []string {
	//Shift elements to the left starting from the given index
	copy(slice[index:], slice[index+1:])

	//Truncate the last element to remove the duplicate
	return slice[:len(slice)-1]
}

func reader() (expenses []expense) {
	// open file
	f, err := os.Open("input2.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// do something with a line
		line := scanner.Text()
		lineSplit := strings.Split(line, ",")
		for i, l := range lineSplit {
			lineSplit[i] = strings.Trim(l, "\"")
		}
		if len(lineSplit) > 6 {
			lineSplit = removeElement(lineSplit, 4)
		}
		//Find and split date from string
		date := strings.Split(lineSplit[1], "-")
		d, err := strconv.Atoi(date[2])
		if err != nil {
			log.Fatal("Error reading day in date")
		}
		m, err := strconv.Atoi(date[1])
		if err != nil {
			log.Fatal("Error reading month in date")
		}

		//find if debit or credit
		isDebit := lineSplit[4] == "Debit"
		if !isDebit && lineSplit[4] != "Credit" {
			log.Fatal("Error in reading debit/credit")
		}

		//get amount of sale
		amount, err := strconv.ParseFloat(lineSplit[5], 32)
		if err != nil {
			log.Fatal("Error reading month in date")
		}

		data := expense{
			date:        day{day: d, month: m},
			description: lineSplit[2],
			isDebit:     isDebit,
			amount:      amount,
		}

		expenses = append(expenses, data)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func printout() string {
	fmt.Println("Enter a: auto, h: housing, f: food, m: medical, e: education, r: recreation, d: donations, o: other, q: QUIT")
	var input string
	i, err := fmt.Scanln(&input)
	if err != nil {
		log.Fatal("Failure in scan")
	}
	if i != 1 {
		fmt.Println("INVALID INPUT")
		return printout()
	}
	return input
}

func balanceTypeSwitcher(char string, expense expense, balancesMap map[day]balance) {
	if startDay.day == 0 {
		startDay = expense.date
	}
	lastDay = expense.date
	switch char {
	case "a":
		tempExpense := balancesMap[expense.date]
		tempExpense.auto += expense.amount
		balancesMap[expense.date] = tempExpense
	case "h":
		tempExpense := balancesMap[expense.date]
		tempExpense.housing += expense.amount
		balancesMap[expense.date] = tempExpense
	case "f":
		tempExpense := balancesMap[expense.date]
		tempExpense.food += expense.amount
		balancesMap[expense.date] = tempExpense
	case "m":
		tempExpense := balancesMap[expense.date]
		tempExpense.medical += expense.amount
		balancesMap[expense.date] = tempExpense
	case "e":
		tempExpense := balancesMap[expense.date]
		tempExpense.education += expense.amount
		balancesMap[expense.date] = tempExpense
	case "r":
		tempExpense := balancesMap[expense.date]
		tempExpense.recreation += expense.amount
		balancesMap[expense.date] = tempExpense
	case "d":
		tempExpense := balancesMap[expense.date]
		tempExpense.donations += expense.amount
		balancesMap[expense.date] = tempExpense
	case "o":
		tempExpense := balancesMap[expense.date]
		tempExpense.other += expense.amount
		balancesMap[expense.date] = tempExpense
	}
}

func readJson() map[string]string {
	file, err := os.ReadFile("history.json")
	if err != nil {
		return nil
	}
	historyMap := make(map[string]string)
	err = json.Unmarshal(file, &historyMap)
	if err != nil {
		return nil
	}
	return historyMap
}

func menu(expenses []expense) map[day]balance {
	balancesMap := make(map[day]balance)
	//Stores history of past identifications of expense types
	historyMap := readJson()
	if historyMap == nil {
		historyMap = make(map[string]string)
	}

	fmt.Println("Enter transaction type to sort each expense/income.")
	//Run through each expense from bank statement
	for _, expense := range expenses {
		//If expense is a credit, automatically add it to that date's 'income'
		if !expense.isDebit {
			if _, ok := balancesMap[expense.date]; !ok {
				tempBalance := balance{income: expense.amount}
				balancesMap[expense.date] = tempBalance
			} else {
				tempExpense := balancesMap[expense.date]
				tempExpense.income += expense.amount
				balancesMap[expense.date] = tempExpense
			}
			continue
		}
		//Based on first 15 characters of description, see if it's been assigned before
		if len(expense.description) > 44 {
			if v, ok := historyMap[expense.description[31:45]]; ok {
				balanceTypeSwitcher(v, expense, balancesMap)
				fmt.Println("Presorted from past input. Category: " + v + "---For " + expense.description + "\n")
				continue
			}
		} else { //Description is too short, use length instead
			if v, ok := historyMap[expense.description[:len(expense.description)]]; ok {
				balanceTypeSwitcher(v, expense, balancesMap)
				fmt.Println("Presorted from past input. Category: " + v + "---For " + expense.description + "\n")
				continue
			}
		}
		//Get user response based on description, sort it, tally it
		fmt.Println("Description: " + expense.description)
		var response = printout()
		if response == "q" {
			break
		}
		if len(expense.description) > 44 {
			historyMap[expense.description[31:45]] = response
		} else { //Description is too short, save in map from start to end for it's key
			historyMap[expense.description[:len(expense.description)]] = response
		}
		balanceTypeSwitcher(response, expense, balancesMap)
	}
	j, err := json.Marshal(historyMap)
	if err == nil {
		err := os.WriteFile("history.json", j, 0777)
		if err != nil {
			fmt.Println("ERROR WRITING JSON")
		}
	}
	return balancesMap
}

func dayToString(d day) string {
	return fmt.Sprintf("%v/%v/23", d.month, d.day)
}

// Returns a string of format "B3" for input x=1, y=4. Treat as zero indexed excel sheet.
func cordToString(x, y int) string {
	return fmt.Sprintf("%c%v", 'A'+x, y+1)
}

func daysInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
func (d day) less(d2 day) bool {
	return d.month < d2.month || (d.month == d2.month && d.day < d2.day)
}
func (d day) lessOrEqual(d2 day) bool {
	return d.month < d2.month || (d.month == d2.month && d.day <= d2.day)
}
func (d *day) increment() {
	d.day++
	if d.day > daysInMonth(time.Month(d.month), 2023) {
		d.month++
		d.day = 1
	}
}
func fileWriter(balances map[day]balance) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	cellMap := make(map[day]int)
	y := 0
	currentDay := startDay
	for currentDay.lessOrEqual(lastDay) {
		f.SetCellValue("Sheet1", cordToString(0, y), dayToString(currentDay))
		cellMap[currentDay] = y
		currentDay.increment()
		y++

	}
	for k, v := range balances {
		y = cellMap[k]
		x := 1
		f.SetCellValue("Sheet1", cordToString(x, y), v.income)
		x++
		f.SetCellValue("Sheet1", cordToString(x, y), v.auto)
		x++
		f.SetCellValue("Sheet1", cordToString(x, y), v.housing)
		x++
		f.SetCellValue("Sheet1", cordToString(x, y), v.food)
		x++
		f.SetCellValue("Sheet1", cordToString(x, y), v.medical)
		x++
		f.SetCellValue("Sheet1", cordToString(x, y), v.recreation)
		x++
		f.SetCellValue("Sheet1", cordToString(x, y), v.donations)
		x++
		f.SetCellValue("Sheet1", cordToString(x, y), v.other)
	}
	f.SaveAs("text.xlsx")
}

var startDay = day{0, 0}
var lastDay = day{0, 0}

func main() {
	expenses := reader()
	balances := menu(expenses)
	fileWriter(balances)
}
