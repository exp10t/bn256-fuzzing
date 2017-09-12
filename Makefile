all: bnfuzzer/target/release fuzzer
bnfuzzer/target/release:
	cd bnfuzzer; cargo rustc --release -- -Cpasses=sancov -Cllvm-args="-sanitizer-coverage-level=3  -sanitizer-coverage-trace-pc-guard" --crate-type=staticlib
bn256_instrumented.a:
	go build -buildmode=c-archive ./bn256_instrumented/
fuzzer: fuzzer.c bn256_instrumented.a
	clang fuzzer.c -c -o fuzzer.o
	clang++ fuzzer.o libFuzzer.a bnfuzzer/target/release/deps/libbnfuzzer*.a bn256_instrumented.a -ldl -lpthread -o fuzzer
clean:
	rm -rf fuzzer fuzzer.o bn256_instrumented.a bn256_instrumented.h
	cd bnfuzzer; cargo clean
