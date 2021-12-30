package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	//"encoding/json"

	"github.com/go-redis/redis"

	//	SimpleStorage "./contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/gofiber/fiber/v2"
	//"github.com/ethereum/go-ethereum/common"
	/*"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"*/)

//var Client ethclient.Client

func main() {

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("KEC App Running!\n")
	})

	app.Get("/api/v1/testFunction", func(c *fiber.Ctx) error {

		testFunction()
		return c.SendString("Called testFunction\n")
	})

	app.Post("/api/v1/deployContract/:id", func(c *fiber.Ctx) error {

		address, err := deployContract(c.Params("id"))

		if err != nil {
			msg := fmt.Sprintf("Contract did not deploy.  Error %s\n", err)
			return c.SendString(msg)
		}
		msg := fmt.Sprintf("Contract Deployed Successfully, Address is %s\n", address)
		return c.SendString(msg)

	})

	app.Post("/api/v1/user/:id", func(c *fiber.Ctx) error {

		err := addUser(c.Params("id"))

		if err == nil {
			fmt.Println(err)
			msg := fmt.Sprintf("User %s Already Exists\n", c.Params("id"))
			return c.SendString(msg)
		}
		msg := fmt.Sprintf("Added User %s\n", c.Params("id"))
		return c.SendString(msg)
	})

	app.Post("/api/v1/awardItem/:id/:item/:contracthex", func(c *fiber.Ctx) error {

		if c.Params("item") != "hammer" {
			msg := fmt.Sprintf("Item %s is not a valid item.  Run api/v1/getitems to see valid options\n", c.Params("item"))
			return c.SendString(msg)
		}

		itemurl := "https://game.example/item-id-8u5h2m.json"
		tx, err := awardItem(c.Params("id"), itemurl, c.Params("contracthex"))

		if err != nil {
			msg := fmt.Sprintf("Item was not awarded to user %s with url %s and hex address %s\n", c.Params("id"), c.Params("itemurl"), c.Params("contracthex"))
			return c.SendString(msg)
		}
		msg := fmt.Sprintf("Awarded item %s to user %s.  Transaction is: %s", c.Params("item"), c.Params("id"), tx)
		return c.SendString(msg)
	})

	app.Get("api/v1/getItems", func(c *fiber.Ctx) error {
		return c.SendString("hammer\n")
	})

	app.Get("api/v1/getOwner/:itemid", func(c *fiber.Ctx) error {

		owner, err := getOwner(c.Params("itemid"))

		if err != nil {
			msg := fmt.Sprintf("Item %s does not exist", c.Params("itemid"))
			return c.SendString(msg)
		}

		msg := fmt.Sprintf("Getowner not implemented yet.  Owner is %s\n", owner)
		return c.SendString(msg)
	})

	app.Listen(":3000")

	//accounts.NewManager()
}

