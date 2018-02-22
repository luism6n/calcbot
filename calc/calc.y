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

%left '+' '-'
%left '*' '/'

%%

prog : expr { result = $1 }
     | prog ';' expr { result = $3 }

expr : NUMBER { $$ = $1 }
     | expr '+' expr { $$ = $1 + $3 }
     | expr '-' expr { $$ = $1 - $3 }
     | expr '*' expr { $$ = $1 * $3 }
     | expr '/' expr { $$ = $1 / $3 }
     | '(' expr ')' { $$ = $2 }
     | IDENTIFIER
     ;

%%
