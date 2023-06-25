package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"testing/slogtest"
)

func main() {
	// 1. 単純な出力
	slog.Info("hello", "count", 3)
	/** 1. 出力結果
	2023/06/25 20:34:28 INFO hello count=3
	*/
	// 2. コンテキストのログ出力
	ctx := context.WithValue(context.Background(), "key", "value")
	slog.InfoCtx(ctx, "メッセージ", "key2", "value2")
	/** 2. 出力結果
	2023/06/25 20:34:28 INFO メッセージ key2=value2
	*/
	// 3. text形式で出力
	logger1 := slog.New(slog.NewTextHandler(os.Stderr, nil))
	logger1.Info("hello", "count", 3)
	/** 3. 出力結果
	time=2023-06-25T20:34:28.878+09:00 level=INFO msg=hello count=3
	*/
	// 4. json形式で出力
	logger2 := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger2.Info("hello1", "name", "sato", "age", 23)
	/** 4. 出力結果
	{"time":"2023-06-25T20:34:28.87839+09:00","level":"INFO","msg":"hello1","name":"sato","age":23}
	*/
	// 5. ログにjsonのkeyとvalueを追加する
	logger2.Info("hello2",
		slog.Bool("isTrue", true),
		slog.String("name2", "suzuki"),
	)
	/** 5. 出力結果
	{"time":"2023-06-25T20:34:28.878413+09:00","level":"INFO","msg":"hello2","isTrue":true,"name2":"suzuki"
	*/
	// 6. 任意のjson形式のログを構造体で追加する
	data := struct {
		Name   string
		Nested struct {
			Depth int
		}
	}{
		Name: "tanaka",
		Nested: struct{ Depth int }{
			Depth: 1,
		},
	}
	logger2.Info("hello3", slog.Any("json_data_key", data))
	/** 6. 出力結果
	{"time":"2023-06-25T20:34:28.878417+09:00","level":"INFO","msg":"hello3","json_data_key":{"Color":"tanaka","Nested":{"Depth":1}}}
	*/
	// 7. 任意のjson形式のログを構造体を利用せずに追加する
	logger2.Info("hello4",
		slog.Group("group1",
			slog.String("name", "kawasaki"),
		),
		slog.Group("group2",
			slog.String("hight", "170"),
		),
	)
	/** 7. 出力結果
	{"time":"2023-06-25T20:34:28.878519+09:00","level":"INFO","msg":"hello4","group1":{"name":"kawasaki"},"group2":{"hight":"170"}}
	*/
	// 8. ログレベルの設定
	logger3 := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
	logger3.Warn("ログレベルがWarn")
	logger3.Info("ログレベルがInfo")
	logger3.Debug("ログレベルがDebug")
	logger3.Error("ログレベルがError")
	/** 8. 出力結果
	{"time":"2023-06-25T20:34:28.878528+09:00","level":"WARN","msg":"ログレベルがWarn"}
	{"time":"2023-06-25T20:34:28.878532+09:00","level":"ERROR","msg":"ログレベルがError"}
	*/
	// 9. ログレベルがerrorのものだけ出力する
	logger4 := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	if logger4.Enabled(context.Background(), slog.LevelError) {
		logger4.Error("9.ログレベルがErrorです")
	}
	if logger4.Enabled(context.Background(), slog.LevelWarn) {
		logger4.Error("9.ログレベルがWarnです")
	}
	if logger4.Enabled(context.Background(), slog.LevelDebug) {
		logger4.Error("9.ログレベルがDebugです")
	}
	if logger4.Enabled(context.Background(), slog.LevelInfo) {
		logger4.Error("9.ログレベルがInfoです")
	}
	/** 9. 出力結果
	{"time":"2023-06-25T20:42:08.558559+09:00","level":"ERROR","msg":"9.ログレベルがErrorです"}
	*/
	// 10 slogtestの利用
	slogTest()
}

func slogTest() {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, nil)

	results := func() []map[string]any {
		var ms []map[string]any
		for _, line := range bytes.Split(buf.Bytes(), []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			var m map[string]any
			if err := json.Unmarshal(line, &m); err != nil {
				panic(err) // In a real test, use t.Fatal.
			}
			ms = append(ms, m)
		}
		return ms
	}
	err := slogtest.TestHandler(h, results)
	if err != nil {
		log.Fatal(err)
	}
}