func testFunction() error {

	//function to test random stuff while I am learning

	var PrivateKey *ecdsa.PrivateKey

	client, err := getEthClient()

	//client, err := ethclient.Dial("http://host.docker.internal:22001")

	if err != nil {
		log.Fatal("Couldn't get Client Connection", err)
	}
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {

		log.Fatal("Headerbynumber call failed", err)
	}

	fmt.Println("Block Header is", header.Number.String())

	//blockNumber := big.NewInt(39)
	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Block number is", block.Number().Uint64())
	fmt.Println("Block Time is", block.Time())
	fmt.Println("Block Difficulty is", block.Difficulty().Uint64())
	fmt.Println("Block Hash is", block.Hash().Hex())
	fmt.Println("Number of Transactions in block is", len(block.Transactions()))

	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Confirmed number of transactions in block is", count)

	fee := big.NewInt(50)

	for _, tx := range block.Transactions() {
		fmt.Println("Transaction Hash is", tx.Hash().Hex())
		fmt.Println("Transaction Value is", tx.Value().String())
		fmt.Println("Transaction Gas is", tx.Gas())
		fmt.Println("Transaction Gasprice is", tx.GasPrice().Uint64())
		fmt.Println("Transaction Nonce is", tx.Nonce())
		fmt.Println("Transaction Data is", tx.Data())
		//fmt.Println("Transaction Recipient Address is", tx.To().Hex())

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID), fee); err != nil {
			fmt.Println("Message is", msg.From().Hex())
		}

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Receipt Status is", receipt.Status) // 1
		fmt.Println("Receipt Logs are ", receipt.Logs)   // ...
	}

	//blockHash := common.HexToHash("0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9")
	//kcount, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal(err)
	}

	//get TX by hash.

	firstTx, err := client.TransactionInBlock(context.Background(), block.Hash(), 0)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Transaction Hash is", firstTx.Hash().Hex())
	fmt.Println("Transaction To address is", firstTx.To())

	//generate private key
	// serializedKeys[i] = hex.EncodeToString(ecrypto.FromECDSA(keys[i]))
	//storeKey("privateKey", privateKey)

	fmt.Println("Attempting to pull private key from DB")

	PrivateKey, err = getKey("privatekey")

	if err != nil {
		fmt.Println("Private Key not found.  Need to generate")

		//generate new key
		PrivateKey, err := crypto.GenerateKey()

		if err != nil {
			fmt.Println("Failed to Generate Private Key")
			log.Fatal(err)
		}

		storeKey("privateKey", PrivateKey)

		/*if err != nil {
		    log.Fatal("Error storing private key")
		}*/

	} else {
		fmt.Println("Found Key!")
	}

	publicKey := PrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	//value := big.NewInt(1000000000000000000)

	//gasLimit := uint64(21000) // in units

	//gasPrice := big.NewInt(30000000000) // in wei (30 gwei)
	gasPrice := big.NewInt(0)
	//gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	chainID := big.NewInt(3543006677)
	auth, err := bind.NewKeyedTransactorWithChainID(PrivateKey, chainID)

	if err != nil {
		log.Fatal("New Keyed Transactor did not get assigned", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)       // in wei
	auth.GasLimit = uint64(30000000) // in units
	auth.GasPrice = gasPrice

	address, tx2, instance, err := DeployGameitem(auth, client)

	if err != nil {
		log.Fatal("Contract GameItem did not deploy", err)
	}

	//toAddress := tx.To()

	//tx := types.NewTransaction(nonce, *toAddress, value, gasLimit, gasPrice, nil)

	//chainid is not networkid in the kaleido example
	//chainID, err := client.NetworkID(context.Background())

	fmt.Println("Hex of Contract address is ", address.Hex())
	fmt.Println("Contract Transaction Hash is ", tx2.Hash().Hex())
	fmt.Println("Contract Instance is ", instance)

	itemId, err := instance.AwardItem(auth, address, "https://game.example/item-id-8u5h2m.json")

	if err != nil {
		log.Fatal("Unable to Award Item", err)
	}

	fmt.Println("Item ID is ", itemId)
	fmt.Println("Item ID data is ", itemId.Data())

	/*if err != nil {
		log.Fatal(err)
	}*/

	fmt.Println("ChainID pulled is", chainID)

	/*signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())*/

	return err

}

func getKey(keyname string) (*ecdsa.PrivateKey, error) {

	redisClient, err := getRedisClient()

	if err != nil {
		log.Fatal("couldn't get RedisClient", err)
	}

	keyString, err := redisClient.Get(keyname).Result()

	if err != nil {
		fmt.Println("key not found:", keyname, err)
		return nil, err
	}

	//takes a key string and returns the ecdsa key.

	fmt.Println("Found Key!")

	//decode string to byte array

	byteKey, err := hex.DecodeString(keyString)

	fmt.Println("Key String we pulled is", keyString)

	if err != nil {
		log.Fatal("Failed to decode key", err)
	}

	privateKey, err := crypto.ToECDSA(byteKey)

	if err != nil {
		log.Fatal("failure to decode privatekey from bytestring", err)
	}

	return privateKey, nil
}

