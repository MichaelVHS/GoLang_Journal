package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func SaveToFile(students map[int]string, journal map[int][]int) {
	jsonData, err := json.MarshalIndent(students, "", "  ")
	if err != nil {
		fmt.Println("Ошибка сериализации:", err)
		return
	}

	err = os.WriteFile("students.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
	}

	jsonData, err = json.MarshalIndent(journal, "", "  ")
	if err != nil {
		fmt.Println("Ошибка сериализации:", err)
		return
	}

	err = os.WriteFile("journal.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
	}
}

func LoadFromFile() (map[int]string, map[int][]int, error) {
	if _, err := os.Stat("students.json"); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("файл %s не существует", "students.json")
	}

	data, err := os.ReadFile("students.json")
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	students := map[int]string{}
	err = json.Unmarshal(data, &students)

	if _, err := os.Stat("journal.json"); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("файл %s не существует", "journal.json")
	}

	data, err = os.ReadFile("journal.json")
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	journal := map[int][]int{}
	err = json.Unmarshal(data, &journal)

	return students, journal, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getName() string {
	var name, surname, otches string
	_, err := fmt.Scan(&name, &surname, &otches)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return ""
	}
	fullName := name + " " + surname + " " + otches
	return fullName
}

func SortKeys(Map map[int]string) []int {
	keys := make([]int, 0, len(Map))
	for id := range Map {
		keys = append(keys, id)
	}
	sort.Ints(keys)
	return keys
}

func printStud(students map[int]string) {
	keys := SortKeys(students)
	for _, id := range keys {
		fmt.Printf("%d. %s\n", id, students[id])
	}
}

func findKeyByValue(data map[int]string, value string) int {
	for key, val := range data {
		if val == value {
			return key
		}
	}
	return 0
}

func averMark(students map[int]string, journal map[int][]int, name string) float64 {
	index := findKeyByValue(students, name)
	var itog int
	var aver float64
	for _, val := range journal[index] {
		itog += val
	}
	aver = float64(itog) / float64(len(journal[index]))
	return aver
}

func main() {
	var err error
	var check, mark int
	var name string
	var marks []int
	var students map[int]string = map[int]string{}
	var journal map[int][]int = map[int][]int{}
	fmt.Println("Добро пожаловать в журнал")
	for {
		fmt.Println(`
Выберите действие:
1. Создать ученика
2. Удалить ученика
3. Ввести оценку
4. Вывести средний балл
5. Сохранить в файл
6. Получить из файла
7.Выход`)
		fmt.Scan(&check)
		if check < 1 || check > 9 {
			fmt.Println("Такого действия не существует")
			continue
		} else if check == 1 {
			fmt.Print("Введите ФИО: ")
			name = getName()
			if len(students) == 0 {
				index := 1
				students[index] = name
				fmt.Println("Ученик добавлен")
			} else {
				isStud := false
				for _, val := range students {
					if val == name {
						isStud = true
					}
				}
				if isStud {
					fmt.Println("Ученик уже существует")
				} else {
					lastKey := 0
					keys := SortKeys(students)
					for _, key := range keys {
						lastKey = key
					}
					index := lastKey + 1
					students[index] = name
					fmt.Println("Ученик добавлен")
				}
			}
			fmt.Println(students)
		} else if check == 2 {
			printStud(students)
			if len(students) == 0 {
				fmt.Println("Учеников нет")
			} else {
				index := -1
				fmt.Print("Выберите ученика: ")
				fmt.Scan(&index)
				if index != -1 && students[index] != "" {
					delete(students, index)
					if journal[index] != nil {
						delete(journal, index)
					}
					fmt.Println("Ученик удален")
				} else {
					fmt.Println("Такого ученика не существует")
				}
			}
		} else if check == 3 {
			printStud(students)
			if len(students) == 0 {
				fmt.Println("Учеников нет")
			} else {
				fmt.Print("Введите номер ученика: ")
				fmt.Scan(&check)
				index := -1
				for i, val := range students {
					if val == students[check] {
						index = i
					}
				}
				if index != -1 {
					marks = journal[index]
					for {
						fmt.Print("Введите оценку(1-5): ")
						_, err := fmt.Scan(&mark)
						if err != nil {
							fmt.Println("Неверный тип данных, повторите", err)
							continue
						}
						if mark < 1 || mark > 5 {
							fmt.Println("Такой оценки нет")
						} else {
							marks = append(marks, mark)
							fmt.Println("Оценка добавлена")
							fmt.Println("Добавить еще одну оценку?(1 - Да/2 - Нет)")
							fmt.Scan(&check)
							if check < 1 || check > 2 {
								fmt.Println("Такого действия не существует")
							} else if check == 2 {
								break
							}
						}
					}
					journal[index] = marks
					marks = nil
				} else {
					fmt.Println("Такого ученика не существует")
				}
			}
		} else if check == 4 {
			if len(students) == 0 {
				fmt.Println("Учеников нет")
			} else {
				fmt.Println(`
Выберите действие:
1. Вывести средний балл
2. Вывести средний балл выше 4
3. Ввести средний балл ниже 4
4. Вывести средний балл студента`)
				fmt.Scan(&check)
				if check < 1 || check > 4 {
					fmt.Println("Такого действия не существует")
				} else if check == 1 {
					for _, val := range students {
						if len(journal[findKeyByValue(students, val)]) == 0 {
							fmt.Println("У ученика", val, "нет оценок")
						} else {
							aver := fmt.Sprintf("%.2f", averMark(students, journal, val))
							fmt.Println(val, "-", aver)
						}
					}
				} else if check == 2 {
					for _, val := range students {
						aver := fmt.Sprintf("%.2f", averMark(students, journal, val))
						av := averMark(students, journal, val)
						if av >= 4 {
							fmt.Println(val, "-", aver)
						}
					}
				} else if check == 3 {
					for _, val := range students {
						aver := fmt.Sprintf("%.2f", averMark(students, journal, val))
						av := averMark(students, journal, val)
						if av <= 4 {
							fmt.Println(val, "-", aver)
						}
					}
				} else if check == 4 {
					printStud(students)
					fmt.Print("Введите номер ученика: ")
					fmt.Scan(&check)
					index := -1
					for i, val := range students {
						if val == students[check] {
							index = i
						}
					}
					if index != -1 {
						if len(journal[findKeyByValue(students, name)]) == 0 {
							fmt.Println("У ученика", name, "нет оценок")
						} else {
							aver := fmt.Sprintf("%.2f", averMark(students, journal, name))
							fmt.Println(name, "-", aver)
						}
					}
				} else {
					fmt.Println("Такого ученика не существует")
				}
			}
		} else if check == 5 {
			SaveToFile(students, journal)
		} else if check == 6 {
			students, journal, err = LoadFromFile()
			if err != nil {
				fmt.Println(err)
			}
		} else if check == 7 {
			fmt.Println("До свидания!")
			break
		} else if check == 9 {
			printStud(students)
			fmt.Println(journal)
		}
	}
}
