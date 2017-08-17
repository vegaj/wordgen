package main

import (
  "fmt"
  "github.com/vegaj/wordgen/gen"
  "time"
)


func main() {
  GEN_WORK := make(chan []string, gen.N_WORKERS)
  COMS := make(chan byte, gen.BATCH)
  WRI_IN := make(chan string, gen.N_WORKERS)

  generator := gen.NewGenerator("source.txt", GEN_WORK)

  w1, w2 := gen.NewWorker(gen.Filter, GEN_WORK, WRI_IN, COMS), gen.NewWorker(gen.Filter, GEN_WORK, WRI_IN, COMS)
  writer := gen.NewWriter("out.txt", WRI_IN, COMS)

  gen.WG.Add(gen.N_WORKERS + 2) //Worker + generator +  writer

  start := time.Now()

  go writer.Write()
  go w1.Work()
  go w2.Work()
  go generator.ExtractAll()

  gen.WG.Wait()
  delta := time.Since(start)
  fmt.Printf("Elapsed %s\n" , delta)
}
