package main

import (
	"branch/config"
	"branch/grpc"
	grpc_client "branch/grpc/client"
	"branch/packages/logger"
	"branch/storage/memory"
	"context"
	"fmt"
	"log"
	"net"
)

func main() {
	cfg := config.Load()
	lg := logger.NewLogger(cfg.Environment, "debug")
	strg, err := memory.NewStorage(context.Background(), *cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	clients, err := grpc_client.New(*cfg)
	if err != nil {
		log.Fatalf("failed to connect to services: %v", err)
	}
	s := grpc.SetUpServer(lg, strg, clients)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

/* import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	pb "main/genproto"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/fake"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStreamServiceServer
	branch []*pb.Branch
}

func (s *server) Create(ctx context.Context, req *pb.CreateBranch) (*pb.CreateBranchResponce, error) {

	id := uuid.NewString()

	s.branch = append(s.branch, &pb.BranchRes{
		Id:      id,
		Name:    req.Name,
		Address: req.Address,
	})

	return &pb.CreateBranchResponce{Id: id}, nil
}

func (s *server) Update(ctx context.Context, req *pb.BranchRes) (*pb.ResponseString, error) {

	for i, v := range s.branch {
		if v.Id == req.Id {
			s.branch[i].Name = req.Name
			s.branch[i].Address = req.Address

			return &pb.ResponseString{Text: "OK"}, nil
		}
	}

	return &pb.ResponseString{Text: ""}, errors.New("not found")

}

func (s *server) Get(ctx context.Context, req *pb.CreateBranchResponce) (*pb.BranchRes, error) {

	for _, v := range s.branch {
		if v.Id == req.GetId() {
			return v, nil
		}
	}

	return &pb.BranchRes{}, errors.New("not Found")

}

func (s *server) Delete(ctx context.Context, req *pb.CreateBranchResponce) (*pb.ResponseString, error) {

	for i, v := range s.branch {
		if v.Id == req.Id {
			if i == (len(s.branch) - 1) {

				s.branch = s.branch[:i]

				return &pb.ResponseString{Text: "deleted suc"}, nil

			} else {

				s.branch = append(s.branch[:i], s.branch[i+1:]...)

				return &pb.ResponseString{Text: "deleted suc"}, nil
			}
		}
	}

	return &pb.ResponseString{Text: ""}, errors.New("not found")

}

func (s *server) GetAll(ctx context.Context, req *pb.GetAllBranchRequest) (resp *pb.GetAllBranchResponse, err error) {
	log.Printf("recieved: %v", req)
	// start := req.Limit * (req.Page - 1)
	// end := start + req.Limit

	response := make([]*pb.Branch, 0)
	fmt.Println(s.branch)
	// if start > int64(len(s.branch)) {

	// 	resp.Branches = []*pb.BranchRes{}
	// 	return resp, nil

	// } else if end > int64(len(s.branch)) {

	// 	return &pb.GetAllBranchResponse{Branches: response}, nil

	// }

	response = append(response, s.branch...)
	return &pb.GetAllBranchResponse{Branches: response}, nil

}

func (s *server) CreateBranch() {
	branch := &pb.BranchRes{
		Id:      uuid.New().String(),
		Name:    fake.Brand(),
		Address: fake.StreetAddress(),
	}
	s.branch = append(s.branch, branch)
}

func (s *server) Count(req *pb.Request, res pb.StreamService_CountServer) error {
	fmt.Println("request:", req.GetNumber())
	var firstNumber, secondNumber, result = 1, 1, 0
	for i := 0; i < int(req.GetNumber()); i++ {

		if i == 0 {
			result = 0
		}
		if i < 3 && i > 0 {
			result = 1
		} else if i != 0 {

			result = firstNumber + secondNumber
			firstNumber = secondNumber
			secondNumber = result

		}

		if result > int(req.GetNumber()) {
			return nil
		}

		err := res.Send(&pb.Response{Count: int32(result)})
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("res sent:", result)
		time.Sleep(time.Second)

	}

	return nil
}

func (s *server) TranslateColor(stream pb.StreamService_TranslateColorServer) error {

	color := make(map[string]string)

	color["qora"] = "Black"
	color["qizil"] = "Red"
	color["oq"] = "White"
	color["yashil"] = "Green"
	color["sariq"] = "Yellow"

	for {
		value, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Println("received color:", value.GetUzbLangColor())

		for key, v := range color {
			if key == value.GetUzbLangColor() {
				if err := stream.Send(&pb.ResponseColor{
					EngLangColor: v,
				}); err != nil {
					return err
				}
			}
		}

	}
}

func main() {

	/* fmt.Println("Start adding...")
	ser := &server{}
	for i := 1; i <= 10; i++ {
		ser.CreateBranch()
		fmt.Printf("created %d\n", i)
	}

	fmt.Println("finished")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBranchServer(s, ser)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// pb.RegisterBranchServiceServer(s, &server{})
	pb.RegisterStreamServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
*/
