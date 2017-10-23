package main

import (
	"fmt"
)

func checkStack(l int) error {
	if len(stack) < l {
		return fmt.Errorf("%d : [error] Few elements stacked", pc)
	}
	return nil
}

// -------- スタック操作 ---------------------------------------------------------
func stackPush(x Num) error {
	stack.Push(StackType(x))
	return nil
}
func stackDup() error {
	stack.Push(stack.Top())
	return nil
}
func stackNthCopy(n Num) error { // TODO : 実装
	panic("this method is not implemented!")
	return nil
}
func stackSwap() error {
	if err := checkStack(2); err != nil {
		return err
	}
	x := stack.Top()
	stack.Pop()
	y := stack.Top()
	stack.Pop()
	stack.Push(x)
	stack.Push(y)
	return nil
}
func stackPop() error {
	stack.Pop()
	return nil
}

// -------- 数理演算 -----------------------------------------------------------
func add() error {
	if err := checkStack(2); err != nil {
		return err
	}
	y := stack.Top()
	stack.Pop()
	x := stack.Top()
	stack.Pop()
	stack.Push(x + y)
	return nil
}
func sub() error {
	if err := checkStack(2); err != nil {
		return err
	}
	y := stack.Top()
	stack.Pop()
	x := stack.Top()
	stack.Pop()
	stack.Push(x - y)
	return nil
}
func mul() error {
	if err := checkStack(2); err != nil {
		return err
	}
	y := stack.Top()
	stack.Pop()
	x := stack.Top()
	stack.Pop()
	stack.Push(x * y)
	return nil
}
func div() error {
	if err := checkStack(2); err != nil {
		return err
	}
	y := stack.Top()
	stack.Pop()
	x := stack.Top()
	stack.Pop()
	if y == 0 {
		return fmt.Errorf("%d : [error] Divide by 0", pc)
	}
	stack.Push(x / y)
	return nil
}
func mod() error {
	if err := checkStack(2); err != nil {
		return err
	}
	y := stack.Top()
	stack.Pop()
	x := stack.Top()
	stack.Pop()
	if y == 0 {
		return fmt.Errorf("%d : [error] Divide by 0", pc)
	}
	stack.Push(x % y)
	return nil
}

// -------- ヒープ操作 ----------------------------------------------------------
func heapStore() error {
	if err := checkStack(2); err != nil {
		return err
	}
	value := stack.Top()
	stack.Pop()
	addr := stack.Top()
	stack.Pop()
	heap[Num(addr)] = Num(value)
	return nil
}
func heapLoad() error {
	if err := checkStack(1); err != nil {
		return err
	}
	addr := stack.Top()
	stack.Pop()
	stack.Push(StackType(heap[Num(addr)]))
	return nil
}

// -------- フロー制御 ----------------------------------------------------------
func jump(label Num) error {
	v, ex := labels[label]
	if !ex {
		return fmt.Errorf("%d : [error] no such label", pc)
	}
	pc = v - 1
	return nil
}
func jumpIfZero(label Num) error {
	if err := checkStack(1); err != nil {
		return err
	}
	x := stack.Top()
	stack.Pop()
	if x == 0 {
		return jump(label)
	}
	return nil
}
func jumpIfNeg(label Num) error {
	if err := checkStack(1); err != nil {
		return err
	}
	x := stack.Top()
	stack.Pop()
	if x < 0 {
		return jump(label)
	}
	return nil
}
func callSubroutine(label Num) error {
	callstack.Push(StackType(pc + 1))
	jump(label)
	return nil
}
func endSubroutine() error {
	if len(callstack) == 0 {
		return fmt.Errorf("%d : [error] cannot do return", pc)
	}
	p := callstack.Top()
	callstack.Pop()
	pc = int(p) - 1
	return nil
}
func exit() error {
	exitFlag = true
	return nil
}

// -------- IO -----------------------------------------------------------------
func putChar() error {
	if err := checkStack(1); err != nil {
		return err
	}
	x := stack.Top()
	stack.Pop()
	fmt.Printf("%c", rune(x))
	return nil
}
func putNum() error {
	if err := checkStack(1); err != nil {
		return err
	}
	x := stack.Top()
	stack.Pop()
	fmt.Print(x)
	return nil
}
func getChar() error {
	if err := checkStack(1); err != nil {
		return err
	}
	var c rune
	fmt.Scanf("%c", &c)
	stack.Push(StackType(c))
	heapStore()
	return nil
}
func getNum() error {
	if err := checkStack(1); err != nil {
		return err
	}
	var x int
	fmt.Scanf("%d\n", &x)
	stack.Push(StackType(x))
	heapStore()
	return nil
}

// -----------------------------------------------------------------------------

// ParamType ... パラメータの型
type ParamType int

const (
	// NoParam ... パラメータを取らないことを示す
	NoParam ParamType = iota
	// Number ... パラメータが数値であることを示す
	Number
	// Label ... パラメータがラベルであることを示す
	Label
)

// Instruction ... 命令を定義するためのインターフェース
type Instruction interface {
	Imp() string
	Name() string
	ParamType() ParamType
	Func() interface{}
}

