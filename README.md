# Virtual Hid Keyboard - Go
Создайте виртуальную клавиатуру и отправляйте события, такие как нажатие клавиши через файлы сценарий.
Это приложение, написанная на C и Golang. Работает только в Linux из-за зависимостей.

# Виртуальный HID
Эмуляция устройства HID, которой можно управлять удаленно по ssh. Вы можете создать устройство, программировать нажатия клавиш клавиатуры. Поэтому оно идеально подходит для совместного использования клавиатуры и мыши.

Серверные вставки написаны на языке C и в настоящее время работает только в Linux. Но клиентская часть не зависит от платформы, поэтому может быть написана на любой платформе и языке с использованием соединения через сокет TCP.

Поддерживаемые коды событий можно найти здесь: [здесь](https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h)
Посмотреть коды событий при нажатии клавиш в своей системе можно командой _$sudo showkey_

## Установка

```ш
git-клон https://github.com/xela07ax/virtual-hid-keyboard.git
cd virtual-hid-keyboard/
```

### Docs
* Virtual HID over TCP [nmelihsensoy/virtual-hid-tcp](https://github.com/nmelihsensoy/virtual-hid-tcp)
* https://blog.golang.org/c-go-cgo
* https://github.com/golang/go/wiki/cgo
* https://gist.github.com/zchee/b9c99695463d8902cd33
* https://dev.to/mattn/call-go-function-from-c-function-1n3
