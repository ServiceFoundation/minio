/*
 * Minio Cloud Storage, (C) 2019 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/minio/minio/pkg/s3select/sql"
	"github.com/tidwall/sjson"
)

// Record - is CSV record.
type Record struct {
	columnNames  []string
	csvRecord    []string
	nameIndexMap map[string]int64
}

// Get - gets the value for a column name.
func (r *Record) Get(name string) (*sql.Value, error) {
	index, found := r.nameIndexMap[name]
	if !found {
		return nil, fmt.Errorf("column %v not found", name)
	}

	if index >= int64(len(r.csvRecord)) {
		// No value found for column 'name', hence return empty string for compatibility.
		return sql.NewString(""), nil
	}

	return sql.NewString(r.csvRecord[index]), nil
}

// Set - sets the value for a column name.
func (r *Record) Set(name string, value *sql.Value) error {
	r.columnNames = append(r.columnNames, name)
	r.csvRecord = append(r.csvRecord, value.CSVString())
	return nil
}

// MarshalCSV - encodes to CSV data.
func (r *Record) MarshalCSV(fieldDelimiter rune) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	w.Comma = fieldDelimiter
	if err := w.Write(r.csvRecord); err != nil {
		return nil, err
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	data := buf.Bytes()
	return data[:len(data)-1], nil
}

// MarshalJSON - encodes to JSON data.
func (r *Record) MarshalJSON() ([]byte, error) {
	data := "{}"

	var err error
	for i := len(r.columnNames) - 1; i >= 0; i-- {
		if i >= len(r.csvRecord) {
			continue
		}

		if data, err = sjson.Set(data, r.columnNames[i], r.csvRecord[i]); err != nil {
			return nil, err
		}
	}

	return []byte(data), nil
}

// NewRecord - creates new CSV record.
func NewRecord() *Record {
	return &Record{}
}
