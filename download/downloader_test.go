package download

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_NewFile(t *testing.T) {
	file, err := os.Create("/tmp/file")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	fileDl, err := NewFileDl("https://d10.baidupcs.com/file/c98e49c3a3b477332c34cc23ebc88fd7?bkt=p3-1400c98e49c3a3b477332c34cc23ebc88fd7b6fedaef000000341e00&xcode=2d2e45fa551b13775d7ceaea994bffde94d0f479fb285bdc&fid=4185473307-250528-477645975447732&time=1527248264&sign=FDTAXGERQBHSK-DCb740ccc5511e5e8fedcff06b081203-cs9ghSC%2BBq0mvxxH6y2PJPjD3iU%3D&to=d10&size=3415552&sta_dx=3415552&sta_cs=160&sta_ft=dmg&sta_ct=1&sta_mt=1&fm2=MH%2CYangquan%2CAnywhere%2C%2Cshanghai%2Ccmnet&vuk=1751568585&iv=2&newver=1&newfm=1&secfm=1&flow_ver=3&pkey=1400c98e49c3a3b477332c34cc23ebc88fd7b6fedaef000000341e00&expires=8h&rt=sh&r=904924522&mlogid=3358735367544085625&vbdid=457171215&fin=Throng_1.11_xclient.info.dmg&fn=Throng_1.11_xclient.info.dmg&rtype=1&dp-logid=3358735367544085625&dp-callid=0.1.1&hps=1&tsl=0&csl=0&csign=1xTqR5%2B0dDn1R3hqGw3PlqzeuPQ%3D&so=0&ut=1&uter=4&serv=0&uc=2152311193&ic=2834468265&ti=8a6c9448563694cbd6ef8bdcb571c1fbb712ee51eb7fcfa5&by=themis", file, -1)
	if err != nil {
		log.Println(err)
	}

	var exit = make(chan bool)
	var resume = make(chan bool)
	var pause bool
	var wg sync.WaitGroup
	wg.Add(1)
	fileDl.OnStart(func() {
		fmt.Println("download started")
		format := "\033[2K\r%v/%v [%s] %v byte/s %v"
		for {
			status := fileDl.GetStatus()
			var i = float64(status.Downloaded) / 50
			h := strings.Repeat("=", int(i)) + strings.Repeat(" ", 50-int(i))

			select {
			case <-exit:
				fmt.Printf(format, status.Downloaded, fileDl.Size, h, 0, "[FINISH]")
				fmt.Println("\ndownload finished")
				wg.Done()
			default:
				if !pause {
					time.Sleep(time.Second * 1)
					fmt.Printf(format, status.Downloaded, fileDl.Size, h, status.Speeds, "[DOWNLOADING]")
					os.Stdout.Sync()
				} else {
					fmt.Printf(format, status.Downloaded, fileDl.Size, h, 0, "[PAUSE]")
					os.Stdout.Sync()
					<-resume
					pause = false
				}
			}
		}
	})

	fileDl.OnPause(func() {
		pause = true
	})

	fileDl.OnResume(func() {
		resume <- true
	})

	fileDl.OnFinish(func() {
		exit <- true
	})

	fileDl.OnError(func(errCode int, err error) {
		log.Println(errCode, err)
	})

	fmt.Printf("%+v\n", fileDl)

	fileDl.Start()
	time.Sleep(time.Second * 2)
	fileDl.Pause()
	time.Sleep(time.Second * 3)
	fileDl.Resume()
	wg.Wait()
}
