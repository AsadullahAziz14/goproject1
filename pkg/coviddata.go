package pkg

import (
	"encoding/csv"
	"io"
	"os"
)

type data struct {
	Date       string
	Positive   string
	Tests      string
	Expired    string
	Admitted   string
	Discharged string
	Region     string
}

func Load(path string) []data {
	table := make([]data, 0)
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		d := data{
			Positive : row[0],
			Tests : row[1],
			Date : row[2],
			Discharged : row[3],
			Expired : row[4],
			Admitted : row[5],
			Region : row[6],

		}
		table = append(table, d)
	}
	return table
}

func Find(table []data, filter string) []data {
	if filter == "" || filter == "*" {
		return table
	}
	result := make([]data, 0)
	//filter = strings.ToUpper(filter)
	for _, record := range table {
		if record.Date == filter || record.Region == filter  {
			result = append(result, record)
		}
	}
	return result
}
