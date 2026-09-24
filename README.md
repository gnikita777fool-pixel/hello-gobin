## 1. Что должно получиться

В конце у тебя будет:

```text
GitHub
└── hello-go
    ├── код Go
    ├── GitHub Actions
    └── Releases
        └── v1.0.0
            ├── hello-go-linux-amd64
            ├── hello-go-linux-arm64
            ├── hello-go-darwin-amd64
            ├── hello-go-darwin-arm64
            └── hello-go-windows-amd64.exe
```

Схема работы:

```text
Изменил код
     ↓
git push
     ↓
GitHub Actions
     ↓
gofmt → go vet → go test → go build
     ↓
всё успешно
     ↓
git tag v1.0.0
     ↓
git push origin v1.0.0
     ↓
GitHub Actions
     ↓
5 бинарников
     ↓
GitHub Release v1.0.0
```

Это соответствует структуре задания и его CI/CD-пайплайну. 

---

# 2. Что установить

Если ты работаешь в **Windows + VS Code**, тебе понадобятся:

* VS Code
* Git
* Docker Desktop
* GitHub-аккаунт

**Go устанавливать необязательно**, потому что в задании тестирование можно выполнять через Docker. 

Но Git и Docker должны работать.

Проверь в PowerShell:

```powershell
git --version
```

и:

```powershell
docker --version
```

Если обе команды показывают версии — всё нормально.

---

# 3. Создай папку проекта

Открой **PowerShell**.

Выполни:

```powershell
cd ~
```

Создай папку:

```powershell
mkdir hello-go
```

Перейди в неё:

```powershell
cd hello-go
```

Открой папку в VS Code:

```powershell
code .
```

---

# 4. Создай структуру проекта

В VS Code слева создай:

```text
hello-go/
│
├── .github/
│   └── workflows/
│       └── ci.yml
│
├── greeting/
│   ├── greeting.go
│   └── greeting_test.go
│
├── .gitignore
├── go.mod
└── main.go
```

Именно такую структуру требует задание. 

---

# 5. Создай `go.mod`

Открой:

```text
go.mod
```

Вставь:

```go
module hello-go

go 1.23
```

Сохрани `Ctrl + S`.

---

# 6. Создай `greeting/greeting.go`

В папке `greeting` создай:

```text
greeting.go
```

Вставь:

```go
package greeting

import "fmt"

func Greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func SumRange(from, to int) int {
	sum := 0

	for i := from; i <= to; i++ {
		sum += i
	}

	return sum
}
```

Здесь две функции:

```go
Greet()
```

выводит приветствие.

А:

```go
SumRange()
```

считает сумму чисел от `from` до `to`.

В задании проверяется, что:

```text
1 + 2 + ... + 10 = 55
```



---

# 7. Создай тесты

Создай:

```text
greeting/greeting_test.go
```

Вставь:

```go
package greeting

import "testing"

func TestGreet(t *testing.T) {
	got := Greet("Docker")
	want := "Hello, Docker!"

	if got != want {
		t.Errorf("Greet() = %q, want %q", got, want)
	}
}

func TestSumRange(t *testing.T) {
	got := SumRange(1, 10)
	want := 55

	if got != want {
		t.Errorf("SumRange(1, 10) = %d, want %d", got, want)
	}
}
```

Здесь два теста:

```text
TestGreet
```

проверяет приветствие.

```text
TestSumRange
```

проверяет сумму от 1 до 10.

Именно такие тесты указаны в исходном задании. 

---

# 8. Создай `main.go`

В корне проекта создай:

```text
main.go
```

Вставь:

```go
package main

import (
	"fmt"
	"os"
	"runtime"

	"hello-go/greeting"
)

var version = "dev"

func main() {
	fmt.Printf("hello-go version %s\n", version)
	fmt.Println("Hello from Go! 🐹")

	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Arch: %s\n", runtime.GOARCH)

	fmt.Println(greeting.Greet("GitHub"))

	fmt.Printf("Sum 1..10 = %d\n", greeting.SumRange(1, 10))

	if len(os.Args) > 1 {
		fmt.Println("Аргументы:")

		for i, arg := range os.Args[1:] {
			fmt.Printf("  %d: %s\n", i+1, arg)
		}
	}
}
```

