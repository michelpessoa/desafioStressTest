# Go Expert - Desfio Stress Tester

Go | REST 


### Overview do projeto 

A ferramenta foi projetada para realizar testes de carga em um determinado serviço web por meio de uma solução CLI. Onde ao final é apresentado um relatório com as informações da execução. 


### Para executar o projeto via linha de comando siga os seguintes passos:

1. `git clone https://github.com/michelpessoa/desafioStressTest`
2. `go mod tidy` para instalar todas as dependências
3. `execute o seguinte comando`
    - go run main.go --url http://"url a ser testada" --requests 100 --concurrency 1000

4. `deve se preencher os parametros --url --requests --concurrency`
    - url: endereço web de sua preferência
    - requests: quantas requisições na url
    - concurrency: número de chamadas simultâneas

### Para executar o projeto via Docker siga os seguintes passos:

1. `git clone https://github.com/michelpessoa/desafioStressTest`
2. `tenha o docker instalado`
3. `docker build -t stresstester .` para realizar o download da imagem
4. `execute o seguinte comando`
    - docker run stresstester --url=http://"url a ser testada" --requests=100 --concurrency=1000

    - url: endereço web de sua preferência
    - requests: quantas requisições na url
    - concurrency: número de chamadas simultâneas