%{
package calc
%}

%union{
    val float64
    name string
}

%type <val> expr

%token <val> NUMBER
%token <val> IDENTIFIER
%token LOG
%token POW

%%

prog : expr { result = $1 }
     | prog ';' expr { result = $3 }

expr : NUMBER { $$ = $1 }
     | IDENTIFIER
     ;

%%
