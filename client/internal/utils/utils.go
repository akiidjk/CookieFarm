package utils

import (
	"os"
	"runtime"
	"time"

	"math/rand"

	"github.com/ByteTheCookies/cookiefarm-client/internal/config"
	"github.com/ByteTheCookies/cookiefarm-client/internal/models"
	"github.com/google/uuid"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

func CleanGC() (uint64, uint64) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	before := mem.Alloc / 1_048_576
	runtime.GC() //Cleaning garbage collector
	runtime.ReadMemStats(&mem)
	after := mem.Alloc / 1_048_576
	return before, after
}

func RandInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func GenerateFakeFlag(flagCode string) models.Flag {
	randomService := config.FAKE_SERVICES[RandInt(1, len(config.FAKE_SERVICES))]
	id, _ := uuid.NewV7()
	return models.Flag{
		ID:           id.String(),
		FlagCode:     flagCode,
		ServiceName:  randomService.Name,
		ServicePort:  randomService.Port,
		SubmitTime:   uint64(time.Now().UnixNano()),
		ResponseTime: 0,
		Status:       "UNSUBMITTED",
		TeamID:       uint16(RandInt(1, 40)),
	}
}
