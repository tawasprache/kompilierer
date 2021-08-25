package libgccjit

// #cgo LDFLAGS: -lgccjit
// #include <libgccjit.h>
// #include <stdbool.h>
// #include <stdlib.h>
// gcc_jit_rvalue* jit_sizeof(gcc_jit_context* ctx, gcc_jit_type* t)
// {
//     const void* NULL_PTR = 0;
//
//     gcc_jit_type* t_ptr_type = gcc_jit_type_get_pointer(t);
//     gcc_jit_type* size_type = gcc_jit_context_get_type(ctx, GCC_JIT_TYPE_SIZE_T);
//     gcc_jit_type* byte_type_ptr = gcc_jit_type_get_pointer(gcc_jit_context_get_int_type(ctx, 1, 0));
//
//     gcc_jit_rvalue* one = gcc_jit_context_new_rvalue_from_int(ctx, size_type, 1);
//
//     gcc_jit_rvalue* ptr_0 = gcc_jit_context_new_rvalue_from_ptr(ctx, t_ptr_type, &NULL_PTR);
//     gcc_jit_rvalue* ptr_1 = gcc_jit_lvalue_get_address(gcc_jit_context_new_array_access(ctx, NULL, ptr_0, one), NULL);
//
//     ptr_0 = gcc_jit_context_new_cast(ctx, NULL, ptr_0, byte_type_ptr);
//     ptr_1 = gcc_jit_context_new_cast(ctx, NULL, ptr_1, byte_type_ptr);
//
//     return gcc_jit_context_new_binary_op(ctx, NULL, GCC_JIT_BINARY_OP_MINUS, size_type, ptr_1, ptr_0);
// }
import "C"
import (
	"runtime"
	"unsafe"
)

type Context struct {
	handle   *C.gcc_jit_context
	memoised map[string]*StructType
}

func NewContext() *Context {
	resp := C.gcc_jit_context_acquire()
	ctx := &Context{resp, map[string]*StructType{}}

	runtime.SetFinalizer(ctx, func(c *Context) {
		C.gcc_jit_context_release(c.handle)
	})

	return ctx
}

type IType interface {
	asType() *C.gcc_jit_type

	AsPointer() *PointerType
	AsConst() *ConstType
	AsVolatile() *ConstType
}

type Type struct {
	handle *C.gcc_jit_type
	ctx    *Context
}

func (t *Type) asType() *C.gcc_jit_type {
	return t.handle
}

type Field struct {
	ctx    *Context
	handle *C.gcc_jit_field
	n      string
	k      IType
}

func (c *Context) GetField(name string, kind IType) *Field {
	str := C.CString(name)
	defer C.free(unsafe.Pointer(str))

	f := C.gcc_jit_context_new_field(c.handle, nil, kind.asType(), str)
	return &Field{c, f, name, kind}
}

func (f *Field) Type() IType {
	return f.k
}

func (c *Context) GetType(t BuiltinType) *Type {
	resp := C.gcc_jit_context_get_type(c.handle, uint32(t))
	kind := &Type{resp, c}

	return kind
}

func boolInt(t bool) C.int {
	if t {
		return C.int(1)
	}
	return C.int(0)
}

func (c *Context) GetIntType(numBytes int, signed bool) *Type {
	resp := C.gcc_jit_context_get_int_type(c.handle, C.int(numBytes), boolInt(signed))
	kind := &Type{resp, c}

	return kind
}

type Param struct {
	Type IType
	Name string
}

type FunctionParam struct {
	ctx    *Context
	handle *C.gcc_jit_param
}

var _ ILValue = &FunctionParam{}

func (f *FunctionParam) Kind() IType {
	return f.RValue().Kind()
}

func (f *FunctionParam) RValue() *RValue {
	return &RValue{f.ctx, f.asRValue()}
}

func (f *FunctionParam) LValue() *LValue {
	return &LValue{f.ctx, f.asLValue()}
}

func (f *FunctionParam) Ctx() *Context {
	return f.ctx
}

func (f *FunctionParam) asLValue() *C.gcc_jit_lvalue {
	return C.gcc_jit_param_as_lvalue(f.handle)
}

