# goexpert-fullcycle-rate-limiter
FullCycle - Pós Go Expert Desafios técnicos - Rate Limiter

## Entregáveis
1. A base do projeto é um sevidor http utilizando o "package" gorilla/mux como roteador para encaminhamento das chamadas http. A principal vantagem
de usar esse pacote é que pode-se anexar com facilidade um middleware que é chamado de maneira encadeada com a chamada da rota solicitada. Com isso
independente da rota que foi geristrada e invocada, o middleware é invocado antes

**Exemplo**:
> *curl http://localhost:8080/* invoca o *middleware RateLimiter* que posteriormente invoca o *handler HelloWord*

2. O *middleware RateLimiter* é um *MiddlewareFunc* que recebe uma *struct* com as configurações existentes no arquivo *.env* e uma *struct* com o
cliente/operador da persistência.

    * Esse *middleware* determina se a request em questão deve ser tratada como uma request por IP ou uma request com API_KEY. Uma vez determinado a chave,
o *usecase RateLimitByKey* é invocado fazendo a análise e o registro da request, devolvendo TRUE para quanto a chave em questão está bloqueada pelas
configurações do rate limiter ou FALSE caso não esteja bloqueada.

3. Para atender o requisito de flexibilidade da persistência, foi implementada uma interface que descreve as assinaturas das funções necessárias para
controle da persistência. Com isso basta criar o "servidor" do tipo desejado no arquivo main.go e utiliza-lo na chamda *handlers.NewRateLimit*

    * Para maiores informações e um exemplo detalhado, favor consultar o arquivo ***main.go***

4. A implementação com o REDIS utiliza o DB0  como base de dados e para cada chave de request, é adicionado ao banco 2 chaves com valores

    1. **{chave}:cnt** => Armazena o contador com a quantidade de chamada realizadas dentro de um segundo
        1. quando essa chave se encontra com o valor -1, significa que essa chave (IP ou API_KEY) está bloqueada para receber requests.
        2. no momento em que a chave de contagem é colocada em -1, o TTL dessa chave é configurado com o valor contido no arquivo de configuração que
        representa o tempo em que esse IP ou API_KEY deve ficar bloqueado
    2. **{chave}:timestamp** => Representa o timestamp em mili segundos referente ao momento em que a primeira request desse contador foi recebida

5. Executando a aplicação:

    1. Clone o repositório com o comando: `git clone https://github.com/felipeksw/goexpert-fullcycle-rate-limiter`
    2. Para iniciar a aplicação, execute o comando na raiz do repositório: `docker-compose up -d`
    3. Para realizar uma chamda ao **ratelimiter** execute o comando:
```sh
curl --request GET \
  --url http://{HOST}:8080/
```



## Requisitos
***Objetivo***: Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

***Descrição***: O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

1. Endereço IP: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
2. Token de Acesso: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
    * API_KEY: <TOKEN>
3. As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

### Requisitos:
* O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
* O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
* O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
* As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
* Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
* O sistema deve responder adequadamente quando o limite é excedido:
    * Código HTTP: 429
    * Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame
* Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
* Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
* A lógica do limiter deve estar separada do middleware.

### Exemplos:
1. Limitação por IP: Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.
2. Limitação por Token: Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.
3. Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.

### Dicas:
* Teste seu rate limiter sob diferentes condições de carga para garantir que ele funcione conforme esperado em situações de alto tráfego.

### Entrega:
* O código-fonte completo da implementação.
* Documentação explicando como o rate limiter funciona e como ele pode ser configurado.
* Testes automatizados demonstrando a eficácia e a robustez do rate limiter.
* Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
* O servidor web deve responder na porta 8080.