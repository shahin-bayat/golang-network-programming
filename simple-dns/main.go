package main

import "simple-dns/server"

func main() {
	// answers, err := r.Resolve(dnsNama, dnsmessage.TypeA)
	// if err != nil {
	// 	slog.Error("IP not found", "error", err)
	// 	os.Exit(1)
	// }
	// for _, a := range answers {
	// 	fmt.Println("found IP:", a.Body)
	// }
	server := server.NewServer("", 530)
	server.Run()
}
