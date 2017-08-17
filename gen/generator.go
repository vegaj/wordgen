package gen

import (
  "bufio"
  "os"
  "fmt"
  "sync"
)

const BATCH int = 5
const WORK_ENDED byte = 0
const N_WORKERS int = 2
var WG sync.WaitGroup

type Generator struct {
  filename string
  out chan []string
  list []string
}

type Worker struct {
  chin chan []string
  chout chan string
  coms chan byte
  Filter FilterFunction
}

type Writer struct {
  list []string
  registeredWorkers int
  ch chan string
  coms chan byte
  filename string
  concluded int
}

type FilterFunction func (str string) bool

func (g *Generator) ExtractAll() {
  f, err := os.Open(g.filename)
  checkFatal(err)
  defer f.Close()

  scanner := bufio.NewScanner(f)
  scanner.Split(bufio.ScanLines)

  var ok bool = true
  for ok {
    g.list = make([]string, BATCH)
    var scanned int = 0
    for ok && scanned < BATCH {
      ok = scanner.Scan()
      g.list[scanned] = scanner.Text()
      scanned++
    }
    if scanned > 0 {
      g.list = g.list[:scanned]
    }
  //fmt.Printf("Vertiendo...")
    g.out <- g.list
  //fmt.Println("Vertido.")
  }
  close(g.out)
  WG.Done()
}

func (w *Worker) Work() {
  for v := range w.chin {
    for i := 0; i < len(v); i++ {
      if w.Filter(v[i]){
        w.chout <- v[i]
        //fmt.Println(v[i])
      }
    }
  }
  w.coms <- WORK_ENDED
  WG.Done()
}

func (w *Writer) Write() {

  go w.supervisor()

  for str := range w.ch {
    w.list = append(w.list, str)
  }
  fmt.Printf("Printing...")
  w.print()
  fmt.Println("Done. See output file")
  WG.Done()
}

//Builders

func NewGenerator(filename string, ch chan []string) *Generator {
  return &Generator {filename: filename, out: ch, list: make([]string, 0)}
}

func NewWorker(filt FilterFunction, in chan []string, out chan string, com chan byte) *Worker {
  return &Worker {chin: in, chout: out, coms : com, Filter : filt}
}

func NewWriter(outfile string, in chan string, coms chan byte) *Writer {
  return &Writer {list: make([]string, 0), ch: in, coms: coms, filename: outfile, concluded: 0}
}

////

func (w *Writer) print() {

  f, err := os.Create(w.filename)
  checkFatal(err)
  defer f.Close()

  writer := bufio.NewWriter(f)
  for i, el := range w.list {
    writer.WriteString(el + "\n")
    if i % BATCH == 0 {
      writer.Flush()
    }
  }
  writer.Flush()
}

func (w *Writer) supervisor() {
  for v := range w.coms {
    if v == WORK_ENDED {
      w.concluded++
      fmt.Println("[SUP] Worker notified its end. ", w.concluded, "/", N_WORKERS, "workers remaining.")
    }
    if w.concluded == N_WORKERS {
      close(w.coms)
      close(w.ch)
    }
  }

}

func Filter(str string) bool {
  return len(str) == 4 && string(str[0]) == "a" && string(str[3]) == "a"
}

//

func checkFatal(err error) {
  if err != nil {
    panic(err)
  }
}