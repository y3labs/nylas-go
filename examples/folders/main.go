package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/y3labs/nylas-go/nylas"
)

func main() {
	apiKey := mustEnv("NYLAS_API_KEY")
	grantID := mustEnv("NYLAS_GRANT_ID")
	apiURI := os.Getenv("NYLAS_API_URI")

	client := nylas.NewClient(apiKey)
	if apiURI != "" {
		client = nylas.NewClient(apiKey, nylas.WithServerURL(apiURI))
	}

	ctx := context.Background()

	// Multi-level (default)
	resp, err := client.Folders().List(ctx, grantID, nil)
	if err != nil {
		log.Fatalf("list folders: %v", err)
	}
	if len(resp.Data) == 0 {
		fmt.Println("No folders found.")
	} else {

		for i, folder := range resp.Data {
			fmt.Printf("[%d], id=%s, name=%q\n",
				i, folder.ID, folder.Name)
		}
	}

	fmt.Println("\n\n=== Create Folder ===")
	create_resp, err := client.Folders().Create(ctx, grantID, nylas.CreateFolderRequest{
		Name: "Example Folder",
	})
	if err != nil {
		log.Fatalf("create folder: %v", err)
	} else {
		fmt.Printf("\nid=%s, name=%q", create_resp.Data.ID, create_resp.Data.Name)
	}

	fmt.Println("\n\n=== Show New Folder ===")
	get_resp, err := client.Folders().Get(ctx, grantID, create_resp.Data.ID, nil)
	if err != nil {
		log.Fatalf("get folder: %v", err)
	} else {
		fmt.Printf("\n === Folder Found ===\n")
		fmt.Printf("id=%s, name=%q\n", get_resp.Data.ID, get_resp.Data.Name)
	}

	fmt.Println("\n\n=== Clean Up Example Folder ===")
	if err := client.Folders().Delete(ctx, grantID, create_resp.Data.ID); err != nil {
		log.Fatalf("delete folder: %v", err)
	}

	new_list, err := client.Folders().List(ctx, grantID, nil)
	if err != nil {
		log.Fatalf("list folders: %v", err)
	}
	if len(new_list.Data) == 0 {
		fmt.Println("No folders found.")
	} else {

		for i, folder := range new_list.Data {
			fmt.Printf("[%d], id=%s, name=%q\n",
				i, folder.ID, folder.Name)
		}
	}
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing %s", k)
	}
	return v
}
