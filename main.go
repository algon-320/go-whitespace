package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Num ... 内部で扱う整数の型
type Num int64

var commands []*Command // 命令列
var pc int              // どの命令を実行しているか
var callstack Stack     // コールスタック
var stack Stack         // スタック
var heap map[Num]Num    // ヒープ
var labels map[Num]int  // ラベルがどの命令を指しているか
var exitFlag bool

func parse(code string) []*Command {
	parseNum := func(s string) Num {
		n := 0
		sign := rune(s[0])
		p := 1
		for i := len(s) - 1; i > 0; i-- {
			if s[i] == 'T' {
				n += p
			}
			p <<= 1
		}
		if sign == 'T' {
			n *= -1
		}
		return Num(n)
	}

	cnt := 0
	l2N := map[string]Num{}
	labelToNum := func(s string) Num {
		v, ex := l2N[s]
		if ex {
			return v
		}
		l2N[s] = Num(cnt)
		cnt++
		return l2N[s]
	}

	var ret []*Command
	var tmp string
	codeRunes := []rune(code)
	codeST := make([]rune, 0, len(codeRunes))
	for i := 0; i < len(codeRunes); i++ {
		switch codeRunes[i] {
		case ' ':
			codeST = append(codeST, 'S')
		case '\t':
			codeST = append(codeST, 'T')
		case '\n':
			codeST = append(codeST, 'L')
		}
	}

	cur := trieRoot
	for i := 0; i < len(codeST); i++ {
		next := cur.Find(codeST[i])
		if next == nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("[error] `%s` unknown instruction was found", tmp))
			os.Exit(1)
		}
		cur = next

		if cur.ins != nil {
			cmd := &Command{ins: cur.ins, param: -1}
			if cur.ins.ParamType() != NoParam {
				str := ""
				i++
				for codeST[i] != 'L' {
					str += string(codeST[i])
					i++
				}
				if cur.ins.ParamType() == Number {
					cmd.param = parseNum(str)
				} else {
					cmd.param = labelToNum(str)
				}
			}
			if cur.ins.Name() == "DefineLabel" {
				labels[cmd.param] = len(ret)
			} else {
				ret = append(ret, cmd)
			}
			cur = trieRoot
		}
	}
	return ret
}

func execute() {
	fmt.Println("======== output ========")

	if len(commands) > 0 {
		for !exitFlag {
			if pc < 0 || len(commands) <= pc {
				fmt.Fprintln(os.Stderr, fmt.Errorf("[error] program counter out of range"))
				os.Exit(1)
				break
			}

			// ダンプ
			// if commands[pp].ins.ParamType() != NoParam {
			// 	fmt.Println(pp, ":", commands[pp].ins.Name(), commands[pp].param)
			// } else {
			// 	fmt.Println(pp, ":", commands[pp].ins.Name())
			// }

			if err := commands[pc].Do(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
				break
			}
			pc++
		}
	}

	fmt.Println("========================")
	fmt.Println("Stack size :", len(stack))
	fmt.Println("Heap size :", len(heap))
}

func init() {
	heap = make(map[Num]Num)
	labels = make(map[Num]int)
	pc = 0
	exitFlag = false
}

func main() {
	initTrie()

	commandDump := true

	data, err := ioutil.ReadFile(`examples/hworld.ws`)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}
	commands = parse(string(data))

	if commandDump {
		fmt.Println("======== Commands ========")
		for i, v := range commands {
			if v.ins.ParamType() != NoParam {
				fmt.Println(i, ":", v.ins.Name(), v.param)
			} else {
				fmt.Println(i, ":", v.ins.Name())
			}
		}
		fmt.Println("-------- labels --------")
		for i, v := range labels {
			fmt.Println(i, "-->", v)
		}
	}

	execute()
}
