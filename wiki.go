package main

// 표준 라이브러리에서 패키지 가져오기(import)
import (
	"html/template"
	"log"
	"net/http"
	"os"
)

// 페이지 구조체(페이지 제목과 본문)
type Page struct {
	Title string
	Body  []byte
}

// 페이지 저장
func (p *Page) save() error {
	//파일명 변수
	filename := p.Title + ".txt"
	//0600은 읽기-쓰기 권한
	return os.WriteFile(filename, p.Body, 0600)
}

// 페이지 저장
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// 페이지 로드(제목을 입력 받아 페이지 포인터 반환)
func loadPage(title string) (*Page, error) {
	//파일명 변수
	filename := title + ".txt"
	//파일 읽기. error나면 err에 저장
	body, err := os.ReadFile(filename)
	//에러가 존재하면 Page 없이 error 반환
	if err != nil {
		return nil, err
	}
	//제목과 본문으로 구성된 Page 반환
	return &Page{Title: title, Body: body}, nil
}

// 페이지 출력
func viewHandler(w http.ResponseWriter, r *http.Request) {
	//페이지 제목 추출 (파이썬의 슬라이싱 사용)
	title := r.URL.Path[len("/view/"):]
	//페이지 로드
	p, err := loadPage(title)
	//존재하지 않는 페이지 처리
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// 페이지 편집
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// 템플릿 반환
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	//해당 파일의 내용을 읽고 html반환
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//템플릿을 실행하여 HTML 생성
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 메인 함수
func main() {
	// view 경로에 대한 요청 처리 핸들러
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	// 8080 port 실행(오류가 생기면 로그를 남기기 위해 log.Fatal 내부에 작성)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//빌드방법
//$ go build wiki.go
//$ ./wiki
