package main

import (
	"context"
	"main/api"
	"main/api/handler"
	"main/config"
	grpc_client "main/grpc"
	"main/packages/logger"
	"main/storage/redis"
)

func main() {

	cfg := config.Load()
	log := logger.NewLogger("mini-project", logger.LevelInfo)
	redisStrg, err := redis.NewCache(context.Background(), *cfg)
	if err != nil {
		return
	}
	client, err := grpc_client.New(*cfg)
	if err != nil {
		return
	}

	h := handler.NewHandler(*cfg, log, redisStrg, client)

	r := api.NewServer(h)
	r.Run(":8080")

	//withCancel()

}

/* func withCancel() {
	ctx := context.Background()
	// ctx, cancel := context.WithCancel(ctx)
	// // cancel()
	// go func() {
	//  time.Sleep(time.Second)
	//  cancel()
	// }()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	timeConsumingFunc(ctx, 5*time.Second, "hello")
}

func timeConsumingFunc(ctx context.Context, d time.Duration, message string) {
	for {
		select {
		case <-time.After(d):
			fmt.Println(message)
		case <-ctx.Done():
			fmt.Println(ctx.Err().Error())
			return
		}
	}
} */

/* cfg := config.Load()
strg, err := memory.NewStorage(context.Background(), *cfg)

if err != nil {
	return
}

h := handler.NewHandler(strg, *cfg)

http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
	w.WriteHeader(http.StatusOK)
})
http.HandleFunc("/branch/", h.BranchHandler)
http.HandleFunc("/staff/", h.StaffHandler)
http.HandleFunc("/staffTarif/", h.StaffTarifHandler)

fmt.Printf("server is running on port %s\n", cfg.Port)
err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), nil)
if err != nil {
	panic(err)
}

func main() {

	cfg := config.Load()
	strg := memory.NewStorage("data/branches.json", "data/staff.json", "data/sales.json", "data/transactions.json", "data/tarif.json")
	h := handler.NewHandler(strg, *cfg)

	fmt.Println("Welcome to my program!")
	fmt.Println("Available methods:")
	for _, m := range cfg.Methods {
		fmt.Println("- ", m)
	}
	fmt.Println("Available objects:")
	for _, o := range cfg.Objects {
		fmt.Println("- ", o)
	}
	for {
		fmt.Println("Enter method and object:")
		method, object := "", ""
		fmt.Scan(&method)

		if method == "exit" {
			return
		}

		fmt.Scan(&object)

		if object == "exit" {
			return
		}

		switch object {
		case "branch":
			switch method {
			case "create":
				fmt.Println("Enter name and adress: ")
				name, adress := "", ""
				fmt.Scan(&name, &adress)
				h.CreateBranch(name, adress)
			case "get":
				fmt.Print("Enter ID: ")
				var id string
				fmt.Scan(&id)
				h.GetBranch(id)
			case "getAll":
				fmt.Print("Enter search text: ")
				var search string
				fmt.Scan(&search)
				h.GetAllBranch(1, 10, search)
			case "update":
				fmt.Println("Enter ID, name, adress and founded year: ")
				id, name, adress := "", "", ""
				fmt.Scan(&id, &name, &adress)
				h.UpdateBranch(id, name, adress)
			case "delete":
				fmt.Print("Enter ID: ")
				id := ""
				fmt.Scan(&id)
				h.DeleteBranch(id)
			}

		case "staff":
			switch method {
			case "create":

				fmt.Println("Enter TypeStaff(int): ")
				var typeStaff int
				fmt.Scan(&typeStaff)

				fmt.Println("Enter BranchID: ")
				var branchId string
				fmt.Scan(&branchId)

				fmt.Println("Enter TarifID: ")
				var tarifId string
				fmt.Scan(&tarifId)

				fmt.Println("Enter Name: ")
				var name string
				fmt.Scan(&name)

				fmt.Println("Enter BirthDate: ")
				var birthDate string
				fmt.Scan(&birthDate)

				fmt.Println("Enter Balance: ")
				var balance float64
				fmt.Scan(&balance)

				h.CreateStaff(typeStaff, branchId, tarifId, name, birthDate, balance)

			}

		case "staffTarif":
			switch method {
			case "create":

				fmt.Println("Enter Type(int): ")
				var typetrafic int
				fmt.Scan(&typetrafic)

				fmt.Println("Enter FoundedAt: ")
				var foundedAt string
				fmt.Scan(&foundedAt)

				fmt.Println("Enter AmountForCard: ")
				var amountForCard float64
				fmt.Scan(&amountForCard)

				fmt.Println("Enter Name: ")
				var name string
				fmt.Scan(&name)

				fmt.Println("Enter AmountForCash")
				var amountForCash float64
				fmt.Scan(&amountForCash)

				h.CreateStaffTarif(typetrafic, name, foundedAt, amountForCard, amountForCash)

			}

		case "sales":
			switch method {
			case "create":

				fmt.Println("Enter PaymentType(int): ")
				var paymenttype int
				fmt.Scan(&paymenttype)

				fmt.Println("Enter Status(suc, can): ")
				var status int
				fmt.Scan(&status)

				fmt.Println("Enter BranchId: ")
				var branchId string
				fmt.Scan(&branchId)

				fmt.Println("Enter ShopAssistentId: ")
				var shopAssistantid string
				fmt.Scan(&shopAssistantid)

				fmt.Println("Enter CashierId: ")
				var cashierId string
				fmt.Scan(&cashierId)

				fmt.Println("Enter ClientName: ")
				var clientName string
				fmt.Scan(&clientName)

				fmt.Println("Enter Price: ")
				var price float64
				fmt.Scan(&price)

				h.CreateSale(paymenttype, status, branchId, shopAssistantid, cashierId, clientName, price)

			case "cancel":

				fmt.Println("Enter SaleId which you need cancel: ")
				var id string
				fmt.Scan(&id)

				h.CancelSale(id)

			}

		}
	}

}

type User struct {
	Id   int
	Name string
}

var users []User

func Delete(w http.ResponseWriter, r *http.Request) {

	var user User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error ioutil.ReadAll:", err.Error())
		w.Write([]byte("Internal server error"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("error Unmarshal:", err.Error())
		w.Write([]byte("Internal server error"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for idx, val := range users {
		if val.Id == user.Id {
			users = append(users[:idx], users[idx+1:]...)

		}

	}
} */