// InstructionBase ... 埋め込み用
type InstructionBase struct {
	imp  string
	name string
}

// Imp ... 命令のコード
func (ib *InstructionBase) Imp() string {
	return ib.imp
}

// Name ... 命令の名前
func (ib *InstructionBase) Name() string {
	return ib.name
}

// Instruction1Param ... 1つのパラメータを取る命令
type Instruction1Param struct {
	InstructionBase
	pType ParamType
	f     func(Num) error
}

// NewInstruction1Param ... Instruction1Paramを生成
func NewInstruction1Param(imp, name string, pType ParamType, f func(Num) error) *Instruction1Param {
	ins := new(Instruction1Param)
	ins.imp = imp
	ins.name = name
	ins.pType = pType
	ins.f = f
	return ins
}

// ParamType ... 命令のパラメータの型を返す
func (ins *Instruction1Param) ParamType() ParamType {
	return ins.pType
}

// Func ... 命令に対応する関数を返す
func (ins *Instruction1Param) Func() interface{} {
	return ins.f
}

// InstructionNoParam ... パラメータを取らない命令
type InstructionNoParam struct {
	InstructionBase
	f func() error
}

// NewInstructionNoParam ... InstructionNoParamを生成
func NewInstructionNoParam(imp, name string, f func() error) *InstructionNoParam {
	ins := new(InstructionNoParam)
	ins.imp = imp
	ins.name = name
	ins.f = f
	return ins
}

// ParamType ... パラメータを取らないのでNoParamを返す
func (ins *InstructionNoParam) ParamType() ParamType {
	return NoParam
}

// Func ... 命令に対応する関数を返す
func (ins *InstructionNoParam) Func() interface{} {
	return ins.f
}

// InstructionList ... 命令一覧
var InstructionList = []Instruction{
	// スタック操作
	NewInstruction1Param("SS", "StackPush", Number, stackPush), // 数値をスタックにプッシュ
	NewInstructionNoParam("SLS", "StackDup", stackDup),         // スタックの1番目を再プッシュ(複製)
	NewInstructionNoParam("SLT", "StackSwap", stackSwap),       // スタックの1番目と2番目を入れ替える
	NewInstructionNoParam("SLL", "StackPop", stackPop),         // スタックの1番目を捨てる
	//NewInstruction1Param("STL", "StackPopN", stackPopN),
	//NewInstruction1Param("STS", "StackNthCopy", stackNthCopy), // スタックのn番目の値をプッシュ

	// 演算
	NewInstructionNoParam("TSSS", "Add", add), // スタックから2つ取り出して和をプッシュ
	NewInstructionNoParam("TSST", "Sub", sub), // スタックから2つ取り出して差をプッシュ
	NewInstructionNoParam("TSSL", "Mul", mul), // スタックから2つ取り出して積をプッシュ
	NewInstructionNoParam("TSTS", "Div", div), // スタックから2つ取り出して商(切り捨て)をプッシュ
	NewInstructionNoParam("TSTT", "Mod", mod), // スタックから2つ取り出して剰余をプッシュ

	// ヒープ操作
	NewInstructionNoParam("TTS", "HeapStore", heapStore), // スタックの一番目の値をスタックの二番目の指す番地にコピー
	NewInstructionNoParam("TTT", "HeapLoad", heapLoad),   // スタックの一番目の指す番地の値をスタックにプッシュ

	// フロー制御
	NewInstruction1Param("LSS", "DefineLabel", Label, func(_ Num) error { return nil }), // ラベルを定義
	NewInstruction1Param("LST", "CallSubroutine", Label, callSubroutine),                // サブルーチン呼び出し
	NewInstruction1Param("LSL", "Jump", Label, jump),                                    // 無条件ジャンプ
	NewInstruction1Param("LTS", "JumpIfZero", Label, jumpIfZero),                        // スタックトップが0ならジャンプ
	NewInstruction1Param("LTT", "JumpIfNeg", Label, jumpIfNeg),                          // スタックトップが負ならジャンプ
	NewInstructionNoParam("LTL", "EndSubroutine", endSubroutine),                        // サブルーチンを抜ける(呼び出し元に戻る)
	NewInstructionNoParam("LLL", "Exit", exit),                                          // プログラムを終了する

	// IO
	NewInstructionNoParam("TLSS", "PutChar", putChar), // スタックトップの値を文字として出力
	NewInstructionNoParam("TLST", "PutNum", putNum),   // スタックトップの値を数値として出力
	NewInstructionNoParam("TLTS", "GetChar", getChar), // 入力を文字としてスタックトップの指す番地(ヒープ)に読み込む
	NewInstructionNoParam("TLTT", "GetNum", getNum),   // 入力を数値としてスタックトップの指す番地(ヒープ)に読み込む
}

// Command ... コマンドの一単位
type Command struct {
	ins   Instruction
	param Num
}

// Do ... コマンドを実行する
func (c *Command) Do() error {
	f := c.ins.Func()
	if f == nil {
		return fmt.Errorf("[error] '%s' command was not defined", c.ins.Name())
	}
	if c.ins.ParamType() != NoParam {
		return f.(func(Num) error)(c.param)
	}
	return f.(func() error)()
}
