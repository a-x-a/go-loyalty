package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceConfig(t *testing.T) {
    // Сохраняем текущие значения флагов и переменных окружения
    oldArgs := os.Args
    oldEnv := map[string]string{}
    for _, e := range os.Environ() {
        pair := strings.SplitN(e, "=", 2)
        oldEnv[pair[0]] = pair[1]
    }

    // Очищаем флаги и переменные окружения
    os.Args = oldArgs[:1]
    os.Clearenv()

    // Задаем тестовые значения флагов и переменных окружения
    os.Setenv("RUN_ADDRESS", "localhost:8081")
    os.Setenv("DATABASE_URI", "postgres://user:password@localhost:5432/db")
    os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://localhost:8080")

    // Запускаем тестовый код
    cfg := NewServiceConfig()

    // Проверяем результат
    assert.Equal(t, "localhost:8081", cfg.RunAddress)
    assert.Equal(t, "postgres://user:password@localhost:5432/db", cfg.DatabaseURI)
    assert.Equal(t, "http://localhost:8080", cfg.AccrualSystemAddress)

    // Восстанавливаем тестовые значения флагов и переменных окружения
	os.Args = oldArgs
	for k, v := range oldEnv {
		os.Setenv(k, v)
	}
}
