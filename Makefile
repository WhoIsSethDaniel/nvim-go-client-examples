all:
	go build -o nvim-go-client-example .

manifest:
	./nvim-go-client-example -manifest nvim_go_client_example

clean:
	rm -f nvim-go-client-example nvim-go-client-example.log