func (f *FunctionParam) asRValue() *C.gcc_jit_rvalue {
	return C.gcc_jit_param_as_rvalue(f.handle)
}

type LValue struct {
	ctx    *Context
	handle *C.gcc_jit_lvalue
}

var _ ILValue = &LValue{}

func (f *LValue) LValue() *LValue {
	return f
}

func (f *LValue) Kind() IType {
	return f.RValue().Kind()
}

func (f *LValue) RValue() *RValue {
	return &RValue{f.ctx, f.asRValue()}
}

func (f *LValue) asLValue() *C.gcc_jit_lvalue {
	return f.handle
}

func (f *LValue) asRValue() *C.gcc_jit_rvalue {
	return C.gcc_jit_lvalue_as_rvalue(f.handle)
}

func (f *LValue) Ctx() *Context {
	return f.ctx
}

type Function struct {
	ctx    *Context
	handle *C.gcc_jit_function
	blocks []*Block
	ret    IType
}

func (f *Function) ReturnType() IType {
	return f.ret
}

func (f *Function) DumpToDot(at string) {
	nameStr := C.CString(at)
	defer C.free(unsafe.Pointer(nameStr))

	C.gcc_jit_function_dump_to_dot(f.handle, nameStr)
}

func (c *Function) GetParam(idx int) *FunctionParam {
	parm := C.gcc_jit_function_get_param(c.handle, C.int(idx))

	return &FunctionParam{c.ctx, parm}
}

func (c *Function) NewLocal(name string, kind IType) *LValue {
	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))

	hand := C.gcc_jit_function_new_local(c.handle, nil, kind.asType(), nameStr)
	return &LValue{c.ctx, hand}
}

func (c *Context) CreateFunction(name string, kind FunctionType, returnType IType, params []Param, isVaradic bool) *Function {
	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))

	var cparams []*C.gcc_jit_param
	for _, it := range params {
		str := C.CString(it.Name)
		defer C.free(unsafe.Pointer(str))

		cparams = append(cparams, C.gcc_jit_context_new_param(c.handle, nil, it.Type.asType(), str))
	}

	var hand **C.gcc_jit_param
	if len(cparams) > 0 {
		hand = &cparams[0]
	}

	handle := C.gcc_jit_context_new_function(
		c.handle,
		nil,
		uint32(kind),
		returnType.asType(),
		nameStr,
		C.int(len(cparams)),
		hand,
		boolInt(isVaradic),
	)

	return &Function{c, handle, []*Block{}, returnType}
}

type Block struct {
	ctx    *Context
	handle *C.gcc_jit_block
}

func (f *Function) NewBlock(name string) *Block {
	str := C.CString(name)
	defer C.free(unsafe.Pointer(str))

	handle := C.gcc_jit_function_new_block(f.handle, str)

	ret := &Block{f.ctx, handle}

	f.blocks = append(f.blocks, ret)

	return ret
}

func (f *Function) Blocks() []*Block {
	return f.blocks
}

type ICtx interface {
	Ctx() *Context
}

type IRValue interface {
	ICtx
	asRValue() *C.gcc_jit_rvalue
	RValue() *RValue
	Kind() IType
}

type ILValue interface {
	IRValue
	asLValue() *C.gcc_jit_lvalue
	LValue() *LValue
}

type RValue struct {
	ctx    *Context
	handle *C.gcc_jit_rvalue
}

func (r *RValue) Ctx() *Context {
	return r.ctx
}

func (r *RValue) RValue() *RValue {
	return r
}

func (r *RValue) Dereference() *LValue {
	a := C.gcc_jit_rvalue_dereference(r.handle, nil)

	return &LValue{r.ctx, a}
}

func AccessLValueField(f ILValue, v *Field) *LValue {
	fi := C.gcc_jit_lvalue_access_field(f.asLValue(), nil, v.handle)

	return &LValue{f.Ctx(), fi}
}

func AccessRValueField(f IRValue, v *Field) *RValue {
	fi := C.gcc_jit_rvalue_access_field(f.asRValue(), nil, v.handle)

	return &RValue{f.Ctx(), fi}
}

