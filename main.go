/*
 * Author: Erdyakov Aleksey Gennadyevich, 2025 year
 * github.com/xela07ax
 * Desc: Virtual HID Keyboard Go for C
 *
 * Input Event codes can be found here
 * https://github.com/nmelihsensoy/virtual-hid-tcp
 * https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h
 *
 * Ex:
 *   Left Click Event -> 11272 -> 1:key press, 1: press val, 272: keycode
 *   X Axis -> 20010 -> 2: pointer, 0: x axis, 010: +10
 *   Y Axis -> 21-10 -> 2: pointer, 1: y axis, -10: -10
 *
 * This is the code side works only on linux because of the dependencies.*
 */

package main

/*
#include <stdlib.h>
#include <unistd.h> //usleep, sleep, close, read, write
#include <stdio.h> //printf, perror
#include <netinet/in.h> //htons, sockaddr_in
#include <string.h> //memset, strcpy, strlen
#include <fcntl.h> //open
#include <stdlib.h> //exit
#include <linux/uinput.h>
#include <errno.h>

int fd;
int bytes_written;
int incoming_code;
short ev_val;
short mode;
char buffer[10];
short neg;

void detach() {
		//Detach HID Device
		ioctl(fd, UI_DEV_DESTROY);
		close(fd);
}

void greet() {
		fd = open("/dev/uinput", O_WRONLY | O_NONBLOCK);
		struct uinput_setup usetup;
		int keys[] = {BTN_LEFT, BTN_RIGHT};

        //Custom key events init
        ioctl(fd, UI_SET_EVBIT, EV_KEY);
        for(int i=0; i<sizeof(keys)/sizeof(int); i++){
            ioctl(fd, UI_SET_KEYBIT, keys[i]);
        }

        //Keyboard init
        for(int i=0; i<227; i++){
            ioctl(fd, UI_SET_KEYBIT, i);
        }

        //Mouse Pointer events init
        ioctl(fd, UI_SET_EVBIT, EV_REL);
        ioctl(fd, UI_SET_RELBIT, REL_X);
        ioctl(fd, UI_SET_RELBIT, REL_Y);

        memset(&usetup, 0, sizeof(usetup));
        usetup.id.bustype = BUS_USB;
        usetup.id.vendor = 0x1357;
		usetup.id.product = 0x5008;
		strcpy(usetup.name, "helloi-go");

		ioctl(fd, UI_DEV_SETUP, &usetup);
		ioctl(fd, UI_DEV_CREATE);
}
// Sends input events
int emit(int type, int code, int val){
    struct input_event data;

    data.type = type;
    data.code = code;
    data.value = val;
    // timestamp values below are ignored
    data.time.tv_sec = 0;
    data.time.tv_usec = 0;

    bytes_written =  write(fd, &data, sizeof(data));
	return bytes_written;
}

int hid_write(int incoming_code, int stat) {
	emit(EV_KEY, incoming_code, stat);
	emit(EV_SYN, SYN_REPORT, stat);
	usleep(5);
	return 1;
}


*/
import "C"
import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.Printf("Virtual HID Device creating device \"helloi-go\"")
	log.Printf("Programmer Erdyakov Aleksey 2025 xela07ax")
	log.Printf("v 1.2.5 - 26.07.2025 12:10 MSK")
	C.greet()
	defer Closer()
	txtch := make(chan [2]string)
	var filepath string
	// Iterate and print individual arguments
	if len(os.Args) > 1 {
		fmt.Println("Scenario file argument:")
		for i, arg := range os.Args[1:] {
			fmt.Printf("Argument %d: %s\n", i+1, arg)
			filepath = arg
		}
	} else {
		filepath = "_scenario.txt"
		log.Println("No additional command-line arguments provided. Use", filepath)
	}
	go func() {
		for msg := range txtch {
			cnt, err := Write(msg[0], msg[1])
			if err != nil {
				log.Printf("error %s|%d", err.Error(), cnt)
				break
			}
		}
	}()
	dat, err := OpenReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	// Поддержка русского
	var src [3]string
	var mdd [3]string
	src[0] = "код клавиши"
	mdd[0] = KeyCode
	src[1] = "нажата"
	mdd[1] = Pressed
	src[2] = "отпущена"
	mdd[2] = Omitted
	var txtfile string = string(dat)
	for i := 0; i < 3; i++ {
		txtfile = strings.ReplaceAll(txtfile, src[i], mdd[i])
		txtfile = strings.ReplaceAll(txtfile, "  ", " ") // если лишние пробелы, можно подкорректировать
	}
	rows := strings.Split(txtfile, "\n")
	for _, row := range rows {
		if row == "" {
			continue
		} else if strings.HasPrefix(row, "//") {
			continue
		}
		log.Println("event-->", row)
		cols := strings.Split(row, " ")
		if cols[0] == KeyCode {
			txtch <- [2]string{cols[1], cols[2]}
		} else if cols[0] == Sleep {
			sptm, err := strconv.Atoi(fmt.Sprintf("%s", cols[1]))
			if err != nil {
				log.Fatalf("error sleep time value %v", err)
			}
			if cols[2] == Second {
				time.Sleep(time.Duration(sptm) * time.Second)
			} else if cols[2] == MicSecond {
				time.Sleep(time.Duration(sptm) * time.Microsecond)
			} else {
				log.Fatalf("sleep time value %v is invalid", sptm)
			}
		} else {
			log.Printf("error unknown command: %s| row: %s", cols[0], row)
			return
		}
	}
	log.Println("Good by")
}

const (
	KeyCode   = "key_code"
	Omitted   = "omitted"
	Pressed   = "pressed"
	Sleep     = "sleep"
	Second    = "second"
	MicSecond = "microsecond"
)

func OpenReadFile(filePath string) (dat []byte, err error) {
	dat, err = ioutil.ReadFile(filePath)
	return
}

func Closer() {
	C.detach()
}

func Write(num, press string) (int, error) {
	key, err := strconv.Atoi(fmt.Sprintf("%s", num))
	if err != nil {
		return -1,
			fmt.Errorf("converting key to int: %s", err.Error())
	}
	var rwTip C.int
	if press == Pressed {
		rwTip = C.int(1)
	}
	res := C.hid_write(C.int(key), C.int(rwTip))
	if res == -1 {
		return int(res),
			fmt.Errorf("failed to write \"Zero buffer/length\"")
	}
	return int(res), nil
}
