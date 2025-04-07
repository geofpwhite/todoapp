package records

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/goccy/go-json"
)

type Record struct {
	RecordID  int    `json:"record_id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

type RecordHandler struct {
	mut sync.Mutex
}

func NewRecordHandler() *RecordHandler {
	return &RecordHandler{
		mut: sync.Mutex{},
	}
}

func (rh *RecordHandler) NumberOfRecords() int {
	rh.mut.Lock()
	defer rh.mut.Unlock()
	f, err := os.ReadFile("num_records.txt")
	if err != nil {
		os.Create("num_records.txt")
		os.WriteFile("num_records.txt", []byte("0"), os.FileMode(os.O_RDWR))
		return 0
	}
	str := string(f)
	num, _ := strconv.Atoi(str)
	return num
}

func (rh *RecordHandler) AddRecord(task string) {
	rh.mut.Lock()
	fmt.Println("Adding record")
	record := Record{
		Task: task,
	}
	// Implementation to add a record
	// This could involve inserting the record into a database or an in-memory data structure.
	rh.mut.Unlock()
	record.RecordID = rh.NumberOfRecords()
	rh.mut.Lock()
	defer rh.mut.Unlock()
	records := make([]Record, record.RecordID+10)
	fmt.Println(record.RecordID, len(records))

	f, err := os.ReadFile("active_records.json")
	if err != nil {
		panic("ahh")
	}
	err = json.Unmarshal(f, &records)
	if err != nil {
		panic("ahh")
	}
	fmt.Println(records, len(records))
	records = append(records, record)

	bytes, err := json.Marshal(records)
	if err != nil {
		panic("ahh")
	}
	if err = os.WriteFile("active_records.json", bytes, os.FileMode(os.O_RDWR)); err != nil {
		fmt.Println(err)
		panic("ahh")
	}

	err = os.WriteFile("num_records.txt", []byte(strconv.Itoa(record.RecordID+1)), 0644)
	if err != nil {
		panic("ahh")
	}

}

func (rh *RecordHandler) MarkRecordAsCompleted(recordID int) {
	rh.mut.Lock()
	defer rh.mut.Unlock()
	file, err := os.ReadFile("active_records.json")
	if err != nil {
		os.Create("active_records.json")
		os.WriteFile("active_records.json", []byte("[]"), os.FileMode(os.O_RDWR))
		return
	}
	records := make([]Record, 100)
	err = json.Unmarshal(file, &records)
	if err != nil {
		panic("ah")
	}
	for i, record := range records {
		if record.RecordID == recordID {
			r := record
			records = append(records[:i], records[i+1:]...)
			f, err := os.ReadFile("completed_records.json")
			if err != nil {
				os.Create("completed_records.json")
				os.WriteFile("completed_records.json", []byte("[]"), os.FileMode(os.O_RDWR))
				f = []byte("[]")
			}
			completed_records := make([]Record, 1000)
			err = json.Unmarshal(f, &completed_records)
			if err != nil {
				fmt.Println(err, completed_records)
				panic("ahh")
			}
			fmt.Println(err, completed_records)
			r.Completed = true
			completed_records = append(completed_records, r)
			jsonCompleted, err := json.Marshal(completed_records)
			fmt.Println("completed_records", completed_records, string(jsonCompleted))
			if err != nil {
				panic("ahh")
			}
			err = os.WriteFile("completed_records.json", jsonCompleted, os.FileMode(os.O_RDWR))
			fmt.Println(recordID)
			if err != nil {
				panic("ahh")
			}
			break
		}
	}
	b, err := json.Marshal(records)
	if err != nil {
		panic("ahh")
	}
	os.WriteFile("active_records.json", b, 0644)
}

func (rh *RecordHandler) GetActiveRecords() []Record {
	rh.mut.Lock()
	defer rh.mut.Unlock()
	records := make([]Record, 100)
	file, err := os.ReadFile("active_records.json")
	if err != nil {
		os.Create("active_records.json")
		os.WriteFile("active_records.json", []byte("[]"), os.FileMode(os.O_RDWR))
		file = []byte("[]")
	}
	err = json.Unmarshal(file, &records)
	if err != nil {
		panic("ahh")
	}
	fmt.Println(records)

	return records
}

func (rh *RecordHandler) GetCompletedRecords() []Record {
	rh.mut.Lock()
	defer rh.mut.Unlock()
	records := make([]Record, 100)
	file, err := os.ReadFile("completed_records.json")
	if err != nil {
		os.Create("completed_records.json")
		os.WriteFile("completed_records.json", []byte("[]"), os.FileMode(os.O_RDWR))
		file = []byte("[]")
	}
	err = json.Unmarshal(file, &records)
	if err != nil {
		panic("ahh")
	}
	fmt.Println(records)

	return records
}
