package main

import (
	"context"
	"errors"
	"fmt"
	"go-iss/internal"
	"go-iss/internal/config"
	handlers "go-iss/internal/rpc/handlers"
	"go-iss/internal/rpc/handlers/iss/iss_fetch_get"
	"go-iss/internal/rpc/handlers/iss/iss_last_get"
	"go-iss/internal/rpc/handlers/iss/iss_trend_get"
	"go-iss/internal/rpc/handlers/osdr/osdr_list_get"
	"go-iss/internal/rpc/handlers/osdr/osdr_sync_get"
	space_refresh_get "go-iss/internal/rpc/handlers/space/space_refresh_get"
	space_src_latest_get "go-iss/internal/rpc/handlers/space/space_src_latest_get"
	space_summary_get "go-iss/internal/rpc/handlers/space/space_summary_get"
	"go-iss/internal/rpc/routes"
	"go-iss/internal/usecase/fetch_and_store_iss"
	"go-iss/internal/usecase/fetch_and_store_osdr"
	"go-iss/internal/usecase/fetch_and_store_space"
	"go-iss/internal/usecase/get_iss_trend"
	"go-iss/internal/usecase/get_last_iss"
	"go-iss/internal/usecase/get_latest_space_cache"
	"go-iss/internal/usecase/get_osdr_list"
	"go-iss/internal/usecase/get_space_summary"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	ctx := context.Background()

	db, err := internal.NewDatabaseConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	issClient, err := internal.NewISSClient()
	if err != nil {
		log.Fatalf("Failed to create ISS client: %v", err)
	}
	defer issClient.Close()

	nasaClient, err := internal.NewNASAClient()
	if err != nil {
		log.Fatalf("Failed to create NASA client: %v", err)
	}
	defer nasaClient.Close()

	spacexClient, err := internal.NewSpaceXClient()
	if err != nil {
		log.Fatalf("Failed to create SpaceX client: %v", err)
	}
	defer spacexClient.Close()

	backgroundConfig, err := config.SetupBackgroundConfig(ctx, db, issClient, nasaClient, spacexClient)
	if err != nil {
		log.Fatalf("Failed to setup background config: %v", err)
	}

	startBackgroundTasks(ctx, backgroundConfig)

	e := echo.New()

	issRepo := backgroundConfig.ISSRepo
	osdrRepo := backgroundConfig.OSDRRepo
	cacheRepo := backgroundConfig.CacheRepo

	getLastIssService := get_last_iss.New(issRepo)
	fetchAndStoreIssService := fetch_and_store_iss.New(issRepo, issClient)
	getIssTrendService := get_iss_trend.New(issRepo)
	getOsdrListService := get_osdr_list.New(osdrRepo)
	fetchAndStoreOsdrService := fetch_and_store_osdr.New(osdrRepo, nasaClient)
	refreshSpaceService := fetch_and_store_space.New(cacheRepo, nasaClient, spacexClient)
	getLatestSpaceCacheService := get_latest_space_cache.New(cacheRepo)
	getSpaceSummaryService := get_space_summary.New(cacheRepo, issRepo, osdrRepo)

	// Создание handlers
	issLastGetHandler := iss_last_get.New(getLastIssService)
	issFetchGetHandler := iss_fetch_get.New(fetchAndStoreIssService)
	issTrendGetHandler := iss_trend_get.New(getIssTrendService)
	osdrListGetHandler := osdr_list_get.New(getOsdrListService)
	osdrSyncGetHandler := osdr_sync_get.New(fetchAndStoreOsdrService)
	spaceRefreshGetHandler := space_refresh_get.New(refreshSpaceService)
	spaceSrcLatestGetHandler := space_src_latest_get.New(getLatestSpaceCacheService)
	spaceSummaryGetHandler := space_summary_get.New(getSpaceSummaryService)

	h := routes.Handlers{
		HandleHealth:       handlers.HealthCheck,
		HandleLastISS:      issLastGetHandler.ISSLastGet,
		HandleTriggerISS:   issFetchGetHandler.ISSFetchGet,
		HandleISSTrend:     issTrendGetHandler.ISSTrendGet,
		HandleOSDRSync:     osdrSyncGetHandler.OSDRSyncGet,
		HandleOSDRList:     osdrListGetHandler.OSDRListGet,
		HandleSpaceLatest:  spaceSrcLatestGetHandler.SpaceSrcLatestGet,
		HandleSpaceRefresh: spaceRefreshGetHandler.SpaceRefreshGet,
		HandleSpaceSummary: spaceSummaryGetHandler.SpaceSummaryGet,
	}

	routes.InitRoutes(e, h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := ":" + port

	go func() {
		fmt.Printf("Starting evently-apico HTTP server on %s\n", addr)
		if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}

func startBackgroundTasks(ctx context.Context, cfg *config.BackgroundConfig) {
	fetchAndStoreIssService := fetch_and_store_iss.New(cfg.ISSRepo, cfg.ISSClient)
	fetchAndStoreOsdrService := fetch_and_store_osdr.New(cfg.OSDRRepo, cfg.NasaClient)
	refreshSpaceService := fetch_and_store_space.New(cfg.CacheRepo, cfg.NasaClient, cfg.SpaceXClient)

	go backgroundTask(ctx, "OSDR", time.Duration(cfg.RetryTimes.EveryOSDR)*time.Second, func() error {
		_, err := fetchAndStoreOsdrService.FetchAndStoreOSDR(ctx)
		return err
	})

	go backgroundTask(ctx, "ISS", time.Duration(cfg.RetryTimes.EveryISS)*time.Second, func() error {
		_, err := fetchAndStoreIssService.FetchAndStoreISS(ctx)
		return err
	})

	go backgroundTask(ctx, "APOD", time.Duration(cfg.RetryTimes.EveryAPOD)*time.Second, func() error {
		_, err := refreshSpaceService.RefreshSpace(ctx, []string{"apod"})
		return err
	})

	go backgroundTask(ctx, "NEO", time.Duration(cfg.RetryTimes.EveryNEO)*time.Second, func() error {
		_, err := refreshSpaceService.RefreshSpace(ctx, []string{"neo"})
		return err
	})

	go backgroundTask(ctx, "DONKI", time.Duration(cfg.RetryTimes.EveryDONKI)*time.Second, func() error {
		_, err := refreshSpaceService.RefreshSpace(ctx, []string{"flr", "cme"})
		return err
	})

	go backgroundTask(ctx, "SpaceX", time.Duration(cfg.RetryTimes.EverySpaceX)*time.Second, func() error {
		_, err := refreshSpaceService.RefreshSpace(ctx, []string{"spacex"})
		return err
	})
}

func backgroundTask(ctx context.Context, name string, interval time.Duration, task func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	if err := task(); err != nil {
		log.Printf("%s err: %v", name, err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := task(); err != nil {
				log.Printf("%s err: %v", name, err)
			}
		}
	}
}

// Handlers

//
//func handleHealth(c echo.Context) error {
//	return c.JSON(http.StatusOK, HealthResponse{
//		Status: "ok",
//		Now:    time.Now().UTC(),
//	})
//}
//
//func makeHandleLastISS(db *sqlx.DB) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		ctx := c.Request().Context()
//
//		query := `SELECT id, fetched_at, source_url, payload
//				  FROM iss_fetch_log
//				  ORDER BY id DESC LIMIT 1`
//
//		var id int64
//		var fetchedAt time.Time
//		var sourceURL string
//		var payload json.RawMessage
//
//		err := db.QueryRowContext(ctx, query).Scan(&id, &fetchedAt, &sourceURL, &payload)
//		if err == sql.ErrNoRows {
//			return c.JSON(http.StatusOK, map[string]string{"message": "no data"})
//		}
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
//		}
//
//		var payloadJSON interface{}
//		if err := json.Unmarshal(payload, &payloadJSON); err != nil {
//			payloadJSON = map[string]interface{}{}
//		}
//
//		return c.JSON(http.StatusOK, map[string]interface{}{
//			"id":         id,
//			"fetched_at": fetchedAt,
//			"source_url": sourceURL,
//			"payload":    payloadJSON,
//		})
//	}
//}

//
//func makeHandleTriggerISS(db *sqlx.DB, issClient *iss.Client) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		ctx := c.Request().Context()
//
//		if err := fetchAndStoreISS(ctx, db, issClient); err != nil {
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
//		}
//
//		// Повторно получить последние данные
//		query := `SELECT id, fetched_at, source_url, payload
//				  FROM iss_fetch_log
//				  ORDER BY id DESC LIMIT 1`
//
//		var id int64
//		var fetchedAt time.Time
//		var sourceURL string
//		var payload json.RawMessage
//
//		err := db.QueryRowContext(ctx, query).Scan(&id, &fetchedAt, &sourceURL, &payload)
//		if err == sql.ErrNoRows {
//			return c.JSON(http.StatusOK, map[string]string{"message": "no data"})
//		}
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
//		}
//
//		var payloadJSON interface{}
//		if err := json.Unmarshal(payload, &payloadJSON); err != nil {
//			payloadJSON = map[string]interface{}{}
//		}
//
//		return c.JSON(http.StatusOK, map[string]interface{}{
//			"id":         id,
//			"fetched_at": fetchedAt,
//			"source_url": sourceURL,
//			"payload":    payloadJSON,
//		})
//	}
//}
//
//func makeHandleISSTrend(db *sqlx.DB) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		ctx := c.Request().Context()
//
//		query := `SELECT fetched_at, payload FROM iss_fetch_log ORDER BY id DESC LIMIT 2`
//		rows, err := db.QueryContext(ctx, query)
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
//		}
//		defer rows.Close()
//
//		var records []struct {
//			FetchedAt time.Time
//			Payload   json.RawMessage
//		}
//
//		for rows.Next() {
//			var r struct {
//				FetchedAt time.Time
//				Payload   json.RawMessage
//			}
//			if err := rows.Scan(&r.FetchedAt, &r.Payload); err != nil {
//				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
//			}
//			records = append(records, r)
//		}
//
//		if len(records) < 2 {
//			return c.JSON(http.StatusOK, Trend{
//				Movement: false,
//				DeltaKm:  0.0,
//				DtSec:    0.0,
//			})
//		}
//
//		t2 := records[0].FetchedAt
//		t1 := records[1].FetchedAt
//
//		var p2, p1 map[string]interface{}
//		json.Unmarshal(records[0].Payload, &p2)
//		json.Unmarshal(records[1].Payload, &p1)
//
//		lat1 := num(p1, "latitude")
//		lon1 := num(p1, "longitude")
//		lat2 := num(p2, "latitude")
//		lon2 := num(p2, "longitude")
//		v2 := num(p2, "velocity")
//
//		var deltaKm float64
//		movement := false
//		if lat1 != nil && lon1 != nil && lat2 != nil && lon2 != nil {
//			deltaKm = haversineKm(*lat1, *lon1, *lat2, *lon2)
//			movement = deltaKm > 0.1
//		}
//
//		dtSec := t2.Sub(t1).Seconds()
//
//		trend := Trend{
//			Movement: movement,
//			DeltaKm:  deltaKm,
//			DtSec:    dtSec,
//			FromTime: &t1,
//			ToTime:   &t2,
//		}
//
//		if v2 != nil {
//			trend.VelocityKmh = v2
//		}
//		if lat1 != nil {
//			trend.FromLat = lat1
//		}
//		if lon1 != nil {
//			trend.FromLon = lon1
//		}
//		if lat2 != nil {
//			trend.ToLat = lat2
//		}
//		if lon2 != nil {
//			trend.ToLon = lon2
//		}
//
//		return c.JSON(http.StatusOK, trend)
//	}
//}
//
//func makeHandleOSDRSync(db *sqlx.DB, nasaClient *nasa.Client) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		ctx := c.Request().Context()
//
//		written, err := fetchAndStoreOSDR(ctx, db, nasaClient, nasaClient.GetOSDRURL())
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
//		}
//
//		return c.JSON(http.StatusOK, map[string]int{"written": written})
//	}
//}
//
//func makeHandleOSDRList(db *sqlx.DB) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		ctx := c.Request().Context()
//
//		limitStr := getEnv("OSDR_LIST_LIMIT", "20")
//		limit, _ := strconv.Atoi(limitStr)
//		if limit <= 0 {
//			limit = 20
//		}
//
//		query := squirrel.Select("id", "dataset_id", "title", "status", "updated_at", "inserted_at", "raw").
//			From("osdr_items").
//			OrderBy("inserted_at DESC").
//			Limit(uint64(limit))
//
//		sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
//		}
//
//		rows, err := db.QueryContext(ctx, sql, args...)
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
//		}
//		defer rows.Close()
//
//		var items []map[string]interface{}
//		for rows.Next() {
//			var id int64
//			var datasetID, title, status *string
//			var updatedAt *time.Time
//			var insertedAt time.Time
//			var raw json.RawMessage
//
//			if err := rows.Scan(&id, &datasetID, &title, &status, &updatedAt, &insertedAt, &raw); err != nil {
//				continue
//			}
//
//			var rawJSON interface{}
//			json.Unmarshal(raw, &rawJSON)
//
//			item := map[string]interface{}{
//				"id":          id,
//				"inserted_at": insertedAt,
//				"raw":         rawJSON,
//			}
//			if datasetID != nil {
//				item["dataset_id"] = *datasetID
//			}
//			if title != nil {
//				item["title"] = *title
//			}
//			if status != nil {
//				item["status"] = *status
//			}
//			if updatedAt != nil {
//				item["updated_at"] = *updatedAt
//			}
//			items = append(items, item)
//		}
//
//		return c.JSON(http.StatusOK, map[string]interface{}{"items": items})
//	}
//}
//
//func makeHandleSpaceLatest(db *sqlx.DB) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		ctx := c.Request().Context()
//		src := c.Param("src")
//
//		query := `SELECT fetched_at, payload FROM space_cache
//				  WHERE source = $1 ORDER BY id DESC LIMIT 1`
//
//		var fetchedAt time.Time
//		var payload json.RawMessage
//
//		err := db.QueryRowContext(ctx, query, src).Scan(&fetchedAt, &payload)
//		if err == sql.ErrNoRows {
//			return c.JSON(http.StatusOK, map[string]interface{}{
//				"source":  src,
//				"message": "no data",
//			})
//		}
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
//		}
//
//		var payloadJSON interface{}
//		json.Unmarshal(payload, &payloadJSON)
//
//		return c.JSON(http.StatusOK, map[string]interface{}{
//			"source":     src,
//			"fetched_at": fetchedAt,
//			"payload":    payloadJSON,
//		})
//	}
//}
//
//func makeHandleSpaceRefresh(db *sqlx.DB, nasaClient *nasa.Client, spacexClient *spacex.SpaceXClient) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		ctx := c.Request().Context()
//
//		srcParam := c.QueryParam("src")
//		if srcParam == "" {
//			srcParam = "apod,neo,flr,cme,spacex"
//		}
//
//		sources := strings.Split(srcParam, ",")
//		var done []string
//
//		for _, s := range sources {
//			s = strings.TrimSpace(strings.ToLower(s))
//			var err error
//			switch s {
//			case "apod":
//				err = fetchAPOD(ctx, db, nasaClient)
//				done = append(done, "apod")
//			case "neo":
//				err = fetchNEOFeed(ctx, db, nasaClient)
//				done = append(done, "neo")
//			case "flr":
//				err = fetchDONKIFLR(ctx, db, nasaClient)
//				done = append(done, "flr")
//			case "cme":
//				err = fetchDONKICME(ctx, db, nasaClient)
//				done = append(done, "cme")
//			case "spacex":
//				err = fetchSpaceXNext(ctx, db, spacexClient)
//				done = append(done, "spacex")
//			}
//			if err != nil {
//				log.Printf("Error refreshing %s: %v", s, err)
//			}
//		}
//
//		return c.JSON(http.StatusOK, map[string]interface{}{"refreshed": done})
//	}
//}
//
//func makeHandleSpaceSummary(db *sqlx.DB) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		ctx := c.Request().Context()
//
//		result := make(map[string]interface{})
//
//		// Кэш источников
//		for _, src := range []string{"apod", "neo", "flr", "cme", "spacex"} {
//			cache := latestFromCache(ctx, db, src)
//			result[src] = cache
//		}
//
//		// ISS последний
//		query := `SELECT fetched_at, payload FROM iss_fetch_log ORDER BY id DESC LIMIT 1`
//		var fetchedAt time.Time
//		var payload json.RawMessage
//		err := db.QueryRowContext(ctx, query).Scan(&fetchedAt, &payload)
//		if err == nil {
//			var payloadJSON interface{}
//			json.Unmarshal(payload, &payloadJSON)
//			result["iss"] = map[string]interface{}{
//				"at":      fetchedAt,
//				"payload": payloadJSON,
//			}
//		} else {
//			result["iss"] = map[string]interface{}{}
//		}
//
//		// OSDR count
//		var osdrCount int64
//		db.QueryRowContext(ctx, "SELECT count(*) FROM osdr_items").Scan(&osdrCount)
//		result["osdr_count"] = osdrCount
//
//		return c.JSON(http.StatusOK, result)
//	}
//}
//
//// Фетчеры
//
//func latestFromCache(ctx context.Context, db *sqlx.DB, src string) map[string]interface{} {
//	query := `SELECT fetched_at, payload FROM space_cache WHERE source=$1 ORDER BY id DESC LIMIT 1`
//	var fetchedAt time.Time
//	var payload json.RawMessage
//
//	err := db.QueryRowContext(ctx, query, src).Scan(&fetchedAt, &payload)
//	if err != nil {
//		return map[string]interface{}{}
//	}
//
//	var payloadJSON interface{}
//	json.Unmarshal(payload, &payloadJSON)
//
//	return map[string]interface{}{
//		"at":      fetchedAt,
//		"payload": payloadJSON,
//	}
//}
//
//func writeCache(ctx context.Context, db *sqlx.DB, source string, payload interface{}) error {
//	payloadJSON, err := json.Marshal(payload)
//	if err != nil {
//		return err
//	}
//
//	query := `INSERT INTO space_cache(source, payload) VALUES ($1, $2)`
//	_, err = db.ExecContext(ctx, query, source, payloadJSON)
//	return err
//}
//
//func fetchAPOD(ctx context.Context, db *sqlx.DB, nasaClient *nasa.Client) error {
//	jsonData, err := nasaClient.FetchAPOD(ctx)
//	if err != nil {
//		return err
//	}
//	return writeCache(ctx, db, "apod", jsonData)
//}
//
//func fetchNEOFeed(ctx context.Context, db *sqlx.DB, nasaClient *nasa.Client) error {
//	today := time.Now().UTC()
//	start := today.AddDate(0, 0, -2)
//
//	jsonData, err := nasaClient.FetchNEOFeed(ctx, start.Format("2006-01-02"), today.Format("2006-01-02"))
//	if err != nil {
//		return err
//	}
//	return writeCache(ctx, db, "neo", jsonData)
//}
//
//func fetchDONKIFLR(ctx context.Context, db *sqlx.DB, nasaClient *nasa.Client) error {
//	from, to := lastDays(5)
//	jsonData, err := nasaClient.FetchDONKIFLR(ctx, from, to)
//	if err != nil {
//		return err
//	}
//	return writeCache(ctx, db, "flr", jsonData)
//}
//
//func fetchDONKICME(ctx context.Context, db *sqlx.DB, nasaClient *nasa.Client) error {
//	from, to := lastDays(5)
//	jsonData, err := nasaClient.FetchDONKICME(ctx, from, to)
//	if err != nil {
//		return err
//	}
//	return writeCache(ctx, db, "cme", jsonData)
//}
//
//func fetchSpaceXNext(ctx context.Context, db *sqlx.DB, spacexClient *spacex.SpaceXClient) error {
//	jsonData, err := spacexClient.FetchNextLaunch(ctx)
//	if err != nil {
//		return err
//	}
//	return writeCache(ctx, db, "spacex", jsonData)
//}
//
//func fetchAndStoreISS(ctx context.Context, db *sqlx.DB, issClient *iss.Client) error {
//	jsonData, err := issClient.FetchISS(ctx)
//	if err != nil {
//		return err
//	}
//
//	payloadJSON, err := json.Marshal(jsonData)
//	if err != nil {
//		return err
//	}
//
//	query := `INSERT INTO iss_fetch_log (source_url, payload) VALUES ($1, $2)`
//	_, err = db.ExecContext(ctx, query, issClient.BaseURL, payloadJSON)
//	return err
//}
//
//func fetchAndStoreOSDR(ctx context.Context, db *sqlx.DB, nasaClient *nasa.Client, nasaURL string) (int, error) {
//	jsonData, err := nasaClient.FetchOSDR(ctx, nasaURL)
//	if err != nil {
//		return 0, err
//	}
//
//	var items []interface{}
//	switch v := jsonData.(type) {
//	case []interface{}:
//		items = v
//	case map[string]interface{}:
//		if vItems, ok := v["items"].([]interface{}); ok {
//			items = vItems
//		} else if vResults, ok := v["results"].([]interface{}); ok {
//			items = vResults
//		} else {
//			items = []interface{}{v}
//		}
//	default:
//		items = []interface{}{v}
//	}
//
//	written := 0
//	for _, item := range items {
//		itemMap, ok := item.(map[string]interface{})
//		if !ok {
//			continue
//		}
//
//		id := sPick(itemMap, []string{"dataset_id", "id", "uuid", "studyId", "accession", "osdr_id"})
//		title := sPick(itemMap, []string{"title", "name", "label"})
//		status := sPick(itemMap, []string{"status", "state", "lifecycle"})
//		updated := tPick(itemMap, []string{"updated", "updated_at", "modified", "lastUpdated", "timestamp"})
//
//		itemJSON, _ := json.Marshal(item)
//
//		if id != nil {
//			// Upsert
//			query := `INSERT INTO osdr_items(dataset_id, title, status, updated_at, raw)
//					  VALUES($1, $2, $3, $4, $5)
//					  ON CONFLICT (dataset_id) DO UPDATE
//					  SET title=EXCLUDED.title, status=EXCLUDED.status,
//						  updated_at=EXCLUDED.updated_at, raw=EXCLUDED.raw`
//			_, err := db.ExecContext(ctx, query, *id, title, status, updated, itemJSON)
//			if err != nil {
//				continue
//			}
//		} else {
//			// Insert без dataset_id
//			query := `INSERT INTO osdr_items(dataset_id, title, status, updated_at, raw)
//					  VALUES($1, $2, $3, $4, $5)`
//			_, err := db.ExecContext(ctx, query, nil, title, status, updated, itemJSON)
//			if err != nil {
//				continue
//			}
//		}
//		written++
//	}
//
//	return written, nil
//}
//
//// Вспомогательные функции
//
//func lastDays(n int) (string, string) {
//	to := time.Now().UTC()
//	from := to.AddDate(0, 0, -n)
//	return from.Format("2006-01-02"), to.Format("2006-01-02")
//}
//
//func sPick(v map[string]interface{}, keys []string) *string {
//	for _, k := range keys {
//		if x, ok := v[k]; ok {
//			if s, ok := x.(string); ok && s != "" {
//				return &s
//			}
//			if n, ok := x.(float64); ok {
//				s := fmt.Sprintf("%.0f", n)
//				return &s
//			}
//		}
//	}
//	return nil
//}
//
//func tPick(v map[string]interface{}, keys []string) *time.Time {
//	for _, k := range keys {
//		if x, ok := v[k]; ok {
//			if s, ok := x.(string); ok {
//				// Попробовать разные форматы
//				formats := []string{
//					time.RFC3339,
//					"2006-01-02T15:04:05Z07:00",
//					"2006-01-02 15:04:05",
//					"2006-01-02T15:04:05",
//				}
//				for _, format := range formats {
//					if t, err := time.Parse(format, s); err == nil {
//						utc := t.UTC()
//						return &utc
//					}
//				}
//			}
//			if n, ok := x.(float64); ok {
//				t := time.Unix(int64(n), 0).UTC()
//				return &t
//			}
//		}
//	}
//	return nil
//}
//
//func num(v map[string]interface{}, key string) *float64 {
//	if x, ok := v[key]; ok {
//		if f, ok := x.(float64); ok {
//			return &f
//		}
//		if s, ok := x.(string); ok {
//			if f, err := strconv.ParseFloat(s, 64); err == nil {
//				return &f
//			}
//		}
//	}
//	return nil
//}
//
//func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
//	const earthRadiusKm = 6371.0
//
//	rlat1 := lat1 * math.Pi / 180.0
//	rlat2 := lat2 * math.Pi / 180.0
//	dlat := (lat2 - lat1) * math.Pi / 180.0
//	dlon := (lon2 - lon1) * math.Pi / 180.0
//
//	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
//		math.Cos(rlat1)*math.Cos(rlat2)*
//			math.Sin(dlon/2)*math.Sin(dlon/2)
//	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
//
//	return earthRadiusKm * c
//}
