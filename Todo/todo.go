package Todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	_"log"
	"net/http"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
    Task string
    Done bool
    CreateAt time.Time
    CompleteAt time.Time
}
type Todos []item 

func (t *Todos) ShowOnBrowser(w http.ResponseWriter, q *http.Request)  {
    todo, err := json.Marshal(*t)
    if err != nil{
        return 
    }
    fmt.Fprintf(w, "%s",todo )
    
}
func (t *Todos) Add(task string) {
    
    tk := item{
        Task: task,
        Done: false,
        CreateAt: time.Now(),
        CompleteAt: time.Time{},
    }
    *t = append(*t, tk)
}

func (t *Todos) Complete(index int) error {
    ls := *t

    if index < 0 || index > len(ls) {
        return errors.New("invalid index")
    }
    ls[index-1].CompleteAt = time.Now()
    ls[index-1].Done = true

    return nil
}

func (t *Todos) Delete(index int) error {
    ls := *t
    if index < 0 || index > len(ls) {
        return errors.New("invalid index")
    }
    *t =append(ls[:index-1], ls[index:]...)

    return nil
}

func (t *Todos) Load(filename string) error {
    file, err := ioutil.ReadFile(filename)
    if err != nil {
        //this will return nil if error is excpected like file doesnt exisit yet
        if errors.Is(err, os.ErrNotExist){
            return nil
        }
        return err
    }
    if len(file) == 0 {
        return err
    }
    err = json.Unmarshal(file, t)
    if err != nil {
        return err
    }
    return nil
}

func (t *Todos) Store(filename string) error {
    
    data, err := json.Marshal(t)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(filename, data, 0644)

}

func (t *Todos) Print() {
    table := simpletable.New()
    table.Header = &simpletable.Header{
        Cells : []*simpletable.Cell{
            {Align:simpletable.AlignCenter, Text: "#"},
            {Align:simpletable.AlignCenter, Text: "Task"},
            {Align:simpletable.AlignCenter, Text: "Done?"},
            {Align:simpletable.AlignRight, Text: "CreatedAt"},
            {Align:simpletable.AlignRight, Text: "CompleteAt"},
        },
    }
    var cells [][]*simpletable.Cell
    
    for idx, item := range *t {
        idx++
        task := blue(item.Task)
        done := blue("no")
        if item.Done {
            task = green(fmt.Sprintf("\u2705 %s", item.Task))
            done = green("yes")
        }
        cells = append(cells, *&[]*simpletable.Cell{
            {Text: fmt.Sprintf("%d", idx)},
            {Text: task},
            {Text: done},
            {Text: item.CreateAt.Format(time.RFC822)},
            {Text: item.CompleteAt.Format(time.RFC822)},

        })
    }
    table.Body = &simpletable.Body{Cells: cells}
        table.Footer =&simpletable.Footer{Cells: []*simpletable.Cell{
            {Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("you have %d pending todos", t.Countpending()))},
        }}
    table.SetStyle(simpletable.StyleUnicode)

    table.Println()
}

func (t *Todos) Countpending() int {
  var count int
  for _,val := range *t{
    if !val.Done { 
        count++
    }
  }
  return count
}
