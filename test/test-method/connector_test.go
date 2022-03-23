package test_method

import (
	"testing"
	"z-common/connector"
)

func TestGetMySQLEngine(t *testing.T) {
	config := connector.NewDefaultMysqlConfig()
	_, err := config.GetMySQLEngine()
	if err != nil {
		t.Errorf("TestGetMySQLEngine.GetMySQLEngine err:%v", err)
	}
	t.Log("DB success")
}

func BenchmarkGetDsnByString(b *testing.B) {
	config := connector.NewDefaultMysqlConfig()
	for i := 0; i < b.N; i++ {
		config.GetDsnByString()
	}
}

func BenchmarkGetDsnByBuffer(b *testing.B) {
	config := connector.NewDefaultMysqlConfig()
	for i := 0; i < b.N; i++ {
		config.GetDsnByString()
	}
}
