
testfiles = tests/cube_deriv.go tests/cube_root.go tests/cube_integral.go tests/cube_math.go

$(testfiles):
	@go install
	@go generate ./...

test: $(testfiles)
	@go test ./tests

bench: $(testfiles)
	@go test ./tests -bench .

clean:
	$(RM) $(testfiles)