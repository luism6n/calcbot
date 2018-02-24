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
%token LOG10
%token LOG2
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
     | LOG '(' expr ',' expr ')' { $$ = log($3, $5) }
     | LOG10 '(' expr ')' { $$ = log(10, $3) }
     | LOG2 '(' expr ')' { $$ = log(2, $3) }
     | IDENTIFIER
     ;

%%
