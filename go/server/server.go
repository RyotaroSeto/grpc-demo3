package main

import (
	"context"
	"grpc-demo3/deepthought"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	// 将来 proto ファイルに RPC が追加されてインタフェースが拡張された際、
	// ビルドエラーになるのを防止する仕組み。
	deepthought.UnimplementedComputeServer
}

// インタフェースが実装できていることをコンパイル時に確認するおまじない
var _ deepthought.ComputeServer = &Server{}

// ストリームでデータを送るには stream.Send を使う
// Deep Thought はクライアントからキャンセルしてくるまで、延々とデータを投げ続ける。
func (s *Server) Boot(req *deepthought.BootRequest, stream deepthought.Compute_BootServer) error {
	for {
		select {
		// クライアントがリクエストをキャンセルしたら終了
		case <-stream.Context().Done():
			return nil
		// それ以外は1秒待機してデータを送信
		case <-time.After(1 * time.Second):
		}
		if err := stream.Send(&deepthought.BootResponse{
			Message: "I THINK THEREFORE I AM.",
		}); err != nil {
			return err
		}
	}
}

// "Life", "Universe", "Everything" に関しては750ミリ秒考えて答えを返す。
// もしクライアントがそれよりも短いデッドラインを指定している場合、エラーを返す。
func (s *Server) Infer(ctx context.Context, req *deepthought.InferRequest) (*deepthought.InferResponse, error) {
	switch req.Query {
	case "Life", "Universe", "Everything":
	default:
		// gRPC は共通で使われるエラーコードを定めているので、基本は定義済みのコードを使う
		// https://grpc.github.io/grpc/core/md_doc_statuscodes.html
		return nil, status.Error(codes.InvalidArgument, "Contemplate your query")
	}

	// クライアントがタイムアウトを指定しているかチェック
	deadline, ok := ctx.Deadline()

	// 指定されていない、もしくは十分な時間があれば回答
	if !ok || time.Until(deadline) > 750*time.Millisecond {
		time.Sleep(750 * time.Millisecond)
		return &deepthought.InferResponse{
			Answer:      42,
			Description: []string{"I checked it"},
		}, nil
	}

	// 時間が足りなければ DEADLINE_EXCEEDED (code 4) エラーを返す
	// https://grpc.github.io/grpc/core/md_doc_statuscodes.html
	return nil, status.Error(codes.DeadlineExceeded, "It would take longer")

}
