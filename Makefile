
tests/cube_deriv.go tests/cube_root.go:
	@go install
	@go generate ./...

test: tests/cube_deriv.go tests/cube_root.go
	@go test ./tests

bench: tests/cube_deriv.go tests/cube_root.go
	@go test ./tests -bench .

clean:
	$(RM) tests/cube_deriv.go tests/cube_root.go