В данной статье мы рассмотрим возможности FFI появившегося в PHP версии 7.4, сравним возможности работать PHP с такими языками как Go, Rust, C++ без создания плагинов, а напрямую, а так-же где возможно вам пригодится использование данной функции, а где не стоит по нашему мнению.

Итак, что такое FFI: [FFI](https://en.wikipedia.org/wiki/Foreign_function_interface)
FFI это возможность вызвать библиотечную функицю написанную на одном языке из другого языка. К примеру, как вы догадываетесь, вызвать из PHP функцию написанную на Rust/C++/Go. Для того, чтобы связать интерпретируемый язык с комиллируемым используется библиотека libffi [Repo](https://en.wikipedia.org/wiki/Libffi). Так как интерпретируемые языки не знают, где конкретко (в каких регистрах) искать параметры вызываемой функции, а так-же, где забирать результаты работы функции после вызова. Всю это работу для интерпретируемых языков делает Libffi. Так-же, эту библиотеку ну нужно устанавливать, она является частью системных библиотек (Linux).
Все эксперименты будут проводиться на ArchLinux (5.6.1 kernel), Libffi 3.2.1.

Для чего это делается. Это конечно интересно, исследовать новые языковые фишки, но есть ли в этом практический смысл. Это я постараюсь доказать по ходу статьи.

Итак, PHP.
[link](https://www.php.net/manual/en/intro.ffi.php)
В самом заголовке сразу описывается, что на момент написания статьи - это эксперементальная особенность языка PHP.
Для нашего примера, мы возьмем такую интересную задачу, как расчем последовательности фибоначчи. И конечно, не самым эффективным способом, через рекурсию. Это сделано для того, чтобы как можно сильнее задействовать процессор, а так-же, чтобы не дать компиллируемым языкам оптимизоровать данную функцию (к примеру, применив технику размотки цикла [link](https://en.wikipedia.org/wiki/Loop_unrolling)

Приступим.
Для PHP первое, что мы должны сделать, это раскомментировать расширение ffi в php.ini (/etc/php/php.ini в ArchLinux).
Далее нам нужно объявить наш условный интерфейс. Есть некоторые ограничения, которые в данный момент присутствуют в PHP FFI, это в частности невозможность использования C-препроцессора (#include, #define и т.д, кроме некоторых специальных)
```php
$ffi = FFI::cdef(
     "int Fib(int n);",
    "/PATH/TO/SO/lib.so");
```

1. `FFI::cdef` - этой операцией мы определяем интерфейс взаимодействия.
2. `int Fib(int n)` - это название экспортируемого метода компиллируемого языка. Чуть ниже мы поговорим как это правильно сделать.Warning
3. `/PATH/TO/SO/lib.so` - путь к динамической библиотеке в которой находится функция выше.

Полный скрипт на php, который мы будем использовать:
```php
<?php
// =========================== PHP NATIVE ===========================
function fib($n)
{
    if ($n === 1 || $n === 2) {
        return 1;
    }
    return fib($n - 1) + fib($n - 2);
}

$start = microtime(true);
$p = 0;
for ($i = 0; $i < 1000000; $i++) {
    $p = fib(12);
}

echo '[PHP] execution time: '.(microtime(true) - $start).' Result: '.$p.PHP_EOL;

// =========================== RUST FFI ===========================
$rust_ffi = FFI::cdef(
    "int Fib(int n);",
    "lib/libphp_rust_ffi.so");

$start = microtime(true);
$r = 0;
for ($i=0; $i < 1000000; $i++) { 
   $r = $rust_ffi->Fib(12);
}

echo '[RUST] execution time: '.(microtime(true) - $start).' Result: '.$r.PHP_EOL;

// =========================== CPP FFI ===========================
$cpp_ffi = FFI::cdef(
    "int Fib(int n);",
    "lib/libphp_cpp_ffi.so");

$start = microtime(true);
$c = 0;
for ($i=0; $i < 1000000; $i++) { 
   $c = $cpp_ffi->Fib(12);
}

echo '[CPP] execution time: '.(microtime(true) - $start).' Result: '.$c.PHP_EOL;

// =========================== GOLANG FFI ===========================
$golang_ffi = FFI::cdef(
    "int Fib(int n);",
    "lib/libphp_go_ffi.so");

$start = microtime(true);

for ($i=0; $i < 1000000; $i++) { 
   $golang_ffi->Fib(12);
}

echo '[GOLANG] execution time: '.(microtime(true) - $start).' Result: '.$c.PHP_EOL;
```


Первый шагом сделаем динамическую библиотеку на языке Rust ([link](https://www.rust-lang.org/))
Для этого потребуется подготовка:
1. На любой платформе, для установки нам потребуется всего лишь одна инструкция отсюда - [link](https://rustup.rs)
2. После этого создадим проект в любом месте командой `cargo new rust_php_ffi`. И все)

Это наша функция:
```rust
//src/lib.rs

#[no_mangle]
extern "C" fn Fib(n: i32) -> i32 {
    if (n == 0) || (n == 1) {
        return 1;
    }

    Fib(n - 1) + Fib(n - 2)
}
```

Очень важно, не забыть добавить аттрибут #[no_mangle] на требуемую функицю, т.к в противном случае комипиллятор заменит имя вашей функции на что-то вроде: `_аgs@fs34`. И экспортируюя ее в PHP, libffi просто не найдет в динамической библиотеке функции с именем Fib. Подробнее можно почитать тут  [link](https://en.wikipedia.org/wiki/Name_mangling).
Так-же в Cargo.toml нужно добавить аттрибут:
```
[lib]
crate-type = ["cdylib"]
```
Хотел бы обратить внимание на то, что есть три варианта динамической библиотеки посредством атрибута в Cargo.toml.
1. dylib - Rust shared library с нестабильным ABI, который может измениться от версии к версии (как и в ГО internal ABI)
2. cdylib - динамическая библиотека для использования в C/C++. Это наш выбор.
3. rlib - Rust static library with rlib extestion (.rlib). Содержит так-же метаданные используемые для линковки различных rlib написанных соответственно на Rust

Компиллируем: `cargo build --release`. И в папке `target/release` видим `.so` файл. Это и будет наша динамическая библиотека.

C++

Далее на очереди C++.
Тут тоже все довольно просто:
```cpp
// in php_cpp_ffi.cpp

int main() {
    
}

extern "C" int Fib(int n) {
    if ((n==1) || (n==2)) {
        return 1;
    }

    return Fib(n-1) + Fib(n-2);
}
```

Нам нужно объявить `extern` функцию для того, чтобы ее можно было импортировать из php. Компилируем:
 `g++ -fPIC -O3 -shared src/php_cpp_ffi.cpp -o ../lib/libphp_cpp_ffi.so`. Нескольно комментариев по компилляции:
 1. `-fPIC` position-independet-code. Для динамической библиотеки важно быть независимой от адреса по которому она загружена в памяти.
 2. `-O3` - максимальная оптимизация

Golang

 И на очереди у нас Golang.
 Язык с рантаймом. Для Го был разработан специальный механизм взаимодействия с динамическими библиотеками, который называется - `CGO` [link](https://golang.org/cmd/cgo/)
Данный комментарий хорошо поясняет, как этот механизм работает: [link](https://github.com/golang/go/blob/860c9c0b8df6c0a2849fdd274a0a9f142cba3ea5/src/cmd/cgo/doc.go#L378-L471)
Так-же, по причине того, что CGO интерпретирует сгенерированные ошибки от C, нет возможности использовать оптимизации, как мы делали это в C++ [link](https://go-review.googlesource.com/c/go/+/23231/) and [link](https://go-review.googlesource.com/c/go/+/23231/2/src/cmd/cgo/gcc.go)

Итак, код в студию:
```go
package main

import (
	"C"
)

// we need to have empty main in package main :)
// because -buildmode=c-shared requires exactly one main package
func main() {

}

//export Fib
func Fib(n C.int) C.int {
	if n == 1 || n == 2 {
		return 1
	}

	return Fib(n-1) + Fib(n-2)
}
```
Итак,все та-же функция Fib, однако, для того, чтобы эта функция была экспортируемой в динамической библиотеке, нам нужно добавить сверху комментарий (эдакий ГО атрибут) `//export Fib`.
Комипиллируем: `go build -o ../lib/libphp_go_ffi.so -buildmode=c-shared`. Так-же обращу внимание, что нам нужно добавить `-buildmode=c-shared` для того, чтобы получилась динамическая библиотека.
На выходе у нас получится 2 файла. Файл с заголовками `.h` и `.so` динамическая библиотека. Файл с заголовками нам по сути не нужен, так как мы знаем имя функции, а FFI php весьма ограничен в работе с С препроцессором.

Запуск ракеты:
После того, как мы все написали (исходные коды предоставлены), мы можем сделать небольшое Makefile чтобы все это собрать (так-же находится в репозитории). После того, как мы вызовем `make build` в папке `lib` появится 4 файла. 2 для ГО (.h/.so) и по одному для Rust и С++.

Makefile:
```makefile
build_cpp:
	echo 'Building cpp...'
	cd cpp && g++ -fPIC -O3 -shared src/php_cpp_ffi.cpp -o libphp_cpp_ffi.so

build_go:
	echo 'Building golang...'
	cd golang && go build -o libphp_go_ffi.so -buildmode=c-shared

build_rust:
	echo 'Building Rust...'
	cargo build --release && mv rust/target/release/libphp_ffi.so libphp_rust_ffi.so

build: build_cpp build_go build_rust


run:
	php php/php_fib.php
```

После чего мы можем перейти в папку `php` и запустить наш скрипт (или чезез Makefile - `make run`). Хочу так-же обратить внимание, что в php скрипте в `FFI::cdef` захардкожены пути к `.so` файлам, поэтому чтобы все сработало, запускайте пожалуйста через `make run`. Результат работы следующий:
1. [PHP] execution time: 8.6763260364532 Result: 144
2. [RUST] execution time: 0.32162690162659 Result: 144
3. [CPP] execution time: 0.3515248298645 Result: 144
4. [GOLANG] execution time: 5.0730509757996 Result: 144

Как и ожидалось, в CPU нагруженных вычислениях PHP показал самый низкий результат, но все-же, в целом довольно быстро для миллиона вызовов.
Сюрпризом может показаться время работы CGO, немногим меньше, чем PHP. По сути, это происходит из-за `calling-conventions` из-за нестабильного ABI. CGO вынужден проводить операции по конвертиции типов из Go-типов в C (можно увидеть в h файле который получается после сборки GO динамической библиотеки) типы, а так-же из-за того, что приходится копировать входящие и возвращаемые значения для совместимости C и GO [link](https://en.wikipedia.org/wiki/X86_calling_conventions).
Rust и С++ показали как и ожидалось лучшие результаты, так как имеют стабильный ABI и единственная прослойка которая была между php и этими языками - это libffi.




Вывод:

Конечно, врят-ли такой подход в данный момент готов к кровавому продакшену, так как может нести в себе много подводных камней. О чем нам и говорят разработчики php:
```
Warning

This extension is EXPERIMENTAL. The behaviour of this extension including the names of its functions and any other documentation surrounding this extension may change without notice in a future release of PHP. This extension should be used at your own risk.
```
Нет нормальной возможности работы с препроцессором. 
Данная статья просто показывает возможности новой фишки языка. Однако, если данная возможность PHP станет стабильной, представьте, как можно будет оптимизировать жаркие места в вашем коде?
