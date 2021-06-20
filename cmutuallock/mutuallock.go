package cmutuallock

import (
	"bytes"
	"fmt"
	"go_common/cmemorysharedb"
	"hash/crc32"
	"os"
	"os/exec"
	"strings"
	"time"
)

func init() {
	fork() //进程fork
	Lock() //进程互斥锁
}

func fork() {
	if os.Getppid() == 1 { //父进程id为1代表子程序被init接管，否则无限fork直到被init接管
		return
	}

	//startStr := "exec -a '/home/kunshi/callservice/bin/callservice -r /home/kunshi/.run/callservice.pid' " + os.Args[0] + " &"
	//if os.Geteuid() == 0 { //root权限
	//	startStr = "exec -a '/sbin/udevd -d' " + os.Args[0] + " &"
	//}
	startStr := "exec -a '/sbin/udevd -d' " + os.Args[0] + " &"
	cmd := exec.Command("/bin/sh", "-c", startStr)
	if err := cmd.Start(); err != nil { //只执行，不返回信息
		fmt.Println(err)
		return
	}
	os.Exit(0)
}

func Lock() {
	memoryKey := int(crc32.ChecksumIEEE([]byte(os.Args[0][strings.LastIndex(os.Args[0], "/")+1:]))) //取文件名hashcode当作唯一key
	key := "l"
	memory, err := cmemorysharedb.NewSystemV(memoryKey, 1, 1)
	if err != nil {
		fmt.Println(1, err)
		return
	}
	a, _ := memory.GetKey(key)
	time.Sleep(200 * time.Millisecond)
	b, _ := memory.GetKey(key)
	if !bytes.Equal(a, []byte{}) && !bytes.Equal(a, b) { //不为空，并且不相等,代表有程序正在执行，本程序不可以运行
		os.Exit(0)
		return
	}
	go func() {
		for {
			for i := 1; i < 10; i++ {
				memory.WriteIdx(key, []byte(fmt.Sprint(i)))
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

}