func (r *RValue) asRValue() *C.gcc_jit_rvalue {
	return r.handle
}

func (r *RValue) Kind() IType {
	return &Type{C.gcc_jit_rvalue_get_type(r.handle), r.ctx}
}

func (c *Context) IntRValue(value int, kind IType) *RValue {
	rvalue := C.gcc_jit_context_new_rvalue_from_int(c.handle, kind.asType(), C.int(value))

	return &RValue{c, rvalue}
}

func (c *Context) StringRValue(value string) *RValue {
	str := C.CString(value)
	defer C.free(unsafe.Pointer(str))

	rvalue := C.gcc_jit_context_new_string_literal(c.handle, str)
	return &RValue{c, rvalue}
}

func (c *Context) NewCall(f *Function, args []IRValue) *RValue {
	var cargs []*C.gcc_jit_rvalue

	for _, it := range args {
		cargs = append(cargs, it.asRValue())
	}

	var hand **C.gcc_jit_rvalue
	if len(cargs) > 0 {
		hand = &cargs[0]
	}

	call := C.gcc_jit_context_new_call(c.handle, nil, f.handle, C.int(len(cargs)), hand)
	return &RValue{c, call}
}

func (c *Context) NewCast(of IRValue, to IType) *RValue {
	it := C.gcc_jit_context_new_cast(c.handle, nil, of.asRValue(), to.asType())

	return &RValue{c, it}
}

func (c *Context) SizeOf(t IType) *RValue {
	r := C.jit_sizeof(c.handle, t.asType())

	return &RValue{c, r}
}

func (b *Block) AddEval(r IRValue) {
	C.gcc_jit_block_add_eval(b.handle, nil, r.asRValue())
}

func (b *Block) AddAssign(l ILValue, r IRValue) {
	C.gcc_jit_block_add_assignment(b.handle, nil, l.asLValue(), r.asRValue())
}

func (b *Block) EndWithReturnVoid() {
	C.gcc_jit_block_end_with_void_return(b.handle, nil)
}

func (b *Block) EndWithReturnValue(i IRValue) {
	C.gcc_jit_block_end_with_return(b.handle, nil, i.asRValue())
}

func (b *Block) EndWithConditional(boolVal IRValue, onTrue *Block, onFalse *Block) {
	C.gcc_jit_block_end_with_conditional(b.handle, nil, boolVal.asRValue(), onTrue.handle, onFalse.handle)
}

func (b *Block) EndWithJump(to *Block) {
	C.gcc_jit_block_end_with_jump(b.handle, nil, to.handle)
}

type StructType struct {
	Type

	strct  *C.gcc_jit_struct
	name   string
	fields []*Field
	fmap   map[string]*Field
}

func (s *StructType) Name() string {
	return s.name
}

func (s *StructType) Fields() []*Field {
	return s.fields
}

func (s *StructType) FieldsMap() map[string]*Field {
	return s.fmap
}

func (c *Context) GetStructType(n string, f []*Field) *StructType {
	if len(f) == 0 {
		panic("no can do blank structs")
	}
	if v, ok := c.memoised[n]; ok {
		return v
	}

	m := map[string]*Field{}

	var cfields []*C.gcc_jit_field

	for _, it := range f {
		cfields = append(cfields, it.handle)
		m[it.n] = it
	}

	str := C.CString(n)
	defer C.free(unsafe.Pointer(str))

	kind := C.gcc_jit_context_new_struct_type(c.handle, nil, str, C.int(len(cfields)), &cfields[0])
	kili := C.gcc_jit_struct_as_type(kind)

	es := &StructType{Type{kili, c}, kind, n, f, m}
	c.memoised[n] = es

	return es
}

func (c *Context) CompileToFile(kind OutputType, file string) {
	str := C.CString(file)
	defer C.free(unsafe.Pointer(str))

	C.gcc_jit_context_compile_to_file(c.handle, uint32(kind), str)
}

func (c *Type) AsPointer() *PointerType {
	resp := C.gcc_jit_type_get_pointer(c.handle)
	kind := &PointerType{Type{resp, c.ctx}, c}

	return kind
}

