package handlers

import (
	"backend/proto/searchpb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func InitSearchService() *searchpb.SearchServiceClient {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", os.Getenv("SEARCH_SERVICE_ADDRESS"), os.Getenv("SEARCH_SERVICE_PORT")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("could not connect to search service: %v", err)
	}

	searchClient := searchpb.NewSearchServiceClient(conn)

	return &searchClient
}
