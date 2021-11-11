# generate.makefile

.PHONY: gen.run
gen.run: gen.clean gen.errcode

.PHONY: gen.errcode
gen.errcode: gen.errcode.code gen.errcode.doc

.PHONY: gen.errcode.code
gen.errcode.code:
	@echo "==========> Generating error code Go source files"
	@go run ${ROOT_DIR}/tools/codegen/codegen.go -type=int ${ROOT_DIR}/internal/pkg/code
	@echo "${ROOT_DIR}/internal/pkg/code/code_generated.go"

.PHONY: gen.errcode.doc
gen.errcode.doc:
	@echo "==========> Generating error code documentation"
	@mkdir -p ${ROOT_DIR}/docs/api/
	@go run ${ROOT_DIR}/tools/codegen/codegen.go -type=int -doc \
		-output ${ROOT_DIR}/docs/api/error_code_generated.md ${ROOT_DIR}/internal/pkg/code
	@echo "${ROOT_DIR}/docs/api/error_code_generated.md"

.PHONY: gen.clean
gen.clean:
	@echo "==========> Clean old generated Go source files"
	@${FIND} -type f -name "*_generated.go" -print -delete
