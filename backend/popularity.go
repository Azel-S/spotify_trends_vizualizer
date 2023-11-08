package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func findPopularity(w http.ResponseWriter, r *http.Request) {
	sql := "SELECT track_id,release_date, popularity FROM Track GROUP BY EXTRACT(YEAR FROM release_date) ORDER BY popularity DESC, EXTRACT(YEAR FROM release_date) DESC"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer func() {
		cancel()
	}()

	if r.Method != http.MethodGet {
		writeError(w, "This requires get")
		return
	}
	stmt, err := conn.PrepareContext(ctx, sql)
	if err != nil {
		panic(err)
	}
	if t, ok := r.URL.Query()["title"]; len(t) == 1 && ok {
		res, err := stmt.QueryContext(ctx, t[0])
		if err != nil {
			panic(err)
		}
		popularityReturn := []PopularityReturn{}
		for res.Next() {
			var track_id string
			var release_date string
			var popularity int
			err := res.Scan(&track_id, &release_date, &popularity)
			if err != nil {
				panic(err)
			}
			popularityReturn = append(popularityReturn, PopularityReturn{
				Track_id:     track_id,
				Release_date: release_date,
				Popularity:   popularity,
			})

		}
		msg := map[string]any{
			"result": popularityReturn,
		}
		b, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, string(b))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		writeError(w, "There must be a title query")
	}
}

type PopularityReturn struct {
	Track_id     string `json:"track_id"`
	Release_date string `json:"release_date"`
	Popularity   int    `json:"popularity"`
}