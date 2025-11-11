package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Marks []int  `json:"marks"`
}

type Journal struct {
	Students []Student `json:"students"`
}

func (j *Journal) SaveToFile(filename string) {
	data, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		fmt.Println("Ошибка сериализации:", err)
		return
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
	}
}

func LoadFromFile(filename string) (*Journal, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл %s не существует", filename)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}
	var j Journal
	err = json.Unmarshal(data, &j)
	if err != nil {
		return nil, fmt.Errorf("ошибка десериализации: %w", err)
	}
	return &j, nil
}

func (j *Journal) findStudentIndex(id int) int {
	for i, s := range j.Students {
		if s.ID == id {
			return i
		}
	}
	return -1
}

func (j *Journal) findStudentIndexByName(name string) int {
	for i, s := range j.Students {
		if s.Name == name {
			return i
		}
	}
	return -1
}

func (j *Journal) addStudent(name string) {
	if j.findStudentIndexByName(name) != -1 {
		fmt.Println("Ученик уже существует")
		return
	}
	newID := 1
	if len(j.Students) > 0 {
		maxID := j.Students[0].ID
		for _, s := range j.Students {
			if s.ID > maxID {
				maxID = s.ID
			}
		}
		newID = maxID + 1
	}
	j.Students = append(j.Students, Student{
		ID:    newID,
		Name:  name,
		Marks: []int{},
	})
	fmt.Println("Ученик добавлен")
}

func (j *Journal) removeStudent(id int) bool {
	idx := j.findStudentIndex(id)
	if idx == -1 {
		return false
	}
	j.Students = append(j.Students[:idx], j.Students[idx+1:]...)
	return true
}

func (j *Journal) addMark(id, mark int) bool {
	idx := j.findStudentIndex(id)
	if idx == -1 {
		return false
	}
	if mark < 1 || mark > 5 {
		fmt.Println("Оценка должна быть от 1 до 5")
		return false
	}
	j.Students[idx].Marks = append(j.Students[idx].Marks, mark)
	return true
}

func (j *Journal) averageMark(id int) (float64, bool) {
	idx := j.findStudentIndex(id)
	if idx == -1 || len(j.Students[idx].Marks) == 0 {
		return 0, false
	}
	sum := 0
	for _, m := range j.Students[idx].Marks {
		sum += m
	}
	return float64(sum) / float64(len(j.Students[idx].Marks)), true
}

func (j *Journal) printStudents() {
	if len(j.Students) == 0 {
		fmt.Println("Учеников нет")
		return
	}
	sort.Slice(j.Students, func(i, k int) bool {
		return j.Students[i].ID < j.Students[k].ID
	})
	for _, s := range j.Students {
		fmt.Printf("%d. %s\n", s.ID, s.Name)
	}
}

func getName() string {
	var name, surname, patronymic string
	_, err := fmt.Scan(&name, &surname, &patronymic)
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return ""
	}
	return name + " " + surname + " " + patronymic
}

func main() {
	var journal Journal
	var check int
	var subCheck int
	filename := "journal.json"
	fmt.Println("Добро пожаловать в журнал")
	for {
		fmt.Println(`
Выберите действие:
1. Создать ученика
2. Удалить ученика
3. Ввести оценку
4. Вывести средний балл
5. Сохранить в файл
6. Загрузить из файла
7. Выход`)
		fmt.Scan(&check)
		switch check {
		case 1:
			fmt.Print("Введите ФИО: ")
			name := getName()
			if name != "" {
				journal.addStudent(name)
			}
		case 2:
			journal.printStudents()
			if len(journal.Students) == 0 {
				continue
			}
			fmt.Print("Выберите ID ученика для удаления: ")
			var id int
			fmt.Scan(&id)
			if journal.removeStudent(id) {
				fmt.Println("Ученик удален")
			} else {
				fmt.Println("Такого ученика не существует")
			}
		case 3:
			journal.printStudents()
			if len(journal.Students) == 0 {
				continue
			}
			fmt.Print("Введите ID ученика: ")
			var id int
			fmt.Scan(&id)
			if journal.findStudentIndex(id) == -1 {
				fmt.Println("Такого ученика нет")
				continue
			}
			for {
				fmt.Print("Введите оценку (1–5): ")
				var mark int
				if _, err := fmt.Scan(&mark); err != nil {
					fmt.Println("Неверный ввод")
					continue
				}
				if mark < 1 || mark > 5 {
					fmt.Println("Оценка должна быть от 1 до 5")
					continue
				}
				journal.addMark(id, mark)
				fmt.Println("Оценка добавлена")
				fmt.Print("Добавить ещё? (1 — да, иное — нет): ")
				var cont int
				fmt.Scan(&cont)
				if cont != 1 {
					break
				}
			}
		case 4:
			if len(journal.Students) == 0 {
				fmt.Println("Учеников нет")
				continue
			}
			fmt.Println(`
Выберите действие:
1. Средний балл всех
2. Средний балл ≥ 4.00
3. Средний балл ≤ 4.00
4. Средний балл конкретного ученика`)
			fmt.Scan(&subCheck)
			switch subCheck {
			case 1:
				for _, s := range journal.Students {
					if avg, ok := journal.averageMark(s.ID); ok {
						fmt.Printf("%s — %.2f\n", s.Name, avg)
					} else {
						fmt.Printf("%s — нет оценок\n", s.Name)
					}
				}
			case 2, 3:
				for _, s := range journal.Students {
					if avg, ok := journal.averageMark(s.ID); ok {
						if (subCheck == 2 && avg >= 4.0) || (subCheck == 3 && avg <= 4.0) {
							fmt.Printf("%s — %.2f\n", s.Name, avg)
						}
					}
				}
			case 4:
				journal.printStudents()
				fmt.Print("Введите ID ученика: ")
				var id int
				fmt.Scan(&id)
				if avg, ok := journal.averageMark(id); ok {
					name := journal.Students[journal.findStudentIndex(id)].Name
					fmt.Printf("%s — %.2f\n", name, avg)
				} else {
					fmt.Println("Ученика не существует или у него нет оценок")
				}
			default:
				fmt.Println("Неверный выбор")
			}
		case 5:
			journal.SaveToFile(filename)
			fmt.Println("Сохранено в", filename)
		case 6:
			loaded, err := LoadFromFile(filename)
			if err != nil {
				fmt.Println("Ошибка загрузки:", err)
			} else {
				journal = *loaded
				fmt.Println("Загружено из", filename)
			}
		case 7:
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Такого действия не существует")
		}
	}
}
