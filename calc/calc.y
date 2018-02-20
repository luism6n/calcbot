%{
    package calc
%}

%union{
    val float64
    name string
}

%token NUMBER
%token IDENTIFIER

%%

expr : NUMBER
     | IDENTIFIER
     ;

%%
