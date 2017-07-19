OUTDIR=out
OUTFILE=$(OUTDIR)/cgnet.o

LINUX_HEADERS="/usr/lib/modules/$(shell uname -r)/build"

.PHONY: all build clean

all: $(OUTFILE) bindata.go

$(OUTFILE):
	mkdir -p $(OUTDIR)
	clang -D__KERNEL__ -D__ASM_SYSREG_H \
	-Wno-unused-value -Wno-pointer-sign -Wno-compare-distinct-pointer-types \
			-O2 -emit-llvm -c "src/cgnet.c" \
			$(foreach path,$(LINUX_HEADERS), -I $(path)/arch/x86/include -I $(path)/arch/x86/include/generated -I $(path)/include -I $(path)/include/generated/uapi -I $(path)/arch/x86/include/uapi -I $(path)/include/uapi) \
			-o - | llc -march=bpf -filetype=obj -o $(OUTFILE)

bindata.go:
	go-bindata -pkg bpf out/

clean:
	rm -rf cgnet.o
	rm -rf bindata.go
	rm -rf $(OUTFILE)
