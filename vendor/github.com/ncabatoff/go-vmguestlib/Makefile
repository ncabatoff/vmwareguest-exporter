SPEC = spec
PACKAGE = vmguestlib
TEMPLATES = templates

gen-accessor-code:
	PACKAGE=$(PACKAGE) envsubst <$(TEMPLATES)/accessor_header.template >$(PACKAGE)/session_accessor.go
	while read line ; do                                                           \
		eval "$$line";                                                             \
		eval `grep "$$FUNCTION " $(SPEC)/accessors.godoc`;                         \
		FUNCTION=$$FUNCTION                                                        \
		NATIVE_FUNCTION=$$NATIVE_FUNCTION                                          \
		VALUE_TYPE_L=$$VALUE_TYPE_L                                                \
		VALUE_TYPE_U=$$VALUE_TYPE_U                                                \
		COMMENT=$$COMMENT                                                          \
		FUNCTION=$$FUNCTION                                                        \
		envsubst <$(TEMPLATES)/accessor.template >>$(PACKAGE)/session_accessor.go; \
	done < $(SPEC)/accessors.list

gen-accessor-test:
	PACKAGE=$(PACKAGE) envsubst <$(TEMPLATES)/test_header.template >$(PACKAGE)/session_accessor_test.go
	for accessor in `egrep -i "Get[^ ]+\(\)" $(PACKAGE)/session_accessor.go | egrep -o "Get[^(]+"`; do  \
		FUNCTION=$$accessor envsubst <$(TEMPLATES)/test.template >>$(PACKAGE)/session_accessor_test.go; \
	done

gen-accessor: gen-accessor-code gen-accessor-test

test:
	go test -v ./...

.PHONY: gen-accessor gen-accessor-code gen-accessor-test test