`version` изначально:

```go
var version = "dev"
```

Но позже GitHub Actions автоматически заменит её на:

```text
v1.0.0
```

при сборке через `-ldflags`. 

---

# 9. Создай GitHub Actions

Это **самая важная часть задания**.

Создай:

```text
.github/workflows/ci.yml
```

Вставь туда:

```yaml
name: Go CI/CD

on:
  push:
    branches: [ main ]
    tags: [ 'v*' ]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v7

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'
          cache: true

      - name: Format check
        run: |
          UNFORMATTED=$(gofmt -l .)

          if [ -n "$UNFORMATTED" ]; then
            echo "Следующие файлы не отформатированы:"
            echo "$UNFORMATTED"
            echo "Запустите локально: gofmt -w ."
            exit 1
          fi

      - name: Lint with go vet
        run: go vet ./...

      - name: Run tests
        run: go test ./... -v

      - name: Build
        run: CGO_ENABLED=0 go build -o hello-go .

  release:
    needs: test

    if: startsWith(github.ref, 'refs/tags/v')

    runs-on: ubuntu-latest

    permissions:
      contents: write

    strategy:
      matrix:
        include:
          - goos: linux
            goarch: amd64
            suffix: linux-amd64

          - goos: linux
            goarch: arm64
            suffix: linux-arm64

          - goos: darwin
            goarch: amd64
            suffix: darwin-amd64

          - goos: darwin
            goarch: arm64
            suffix: darwin-arm64

          - goos: windows
            goarch: amd64
            suffix: windows-amd64.exe

    steps:
      - uses: actions/checkout@v7

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Build binary
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
          CGO_ENABLED: 0

        run: |
          go build \
            -ldflags="-s -w -X main.version=${{ github.ref_name }}" \
            -o hello-go-${{ matrix.suffix }} .

      - name: Upload to Release
        uses: softprops/action-gh-release@v2
        with:
          files: hello-go-${{ matrix.suffix }}
          generate_release_notes: true
```

В исходном файле этот workflow разделён на два job:

```text
test
release
```

`release` зависит от `test`, поэтому публикация не произойдёт, если проверки не прошли. 

---

# 10. Создай `.gitignore`

В корне создай:

```text
.gitignore
```

Вставь:

```gitignore
/hello-go
/hello-go-*
.env
.idea/
.vscode/
*.iml
```

Это нужно, чтобы скомпилированные бинарники и настройки IDE не попали в Git. 

---

# 11. Проверь структуру

В VS Code должно быть примерно так:

```text
hello-go
│
├── .github
│   └── workflows
│       └── ci.yml
│
├── greeting
│   ├── greeting.go
│   └── greeting_test.go
│
├── .gitignore
├── go.mod
└── main.go
```

Если у тебя получилось именно так — идём дальше.

---

# 12. Проверяем проект через Docker

Так как Go можно не устанавливать, будем запускать тесты в Docker.

В PowerShell:

```powershell
cd ~/hello-go
```

Запусти:

```powershell
docker run --rm `
  -e GOPATH=/tmp/go `
  -e GOCACHE=/tmp/go-cache `
  -v "${PWD}:/app" `
  -w /app `
  golang:1.23-alpine `
  go test ./... -v
```

Эта команда берёт официальный Go-контейнер и запускает тесты внутри него. Такой способ предусмотрен заданием. 

Если всё правильно, увидишь примерно:

```text
=== RUN   TestGreet
--- PASS: TestGreet (0.00s)

=== RUN   TestSumRange
--- PASS: TestSumRange (0.00s)

PASS
ok      hello-go/greeting
```

