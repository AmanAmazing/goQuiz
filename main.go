package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
    q string 
    a string 
}


func main(){
//    var answers []bool
    var filepath string 
    var timelimit int
    flag.StringVar(&filepath,"csv","","Please enter filepath for csv file that has the following format 'question,answer'")    
    flag.IntVar(&timelimit,"time",30,"Please enter the time in seconds you would like the quiz to run for?")
    flag.Parse();
    if filepath == ""{
        fmt.Println("You did not provide any data");
        os.Exit(1)
    }
    file,err := os.Open(filepath);
    if err != nil {
        fmt.Printf("Error: %s\n",err)
        os.Exit(1)
    }
    defer file.Close()
    
    r := csv.NewReader(file);
    r.FieldsPerRecord = 2;
    lines, err := r.ReadAll()
    if err != nil{
        fmt.Printf("An error occurred while parsing the csv file\nError: %s\n",err)
        os.Exit(2)
    }
    problems := parseLines(lines)
    
    // creating timer 
    timer := time.NewTimer(time.Duration(timelimit) * time.Second)

    correct := 0 
    for i, p := range problems{
        fmt.Printf("Problem %d: %s = \n",i+1, p.q)
        answerCh := make(chan string )
        go func(){
            // creating an annoymous go func since Scanf is blocking 
            var answer string 
            fmt.Scanf("%s\n",&answer)
            answerCh <- answer
        }()
        select {
        case <-timer.C:
            fmt.Printf("\nYou scored %d out of %d.\n",correct, len(problems))
            return 
        case answer := <-answerCh:
            if answer == p.a {
                fmt.Println("Correct!")
                correct++ 
            }
        }
    }

    fmt.Printf("You scored %d out of %d\n", correct, len(problems))

}

func parseLines(lines [][]string) []problem {
    ret := make([]problem, len(lines))
    for i, line := range lines{
        ret[i] = problem{
            q: line[0],
            a: strings.TrimSpace(line[1]), // trims space. Makes matching user input easier.
        }
    }
    return ret 
}
