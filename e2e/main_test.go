package e2e

import (
	"fmt"
	"la-warehouse/config"
	"la-warehouse/internal/model"
	"la-warehouse/internal/repo"
	"la-warehouse/internal/service"
	"la-warehouse/internal/transport"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog"
)

type ReqReservation struct {
	Method string                 `json:"method"`
	Params []model.ReqReservation `json:"params"`
	Id     int                    `json:"id"`
}

type RespReservation struct {
	Result struct {
		Products []model.RespReservation `json:"products"`
	} `json:"result"`
	Error interface{} `json:"error"`
	Id    int         `json:"id"`
}

const (
	POSTGRES_PASSWORD = "postgres"
	POSTGRES_USER     = "postgres"
	POSTGRES_DB       = "postgres"
	TEST_URL          = ":8012"

	RESERVATION_METHOD     = "Warehouse.Reservation"
	RESERVATIONFREE_METHOD = "Warehouse.ReservationFree"
)

var (
	pool     *dockertest.Pool
	resource *dockertest.Resource
	pg       *repo.Pg
	logger   *zerolog.Logger
)

func TestMain(m *testing.M) {
	defer func() {

	}()
	host, post := initTestContainer()

	initMigration(host, post)

	tmpLog := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "02.01.2006 -- 15:04:05.9999",
	}).With().Timestamp().Logger()
	logger = &tmpLog
	serv := service.NewWarehouse(logger, pg)

	transportWh := transport.NewWarehouse(serv)
	go func() {
		err := transportWh.ListenAndServe(TEST_URL)
		if err != nil {
			return
		}
	}()
	//time.Sleep(time.Second * 5)
	m.Run()
}

func initTestContainer() (string, string) {
	var err error
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_PASSWORD=" + POSTGRES_PASSWORD,
			"POSTGRES_USER=" + POSTGRES_USER,
			"POSTGRES_DB=" + POSTGRES_DB,
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// остановленный контейнер сам удаляется
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	timeToKillContainerSec := uint(60 * 5)
	_ = resource.Expire(timeToKillContainerSec) // hard kill контейнера спустя это время

	host := "localhost"
	port := resource.GetPort("5432/tcp")

	initDB(host, port)
	return host, port
}

func initDB(host, port string) {
	var err error
	cfg := config.Postgres{
		Host:     host,
		Port:     port,
		User:     POSTGRES_USER,
		Password: POSTGRES_PASSWORD,
		DBName:   POSTGRES_DB,
		SSLMode:  "disable",
		Settings: struct {
			MaxOpenConns    int           `yaml:"MaxOpenConns"`
			ConnMaxLifeTime time.Duration `yaml:"ConnMaxLifeTime"`
			MaxIdleConns    int           `yaml:"MaxIdleConns"`
			MaxIdleLifeTime time.Duration `yaml:"MaxIdleLifeTime"`
		}{
			MaxOpenConns:    10,
			ConnMaxLifeTime: time.Second * 10,
			MaxIdleConns:    10,
			MaxIdleLifeTime: time.Second * 5,
		},
	}
	pg, err = repo.NewPostgres(cfg)

	if err = pool.Retry(func() error {
		err = pg.Ping()
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
}

// initMigration поднятие миграций
func initMigration(host, port string) {
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		POSTGRES_USER, POSTGRES_PASSWORD, host, port, POSTGRES_DB)
	migrator, err := migrate.New("file://../migrations/postgres", databaseUrl)
	if err != nil {
		log.Fatalf("Could not create migrator for DB: %s", err)
	}

	if err := migrator.Up(); err != nil {
		log.Println(err)
		return
	}
}
