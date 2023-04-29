package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)


type task struct {
	ID int
	Description string
	URL string
	Completed bool
	Created time.Time
}


type list struct {
	items []task
}

func (l *list) add(taskDescription string) {
	newID := len(l.items)
	todo := task{ID: newID, Description: taskDescription, Completed: false, Created: time.Now()}
	l.items = append(l.items, todo)
}

func (l *list) delete(id int) error {
	for k, v := range(l.items) {
		if v.ID == id {
			l.items = append(l.items[:k], l.items[k+1:]...)
		} else {
			return fmt.Errorf("ID not found in task list")
		}
	}
	return nil
}

func (l *list) save(p string) error {
	listJSON, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("could not create JSON: %w", err)
	}

	err = os.WriteFile(p, listJSON, 0644)
	if err != nil {
		return fmt.Errorf("Could not dave to datafile: %w", err)
	}

	return nil
}

func (l *list) load(p string) error {
	dataJSON, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("Could not read from datafile: %w", err)
	}
	
	err = json.Unmarshal(dataJSON, &l.items)
	if err != nil {
		return fmt.Errorf("Could not parse JSON: %w", err)
	}

	return nil
}

func (l *list) complete(id int) error {
	for _, v := range(l.items) {
		if v.ID == id {
			if v.Completed == false {
				v.Completed = true
				return nil
			} else {
				v.Completed = false
				return nil
			}
		}
	}
	return fmt.Errorf("Task ID not found")
}

func buildDataPath() string {

	return filepath.Join(getHomeDir(), ".local", "share", "tasker")
}

func getHomeDir() string {
	homedir, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}
	
	return homedir
}

func makeDir(p string) error {

	_, err := os.Stat(p)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Directory already exists")
		}
	} else {
		err := os.Mkdir(p, 755)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func makeDataFile(p string) (string, error) {
	dataFilePath := filepath.Join(p, "tasker.txt")
	_, err := os.Stat(dataFilePath)

	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(dataFilePath)
			if err != nil {
				return "", fmt.Errorf("File create error: %w", err)
			}
			 defer f.Close()

		} else {
			return "", fmt.Errorf("Unexpected error locating datafile: %w", err)
		}
	}
	return dataFilePath, nil
}

func main() {
	flag.String("add", "", "add a new task")
	flag.Int("del", -1, "delete a task")
	flag.Int("done", -1, "complete a task")

	flag.Parse()

	p := buildDataPath()
	
	err := makeDir(p)
	if err != nil {
		log.Println(err)
	}
	
	dataFilePath, err := makeDataFile(p)
	if err != nil {
		log.Fatal(err)
	}
	
	data, err:= os.ReadFile(dataFilePath)
	if err != nil {
		log.Fatal(err)
	}

	l := &list{}

	if len(data) != 0 {
		json.Unmarshal(data, l)
	}

	fmt.Println(data)



	fmt.Println("hello")
}