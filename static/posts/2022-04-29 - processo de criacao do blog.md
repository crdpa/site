Title: Criando o site do zero - parte 1
Description: Como criei este site utilizando Go no backend.
Date: 2022-04-29
Tags: go, website
---
# Criando o site do zero - parte 1

Resolvi refazer meu site como um projeto de aprendizado, portanto decidi que faria do zero. Decidi que faria o backend em Go, que é a linguagem que estou aprendendo e que tenho mais afinidade no momento, e aproveitei para aprender mais sobre html e css.

## Por onde começar?

Eu já programei utilizando Go antes, mas nunca trabalhei profissionalmente (até agora) com programação. Inclusive um dos meus projetos foi para um processo seletivo para vaga de desenvolvedor. Portanto me senti um pouco perdido em por onde começar essa empreitada, então resolvi documentar o processo.

O que eu quero é bem básico. Um site com informações pessoais minhas para servir de portfólio e um local para escrever sobre assuntos variados do meu interesse.
Meu primeiro passo foi criar um servidor para a página. Algo simples de se implementar em Go.

```go
// Esta é a função de request para ser usado na http.HandleFunc, ela irá direcionar
// o endereço raiz do site para index.html e executar a função que irá gerar os posts
// da página inicial. Também direcionará o endereço crdpa.net/blog para o arquivo blog.html
// e caso haja uma tag no endereço (?tag=), irá extrair e utilizar na função para exibir os
// posts com aquela tag
func httpFunc(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
		case "/", "/index.html":
		executeTemplate(w, "index.html", blogposts.FrontPage(posts))
		return
		case "/blog":
		tag = r.URL.Query().Get("tag")
		executeTemplate(w, "blog.html", blogposts.Archive(posts, tag))
	}
}

// Função para executar os templates das páginas html
func executeTemplate(w http.ResponseWriter, templ string, content interface{}) {
	templates := template.Must(template.ParseGlob("./static/*.html"))
	err := templates.ExecuteTemplate(w, templ, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Aqui iremos servir os arquivos necessários para a página
	stylesheets := http.FileServer(http.Dir("./static/css/"))
	http.Handle("/css/", http.StripPrefix("/css/", stylesheets))
	images := http.FileServer(http.Dir("./static/img/"))
	http.Handle("/img/", http.StripPrefix("/img/", images))

	// Executamos as funções para directionamento e exibição das páginas
	http.HandleFunc("/", httpFunc)
	http.HandleFunc("/blog", httpFunc)

	port := ":8000"
	log.Println("Server is running on port" + port)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
```

Rodando o aplicativo, é só digitar *localhost:8000* no navegador e a página inicial (/static/index.html) irá abrir.

Meu próximo passo seria mostrar uma lista de links para os posts do blog. A estrutura dos posts ficou assim:

```go
type Post struct {
    Title       string
    Description string
    Date        time.Time
    Tags        []string
    Body        string
}
```

O campo *Date* ficou no formato time.Time para eu poder organizar os posts por data com os mais recentes primeiro e também poder exibir a data de várias maneiras possíveis.

Arquivos markdown não contém metadados, porém eu preciso de algo similar. Minha solução foi colocar os metadados no topo do arquivo markdown, antes do conteúdo do post, ficando assim:

```
Title: Título do post
Description: Descrição do post
Date: 2006-01-02
Tags: go, website
---
Conteúdo do post.
```

O programa deve ler a primeira linha e atribuir ao título, a segunda à data e assim por diante. No final, pulará os "---" e o restante será o conteúdo do post.

A função para ler os arquivos e criar a estrutura Post ficou assim:

```go
// Separadores para definir o que cada linha representa
const (
	titleSeparator = "Title: "
	descSeparator  = "Description: "
	dateSeparator  = "Date: "
	tagsSeparator  = "Tags: "
)

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)

	// Função anônima para ler o conteúdo do arquivo removendo os separadores
	readLines := func(tag string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tag)
	}

	// Atribuindo variáveis título, descrição, data, tags e corpo
	title := readLines(titleSeparator)
	desc := readLines(descSeparator)
	date := readLines(dateSeparator)
	tags := strings.Split(readLines(tagsSeparator), ", ")
	body := strings.TrimSuffix(readBody(scanner), "\n")

	// Aqui analiso a data declarando o formato em que ela foi digitada
	// e convertendo para time.Time
	const dateForm = "2006-01-02"
	parsedDate, err := time.Parse(dateForm, date)
	if err != nil {
		return Post{}, nil
	}

	return Post{
		Title:       title,
		Description: desc,
		Date:        parsedDate,
		Tags:        tags,
		Body:        body,
	}, nil
}

// função para ler o conteúdo do post
func readBody(scanner *bufio.Scanner) string {
	scanner.Scan()
	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}

	newBuf := buf.String()
	// blackfriday converte markdown para html
	// https://github.com/russross/blackfriday/tree/v2
	content := blackfriday.Run([]byte(newBuf))
	return string(content)
}
```

Há várias maneiras de se ler os arquivos de uma pasta em Go. Não irei entrar nestes detalhes aqui. Criei uma função chamada NewPostsFromFS que lê a pasta onde estão os arquivos markdown e retorna um slice do struct Post ordenado por data utilizando uma função anônima que ficou assim:

```go
sort.Slice(posts, func(i, j int) bool {
	return posts[i].Date.After(posts[j].Date)
})
```

No próximo post irei mostrar como implementei as tags e filtrei a página de posts para exibir somente os posts com a tag que o visitante seleciona e como implementei os links para cada post.
