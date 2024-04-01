package main

import (
	"encoding/binary"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func uint64ToBytes(i uint64) []byte {
	n1 := make([]byte, 8)
	binary.BigEndian.PutUint64(n1, i)
	return n1
}

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page <= 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