func storeKey(keyname string, privateKey *ecdsa.PrivateKey) error {

	redisClient, err := getRedisClient()

	if err != nil {
		log.Fatal("couldn't get RedisClient", err)
	}

	// store key in DB
	keyString := hex.EncodeToString(crypto.FromECDSA(privateKey))

	fmt.Println("KeyString we stored is", keyString)
	//jsonKey, err := json.Marshal(serializedKey)

	/*if err != nil {
	    log.Fatal("Failed to Marshal Key.  Exiting.", err)
	}*/
	redisClient.Set(keyname, keyString, 0)

	//takes a key and stores it in the DB under keyname
	return nil

}

func getEthClient() (*ethclient.Client, error) {

	hostURL := "http://localhost:22001"
	dockerURL := "http://host.docker.internal:22001"

	fmt.Println("Attempting to connect to Ethereum Node at", dockerURL)
	client, err := ethclient.Dial(dockerURL)

	// so the err value isn't non nil on error so I'll need to actually get a decent error or maybe I can just check client.

	//client.

	if err != nil {
		fmt.Printf("Could not coonect to DB at %s %s.  Trying docker URL", dockerURL, err)
	} else {
		fmt.Println("Returning client for docker URL")
		return client, err
	}

	fmt.Println("trying to connect to hostURL")
	client, err = ethclient.Dial(dockerURL)

	if err != nil {
		log.Fatal("Couldn't connect to Anything", hostURL)
	}

	fmt.Println("we have a connection")

	return client, err

}

func addUser(userName string) error {

	//add a user userName and store their address in the DB as userName-address
	//should change this to a struct later probably to better store data

	//see if user already exists
	keytoGet := userName + "-key"
	fmt.Println("keytoGet is ", keytoGet)
	userKey, err := getKey(keytoGet)

	fmt.Println("UserKey and Error are", keytoGet, err)

	if err != redis.Nil {
		fmt.Println("User already exists", userName)
		fmt.Println("Userkey is ", userKey)
		return err

	} else {
		fmt.Println("User does not exist, creating", userName)
		privateKey, err := crypto.GenerateKey()

		if err != nil {
			fmt.Println("Failed to Generate Private Key")
			log.Fatal(err)
		}

		err2 := storeKey(userName+"-key", privateKey)

		if err2 != nil {
			fmt.Println("Error Storing Key!")
		}

		/*publicKey := privateKey.Public()
				publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
				if !ok {
		  			log.Fatal("error casting public key to ECDSA")
				}

				fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)*/
	}

	return err
}

func getUserKey(userName string) (*ecdsa.PrivateKey, error) {

	keytoGet := userName + "-key"
	fmt.Println("keytoGet is ", keytoGet)
	userKey, err := getKey(keytoGet)

	fmt.Println("UserKey and Error are", keytoGet, err)

	if err != redis.Nil {
		fmt.Println("User already exists", userName)
		fmt.Println("Userkey is ", userKey)
	}
	return userKey, err
}

func getRedisClient() (redis.Client, error) {

	//function that returns a connection to Redis
	//could make this a global later

	redisURL := "localhost:6379"
	redisdockerURL := "redis:6379"

	fmt.Println("Attemptig to connect to Redis DB at", redisURL)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping().Result()

	if err != nil {
		fmt.Println("Error Connecting to DB, trying by name")
	}

	fmt.Println(pong)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisdockerURL,
		Password: "",
		DB:       0,
	})

	pong, err = redisClient.Ping().Result()

	if err != nil {
		log.Fatal("Error Connecting to DB by name.  Exiting", err)
	}

	fmt.Println("Good Response from DB", pong)
	return *redisClient, err
}