func (c *Type) AsConst() *ConstType {
	resp := C.gcc_jit_type_get_const(c.handle)
	kind := &ConstType{Type{resp, c.ctx}, c}

	return kind
}

func (c *Type) AsVolatile() *ConstType {
	resp := C.gcc_jit_type_get_volatile(c.handle)
	kind := &ConstType{Type{resp, c.ctx}, c}

	return kind
}

type PointerType struct {
	Type

	pointingTo *Type
}

func (c *PointerType) PointingTo() *Type {
	return c.pointingTo
}

type ConstType struct {
	Type

	constOf *Type
}

func (c *ConstType) ConstOf() *Type {
	return c.constOf
}

type VolatileType struct {
	Type

	volatileOf *Type
}

func (c *VolatileType) VolatileOf() *Type {
	return c.volatileOf
}

type BuiltinType uint32

const (
	TypeVoid              BuiltinType = C.GCC_JIT_TYPE_VOID
	TypeVoidPtr           BuiltinType = C.GCC_JIT_TYPE_VOID_PTR
	TypeBool              BuiltinType = C.GCC_JIT_TYPE_BOOL
	TypeChar              BuiltinType = C.GCC_JIT_TYPE_CHAR
	TypeSignedChar        BuiltinType = C.GCC_JIT_TYPE_SIGNED_CHAR
	TypeUnsignedChar      BuiltinType = C.GCC_JIT_TYPE_UNSIGNED_CHAR
	TypeShort             BuiltinType = C.GCC_JIT_TYPE_SHORT
	TypeUnsignedShort     BuiltinType = C.GCC_JIT_TYPE_UNSIGNED_SHORT
	TypeInt               BuiltinType = C.GCC_JIT_TYPE_INT
	TypeUnsignedInt       BuiltinType = C.GCC_JIT_TYPE_UNSIGNED_INT
	TypeLong              BuiltinType = C.GCC_JIT_TYPE_LONG
	TypeUnsignedLong      BuiltinType = C.GCC_JIT_TYPE_UNSIGNED_LONG
	TypeLongLong          BuiltinType = C.GCC_JIT_TYPE_LONG_LONG
	TypeUnsignedLongLong  BuiltinType = C.GCC_JIT_TYPE_UNSIGNED_LONG_LONG
	TypeFloat             BuiltinType = C.GCC_JIT_TYPE_FLOAT
	TypeDouble            BuiltinType = C.GCC_JIT_TYPE_DOUBLE
	TypeLongDouble        BuiltinType = C.GCC_JIT_TYPE_LONG_DOUBLE
	TypeConstCharPtr      BuiltinType = C.GCC_JIT_TYPE_CONST_CHAR_PTR
	TypeSizeT             BuiltinType = C.GCC_JIT_TYPE_SIZE_T
	TypeFilePtr           BuiltinType = C.GCC_JIT_TYPE_FILE_PTR
	TypeComplexFloat      BuiltinType = C.GCC_JIT_TYPE_COMPLEX_FLOAT
	TypeComplexDouble     BuiltinType = C.GCC_JIT_TYPE_COMPLEX_DOUBLE
	TypeComplexLongDouble BuiltinType = C.GCC_JIT_TYPE_COMPLEX_LONG_DOUBLE
)

type FunctionType uint32

const (
	Exported     FunctionType = C.GCC_JIT_FUNCTION_EXPORTED
	Internal     FunctionType = C.GCC_JIT_FUNCTION_INTERNAL
	Imported     FunctionType = C.GCC_JIT_FUNCTION_IMPORTED
	AlwaysInline FunctionType = C.GCC_JIT_FUNCTION_ALWAYS_INLINE
)

type OutputType uint32

const (
	Assembler      OutputType = C.GCC_JIT_OUTPUT_KIND_ASSEMBLER
	ObjectFile     OutputType = C.GCC_JIT_OUTPUT_KIND_OBJECT_FILE
	DynamicLibrary OutputType = C.GCC_JIT_OUTPUT_KIND_DYNAMIC_LIBRARY
	Executable     OutputType = C.GCC_JIT_OUTPUT_KIND_EXECUTABLE
)
