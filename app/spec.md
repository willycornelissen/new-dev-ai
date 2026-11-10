# Especificação: new-dev-ai

## 1. Objetivo & Contexto
* **O que faz**: Aplicação CLI que gera um scaffold para um projeto de desenvolvimento usando IA. A aplicação vai receber o nome do diretório que será criado, contendo o conteúdo do repositório https://github.com/willycornelissen/new-dev-ai menos o diretório app e o arquivo quickstart.md. O Git deve ser inicializado neste novo diretório criado. O novo arquivo README.md deve conter apenas o nome passado como parâmetro como título. 
* **Para quem**: Desenvolver solo
* **O que NÃO faz (Não-Objetivos / Out of Scope)**: Não precisa clonar o diretório app nem o arquivo quickstart.md no novo diretório criado

## 2. Entradas, Saídas e Interfaces (I/O)
* **Entradas**: O nome do diretório a ser criado com o clone de https://github.com/willycornelissen/new-dev-ai.
* **Saídas**: O diretório criado, com o git inicializado e o README.md inicializado com o nome do projeto.
* **Exemplo de Uso**: new-dev-ai hello-world
  ```bash
  ./new-dev-ai <nome-do-projeto>