Главное:

```text
PASS
```

---

# 13. Создай репозиторий GitHub

Теперь зайди на GitHub.

Создай новый репозиторий:

```text
hello-go
```

**Важно:** при создании репозитория не ставь галочки:

```text
☐ Add a README file
☐ Add .gitignore
☐ Choose a license
```

Репозиторий должен быть пустым.

Это прямо указано в задании, потому что иначе первый `push` может вызвать конфликт. 

---

# 14. Подключи Git

Вернись в PowerShell:

```powershell
cd ~/hello-go
```

Инициализируй Git:

```powershell
git init
```

Добавь файлы:

```powershell
git add .
```

Создай первый commit:

```powershell
git commit -m "Initial commit: Go app with CI/CD to GitHub Releases"
```

Переименуй ветку:

```powershell
git branch -M main
```

---

# 15. Подключи GitHub

Допустим, твой GitHub username:

```text
SherKron
```

Тогда:

```powershell
git remote add origin https://github.com/SherKron/hello-go.git
```

Проверь:

```powershell
git remote -v
```

Должно быть примерно:

```text
origin  https://github.com/SherKron/hello-go.git (fetch)
origin  https://github.com/SherKron/hello-go.git (push)
```

---

# 16. Первый Push

Теперь:

```powershell
git push -u origin main
```

GitHub может попросить авторизацию.

После успешного push открой свой репозиторий на GitHub.

---

# 17. Проверяем GitHub Actions

В репозитории сверху нажми:

**Actions**

Там должен появиться:

```text
Go CI/CD
```

Нажми на него.

Workflow должен запустить:

```text
test
 ├── Format check
 ├── Lint with go vet
 ├── Run tests
 └── Build
```

Поскольку сейчас мы сделали:

```text
git push origin main
```

job:

```text
release
```

запускаться не должен.

Это **нормально**.

В задании Release запускается только для тегов `v*`. 

---

# 18. Если `test` зелёный

Если возле workflow стоит зелёная галочка:

```text
✓
```

значит CI работает.

Теперь можно создавать Release.

---

# 19. Создай тег `v1.0.0`

В PowerShell:

```powershell
git tag v1.0.0
```

Проверь:

```powershell
git tag
```

Должно появиться:

```text
v1.0.0
```

Теперь самое важное:

```powershell
git push origin v1.0.0
```

Именно эта команда запускает Release job. 

---

# 20. Что произойдёт после `git push origin v1.0.0`

GitHub Actions снова запустит:

```text
test
```

После успешного `test` запустится:

```text
release
```

И GitHub Actions параллельно соберёт:

```text
Linux AMD64
Linux ARM64
macOS AMD64
macOS ARM64
Windows AMD64
```

То есть получится 5 файлов. 

---

# 21. Открой Releases

В GitHub зайди:

```text
hello-go
↓
Releases
```

Там должен появиться:

```text
v1.0.0
```

И файлы:

```text
hello-go-linux-amd64
hello-go-linux-arm64
hello-go-darwin-amd64
hello-go-darwin-arm64
hello-go-windows-amd64.exe
```

Именно такой результат должен получиться согласно заданию. 

---

# 22. Проверяем Windows-бинарник

Поскольку ты, судя по работе в VS Code/PowerShell, вероятно выполняешь задание на Windows, скачай:

```text
hello-go-windows-amd64.exe
```

Положи его, например, в:

```text
Downloads
```

В PowerShell перейди туда:

```powershell
cd ~/Downloads
```

Запусти:

```powershell
.\hello-go-windows-amd64.exe
```

Должно получиться примерно:

```text
hello-go version v1.0.0
Hello from Go! 🐹
OS: windows
Arch: amd64
Hello, GitHub!
Sum 1..10 = 55
```

Первая строка особенно важна:

```text
hello-go version v1.0.0
```

Она показывает, что версия была внедрена в бинарник через:

```text
-ldflags
```



---
