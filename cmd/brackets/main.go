package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/yaBliznyk/brackets/internal/brackets"
	"github.com/yaBliznyk/brackets/internal/config"
	"github.com/yaBliznyk/brackets/internal/server"
	gw "github.com/yaBliznyk/brackets/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

// Main function
// Run service and listen for errors
func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	cfg, err := config.MakeConfig()
	if err != nil {
		return errors.Wrap(err, "инициализация конфига")
	}

	err, bktSvc := brackets.NewBaseService(&brackets.Config{
		Logger: nil,
		Bkts:   "{}[]()",
	})
	if err != nil {
		return errors.Wrap(err, "создание base сервиса")
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	srv := server.NewBracketServer(bktSvc)

	wg := sync.WaitGroup{}
	// Запускаем grpc сервер
	{
		lis, err := net.Listen("tcp", cfg.GRPCPort)
		if err != nil {
			return errors.Wrap(err, "невозможно запустить listen на порте "+cfg.GRPCPort)
		}
		s := grpc.NewServer()
		gw.RegisterBracketsServer(s, srv)

		wg.Add(1)
		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)

			}
		}()
	}

	// Запускаем http прокси
	{
		opts := []grpc.DialOption{grpc.WithInsecure()}
		err = gw.RegisterBracketsHandlerFromEndpoint(ctx, mux, "localhost"+cfg.GRPCPort, opts)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			if err := http.ListenAndServe(cfg.HTTPPort, mux); err != nil {
				log.Fatalf("Ошибка запуска grpc прокси: %v", err)
			}
		}()

	}
	wg.Wait()
	return nil
}
