package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sonr-io/highway-go/pkg/client"
	"github.com/sonr-io/sonr/x/registry/types"
)

func main() {
	c, err := client.NewClient(context.Background(), "http://localhost:26657", "test", "bad-password")
	if err != nil {
		log.Fatal(err)
	}

	accName := c.Account.Name
	log.Println("Account name:", accName)
	// define a message to create a post
	msg := &types.MsgRegisterName{
		Creator:        accName,
		NameToRegister: "prad.snr",
	}

	// broadcast a transaction from account `alice` with the message to create a post
	// store response in txResp
	txResp, err := c.BroadcastTx(accName, msg)
	if err != nil {
		log.Fatal(err)
	}

	// print response from broadcasting a transaction
	fmt.Print("MsgCreatePost:\n\n")
	fmt.Println(txResp)

	// instantiate a query client for your `blog` blockchain
	queryClient := types.NewQueryClient(c.Context)

	// query the blockchain using the client's `PostAll` method to get all posts
	// store all posts in queryResp
	queryResp, err := queryClient.WhoIsAll(context.Background(), &types.QueryAllWhoIsRequest{})
	if err != nil {
		log.Fatal(err)
	}

	// print response from querying all the posts
	fmt.Print("\n\nAll posts:\n\n")
	fmt.Println(queryResp)
}
