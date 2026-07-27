package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Luxurymann/configman/internal/storage"
)

type stringVarFlag struct {
	val *string
}

func (f *stringVarFlag) Set(value string) error {
	*f.val = value
	return nil
}
func (f *stringVarFlag) String() string { return *f.val }
func parseArgs() (cmd string, prefix string, args []string) {
	flag.Var(&stringVarFlag{val: &prefix}, "prefix", "Префикс для ключей (пример: MYAPP_)")
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("Не указана команда. Используйте: set, get, list, delete, export")
	}
	cmd = flag.Arg(0)
	args = flag.Args()[1:]
	return cmd, prefix, args
}
func Set(st *storage.Storage, args []string) {
	if len(args) < 2 {
		log.Fatal("Не достаточно аргументов")
	}
	key, value := args[0], args[1]
	err := st.Set(key, value)
	if err != nil {
		log.Fatal(err)
	}
}
func Get(st *storage.Storage, args []string) {
	if len(args) < 1 {
		log.Fatal("Не достаточно аргументов")
	}
	key := args[0]
	value, ok := st.Get(key)
	if ok {
		fmt.Println(value)
	} else {
		fmt.Fprintf(os.Stderr, "Не найдено значения по ключу %s\n", key)
	}
}
func List(st *storage.Storage) {
	list := st.List()
	fmt.Printf("Список всех элементов: \n")
	for _, val := range list {
		fmt.Printf("%s\n", val)
	}
}
func Delete(st *storage.Storage, args []string) {
	if len(args) < 1 {
		log.Fatal("Не достаточно аргументов")
	}
	key := args[0]
	err := st.Delete(key)
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	st, err := storage.NewStorage("storage.json")
	if err != nil {
		log.Fatal(err)
	}
	cmd, prefix, args := parseArgs()
	switch cmd {
	case "set":
		Set(st, args)
	case "get":
		Get(st, args)
	case "list":
		List(st)
	case "delete":
		Delete(st, args)
	case "export":
		st.Export(prefix)
	default:
		fmt.Fprintf(os.Stderr, "Не известная команда\n")
	}
}
