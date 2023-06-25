package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {
	// 単純な出力
	slog.Info("hello", "count", 3)
	// コンテキストのログ出力
	ctx := context.WithValue(context.Background(), "key", "value")
	slog.InfoCtx(ctx, "メッセージ", "key2", "value2")
	// text形式で出力
	logger5 := slog.New(slog.NewTextHandler(os.Stderr, nil))
	logger5.Info("hello", "count", 3)
	// json形式で出力
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("hello1", "name", "blue", "No", 5)
	// ログにjsonのkeyとvalueを追加する
	logger.Info("hello2",
		slog.Bool("allowed", true),
		slog.String("name2", "red"),
	)
	// 任意のjson形式のログを構造体で追加する
	data := struct {
		Color  string
		Nested struct {
			Depth int
		}
	}{
		Color: "yellow",
		Nested: struct{ Depth int }{
			Depth: 1,
		},
	}
	logger.Info("hello3", slog.Any("json_data_key", data))
	// 任意のjson形式のログを構造体を利用せずに追加する
	logger.Info("hello4",
		slog.Group("group1",
			slog.String("color", "green"),
		),
		slog.Group("group2",
			slog.String("hight", "170"),
		),
	)
	// ログレベルの設定
	logger4 := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
	logger4.Warn("ログレベルがWarn")
	logger4.Info("ログレベルがInfo")
	logger4.Debug("ログレベルがDebug")
	logger4.Error("ログレベルがError")
}
