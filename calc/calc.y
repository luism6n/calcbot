%{
    package calc
%}

%union{
    val float64
    name string
}

%token NUMBER
%token IDENTIFIER
%token LOG
%token POW

%%

expr : NUMBER
     | IDENTIFIER
     ;

%%