func deployContract(deployer string) (string, error) {

	//Deploys an instance of the contract and returns the address so we can use it later to interact with the depoyed contract

	client, _ := getEthClient()

	privateKey, err := getUserKey(deployer)

	if err != nil {
		fmt.Println("Unable to get Private Key!!")
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	chainID := big.NewInt(3543006677)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)

	if err != nil {
		log.Fatal("New Keyed Transactor did not get assigned", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)       // in wei
	auth.GasLimit = uint64(30000000) // in units
	auth.GasPrice = big.NewInt(0)

	address, tx2, instance, err := DeployGameitem(auth, client)

	if err != nil {
		log.Fatal("Contract GameItem did not deploy", err)
	}

	fmt.Println("Hex of Contract address is ", address.Hex())
	fmt.Println("Contract Transaction Hash is ", tx2.Hash().Hex())
	fmt.Println("Contract Instance is ", instance)

	/*itemId, err := instance.AwardItem(auth, address, "https://game.example/item-id-8u5h2m.json")

	if err != nil {
		log.Fatal("Unable to Award Item", err)
	}

	fmt.Println("Item ID is ", itemId)
	fmt.Println("Item ID data is ", itemId.Data())*/

	//deploys a contract

	return address.Hex(), nil
}

func getBalance(username string) uint16 {
	//gets ether balance for a user

	return 50
}

func awardItem(awardtoplayer string, tokenURI string, contractHex string) (string, error) {

	// awards an item to a player
	//takes playerID and item type and mayb contract instance

	client, _ := getEthClient()

	// get address of player getting item
	privateKey, err := getUserKey(awardtoplayer)

	if err != nil {
		fmt.Println("Unable to get Private Key!!")
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	playerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//get instance of deployed contract from address

	contractAddress := common.HexToAddress(contractHex)

	instance, err := NewGameitem(contractAddress, client)

	if err != nil {
		log.Fatal("could not get instance of NewGameitem")
	}

	fmt.Println("Contract is loaded!")

	//award the item

	nonce, err := client.PendingNonceAt(context.Background(), playerAddress)
	if err != nil {
		log.Fatal(err)
	}

	chainID := big.NewInt(3543006677)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)

	if err != nil {
		log.Fatal("New Keyed Transactor did not get assigned", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)       // in wei
	auth.GasLimit = uint64(30000000) // in units
	auth.GasPrice = big.NewInt(0)

	tx, err := instance.AwardItem(auth, playerAddress, tokenURI)

	if err != nil {
		log.Fatal("Unable to Award Item", err)
	}

	//fmt.Println("Item ID is ", item.)
	fmt.Println("Transaction Sent: ", tx.Hash().Hex())
	fmt.Println("Item ID data is ", tx.Data())
	fmt.Println("Player Address receiving item is ", playerAddress)

	//sleep because I'm lame

	fmt.Println("Waiting 10s to be able to mine transaction..")

	time.Sleep(10 * time.Second)

	//get latest blocknumber

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {

		log.Fatal("Headerbynumber call failed", err)
	}

	fmt.Println("Block Header is", header.Number.String())

	/*blockNumber := big.NewInt(39)
	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}*/

	callopts := bind.CallOpts{false, playerAddress, header.Number, nil}
	itemnum := big.NewInt(1)

	owner, err2 := instance.OwnerOf(&callopts, itemnum)

	if err2 != nil {
		fmt.Printf("Unable to get owner of item %d.  Error is %s\n", itemnum, err2)
	} else {

		fmt.Printf("Owner of item %d is address %s", itemnum, owner)
		itemURI, _ := instance.TokenURI(&callopts, itemnum)
		fmt.Println("Item 1 details are ", itemURI)

	}

	return tx.Hash().Hex(), err

}

func getOwner(id string) (string, error) {
	//input: takes in an item ID
	//output: returns the owner of that item.

	//Get owner of item 1
	//item.UnmarshalJSON(item.Data())

	return "", nil
}

func tradeItem() {
	//trades item from one player to another
	//the player needs to own the item to transfer, takes player 1 and 2 and item id.
}

func destroyitem() {
	//destroys an item the player owns.
	//takes player and item id.
}
