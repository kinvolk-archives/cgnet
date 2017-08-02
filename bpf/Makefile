OUTDIR=out
OUTFILE=$(OUTDIR)/cgnet.o

.PHONY: all build clean

all: $(OUTFILE) bindata.go

$(OUTFILE):
	mkdir -p $(OUTDIR)
	clang -D__KERNEL__ -D__ASM_SYSREG_H \
	-Wno-unused-value -Wno-pointer-sign -Wno-compare-distinct-pointer-types \
			-fno-stack-protector \
			-O2 -emit-llvm -c "src/cgnet.c" \
			-o - | llc -march=bpf -filetype=obj -o $(OUTFILE)

bindata.go:
	go-bindata -pkg bpf out/

clean:
	rm -rf bindata.go
	rm -rf $(OUTFILE)
