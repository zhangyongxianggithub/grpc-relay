package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/zhangyongxianggithub/grpc-relay/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// main main 函数是程序的入口，打印了一条信息，然后创建了一个 gRPC 客户端连接到服务器，并获取了一个 pb.GRPCInferenceServiceClient 类型的客户端对象。
// 使用该客户端对象调用了 ModelStreamInfer 方法，传入了一个 pb.ModelInferRequest 类型的请求参数，包括模型 ID、请求 ID、跟踪 ID、输入内容等信息。
// 如果发生错误，则打印出错误信息；否则，循环接收结果，直到接收到 io.EOF 或者其他错误为止，并打印每次接收到的结果。
func main() {
	fmt.Println("grpc access test")
	conn, err := grpc.NewClient("10.63.41.18:8671",
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}...)
	if err != nil {
		slog.Error(fmt.Sprintf("connect to grpc server failed, reason: %v", err))
		return
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("close connection error: %v", err))
		}
	}(conn)
	client := pb.NewGRPCInferenceServiceClient(conn)
	result, err := client.ModelStreamInfer(context.Background(), &pb.ModelInferRequest{
		ModelId:   "ernie-code2-sft",
		RequestId: "0744a3cc-de7d-4158-bea4-32c24a0e1869",
		TraceId:   "",
		Input:     `{"text":"ERNIE User: 给下面这段代码加上中文的文档注释\njava\n    public void trigger(LocalDate date){\n        RedisLock lock \u003d new RedisLock(LOCK_KEY, UUID.randomUUID().toString(), \n                60L, TimeUnit.MINUTES, stringRedisTemplate); \n        try{\n            if (lock.tryLock()) {\n                log.info(\"RecordSchedule calcRecordTokens start\");\n                processRecord.calcRecordTokens(date);\n            }\n            log.info(\"RecordSchedule calcRecordTokens finished.\");\n        } finally {\n            lock.unlock();\n        }\n    }\n", "top_p":0.7, "temperature":0.2, "penalty_score":1.0, "is_result_all":0, "req_id":"0744a3cc-de7d-4158-bea4-32c24a0e1869", "model_id":"ernie-code2-sft"}`,
		TenantId:  "",
		InputType: 0,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("access model stream infer failed: %v", err))
		return
	}
	for {
		recv, err := result.Recv()
		if err != nil {
			slog.Error(fmt.Sprintf("receive error: %v", err))
			return
		} else {
			fmt.Printf("%v\n\n", recv)
		}
	}

}
